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
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gallery-daemon",
	Short: "",
	Long: `
The daemon that index, and stream all the videos to the backend.`,
}

var version = "v0.0.0"

var localDb bool
var verbose bool
var folder string

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Version = version

	rootCmd.PersistentFlags().BoolVar(&localDb, "local-db", false, "Place the database in the current folder")
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "Verbose display")
	rootCmd.PersistentFlags().StringVarP(&folder, "folder", "f", ".", "The folder containing the photos")
}
