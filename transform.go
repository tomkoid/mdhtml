package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/gomarkdown/markdown"
)

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

func transformMarkdownToHTML(args Args) bool {
	content, err := os.ReadFile(args.file)
	if err != nil {
		log.Fatalf("Error reading file: %s", err)
		os.Exit(1)
	}

	html := markdown.ToHTML(content, nil, nil)

	// apply styling if provided
	if args.style != "" {
		style, err := os.ReadFile(args.style)
		if err != nil {
			log.Fatalf("Error reading style file: %s", err)
			os.Exit(1)
		}

		html = append(html, fmt.Sprintf("<style>\n%s\n</style>", style)...)
	}

	err = os.WriteFile(args.out, html, 0644)

	if err != nil {
		log.Fatalf("Error writing to file: %s", err)
		os.Exit(1)
	}

	return true
}
