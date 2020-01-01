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
	// aperture & imovie create (i'm guessing) thumbnails or smaller res caches
	// we don't want them
	if strings.HasPrefix(path, ".") {
		return nil
	}
	hc := make(chan int)
	exifc := make(chan int)
	uc := make(chan int)
	fupc := make(chan int)

	switch strings.ToLower(filepath.Ext(path)) {
	case ".jpg", ".jpeg", ".gif", ".png", ".mov", ".mp4", ".nef", ".cr2":
		abs, _ := filepath.Abs(path)
		mf := &MediaFile{
			AbsolutePath: abs,
			FileName:     info.Name(),
			Extension:    filepath.Ext(path),
			Drivename:    drivename,
			Size:         info.Size() >> 10,
		}
		go func(mf *MediaFile, path string) {
			mf.Sha256 = HashFile(path)
			hc <- 1
		}(mf, path)

		go func(mf *MediaFile, path string) {
			mf.ExifData = ParseExifData(path)
			exifc <- 1
		}(mf, path)

		go func(path string) {
			fmt.Println("Begin file upload")
			uploadFile(path)
			fmt.Println("Finish file upload")
			fupc <- 1
		}(path)
		<-hc
		<-exifc
		go func(mf *MediaFile) {
			uploadMetadata(mf)
			uc <- 1
		}(mf)
		<-fupc
		<-uc
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
