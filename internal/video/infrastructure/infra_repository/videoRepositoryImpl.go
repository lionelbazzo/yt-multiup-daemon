/*
 * Copyright (c) 2022 Michaël COLL.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package infra_repository

import (
	"bufio"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/michaelcoll/gallery-daemon/internal/video/domain/consts"
	"github.com/michaelcoll/gallery-daemon/internal/video/domain/model"
	"github.com/michaelcoll/gallery-daemon/internal/video/domain/repository"
	"github.com/michaelcoll/gallery-daemon/internal/video/infrastructure/db"
	"github.com/michaelcoll/gallery-daemon/internal/video/infrastructure/sqlc"
)

const BufferSize = 1024 * 1024 * 2
const webPContentType = "image/webp"

type VideoDBRepository struct {
	repository.VideoRepository

	databaseLocation string

	c *sql.DB
	q *sqlc.Queries
}

func New(localDb bool, videosPath string) *VideoDBRepository {
	var databaseLocation string
	if localDb {
		databaseLocation = "."
	} else {
		databaseLocation = videosPath
	}

	connection := db.Connect(false, databaseLocation)
	db.New(connection).Migrate()

	return &VideoDBRepository{databaseLocation: databaseLocation, q: sqlc.New(), c: connection}
}

func (r *VideoDBRepository) Connect(readOnly bool) {
	r.c = db.Connect(readOnly, r.databaseLocation)
}

func (r *VideoDBRepository) Close() {
	r.c.Close()
}

func (r *VideoDBRepository) CreateOrReplace(ctx context.Context, video model.Video) error {
	params, err := r.toInfra(video)
	if err != nil {
		return err
	}

	if err := r.q.CreateOrReplaceVideo(ctx, r.c, params); err != nil {
		return err
	}

	return nil
}

func (r *VideoDBRepository) Get(ctx context.Context, hash string) (model.Video, error) {
	video, err := r.q.GetVideo(ctx, r.c, hash)
	if err == sql.ErrNoRows {
		return model.Video{}, status.Error(codes.NotFound, "media not found")
	}
	if err != nil {
		return model.Video{}, err
	}
	domain, err := r.toDomain(video)
	if err != nil {
		return model.Video{}, err
	}

	return domain, nil
}

func (r *VideoDBRepository) Exists(ctx context.Context, hash string) bool {
	count, err := r.q.CountVideoByHash(ctx, r.c, hash)
	if err != nil {
		return false
	}

	return count == 1
}

func (r *VideoDBRepository) List(ctx context.Context, offset uint32, limit uint32) ([]model.Video, error) {
	list, err := r.q.List(ctx, r.c, sqlc.ListParams{
		Limit:  int64(limit),
		Offset: int64(offset),
	})
	if err != nil {
		return nil, err
	}

	videos := make([]model.Video, len(list))
	for i, video := range list {
		domain, err := r.toDomain(video)
		if err != nil {
			return nil, err
		}
		videos[i] = domain
	}

	return videos, nil
}

func (r *VideoDBRepository) ReadContent(ctx context.Context, hash string, reader repository.VideoReader) error {
	video, err := r.q.GetVideo(ctx, r.c, hash)
	if err == sql.ErrNoRows {
		return status.Error(codes.NotFound, "media not found")
	}
	if err != nil {
		return err
	}

	contentType, err := detectContentType(video.Path)
	if err != nil {
		return err
	}

	f, err := os.Open(fmt.Sprintf("%s%s", r.databaseLocation, video.Path))
	if err != nil {
		return err
	}
	defer f.Close()

	fReader := bufio.NewReader(f)
	buf := make([]byte, BufferSize)

	for {
		n, err := fReader.Read(buf)

		if err != nil {
			if err != io.EOF {
				return err
			}

			break
		}

		err = reader.ReadChunk(buf[0:n], contentType)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *VideoDBRepository) ReadThumbnail(ctx context.Context, hash string, width uint32, height uint32, reader repository.VideoReader) error {

	var w, h uint32
	if width > 0 && height > 0 {
		w = width
	} else if width == 0 && height == 0 {
		w = 200
	} else {
		w, h = width, height
	}

	thumbnailBytes, err := r.q.GetThumbnail(ctx, r.c, sqlc.GetThumbnailParams{
		Hash:   hash,
		Width:  int64(w),
		Height: int64(h),
	})
	if err == sql.ErrNoRows {
		thumbnailBytes, err = r.createAndUpdateThumbnail(ctx, hash, w, h)
		if err != nil {
			return err
		}
	}
	if err != nil {
		return err
	}

	if err = reader.ReadChunk(thumbnailBytes, webPContentType); err != nil {
		return err
	}

	return nil
}

func (r *VideoDBRepository) createAndUpdateThumbnail(ctx context.Context, hash string, width uint32, height uint32) ([]byte, error) {
	video, err := r.q.GetVideo(ctx, r.c, hash)
	if err == sql.ErrNoRows {
		return nil, status.Error(codes.NotFound, "media not found")
	}

	path := fmt.Sprintf("%s%s", r.databaseLocation, video.Path)
	orientation := getOrientation(video)

	if thumbnail, err := webpEncoder(path, orientation); err != nil {
		return nil, status.Errorf(codes.Internal, "Error while creating the thumbnail of the file %s : %v\n", video.Path, err)
	} else {
		if err := r.createOrReplaceThumbnail(ctx, video.Hash, width, height, thumbnail); err != nil {
			return nil, status.Errorf(codes.Internal, "Error save thumbnail in database (%v).\n", err)
		}

		return thumbnail, nil
	}
}

func getOrientation(video sqlc.Video) uint {
	if video.Orientation.Valid {
		return uint(video.Orientation.Int64)
	}
	return 1
}

func (r *VideoDBRepository) createOrReplaceThumbnail(ctx context.Context, hash string, width uint32, height uint32, thumbnail []byte) error {
	if err := r.q.CreateOrReplaceThumbnail(ctx, r.c, sqlc.CreateOrReplaceThumbnailParams{
		Hash:      hash,
		Width:     int64(width),
		Height:    int64(height),
		Thumbnail: thumbnail,
	}); err != nil {
		return err
	}

	return nil
}

func detectContentType(videoPath string) (string, error) {
	for ext, contentType := range consts.ExtensionsAndContentTypesMap {
		if strings.HasSuffix(videoPath, ext) {
			return contentType, nil
		}
	}

	return "", errors.New("content type not supported")
}

func (r *VideoDBRepository) Delete(ctx context.Context, path string) error {
	err := r.q.DeleteVideoByPath(ctx, r.c, path)
	if err != nil {
		return err
	}

	return nil
}

func (r *VideoDBRepository) DeleteAllVideoInPath(ctx context.Context, path string) error {
	err := r.q.DeleteAllVideoInPath(ctx, r.c, fmt.Sprintf("'%s%%'", strings.ReplaceAll(path, r.databaseLocation, "")))
	if err != nil {
		return err
	}

	return nil
}

func (r *VideoDBRepository) DeleteAll(ctx context.Context) error {
	err := r.q.DeleteAllVideos(ctx, r.c)
	if err != nil {
		return err
	}

	return nil
}

func (r *VideoDBRepository) CountVideos(ctx context.Context) (uint32, error) {
	count, err := r.q.CountVideos(ctx, r.c)
	if err != nil {
		return 0, err
	}

	return uint32(count), nil
}
