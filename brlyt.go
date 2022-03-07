package main

import (
	"bytes"
	"encoding/binary"
	"encoding/xml"
	"io/ioutil"
)

// Header represents the header of our BRLYT
type Header struct {
	Magic        [4]byte
	BOM          uint32
	FileSize     uint32
	HeaderLen    uint16
	SectionCount uint16
}

// SectionTypes are known parts of a BRLYT.
type SectionTypes [4]byte

var (
	headerMagic    SectionTypes = [4]byte{'R', 'L', 'Y', 'T'}
	SectionTypeLYT SectionTypes = [4]byte{'l', 'y', 't', '1'}
	SectionTypeTXL SectionTypes = [4]byte{'t', 'x', 'l', '1'}
	SectionTypeFNL SectionTypes = [4]byte{'f', 'n', 'l', '1'}
	SectionTypeMAT SectionTypes = [4]byte{'m', 'a', 't', '1'}
	SectionTypePAN SectionTypes = [4]byte{'p', 'a', 'n', '1'}
	SectionTypePAS SectionTypes = [4]byte{'p', 'a', 's', '1'}
	SectionTypePAE SectionTypes = [4]byte{'p', 'a', 'e', '1'}
	SectionTypeBND SectionTypes = [4]byte{'b', 'n', 'd', '1'}
	SectionTypePIC SectionTypes = [4]byte{'p', 'i', 'c', '1'}
	SectionTypeTXT SectionTypes = [4]byte{'t', 'x', 't', '1'}
	SectionTypeWND SectionTypes = [4]byte{'w', 'n', 'd', '1'}
	SectionTypeGRP SectionTypes = [4]byte{'g', 'r', 'p', '1'}
	SectionTypeGRS SectionTypes = [4]byte{'g', 'r', 's', '1'}
	SectionTypeGRE SectionTypes = [4]byte{'g', 'r', 'e', '1'}
)

type SectionHeader struct {
	Type SectionTypes
	Size uint32
}

func ParseBRLYT(fileName string) ([]byte, error) {
	contents, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	// Create a new reader
	readable := bytes.NewReader(contents)

	var header Header
	err = binary.Read(readable, binary.BigEndian, &header)
	if err != nil {
		panic(err)
	}

	if !bytes.Equal(headerMagic[:], header.Magic[:]) {
		return nil, ErrInvalidFileMagic
	}

	if readable.Size() != int64(header.FileSize) {
		return nil, ErrFileSizeMismatch
	}

	root := Root{
		XMLName: xml.Name{},
		LYT:     LYTNode{},
		FNL:     nil,
		TXL:     nil,
		Panes:   nil,
	}

	for count := header.SectionCount; count != 0; count-- {
		var sectionHeader SectionHeader
		err = binary.Read(readable, binary.BigEndian, &sectionHeader)
		if err != nil {
			return nil, err
		}

		// Subtract the header size
		sectionSize := int(sectionHeader.Size) - 8
		if readable.Len() == 0 {
			continue
		}

		temp := make([]byte, sectionSize)
		_, err = readable.Read(temp)
		if err != nil {
			return nil, err
		}

		switch sectionHeader.Type {
		case SectionTypeLYT:
			root.ParseLYT(temp)
		case SectionTypeTXL:
			root.ParseTXL(temp, sectionHeader.Size)
		case SectionTypeFNL:
			root.ParseFNL(temp, sectionHeader.Size)
		case SectionTypeMAT:
			root.ParseMAT(temp, sectionHeader.Size)
		case SectionTypePAN:
			root.ParsePAN(temp)
		case SectionTypeBND:
			root.ParseBND(temp)
		case SectionTypePIC:
			root.ParsePIC(temp)
		case SectionTypeTXT:
			root.ParseTXT(temp, sectionHeader.Size)
		case SectionTypeWND:
			root.ParseWND(temp)
		case SectionTypePAS:
			root.ParsePAS()
		case SectionTypePAE:
			root.ParsePAE()
		case SectionTypeGRP:
			root.ParseGRP(temp)
		case SectionTypeGRS:
			root.ParseGRS()
		case SectionTypeGRE:
			root.ParseGRE()
		}
	}

	data, err := xml.MarshalIndent(root, "", "\t")
	if err != nil {
		return nil, err
	}

	return data, nil
}
