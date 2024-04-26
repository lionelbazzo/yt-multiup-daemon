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

package video

import (
	"github.com/michaelcoll/gallery-daemon/internal/video/domain/model"
	"github.com/michaelcoll/gallery-daemon/internal/video/domain/service"
	"github.com/michaelcoll/gallery-daemon/internal/video/infrastructure/caller"
	"github.com/michaelcoll/gallery-daemon/internal/video/infrastructure/infra_repository"
	"github.com/michaelcoll/gallery-daemon/internal/video/presentation"
)

type Module struct {
	videoService    service.VideoService
	registerService service.RegisterService
	controller      presentation.VideoController
}

func (m Module) GetVideoService() *service.VideoService {
	return &m.videoService
}

func (m Module) GetRegisterService() *service.RegisterService {
	return &m.registerService
}

func (m Module) GetController() *presentation.VideoController {
	return &m.controller
}

func NewForServe(localdb bool, videosPath string, param model.ServeParameters) Module {
	repository := infra_repository.New(localdb, videosPath)
	registerCaller := caller.New(param)

	return Module{
		videoService:    service.New(repository),
		controller:      presentation.New(repository, param),
		registerService: service.NewRegisterService(registerCaller, param),
	}
}
func NewForIndex(localdb bool, videosPath string) Module {
	repository := infra_repository.New(localdb, videosPath)

	return Module{videoService: service.New(repository)}
}
