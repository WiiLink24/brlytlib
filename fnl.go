package brlytlib

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// FNLHeader represents the header of the fnl1 section
type FNLHeader struct {
	Magic       [4]byte
	SectionSize uint32
	NumOfFonts  uint16
	_           uint16
}

type FNLOffset struct {
	// OffSet is relative to the beginning of the fnl1 section
	Offset uint32
	_      uint32
}

func ParseFNL(contents []byte) ([]FNLNames, error) {
	var fnl FNLHeader
	var fnlNamesFormat []FNLNamesFormat
	var fnlNames []FNLNames

	fnlOffset := findAllOccurrences(contents, []string{"fnl1"})
	var textOffsets []uint32

	// There is only one instance of fnl1 per BRLYT
	err := binary.Read(bytes.NewReader(contents[fnlOffset[0]:]), binary.BigEndian, &fnl)
	if err != nil {
		return nil, err
	}

	for i := 0; i < int(fnl.NumOfFonts); i++ {
		var fnlOffsets FNLOffset
		offset := (fnlOffset[0] + 12) + (i * 8)

		fmt.Println(offset)
		err := binary.Read(bytes.NewReader(contents[offset:]), binary.BigEndian, &fnlOffsets)
		if err != nil {
			return nil, err
		}

		textOffsets = append(textOffsets, fnlOffsets.Offset+uint32(offset))

		// If we have reached the last index, append the section size to the slice.
		if i == int(fnl.NumOfFonts)-1 {
			textOffsets = append(textOffsets, uint32(fnl.SectionSize)+uint32(fnlOffset[0]))
		}
	}

	for i := 0; i < int(fnl.NumOfFonts); i++ {
		fontName := string(contents[textOffsets[i]:textOffsets[i+1]])

		xmlNode := FNLNamesFormat{
			Index:  i,
			String: fontName,
		}

		fnlNamesFormat = append(fnlNamesFormat, xmlNode)
	}

	fnlNames = append(fnlNames, FNLNames{
		FNLName: fnlNamesFormat,
	})

	return fnlNames, err
}
