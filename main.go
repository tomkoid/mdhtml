package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/briandowns/spinner"
)

func getAbsolutePath(path string) string {
	absPath, err := filepath.Abs(path)
	if err != nil {
		log.Fatalf("Error getting absolute path: %s", err)
		os.Exit(1)
	}

	return absPath
}

func main() {
	// parse arguments
	args := parseArgs()

	checkFilesExist(args)

	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)

	s.Suffix = fmt.Sprintf(" Transforming %s to HTML...", args.file)
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

	fmt.Printf("== Successfully wrote to %s!\n", destPath)
	fmt.Printf("   Source file: %s\n", filePath)

	if stylePath != "" {
		fmt.Printf("   Style file: %s\n", stylePath)
	}

	fmt.Println()

	fmt.Printf("View in browser at: file://%s\n", destPath)
}
