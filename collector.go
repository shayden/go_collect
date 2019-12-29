package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// MediaFile describes a media file
type MediaFile struct {
	AbsolutePath string     `json:"absolute_path"`
	FileName     string     `json:"file_name"`
	Extension    string     `json:"extension"`
	ExifData     []IfdEntry `json:"exif_data"`
	Drivename    string     `json:"drive_name"`
	Sha256       []byte     `json:"sha256"`
	Size         int64      `json:"size_in_kb"`
}

var files []MediaFile
var drivename string

func walk(path string, info os.FileInfo, err error) error {
	if err != nil {
		fmt.Println(err)
		return err
	}
	if strings.HasPrefix(path, ".") {
		return nil
	}

	switch strings.ToLower(filepath.Ext(path)) {
	case ".jpg", ".jpeg", ".gif", ".png", ".mov", ".mp4", ".nef", ".cr2":
		files = append(files, MediaFile{
			AbsolutePath: path,
			FileName:     info.Name(),
			Extension:    filepath.Ext(path),
			Drivename:    drivename,
			Size:         info.Size() >> 10,
			Sha256:       HashFile(path),
			ExifData:     ParseExifData(path),
		})
	}
	return nil
}

// WalkPath Returns a list of media files starting at root
func WalkPath(root, dname string) []MediaFile {
	drivename = dname
	err := filepath.Walk(root, walk)
	if err != nil {
		fmt.Println(err)
	}
	return files
}
