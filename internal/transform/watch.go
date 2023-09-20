package transform

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"os"

	"codeberg.org/Tomkoid/mdhtml/internal/models"
	"codeberg.org/Tomkoid/mdhtml/internal/utils"
	"github.com/fsnotify/fsnotify"
)

func GenerateSourceFileChecksum(args models.Args) string {
	file, err := os.Open(args.File)
	if err != nil {
		log.Fatalf("Error opening file: %s", err)
		os.Exit(1)
	}

	defer file.Close()

	hash := md5.New()
	_, err = io.Copy(hash, file)

	if err != nil {
		log.Fatalf("Error copying file: %s", err)
		os.Exit(1)
	}

	return string(hash.Sum(nil))
}

func TransformWatch(args models.Args, debug bool) {
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

	done := make(chan bool)

	go func() {
		oldHash := GenerateSourceFileChecksum(args)

		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					continue
				}

				newHash := GenerateSourceFileChecksum(args)

				if oldHash != newHash {
					Transform(args, false)

					fmt.Println("== Successfully transformed to markdown...")

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
