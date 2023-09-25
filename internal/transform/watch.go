package transform

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"os"

	"codeberg.org/Tomkoid/mdhtml/internal/httpserver"
	"codeberg.org/Tomkoid/mdhtml/internal/models"
	"codeberg.org/Tomkoid/mdhtml/internal/utils"
	"github.com/fsnotify/fsnotify"
)

func GenerateSourceFileChecksum(args models.Args, oldHash string) string {
	file, err := os.Open(args.File)
	if err != nil {
		log.Printf("Error opening file: %s", err)
		return oldHash
	}

	defer file.Close()

	hash := md5.New()
	_, err = io.Copy(hash, file)

	if err != nil {
		log.Fatalf("Error copying file: %s", err)
		return oldHash
	}

	return string(hash.Sum(nil))
}

func checkEventType(event fsnotify.Event) bool {
	if event.Op&fsnotify.Remove == fsnotify.Remove {
		return false
	}

	return true
}

func TransformWatch(args models.Args, debug bool, httpServer bool) {
	// use fsnotify to watch for changes
	watcher, err := fsnotify.NewWatcher()

	if err != nil {
		log.Fatalf("Error creating watcher: %s", err)
		os.Exit(1)
	}

	defer watcher.Close()

	sourceIsDir, _ := utils.IsDirectory(args.File)

	if sourceIsDir {
		log.Fatalf("Error: %s is a directory", args.File)
		os.Exit(1)
	}

	if httpServer {
		go func() {
			httpserver.HttpServer(args)
		}()
	}

	done := make(chan bool)

	go func() {
		oldHash := GenerateSourceFileChecksum(args, "")

		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					continue
				}

				eventOk := checkEventType(event)
				if !eventOk {
					err := watcher.Add(event.Name)

					if err != nil {
						log.Fatalf("Error adding file to watcher: %s", err)
						os.Exit(1)
					}

					continue
				}

				newHash := GenerateSourceFileChecksum(args, oldHash)

				if oldHash != newHash {
					Transform(args, false)
					fmt.Println("== Successfully transformed to markdown...")

					httpserver.SetReload()
					oldHash = newHash
				}

				watcher.Add(event.Name)
			case err, ok := <-watcher.Errors:
				if !ok {
					continue
				}

				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(args.File)
	if err != nil {
		log.Fatalf("Error adding file to watcher: %s", err)
		os.Exit(1)
	}

	<-done
}
