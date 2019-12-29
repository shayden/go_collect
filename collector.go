package main

import (
	"encoding/json"
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
	TagId       uint16                `json:"tag_id"`
	TagName     string                `json:"tag_name"`
	TagTypeId   exif.TagTypePrimitive `json:"tag_type_id"`
	TagTypeName string                `json:"tag_type_name"`
	UnitCount   uint32                `json:"unit_count"`
	Value       interface{}           `json:"value"`
	ValueString string                `json:"value_string"`
}

// MediaFile describes a media file
type MediaFile struct {
	AbsolutePath string `json:"absolute_path"`
	FileName     string `json:"file_name"`
	Extension    string `json:"extension"`
	ExifData     string `json:"exif_data"`
	Drivename    string `json:"drive_name"`
	Sha256       []byte `json:"sha256"`
	Size         int64  `json:"size_in_kb"`
}

var files []MediaFile
var drivename string

func walk(path string, info os.FileInfo, err error) error {
	if err != nil {
		fmt.Println(err)
		return err
	}

	switch strings.ToLower(filepath.Ext(path)) {
	case ".jpg", ".jpeg", ".gif", ".png", ".mov", ".mp4":
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

var entries []ifdEntry

//basically stolen from main.go in go-exif
func parseExifData(path string) string {
	im := exif.NewIfdMappingWithStandard()
	ti := exif.NewTagIndex()

	entries := make([]ifdEntry, 0)
	visitor := func(fqIfdPath string, ifdIndex int, tagId uint16, tagType exif.TagType, valueContext exif.ValueContext) (err error) {

		ifdPath, err := im.StripPathPhraseIndices(fqIfdPath)

		it, _ := ti.Get(ifdPath, tagId)

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
			TagId:       tagId,
			TagName:     it.Name,
			TagTypeId:   tagType.Type(),
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
		return "{}"
	}
	data, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Println(err)
		return "{}"
	}
	rawExif, _ := exif.SearchAndExtractExif(data)
	if err != nil {
		fmt.Println(err)
		return "{}"
	}
	exif.Visit(exif.IfdStandard, im, ti, rawExif, visitor)

	exifData, _ := json.Marshal(entries)
	return string(exifData)
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
