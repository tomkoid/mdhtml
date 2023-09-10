package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
)

func getAbsolutePath(path string) string {
	absPath, err := filepath.Abs(path)
	if err != nil {
		log.Fatalf("Error getting absolute path: %s", err)
		os.Exit(1)
	}

	return absPath
}

func transform(args Args, debug bool) {
	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)

	s.Suffix = fmt.Sprintf(" Transforming %s to HTML...", color.BlueString(args.file))
	s.Start()

	err := transformMarkdownToHTML(args)
	if !err {
		s.Stop()
		log.Fatalf("Error transforming markdown to HTML")
		os.Exit(1)
	}

	s.Stop()

	filePath := getAbsolutePath(args.file)
	destPath := getAbsolutePath(args.out)

	stylePath := ""

	if args.style != "" {
		stylePath = getAbsolutePath(args.style)
	}

	if debug {
		color.Set(color.FgGreen)
		fmt.Printf("==")
		color.Unset()

		fmt.Printf(" Successfully wrote to %s!\n", destPath)
		fmt.Printf("   Source file: %s\n", filePath)

		if stylePath != "" {
			fmt.Printf("   Style file: %s\n", stylePath)
		}

		fmt.Println()

		color.Set(color.FgBlue, color.Bold)
		fmt.Printf("View in browser at: ")
		color.Unset()

		fmt.Printf("file://%s\n", destPath)
	}
}

func shaString(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func transformWatch(args Args, debug bool) {
	// get initial content
	content, err := os.ReadFile(args.file)
	if err != nil {
		log.Fatalf("Error reading file: %s", err)
		os.Exit(1)
	}

	// get the hash
	oldHash := shaString(string(content))

	for {
		// get the new hash
		content, err := os.ReadFile(args.file)
		if err != nil {
			log.Fatalf("Error reading file: %s", err)
			os.Exit(1)
		}

		newHash := shaString(string(content))

		// compare the hashes
		if oldHash != newHash {
			fmt.Println("== Detected change in file, transforming...")

			// if they're different, transform the file
			transform(args, false)

			// update the old hash
			oldHash = newHash
		}

		// sleep for 500ms
		time.Sleep(500 * time.Millisecond)
	}
}

func main() {
	// parse arguments
	args := parseArgs()

	if args.watch {
		fmt.Printf("== Watching %s for changes...\n", args.file)
		transformWatch(args, true)
		return
	}

	checkFilesExist(args)

	transform(args, true)
}
