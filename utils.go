package main

import (
	"crypto/sha1"
	"encoding/hex"
	"log"
	"os"
	"path/filepath"
)

func getAbsolutePath(path string) string {
	absPath, err := filepath.Abs(path)
	if err != nil {
		log.Fatalf("Error getting absolute path: %s", err)
		os.Exit(1)
	}

	return absPath
}

func shaString(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}
