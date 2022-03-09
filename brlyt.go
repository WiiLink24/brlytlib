package main

import (
	"bytes"
	"encoding/binary"
	"encoding/xml"
	"io"
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

type BRLYTWriter struct {
	*bytes.Buffer
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
			// If our type is one of the section ending types, we can write then finish.
			switch sectionHeader.Type {
			case SectionTypePAE:
				root.ParsePAE()
			case SectionTypeGRE:
				root.ParseGRE()
			}
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

func WriteBRLYT(data []byte) ([]byte, error) {
	var root Root
	err := xml.Unmarshal(data, &root)
	if err != nil {
		return nil, err
	}

	writer := BRLYTWriter{bytes.NewBuffer(nil)}

	// First write the header
	header := Header{
		Magic:        headerMagic,
		BOM:          0xFEFF000A,
		FileSize:     0,
		HeaderLen:    16,
		SectionCount: 0,
	}

	err = binary.Write(writer, binary.BigEndian, header)
	if err != nil {
		return nil, err
	}

	sectionCount := len(root.Panes)

	// Write the LYT1 section
	writer.WriteLYT(root)
	sectionCount += 1

	if root.TXL != nil {
		// Write TXL section
		writer.WriteTXL(root)
		sectionCount += 1
	}

	if root.FNL != nil {
		// Write FNL section
		writer.WriteFNL(root)
		sectionCount += 1
	}

	// Write MAT section
	writer.WriteMAT(root)
	sectionCount += 1

	for _, pane := range root.Panes {
		// Please bear with me, we must check which pane is not nil
		if pane.Pane != nil {
			writer.WritePane(*pane.Pane)
		}
		if pane.PAS != nil {
			writer.WritePAS()
		}
		if pane.PAE != nil {
			writer.WritePAE()
		}
		if pane.BND != nil {
			writer.WriteBND(*pane.BND)
		}
		if pane.PIC != nil {
			writer.WritePIC(*pane.PIC)
		}
		if pane.TXT != nil {
			writer.WriteTXT(*pane.TXT)
		}
		if pane.WND != nil {
			writer.WriteWND(*pane.WND)
		}
		if pane.GRP != nil {
			writer.WriteGRP(*pane.GRP)
		}
		if pane.GRS != nil {
			writer.WriteGRS()
		}
		if pane.GRE != nil {
			writer.WriteGRE()
		}
	}

	binary.BigEndian.PutUint32(writer.Bytes()[8:12], uint32(writer.Len()))
	binary.BigEndian.PutUint16(writer.Bytes()[14:16], uint16(sectionCount))

	return writer.Bytes(), nil
}

func write(writer io.Writer, data interface{}) {
	err := binary.Write(writer, binary.BigEndian, data)
	if err != nil {
		panic(err)
	}
}
