package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

func transformWatch(args Args, debug bool) {
	// get initial content
	content, err := os.ReadFile(args.file)
	if err != nil {
		log.Fatalf("Error reading file: %s", err)
		os.Exit(1)
	}

	// get the hash
	oldHash := shaString(string(content))

	for {
		// get the new hash
		content, err := os.ReadFile(args.file)
		if err != nil {
			log.Fatalf("Error reading file: %s", err)
			os.Exit(1)
		}

		newHash := shaString(string(content))

		// compare the hashes
		if oldHash != newHash {
			fmt.Println("== Detected change in file, transforming...")

			// if they're different, transform the file
			transform(args, false)

			// update the old hash
			oldHash = newHash
		}

		// sleep for 500ms
		time.Sleep(500 * time.Millisecond)
	}
}
