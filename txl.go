package main

import (
	"bytes"
	"encoding/binary"
	"strings"
)

type TXL struct {
	NumOfTPL uint16
	Unknown  uint16
}

type TPLOffSet struct {
	// Offset is relative to the beginning of the txl1 section
	Offset  uint32
	Padding uint32
}

func (r *Root) ParseTXL(data []byte, sectionSize uint32) {
	var tplOffsets []uint32
	var tplNames []string

	var txl TXL
	err := binary.Read(bytes.NewReader(data), binary.BigEndian, &txl)
	if err != nil {
		panic(err)
	}

	for i := 0; i < int(txl.NumOfTPL); i++ {
		// By now we have only read the header.
		// We will read the TPLOffset table in order to get our names.
		var tplTable TPLOffSet
		offset := 4 + (i * 8)

		err = binary.Read(bytes.NewReader(data[offset:]), binary.BigEndian, &tplTable)
		if err != nil {
			panic(err)
		}

		tplOffsets = append(tplOffsets, tplTable.Offset+4)

		// If we have reached the last index, append the section size to the slice.
		if i == int(txl.NumOfTPL)-1 {
			tplOffsets = append(tplOffsets, sectionSize-8)
		}
	}

	// Now that we have the offsets, retrieve the TPL names.
	for i := 0; i < int(txl.NumOfTPL); i++ {
		tplName := string(data[tplOffsets[i]:tplOffsets[i+1]])

		// Strip the null terminator
		tplName = strings.Replace(tplName, "\x00", "", -1)

		tplNames = append(tplNames, tplName)
	}

	r.TXL = &TPLNames{TPLName: tplNames}
}
