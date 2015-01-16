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

	"github.com/go-fsnotify/fsnotify"
)

var (
	verboseFlag = flag.Bool("v", false, "Verbose output")
	namePattern = flag.String("name", "", "Only detect changes on files matching this pattern. For example, -name \"*.go\"")
	opFilter    = fsnotify.Create | fsnotify.Write | fsnotify.Remove | fsnotify.Rename
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
		for {
			select {
			case event := <-watcher.Events:
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
	} else {
		return true
	}
}
