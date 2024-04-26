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
	"log"
	"os"
	"path/filepath"
	"strings"
)

// findFiles scans the given path and return a list of folder or files, depending on the value of dirOnly parameter,
// the extensions parameter is a filter for extensions
func findFiles(path string, dirOnly bool, extensions []string) []string {
	files, err := os.ReadDir(path)
	if err != nil {
		log.Fatalf("Can't open folder : %s (%v)\n", path, err)
	}

	var pathList []string

	for _, file := range files {
		filePath := filepath.Join(path, file.Name())
		if file.IsDir() {
			if dirOnly {
				pathList = append(pathList, filePath)
			}
			pathList = append(pathList, findFiles(filePath, dirOnly, extensions)...)
		} else if extensions != nil && hasExtension(file.Name(), extensions) {
			pathList = append(pathList, filePath)
		} else if extensions == nil && !dirOnly {
			pathList = append(pathList, filePath)
		}
	}

	return pathList
}

// hasExtension returns true if the filename has one of the extensions given in parameter.
func hasExtension(filename string, extensions []string) bool {
	for _, ext := range extensions {
		if strings.HasSuffix(filename, ext) {
			return true
		}
	}

	return false
}

func absPath(path string) *string {
	absPath, err := filepath.Abs(path)
	if err != nil {
		log.Fatalf("Can't determine the absolute path of '%s' (%v)\n", path, err)
	}

	return &absPath
}

func isIgnoredFile(path string, ignoredFiles []string) bool {
	for _, ignoredFile := range ignoredFiles {
		if strings.HasSuffix(path, ignoredFile) {
			return true
		}
	}

	return false
}
