package brlytlib

import (
	"bytes"
	"encoding/binary"
	"encoding/xml"
	"strings"
)

type TXL struct {
	Magic 			[4]byte
	SectionSize		uint32
	NumOfTPL		uint16
	_				uint16
}

type TPLOffSet struct {
	// OffSet is relative to the beginning of the txl1 section
	OffSet			uint32
	_				uint32
}

func ParseTXL(contents []byte) ([]TPLNames, error) {
	var txl TXL
	var tplOffsets []uint32
	var tplFormat []TPLNamesFormat
	var tplNames []TPLNames

	// All appearances of txl1 are at 0x24.
	err := binary.Read(bytes.NewReader(contents[0x24:]), binary.BigEndian, &txl)
	if err != nil {
		return nil, err
	}

	if txl.Magic != txlMagic {
		return nil, ErrInvalidTXLHeader
	}

	// Grab the TPL offsets
	for i := 0; i < int(txl.NumOfTPL); i++ {
		var tpl TPLOffSet
		offset := 48 + (i * 8)
		err = binary.Read(bytes.NewReader(contents[offset:txl.SectionSize + 36]), binary.BigEndian, &tpl)
		if err != nil {
			return nil, err
		}

		tplOffsets = append(tplOffsets, tpl.OffSet + 48)

		// If we have reached the last index, append the section size to the slice.
		if i == int(txl.NumOfTPL) - 1 {
			tplOffsets = append(tplOffsets, txl.SectionSize + 36)
		}
	}

	// Now that we have the TPL offsets, we can properly get the strings
	for i := 0; i < int(txl.NumOfTPL); i++ {
		tplName := string(contents[tplOffsets[i]:tplOffsets[i+1]])
		// Strip null bytes
		tplName = strings.Replace(tplName, "\x00", "", -1)

		xmlNode := TPLNamesFormat{
			Index:  i,
			String: tplName,
		}

		tplFormat = append(tplFormat, xmlNode)
	}

	tplNames = append(tplNames, TPLNames{
		XMLName: xml.Name{},
		TPLName: tplFormat,
	})

	return tplNames, nil
}