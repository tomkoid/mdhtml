package utils

import (
	"crypto/sha1"
	"encoding/hex"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func IsDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	return fileInfo.IsDir(), err
}

func GetAbsolutePath(path string) string {
	absPath, err := filepath.Abs(path)
	if err != nil {
		log.Fatalf("Error getting absolute path: %s", err)
		os.Exit(1)
	}

	return absPath
}

func ShaString(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func OpenInBrowser(path string) {
	cmd := exec.Command("xdg-open", path)
	cmd.Run()
}
