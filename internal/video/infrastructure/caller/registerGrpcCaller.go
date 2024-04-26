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

package caller

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	daemonv1 "github.com/michaelcoll/gallery-proto/gen/proto/go/daemon/v1"

	"github.com/michaelcoll/gallery-daemon/internal/video/domain/model"
)

const (
	webHost = "localhost"
	webPort = 9000
)

type RegisterGrpcCaller struct {
	param model.ServeParameters
}

func New(param model.ServeParameters) *RegisterGrpcCaller {
	return &RegisterGrpcCaller{param: param}
}

func (c *RegisterGrpcCaller) Register() (*model.RegisterResponse, error) {

	client, conn := createClient(webHost, webPort)
	defer closeConnection(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.Register(ctx, &daemonv1.RegisterRequest{
		DaemonName:    c.param.DaemonName,
		DaemonHost:    c.param.ExternalHost,
		DaemonPort:    c.param.GrpcPort,
		DaemonVersion: c.param.DaemonVersion,
		DaemonOwner:   c.param.DaemonOwner,
	})
	if err != nil {
		return nil, err
	}

	id, err := uuid.Parse(resp.Uuid)
	if err != nil {
		return nil, err
	}

	return &model.RegisterResponse{
		Id:    id,
		ExpIn: resp.ExpIn,
	}, nil
}

func (c *RegisterGrpcCaller) HeartBeat(id uuid.UUID) error {
	client, conn := createClient(webHost, webPort)
	defer closeConnection(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := client.HeartBeat(ctx, &daemonv1.HeartBeatRequest{
		Uuid:          id.String(),
		DaemonName:    c.param.DaemonName,
		DaemonHost:    c.param.ExternalHost,
		DaemonPort:    c.param.GrpcPort,
		DaemonVersion: c.param.DaemonVersion,
		DaemonOwner:   c.param.DaemonOwner,
	})
	if err != nil {
		return err
	}

	return nil
}

func createClient(webHost string, webPort int) (daemonv1.DaemonServiceClient, *grpc.ClientConn) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	daemonAddr := fmt.Sprintf("%s:%d", webHost, webPort)

	conn, err := grpc.Dial(daemonAddr, opts...)
	if err != nil {
		log.Fatalf("fail to contact the daemon : %v", err)
	}
	client := daemonv1.NewDaemonServiceClient(conn)

	return client, conn
}

func closeConnection(conn *grpc.ClientConn) {
	err := conn.Close()
	if err != nil {
		log.Fatalf("fail to close the daemon connection : %v", err)
	}
}
