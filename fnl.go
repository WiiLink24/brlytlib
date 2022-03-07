package main

import (
	"bytes"
	"encoding/binary"
	"strings"
)

// FNL represents the header of the fnl1 section
type FNL struct {
	NumOfFonts uint16
	_          uint16
}

type FNLTable struct {
	// OffSet is relative to the beginning of the fnl1 section
	Offset uint32
	_      uint32
}

func (r *Root) ParseFNL(data []byte, sectionSize uint32) {
	var fontOffsets []uint32
	var fontNames []string

	var fnl FNL
	err := binary.Read(bytes.NewReader(data), binary.BigEndian, &fnl)
	if err != nil {
		panic(err)
	}

	for i := 0; i < int(fnl.NumOfFonts); i++ {
		// By now we have only read the header.
		// We will read the FNLOffset table in order to get our names.
		var fnlTable FNLTable
		offset := 4 + (i * 8)

		err = binary.Read(bytes.NewReader(data[offset:]), binary.BigEndian, &fnlTable)
		if err != nil {
			panic(err)
		}

		fontOffsets = append(fontOffsets, fnlTable.Offset+4)

		// If we have reached the last index, append the section size to the slice.
		if i == int(fnl.NumOfFonts)-1 {
			fontOffsets = append(fontOffsets, sectionSize-8)
		}
	}

	// Now that we have the offsets, retrieve the TPL names.
	for i := 0; i < int(fnl.NumOfFonts); i++ {
		fontName := string(data[fontOffsets[i]:fontOffsets[i+1]])

		// Strip the null terminator
		fontName = strings.Replace(fontName, "\x00", "", -1)

		fontNames = append(fontNames, fontName)
	}

	r.FNL = &FNLNames{FNLName: fontNames}
}
