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

package cmd

import (
	"context"
	"fmt"
	"os"
	"regexp"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/michaelcoll/gallery-daemon/internal/video"
	"github.com/michaelcoll/gallery-daemon/internal/video/domain/banner"
	"github.com/michaelcoll/gallery-daemon/internal/video/domain/model"
)

const (
	ownerEnvVarName  = "OWNER"
	daemonEnvVarName = "DAEMON_NAME"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "",
	Long: `
Starts the daemon in server mode.

In this mode it will :
 - index the images if the database is not up-to-date
 - register the daemon to the backend
 - watch for file changes
 - serve backend requests`,
	Run: func(cmd *cobra.Command, args []string) {
		owner := getOwner()
		if !isEmailValid(owner) {
			fmt.Printf("%s Invalid owner email : %s\n", color.RedString("✗"), owner)
			return
		}

		banner.Print(rootCmd.Version, owner, banner.Serve)

		module := video.NewForServe(localDb, folder, model.ServeParameters{
			GrpcPort:      grpcPort,
			ExternalHost:  externalHost,
			DaemonName:    getDaemonName(),
			DaemonVersion: version,
			DaemonOwner:   owner,
		})

		// Indexation
		videoService := module.GetVideoService()
		if reIndex {
			videoService.ReIndex(context.Background(), folder)
		} else {
			videoService.Index(context.Background(), folder)
		}
		videoService.CloseDb()

		// Registration
		go module.GetRegisterService().Register()

		// Watch for file changes
		go videoService.Watch(folder)

		// Serving backend requests
		module.GetController().Serve()
	},
}

var grpcPort uint32
var externalHost string
var name string
var owner string
var reIndex bool

func init() {
	serveCmd.Flags().Uint32VarP(&grpcPort, "port", "p", 9000, "Grpc Port")
	serveCmd.Flags().StringVarP(&externalHost, "external-host", "H", "localhost", "External host")
	serveCmd.Flags().StringVarP(&name, "name", "n", "localhost-daemon", "Daemon name")
	serveCmd.Flags().StringVarP(&owner, "owner", "o", "no@name.com", "Daemon owner email")
	serveCmd.Flags().BoolVar(&reIndex, "re-index", false, "Launch a full re-indexation")

	rootCmd.AddCommand(serveCmd)
}

func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(e)
}

func getOwner() string {
	env, present := os.LookupEnv(ownerEnvVarName)
	if present {
		return env
	}

	return owner
}

func getDaemonName() string {
	env, present := os.LookupEnv(daemonEnvVarName)
	if present {
		return env
	}

	return name
}
