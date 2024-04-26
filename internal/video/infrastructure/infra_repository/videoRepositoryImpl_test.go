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
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	expectedContentType = "image/jpeg"
)

func TestDetectContentType(t *testing.T) {

	contentType1, _ := detectContentType("exif_sample.jpg")
	contentType2, _ := detectContentType("exif_sample.jpeg")
	contentType3, _ := detectContentType("exif_sample.JPG")
	contentType4, _ := detectContentType("exif_sample.JPEG")

	assert.Equal(t, expectedContentType, contentType1, "Invalid content type")
	assert.Equal(t, expectedContentType, contentType2, "Invalid content type")
	assert.Equal(t, expectedContentType, contentType3, "Invalid content type")
	assert.Equal(t, expectedContentType, contentType4, "Invalid content type")

	_, err := detectContentType("exif_sample.png")
	assert.Errorf(t, err, "content type not supported")
}
