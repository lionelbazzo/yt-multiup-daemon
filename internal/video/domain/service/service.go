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

package service

import (
	"context"
	"log"

	"github.com/michaelcoll/gallery-daemon/internal/video/domain/model"
	"github.com/michaelcoll/gallery-daemon/internal/video/domain/repository"
)

type VideoService struct {
	videoPath *string

	r repository.VideoRepository

	watcherStats *stats
}

func New(r repository.VideoRepository) VideoService {
	return VideoService{r: r, watcherStats: &stats{}}
}

func (s *VideoService) indexVideo(ctx context.Context, imagePath string) *model.Video {
	video := &model.Video{Path: imagePath, Orientation: 1}
	extractData(video)

	if err := s.r.CreateOrReplace(ctx, *video); err != nil {
		log.Fatalf("Can't insert video located at '%s' into database (%v)\n", imagePath, err)
	}

	return video
}

func (s *VideoService) deleteImage(ctx context.Context, imagePath string) {
	if err := s.r.Delete(ctx, imagePath); err != nil {
		log.Fatalf("Can't delete video with path '%s' (%v)\n", imagePath, err)
	}
}

func (s *VideoService) deleteAllImageInPath(ctx context.Context, path string) {
	if err := s.r.DeleteAllVideoInPath(ctx, path); err != nil {
		log.Fatalf("Can't delete all video in path '%s' (%v)\n", path, err)
	}
}

func (s *VideoService) CloseDb() {
	s.r.Close()
}
