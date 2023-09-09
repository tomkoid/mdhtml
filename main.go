package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
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

	fmt.Println("Transforming markdown file to HTML...")
	err := transformMarkdownToHTML(args)
	if !err {
		log.Fatalf("Error transforming markdown to HTML")
		os.Exit(1)
	}

	fmt.Println()

	destPath := getAbsolutePath(args.out)
	fmt.Printf("== Successfully wrote to %s!\n", destPath)

	fmt.Printf("View in browser at: file://%s\n", destPath)
}
