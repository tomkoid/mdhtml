package utils

import (
	"os"
	"testing"
)

func TestIsDirectory(t *testing.T) {
	// make a new directory
	err := os.Mkdir("test", 0755)
	if err != nil {
		t.Fatal(err)
	}

	// check if it is a directory
	isDir, err := IsDirectory("test")

	if err != nil {
		t.Fatal(err)
	}

	if !isDir {
		t.Fatal("test is not a directory")
	}

	// remove the directory
	err = os.Remove("test")
	if err != nil {
		t.Fatal(err)
	}
}
