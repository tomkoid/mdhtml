package utils

import (
	"fmt"
	"os"

	"codeberg.org/Tomkoid/mdhtml/internal/models"
)

func CheckFilesExist(args models.Args) bool {
	if _, err := os.Stat(args.File); os.IsNotExist(err) {
		fmt.Printf("Error: file %s does not exist\n", args.File)
		os.Exit(1)
	}

	if args.Style != "" {
		if _, err := os.Stat(args.Style); os.IsNotExist(err) {
			fmt.Printf("Error: style file %s does not exist\n", args.Style)
			os.Exit(1)
		}
	}

	return true
}
