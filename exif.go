package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/dsoprea/go-exif"
)

type IfdEntry struct {
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

// ParseExifData uses go-exif to extract any exif data
//basically stolen from main.go in go-exif
func ParseExifData(path string) []IfdEntry {
	im := exif.NewIfdMappingWithStandard()
	ti := exif.NewTagIndex()

	entries := make([]IfdEntry, 0)
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

		entry := IfdEntry{
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
