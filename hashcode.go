package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

// HashFile creates a SHA256 hash of the path given
func HashFile(path string) []byte {
	f, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		fmt.Println(err)
	}
	return h.Sum(nil)
}
