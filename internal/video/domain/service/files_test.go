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
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var extensions = []string{
	".jpg",
	".JPG",
}

func TestFindFilesDirOnly(t *testing.T) {
	files := findFiles("./../..", true, nil)

	for _, filePath := range files {
		if strings.HasSuffix(filePath, ".go") {
			assert.Fail(t, "Not a dir", filePath)
		}
	}
}

func TestFindFilesFileOnly(t *testing.T) {
	files := findFiles("./..", false, nil)

	for _, filePath := range files {
		if !strings.HasSuffix(filePath, ".go") {
			assert.Fail(t, "Not a file", filePath)
		}
	}
}

func TestFindFilesImageOnly(t *testing.T) {
	files := findFiles("./../../../..", false, extensions)

	for _, filePath := range files {
		if !strings.HasSuffix(filePath, ".jpg") && !strings.HasSuffix(filePath, ".JPG") {
			assert.Fail(t, "Not a JPG", filePath)
		}
	}
}

func TestFindFilesNoFile(t *testing.T) {
	files := findFiles("./", false, extensions)

	assert.Equal(t, 0, len(files), "Should find no file")
}
