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

package consts

var ExtensionsAndContentTypesMap = map[string]string{
	".avi":  "video/x-msvideo",
	".AVI":  "video/x-msvideo",
	".mkv":  "video/x-matroska",
	".MKV":  "video/x-matroska",
	".mp4":  "video/mp4",
	".MP4":  "video/mp4",
	".mpeg": "video/mpeg",
	".MPEG": "video/mpeg",
	".webm": "video/webm",
	".WEBM": "video/webm",
}

var SupportedExtensions = []string{
	".avi",
	".AVI",
	".mkv",
	".MKV",
	".mp4",
	".MP4",
	".mpeg",
	".MPEG",
	".webm",
	".WEBM",
}

var IgnoredFiles = []string{
	DatabaseName,
	DatabaseName + "-shm",
	DatabaseName + "-wal",
}

const (
	DatabaseName = "videos.db"
)
