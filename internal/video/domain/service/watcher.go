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
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/fatih/color"
	"github.com/fsnotify/fsnotify"

	"github.com/michaelcoll/rfsnotify"

	"github.com/michaelcoll/gallery-daemon/internal/video/domain/consts"
)

type stats struct {
	inserted      atomic.Int64
	deleted       atomic.Int64
	deletedFolder atomic.Int64
}

func (s *VideoService) Watch(path string) {

	s.r.Connect(false)
	defer s.r.Close()

	watcher, err := rfsnotify.NewBufferedWatcher(2000)
	if err != nil {
		log.Fatalf("Could not create the watcher : %v", err)
	}
	defer watcher.Close()

	err = watcher.AddRecursive(path, nil)
	if err != nil {
		log.Fatalf("Could not add the folder : %v", err)
	}

	fmt.Printf("%s Watching folder %s \n", color.GreenString("✓"), color.GreenString(path))

	go s.displayStats()

	quit := make(chan os.Signal, 1) // we need to reserve to buffer size 1, so the notifier are not blocked
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case <-quit:
			fmt.Printf("%s Stoping watcher ...\n", color.RedString("!"))
			s.r.Close()
			os.Exit(1)

		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			if isCreateEvent(event) || isDeleteEvent(event) {
				go s.handleEvent(event)
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("error:", err)
		}
	}

}

func isCreateEvent(event fsnotify.Event) bool {
	return event.Op == fsnotify.Write
}

func isDeleteEvent(event fsnotify.Event) bool {
	return event.Op == fsnotify.Rename ||
		event.Op == fsnotify.Remove
}

func (s *VideoService) handleEvent(event fsnotify.Event) {
	if isIgnoredFile(event.Name, consts.IgnoredFiles) {
		return
	}

	if isCreateEvent(event) && hasExtension(event.Name, consts.SupportedExtensions) {
		s.indexVideo(context.Background(), event.Name)
		s.watcherStats.inserted.Add(1)
	} else if isDeleteEvent(event) && hasExtension(event.Name, consts.SupportedExtensions) {
		s.deleteImage(context.Background(), event.Name)
		s.watcherStats.deleted.Add(1)
	} else if isDeleteEvent(event) {
		s.deleteAllImageInPath(context.Background(), event.Name)
		s.watcherStats.deletedFolder.Add(1)
	}
}

func (s *VideoService) displayStats() {
	var inserted, deleted, deletedFolder int64

	for {
		time.Sleep(time.Duration(2) * time.Second)
		if inserted != s.watcherStats.inserted.Load() ||
			deleted != s.watcherStats.deleted.Load() ||
			deletedFolder != s.watcherStats.deletedFolder.Load() {

			deltaInsert := s.watcherStats.inserted.Load() - inserted
			deltaDelete := s.watcherStats.deleted.Load() - deleted
			deltaDeleteFolder := s.watcherStats.deletedFolder.Load() - deletedFolder

			if deltaInsert > 0 && deltaDelete > 0 {
				fmt.Printf("%s Indexed %d image(s) and deleted %d image(s)\n", color.GreenString("!"), deltaInsert, deltaDelete)
			} else if deltaInsert > 0 {
				fmt.Printf("%s Indexed %d image(s)\n", color.GreenString("!"), deltaInsert)
			} else if deltaDelete > 0 {
				fmt.Printf("%s Deleted %d image(s)\n", color.GreenString("!"), deltaDelete)
			}
			if deltaDeleteFolder > 0 {
				fmt.Printf("%s Deleted %d folder(s)\n", color.GreenString("!"), deltaDeleteFolder)
			}

			inserted = s.watcherStats.inserted.Load()
			deleted = s.watcherStats.deleted.Load()
			deletedFolder = s.watcherStats.deletedFolder.Load()
		}
	}
}
