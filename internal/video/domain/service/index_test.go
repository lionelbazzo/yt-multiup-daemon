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
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/michaelcoll/gallery-daemon/internal/video/domain/model"
)

const (
	expectedIso          = 80
	expectedDateTime     = "2010-12-27T11:17:34"
	expectedExposureTime = "1/4309"
	expectedXDimension   = 2592
	expectedYDimension   = 1936
	expectedModel        = "iPhone 4"
	expectedFNumber      = "f/2.8"

	expectedSha = "4a1e1af6a9a204a36018c3051f6221e9e000690f"
)

func TestExtractExif(t *testing.T) {
	video := &model.Video{Path: "../../../../test/exif_sample.jpg"}

	err := extractExif(video)
	if err != nil {
		t.Errorf("Error while extracting EXIF data : %v\n", err)
	}

	assert.Equal(t, expectedDateTime, video.DateTime, "Invalid DateTime")
	assert.Equal(t, expectedIso, video.Iso, "Invalid Iso")
	assert.Equal(t, expectedExposureTime, video.ExposureTime, "Invalid ExposureTime")
	assert.Equal(t, expectedXDimension, video.XDimension, "Invalid XDimension")
	assert.Equal(t, expectedYDimension, video.YDimension, "Invalid YDimension")
	assert.Equal(t, expectedModel, video.Model, "Invalid Model")
	assert.Equal(t, expectedFNumber, video.FNumber, "Invalid FNumber")

}

func TestSha(t *testing.T) {
	sha, err := sha("../../../../test/exif_sample.jpg")
	if err != nil {
		t.Errorf("Error while calculating sha : %v\n", err)
	}

	assert.Equal(t, expectedSha, sha, "Invalid Sha")
}
