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
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/fatih/color"
	"github.com/schollz/progressbar/v3"

	"github.com/michaelcoll/gallery-daemon/internal/video/domain/consts"
	"github.com/michaelcoll/gallery-daemon/internal/video/domain/model"
)

// Index tests if the number of images present in the database is different from the number present in the filesystem.
// If this is the case it launches a re-indexation otherwise it do nothing
func (s *VideoService) Index(ctx context.Context, path string) {
	files := findFiles(path, false, consts.SupportedExtensions)

	count, err := s.r.CountVideos(ctx)
	if err != nil {
		log.Fatalf("Can't read the count of videos in the database (%v)\n", err)
	}
	if len(files) != int(count) {
		s.ReIndex(ctx, path)
	} else {
		fmt.Printf("%s Up-to-date. \n", color.GreenString("✓"))
	}
}

// ReIndex scans the folder given in parameter and fill the database with image info and EXIF data found on JPEGs
func (s *VideoService) ReIndex(ctx context.Context, path string) {
	fmt.Printf("%s Re-indexing folder %s \n", color.GreenString("✓"), color.GreenString(path))

	s.videoPath = absPath(path)

	bar := progressbar.Default(-1, "Clearing database... ")
	err := s.r.DeleteAll(ctx)
	if err != nil {
		log.Fatalf("Can't delete all videos in the database (%v)\n", err)
	}
	_ = bar.Clear()

	bar = progressbar.Default(-1, "Searching all the images... ")
	var indexedImages []*model.Video
	// Find images in the folder
	for _, imagePath := range findFiles(*s.videoPath, false, consts.SupportedExtensions) {
		indexedImages = append(indexedImages, s.indexVideo(ctx, imagePath))
		_ = bar.Add(1)
	}
	_ = bar.Clear()

	fmt.Printf("%s Done. \n", color.GreenString("✓"))

	printIndexationStats(indexedImages)
}

func extractData(video *model.Video) {
	hash, err := sha(video.Path)
	if err != nil {
		log.Printf("Can't calculate hash for file : %s : %v\n", video.Path, err)
	}

	video.Hash = hash

	if err = extractExif(video); err != nil {
		log.Printf("Error while extracting EXIF from file %s : %v\n", video.Path, err)
	}

	fixVideoAttributes(video)
}

// sha calculate the SHA256 of a file
func sha(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		//log.Fatalf("SHA Can't open file : %s (%v)\n", path, err)
		panic(err)
	}
	defer f.Close()

	h := sha1.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

func printIndexationStats(indexedImages []*model.Video) {
	var mapIndexedImages = make(map[string][]string)

	for _, image := range indexedImages {
		paths, exists := mapIndexedImages[image.Hash]
		if exists {
			paths = append(paths, image.Path)
			mapIndexedImages[image.Hash] = paths
		} else {
			mapIndexedImages[image.Hash] = []string{image.Path}
		}
	}

	for _, paths := range mapIndexedImages {
		if len(paths) > 1 {
			fmt.Println("Duplicated images detected :")
			for i, path := range paths {
				if i == len(paths)-1 {
					fmt.Printf(" using => %s\n", path)
				} else {
					fmt.Printf("          %s\n", path)
				}
			}
		}
	}
}
