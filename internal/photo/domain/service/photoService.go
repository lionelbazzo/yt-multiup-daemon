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
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/cozy/goexif2/exif"
	"github.com/cozy/goexif2/tiff"
	"github.com/schollz/progressbar/v3"

	"github.com/michaelcoll/gallery-daemon/internal/photo/domain/model"
	"github.com/michaelcoll/gallery-daemon/internal/photo/domain/repository"
)

var supportedExtensions = []string{".jpg", ".jpeg", ".JPG", ".JPEG"}

type PhotoService struct {
	r repository.PhotoRepository
}

func New(r repository.PhotoRepository) PhotoService {
	return PhotoService{r: r}
}

// Index scans the folder given in parameter and fill the database with image info and EXIF data found on JPEGs
func (s *PhotoService) Index(ctx context.Context, path string) {

	s.r.Connect(false)
	defer s.r.Close()

	imagesToInsert := make(chan *model.Photo)

	bar := progressbar.Default(-1, fmt.Sprintf("Finding all images in folder %s ... ", path))
	go func() {
		getImageFiles(path, imagesToInsert)
		close(imagesToInsert)
	}()

	for photo := range imagesToInsert {
		err := s.r.CreateOrReplace(ctx, *photo)
		if err != nil {
			log.Fatalf("Can't insert photo into database (%v)", err)
		}
		_ = bar.Add(1)
	}
	_ = bar.Clear()
}

func getImageFiles(path string, imagesToInsert chan *model.Photo) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatalf("Can't open folder : %s (%v)\n", path, err)
	}

	for _, file := range files {
		imagePath := filepath.Join(path, file.Name())
		if file.IsDir() {
			getImageFiles(imagePath, imagesToInsert)
		} else if hasSupportedExtension(file.Name()) {

			photo := &model.Photo{Path: imagePath}

			extractData(photo)
			imagesToInsert <- photo
		}
	}
}

func hasSupportedExtension(filename string) bool {
	for _, ext := range supportedExtensions {
		if strings.HasSuffix(filename, ext) {
			return true
		}
	}

	return false
}

func extractData(photo *model.Photo) {
	hash, err := sha(photo.Path)
	if err != nil {
		log.Printf("Can't calculate hash for file : %s : %v\n", photo.Path, err)
	}

	photo.Hash = hash

	err = extractExif(photo)
	if err != nil {
		log.Printf("Error while extracting EXIF from file %s : %v\n", photo.Path, err)
	}
}

// sha calculate the SHA256 of a file
func sha(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		//log.Fatalf("SHA Can't open file : %s (%v)\n", path, err)
		panic(err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

// extractExif extracts the EXIF data of a photo
func extractExif(photo *model.Photo) error {
	f, err := os.Open(photo.Path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	x, err := exif.Decode(f)
	if err != nil && !errors.Is(err, io.EOF) {
		fmt.Printf("Error reading EXIF data in file : %s (%v)\n", photo.Path, err)
	} else if err == nil {
		err := x.Walk(&walker{p: photo})
		if err != nil {
			panic(err)
		}
	}

	return nil
}

type walker struct {
	p *model.Photo
}

func (w *walker) Walk(name exif.FieldName, tag *tiff.Tag) error {
	if name == "DateTime" {
		dateTimeStr, _ := tag.StringVal()
		date, _ := time.Parse("2006:01:02 15:04:05", dateTimeStr)

		w.p.DateTime = date.Format("2006-01-02T15:04:05")
	} else if name == "ISOSpeedRatings" {
		w.p.Iso, _ = tag.Int(0)
	} else if name == "ExposureTime" {
		if value, err := tag.Rat(0); err == nil {
			w.p.ExposureTime = toString(value)
		}
	} else if name == "PixelXDimension" {
		w.p.XDimension, _ = tag.Int(0)
	} else if name == "PixelYDimension" {
		w.p.YDimension, _ = tag.Int(0)
	} else if name == "Model" {
		value, _ := tag.StringVal()
		w.p.Model = strings.Trim(value, " ")
	} else if name == "FNumber" {
		if ratValue, err := tag.Rat(0); err == nil {
			floatValue, _ := ratValue.Float32()
			w.p.FNumber = fmt.Sprintf("f/%s", strconv.FormatFloat(float64(floatValue), 'f', -1, 32))
		}
	}
	return nil
}

func toString(rat *big.Rat) string {
	return fmt.Sprintf("%d/%d", rat.Num(), rat.Denom())
}