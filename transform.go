package main

import (
	"fmt"
	"log"
	"os"

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
