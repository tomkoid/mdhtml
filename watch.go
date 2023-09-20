package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
)

func GenerateSourceFileChecksum(args Args) string {
	file, err := os.Open(args.file)
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

func transformWatch(args Args, debug bool) {
	// use fsnotify to watch for changes
	watcher, err := fsnotify.NewWatcher()

	if err != nil {
		log.Fatalf("Error creating watcher: %s", err)
		os.Exit(1)
	}

	defer watcher.Close()

	sourceIsDir, _ := isDirectory(args.file)

	if sourceIsDir {
		log.Fatalf("Error: %s is a directory", args.file)
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
					transform(args, false)

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

	err = watcher.Add(args.file)
	if err != nil {
		log.Fatalf("Error adding file to watcher: %s", err)
		os.Exit(1)
	}

	<-done
}
