package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/dsoprea/go-exif"
)

type ifdEntry struct {
	IfdPath     string                `json:"ifd_path"`
	FqIfdPath   string                `json:"fq_ifd_path"`
	IfdIndex    int                   `json:"ifd_index"`
	TagID       uint16                `json:"tag_id"`
	TagName     string                `json:"tag_name"`
	TagTypeID   exif.TagTypePrimitive `json:"tag_type_id"`
	TagTypeName string                `json:"tag_type_name"`
	UnitCount   uint32                `json:"unit_count"`
	Value       interface{}           `json:"value"`
	ValueString string                `json:"value_string"`
}

// MediaFile describes a media file
type MediaFile struct {
	AbsolutePath string     `json:"absolute_path"`
	FileName     string     `json:"file_name"`
	Extension    string     `json:"extension"`
	ExifData     []ifdEntry `json:"exif_data"`
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

	switch strings.ToLower(filepath.Ext(path)) {
	case ".jpg", ".jpeg", ".gif", ".png", ".mov", ".mp4", ".nef", ".cr2":
		files = append(files, MediaFile{
			AbsolutePath: path,
			FileName:     info.Name(),
			Extension:    filepath.Ext(path),
			Drivename:    drivename,
			Size:         info.Size() >> 10,
			Sha256:       HashFile(path),
			ExifData:     parseExifData(path),
		})
	}
	return nil
}

//basically stolen from main.go in go-exif
func parseExifData(path string) []ifdEntry {
	im := exif.NewIfdMappingWithStandard()
	ti := exif.NewTagIndex()

	entries := make([]ifdEntry, 0)
	visitor := func(fqIfdPath string, ifdIndex int, tagId uint16, tagType exif.TagType, valueContext exif.ValueContext) (err error) {

		ifdPath, err := im.StripPathPhraseIndices(fqIfdPath)

		it, err := ti.Get(ifdPath, tagId)
		if err != nil {
			fmt.Println(err)
			return
		}

		valueString := ""
		var value interface{}
		if tagType.Type() == exif.TypeUndefined {
			value, _ = exif.UndefinedValue(ifdPath, tagId, valueContext, tagType.ByteOrder())
			valueString = fmt.Sprintf("%v", value)
		} else {
			valueString, _ = tagType.ResolveAsString(valueContext, true)
			value = valueString
		}

		entry := ifdEntry{
			IfdPath:     ifdPath,
			FqIfdPath:   fqIfdPath,
			IfdIndex:    ifdIndex,
			TagID:       tagId,
			TagName:     it.Name,
			TagTypeID:   tagType.Type(),
			TagTypeName: tagType.Name(),
			UnitCount:   valueContext.UnitCount(),
			Value:       value,
			ValueString: valueString,
		}

		entries = append(entries, entry)

		return nil
	}

	f, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	rawExif, _ := exif.SearchAndExtractExif(data)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	exif.Visit(exif.IfdStandard, im, ti, rawExif, visitor)

	if err != nil {
		fmt.Print(err)
	}
	return entries
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
