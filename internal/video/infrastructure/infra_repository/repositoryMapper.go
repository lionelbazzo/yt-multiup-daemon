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
	"strings"

	"github.com/michaelcoll/gallery-daemon/internal/video/domain/model"
	"github.com/michaelcoll/gallery-daemon/internal/video/infrastructure/sqlc"
)

func (r *VideoDBRepository) toInfra(video model.Video) (sqlc.CreateOrReplacePhotoParams, error) {
	params := sqlc.CreateOrReplacePhotoParams{
		Hash: video.Hash,
		Path: strings.ReplaceAll(video.Path, r.databaseLocation, ""),
	}

	if err := params.DateTime.Scan(video.DateTime); err != nil {
		return sqlc.CreateOrReplacePhotoParams{}, err
	}
	if video.Iso != 0 {
		if err := params.Iso.Scan(video.Iso); err != nil {
			return sqlc.CreateOrReplacePhotoParams{}, err
		}
	}
	if err := params.ExposureTime.Scan(video.ExposureTime); err != nil {
		return sqlc.CreateOrReplacePhotoParams{}, err
	}
	if err := params.XDimension.Scan(video.XDimension); err != nil {
		return sqlc.CreateOrReplacePhotoParams{}, err
	}
	if err := params.YDimension.Scan(video.YDimension); err != nil {
		return sqlc.CreateOrReplacePhotoParams{}, err
	}
	if err := params.Model.Scan(video.Model); err != nil {
		return sqlc.CreateOrReplacePhotoParams{}, err
	}
	if err := params.FNumber.Scan(video.FNumber); err != nil {
		return sqlc.CreateOrReplacePhotoParams{}, err
	}
	if err := params.Orientation.Scan(video.Orientation); err != nil {
		return sqlc.CreateOrReplacePhotoParams{}, err
	}

	return params, nil
}

func (r *VideoDBRepository) toDomain(video sqlc.Video) (model.Video, error) {

	m := &model.Video{
		Hash: video.Hash,
		Path: video.Path,
	}

	if video.DateTime.Valid {
		m.DateTime = video.DateTime.String
	}
	if video.Iso.Valid {
		m.Iso = int(video.Iso.Int64)
	}
	if video.ExposureTime.Valid {
		m.ExposureTime = video.ExposureTime.String
	}
	if video.XDimension.Valid {
		m.XDimension = int(video.XDimension.Int64)
	}
	if video.YDimension.Valid {
		m.YDimension = int(video.YDimension.Int64)
	}
	if video.Model.Valid {
		m.Model = video.Model.String
	}
	if video.FNumber.Valid {
		m.FNumber = video.FNumber.String
	}
	if video.Orientation.Valid {
		m.Orientation = int(video.Orientation.Int64)
	}

	return *m, nil
}
