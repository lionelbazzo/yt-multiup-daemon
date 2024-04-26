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
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/michaelcoll/gallery-daemon/internal/video/domain/model"

	"github.com/google/uuid"
)

type RegisterService struct {
	c        RegisterCaller
	param    model.ServeParameters
	daemonId uuid.UUID
	expIn    uint32

	connectionUp bool
	hasConnected bool
	canConnect   bool
}

func NewRegisterService(c RegisterCaller, param model.ServeParameters) RegisterService {
	return RegisterService{c: c, param: param}
}

func (s *RegisterService) Register() {
	for {
		if !s.hasConnected {
			response, err := s.c.Register()
			if err != nil {
				s.connectionProblem(err)
				s.expIn = 5
			} else {
				s.daemonId = response.Id
				s.expIn = response.ExpIn
				s.hasConnected = true
				s.connectionUp = true
				fmt.Printf("%s Daemon registered.\n", color.GreenString("✓"))
			}
		}

		time.Sleep(time.Duration(s.expIn) * time.Second)

		err := s.c.HeartBeat(s.daemonId)
		if err != nil {
			s.connectionProblem(err)
		}
	}
}

func (s *RegisterService) connectionProblem(err error) {
	st, _ := status.FromError(err)
	if s.connectionUp && s.hasConnected && st.Code() != codes.NotFound {
		fmt.Printf("%s Connection lost : %s\n", color.RedString("✗"), getErrorMessage(err))
		s.connectionUp = false
	}
	if !s.hasConnected && !s.connectionUp && !s.canConnect {
		fmt.Printf("%s Can't connect : %s\n", color.RedString("✗"), getErrorMessage(err))
		s.canConnect = true
	}
}

func getErrorMessage(err error) string {
	if strings.Contains(err.Error(), "connection refused") {
		return "connection refused"
	}

	return err.Error()
}
