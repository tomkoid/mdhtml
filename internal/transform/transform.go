package transform

import (
	"fmt"
	"log"
	"os"
	"time"

	"codeberg.org/Tomkoid/mdhtml/internal/models"
	"codeberg.org/Tomkoid/mdhtml/internal/utils"
	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/gomarkdown/markdown"
)

func Transform(args models.Args, debug bool) {
	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)

	s.Suffix = fmt.Sprintf(" Transforming %s to HTML...", color.BlueString(args.File))
	s.Start()

	err := transformMarkdownToHTML(args)
	if !err {
		s.Stop()
		log.Fatalf("Error transforming markdown to HTML")
		os.Exit(1)
	}

	s.Stop()

	filePath := utils.GetAbsolutePath(args.File)
	destPath := utils.GetAbsolutePath(args.Out)

	stylePath := ""

	if args.Style != "" {
		stylePath = utils.GetAbsolutePath(args.Style)
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
		fmt.Printf("View in bro/Projects/Gnomewser at: ")
		color.Unset()

		fmt.Printf("file://%s\n", destPath)
	}
}

func transformMarkdownToHTML(args models.Args) bool {
	content, err := os.ReadFile(args.File)
	if err != nil {
		log.Fatalf("Error reading file: %s", err)
		os.Exit(1)
	}

	html := markdown.ToHTML(content, nil, nil)

	// apply styling if provided
	if args.Style != "" {
		style, err := os.ReadFile(args.Style)
		if err != nil {
			log.Fatalf("Error reading style file: %s", err)
			os.Exit(1)
		}

		html = append(html, fmt.Sprintf("<style>\n%s\n</style>", style)...)
	}

	err = os.WriteFile(args.Out, html, 0644)

	if err != nil {
		log.Fatalf("Error writing to file: %s", err)
		os.Exit(1)
	}

	return true
}
