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

package presentation

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/fatih/color"
	"google.golang.org/grpc"

	"github.com/michaelcoll/gallery-daemon/internal/photo/domain/repository"
	photov1 "github.com/michaelcoll/gallery-proto/gen/proto/go/photo/v1"
)

const port = ":9000"

type PhotoController struct {
	photov1.UnimplementedPhotoServiceServer

	r repository.PhotoRepository
}

func New(r repository.PhotoRepository) PhotoController {
	return PhotoController{r: r}
}

func (c *PhotoController) Serve() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// Creates a new gRPC server
	grpcServer := grpc.NewServer()
	photov1.RegisterPhotoServiceServer(grpcServer, c)

	fmt.Printf("Listening on port %s\n", color.GreenString(port))
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

// GetPhotos returns all photos by given filter
func (c *PhotoController) GetPhotos(ctx context.Context, filter *photov1.GetPhotosRequest) (*photov1.GetPhotosResponse, error) {

	c.r.Connect(true)
	defer c.r.Close()

	list, err := c.r.List(ctx)
	if err != nil {
		return nil, err
	}

	responseList := make([]*photov1.Photo, len(list))
	for i, photo := range list {
		responseList[i] = toGrpc(photo)
	}

	return &photov1.GetPhotosResponse{Photos: responseList}, nil
}

func (c *PhotoController) GetByHash(ctx context.Context, request *photov1.GetByHashRequest) (*photov1.GetByHashResponse, error) {

	c.r.Connect(true)
	defer c.r.Close()

	photo, err := c.r.Get(ctx, request.Hash)
	if err != nil {
		return nil, err
	}

	return &photov1.GetByHashResponse{Photo: toGrpc(photo)}, nil
}

func (c *PhotoController) ContentByHash(filter *photov1.ContentByHashRequest, stream photov1.PhotoService_ContentByHashServer) error {

	c.r.Connect(true)
	defer c.r.Close()

	err := c.r.ReadContent(stream.Context(), filter.Hash, streamReader{stream: stream})
	if err != nil {
		return err
	}

	return nil
}

type streamReader struct {
	stream photov1.PhotoService_ContentByHashServer
}

func (r streamReader) ReadChunk(bytes []byte) error {
	err := r.stream.Send(&photov1.PhotoServiceContentByHashResponse{
		Data: bytes,
	})
	if err != nil {
		return err
	}

	return nil
}