package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// MediaFile describes a media file
type MediaFile struct {
	absolutePath string
	fileName     string
	extension    string
	exifData     string
	drivename    string
	sha256       []byte
	size         int64
}

var files []MediaFile

func walk(path string, info os.FileInfo, err error) error {
	if err != nil {
		fmt.Println(err)
		return err
	}

	switch filepath.Ext(path) {
	case ".jpg", ".jpeg", ".gif", ".png", ".mov", ".mp4":
		files = append(files, MediaFile{
			absolutePath: path,
			fileName:     info.Name(),
			extension:    filepath.Ext(path),
			drivename:    filepath.VolumeName(path),
			size:         info.Size() >> 10,
			sha256:       HashFile(path),
		})
	}
	return nil
}

// WalkPath Returns a list of media files starting at root
func WalkPath(root string) []MediaFile {
	err := filepath.Walk(root, walk)
	if err != nil {
		fmt.Println(err)
	}
	return files
}
