package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

func GenerateSourceFileChecksum(args Args) string {
	file, err := os.Open(args.file)
	if err != nil {
		log.Fatalf("Error opening file: %s", err)
		os.Exit(1)
	}

	defer file.Close()

	hash := md5.New()
	_, err = io.Copy(hash, file)

	if err != nil {
		log.Fatalf("Error copying file: %s", err)
		os.Exit(1)
	}

	return string(hash.Sum(nil))
}

func transformWatch(args Args, debug bool) {
	// get the initial hash
	oldHash := GenerateSourceFileChecksum(args)

	for {
		// get the new hash
		newHash := GenerateSourceFileChecksum(args)

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
