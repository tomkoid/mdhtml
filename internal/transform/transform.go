package transform

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"codeberg.org/Tomkoid/mdhtml/internal/models"
	"codeberg.org/Tomkoid/mdhtml/internal/utils"
	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
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
		fmt.Printf("View in browser at: ")
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

	html := mdToHTML(content)

	prismData := `
    <script src="/prism.js" defer></script>
    <link rel="stylesheet" href="/prism.css">
    `

	// if content doesnt include any code block or inline code, dont include prism
	if !strings.Contains(string(content), "```") || args.NoExternalLibs {
		prismData = ``
	}

	headData := `
    <meta name="viewport" content="width=device-width,initial-scale=1">
    <link rel="stylesheet" href="/default.css">
    <script src="/reload.js" defer></script>
  ` + prismData

	// apply styling if provided
	if args.Style != "" {
		style, err := os.ReadFile(args.Style)
		if err != nil {
			log.Fatalf("Error reading style file: %s", err)
			os.Exit(1)

		}

		html = append(html, fmt.Sprintf("<style>\n%s\n</style>\n%s", style, headData)...)

	} else {
		html = append(html, fmt.Sprintf("\n%s", headData)...)
	}

	err = os.WriteFile(args.Out, html, 0644)

	if err != nil {
		log.Fatalf("Error writing to file: %s", err)
		os.Exit(1)
	}

	return true
}

func mdToHTML(md []byte) []byte {
	// create markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	// create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return markdown.Render(doc, renderer)
}
