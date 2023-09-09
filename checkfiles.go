package main

import (
	"fmt"
	"os"
)

func checkFilesExist(args Args) bool {
	if _, err := os.Stat(args.file); os.IsNotExist(err) {
		fmt.Printf("Error: file %s does not exist\n", args.file)
		os.Exit(1)
	}

	if args.style != "" {
		if _, err := os.Stat(args.style); os.IsNotExist(err) {
			fmt.Printf("Error: style file %s does not exist\n", args.style)
			os.Exit(1)
		}
	}

	return true
}
