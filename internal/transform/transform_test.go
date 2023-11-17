package transform

import (
	"os"
	"strings"
	"testing"

	"codeberg.org/Tomkoid/mdhtml/internal/models"
)

func TestTransform(t *testing.T) {
	// create a new file and write some markdown to it
	file := "test.md"
	err := os.WriteFile(file, []byte("# Test"), 0644)
	if err != nil {
		t.Errorf("Error creating test file: %s", err)
	}

	// create a new style
	style := "test.css"
	err = os.WriteFile(style, []byte("body { background-color: red; }"), 0644)
	if err != nil {
		t.Errorf("Error creating test style: %s", err)
	}

	outFile := "test.html"

	// create a new instance of Args
	args := models.Args{
		File:  file,
		Out:   outFile,
		Style: style,
	}

	// transform the markdown to HTML
	startTransform(args)

	// check if the HTML file exists
	if _, err := os.Stat(outFile); os.IsNotExist(err) {
		t.Errorf("Error: %s", err)
	}

	// read the HTML file
	content, err := os.ReadFile(outFile)
	if err != nil {
		t.Errorf("Error reading HTML file: %s", err)
	}

	// remove all whitespaces and newlines from the HTML
	html := strings.ReplaceAll(string(content), " ", "")
	html = strings.ReplaceAll(html, "\n", "")

	t.Logf("HTML: %s", html)
	// check if the HTML contains the markdown
	if !strings.Contains(html, `<h1id="test">Test</h1>`) {
		t.Errorf("Error: HTML does not contain markdown")
	}

	// if HTML is the same
	if html != `<h1id="test">Test</h1><style>body{background-color:red;}</style><metaname="viewport"content="width=device-width,initial-scale=1"><linkrel="stylesheet"href="/default.css"><scriptsrc="/reload.js"defer></script>` {
		t.Errorf("Error: HTML is not the same as expected: %s", html)
	}

	// remove the files
	err = os.Remove(file)
	if err != nil {
		t.Errorf("Error removing test file: %s", err)
	}

	err = os.Remove(style)
	if err != nil {
		t.Errorf("Error removing test style: %s", err)
	}
}
