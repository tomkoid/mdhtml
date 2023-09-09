package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/gomarkdown/markdown"
)

func transformMarkdownToHTML(args Args) bool {
	content, err := os.ReadFile(args.file)
	if err != nil {
		log.Fatalf("Error reading file: %s", err)
		os.Exit(1)
	}

	html := markdown.ToHTML(content, nil, nil)

	// apply styling if provided
	if args.style != "" {
		fmt.Printf("config: applying the provided %s styling file..\n", args.style)
		style, err := os.ReadFile(args.style)
		if err != nil {
			log.Fatalf("Error reading style file: %s", err)
			os.Exit(1)
		}

		html = append(html, fmt.Sprintf("<style>%s</style>", style)...)
	}

	err = os.WriteFile(args.dest, html, 0644)

	if err != nil {
		log.Fatalf("Error writing to file: %s", err)
		os.Exit(1)
	}

	return true
}

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

	fmt.Println("Transforming markdown file to HTML...")
	err := transformMarkdownToHTML(args)
	if !err {
		log.Fatalf("Error transforming markdown to HTML")
		os.Exit(1)
	}

	fmt.Println()

	destPath := getAbsolutePath(args.dest)
	fmt.Printf("== Successfully wrote to %s!\n", destPath)

	fmt.Printf("View in browser at: file://%s\n", destPath)
}
