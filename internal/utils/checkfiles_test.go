package utils

import (
	"os"
	"testing"

	"codeberg.org/Tomkoid/mdhtml/internal/models"
)

func TestCheckFilesExist(t *testing.T) {
	// create a new file
	_, err := os.Create("test.md")
	if err != nil {
		t.Fatal(err)
	}

	// create a new style file
	_, err = os.Create("style.css")
	if err != nil {
		t.Fatal(err)
	}

	// create a new args instance
	args := models.Args{
		File:  "test.md",
		Style: "style.css",
	}

	// check if the files exist
	exists := CheckFilesExist(args)

	if !exists {
		t.Fatal("files do not exist")
	}

	// remove the files
	err = os.Remove("test.md")
	if err != nil {
		t.Fatal(err)
	}

	err = os.Remove("style.css")
	if err != nil {
		t.Fatal(err)
	}
}
