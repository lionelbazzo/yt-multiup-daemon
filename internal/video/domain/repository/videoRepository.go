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

package repository

import (
	"context"

	"github.com/michaelcoll/gallery-daemon/internal/video/domain/model"
)

type VideoReader interface {
	ReadChunk([]byte, string) error
}

type VideoRepository interface {
	// Connect Opens a database connection
	Connect(readOnly bool)
	// Close Closes the database connection
	Close()

	CreateOrReplace(context.Context, model.Video) error
	Get(ctx context.Context, hash string) (model.Video, error)
	ReadContent(ctx context.Context, hash string, reader VideoReader) error
	ReadThumbnail(ctx context.Context, hash string, width uint32, height uint32, reader VideoReader) error
	Exists(ctx context.Context, hash string) bool
	List(ctx context.Context, offset uint32, limit uint32) ([]model.Video, error)
	Delete(ctx context.Context, path string) error
	DeleteAllVideoInPath(ctx context.Context, path string) error
	DeleteAll(ctx context.Context) error
	CountVideos(ctx context.Context) (uint32, error)
}
