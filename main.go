// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !plan9,!solaris

package main

import (
	"flag"
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
)

var (
	verboseFlag  = flag.Bool("v", false, "Verbose output")
	namePattern  = flag.String("name", "", "Only detect changes on files matching this pattern. For example, -name \"*.go\"")
	opFilter     = fsnotify.Create | fsnotify.Write | fsnotify.Remove | fsnotify.Rename
	quietTimeSec = flag.Float64("quiet", 5.0, "Quiet time after an execution, in seconds")
)

func main() {
	flag.Parse()

	if *namePattern != "" {
		log.Println("Filtering on name pattern '" + *namePattern + "'")
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)

	go func() {
		var lastTime time.Time
		quietInterval := time.Duration(time.Duration(*quietTimeSec*1000) * time.Millisecond)
		log.Println("Quiet period: ", quietInterval)
		for {
			select {
			case event := <-watcher.Events:
				if time.Since(lastTime) < quietInterval {
					if *verboseFlag {
						log.Println("within quiet interval, skipped event:", event, "(", event.Name, ")")
					}
					break // still within the quiet period
				}

				if *verboseFlag {
					log.Println("event:", event, "(", event.Name, ")")
				}

				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file:", event.Name)
				}
				if !matching(event.Name) {
					break // file name doesn't match pattern
				}
				if event.Op&opFilter == 0 {
					break // operation is not among those monitored for change
				}

				args := flag.Args()
				if len(args) > 0 {
					log.Println("change detected...")
					lastTime = time.Now()
					out, err := exec.Command(args[0], args[1:]...).CombinedOutput()
					fmt.Println(string(out))
					if err != nil {
						log.Println(err)
					} else {
						log.Println("ok")
					}
				}
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(".")
	if err != nil {
		log.Fatal(err)
	}
	<-done
}

func matching(name string) bool {
	if namePattern != nil && *namePattern != "" {
		match, err := filepath.Match(*namePattern, name)
		return match || err != nil
	}
	return true
}
