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
	"github.com/fatih/color"
	"github.com/fsnotify/fsnotify"
)

func GenerateChecksum(filePath string, oldHash string) string {
	file, err := os.Open(filePath)
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

func GenerateSourceFileChecksum(args models.Args, oldHash string) string {
	var srcFileChecksum string = ""
	var styleFileChecksum string = ""

	var hashString string

	srcFileChecksum = GenerateChecksum(args.File, oldHash)

	if args.Style != "" {
		styleFileChecksum = GenerateChecksum(args.Style, oldHash)
	}

	if styleFileChecksum != "" {
		hashString = srcFileChecksum + styleFileChecksum
	} else {
		hashString = srcFileChecksum
	}

	return hashString
}

func checkEventType(event fsnotify.Event) bool {
	if event.Op&fsnotify.Remove == fsnotify.Remove {
		return false
	}

	return true
}

func TransformWatch(args models.Args, debug bool, httpServer bool) {
	// check if the out file exists
	if _, err := os.Stat(args.Out); os.IsNotExist(err) {
		Transform(args, false)
	}

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
		go httpserver.HttpServer(args)
	}

	var watchFiles []string
	watchFiles = append(watchFiles, args.File)

	if args.Style != "" {
		watchFiles = append(watchFiles, args.Style)
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
					httpserver.BroadcastMessage("transforming")
					Transform(args, false)

					color.Set(color.FgGreen)
					fmt.Print("==")
					color.Unset()

					fmt.Println(" Successfully transformed to markdown!")

					httpserver.BroadcastMessage("reload")
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

	for _, element := range watchFiles {
		err = watcher.Add(element)
		if err != nil {
			log.Fatalf("Error adding file to watcher: %s", err)
			os.Exit(1)
		}

	}
	<-done
}
