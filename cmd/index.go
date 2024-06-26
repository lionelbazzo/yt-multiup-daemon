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
	"github.com/michaelcoll/gallery-daemon/internal/video"
	"github.com/michaelcoll/gallery-daemon/internal/video/domain/banner"

	"github.com/spf13/cobra"
)

// indexCmd represents the index command
var indexCmd = &cobra.Command{
	Use:   "index",
	Short: "",
	Long: `
Starts the daemon in index mode only.

Indexes the given folder and create a database file.`,
	Run: func(cmd *cobra.Command, args []string) {
		banner.Print(rootCmd.Version, "", banner.Index)

		service := video.NewForIndex(localDb, folder).GetVideoService()
		defer service.CloseDb()

		service.ReIndex(context.Background(), folder)
	},
}

func init() {
	rootCmd.AddCommand(indexCmd)
}
