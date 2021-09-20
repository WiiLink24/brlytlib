package brlytlib

import (
	"bytes"
	"encoding/binary"
	"encoding/xml"
	"io/ioutil"
)

// Header represents the header of our BRLYT
type Header struct {
	Magic            [4]byte
	_                uint32
	FileLength       uint32
	_                uint16
	NumberOfSections uint16
}

// LYT defines the main layout of the BRLYT
type LYT struct {
	Magic [4]byte
	_     uint32
	_     [4]byte
	_     float32
	_     float32
}

// BRLYT is the internal structure of our BRLYT
type BRLYT struct {
	Header Header
	LYT1   LYT
}

// SectionMagic are the known parts of the BRLYT
type SectionMagic [4]byte

var (
	headerMagic SectionMagic = [4]byte{'R', 'L', 'Y', 'T'}
	txlMagic    SectionMagic = [4]byte{'t', 'x', 'l', '1'}
)

func ParseBRLYT(fileName string) ([]byte, error) {
	contents, err := ioutil.ReadFile(fileName)
	var brlyt BRLYT
	err = binary.Read(bytes.NewReader(contents), binary.BigEndian, &brlyt)
	if err != nil {
		panic(err)
	}

	if brlyt.Header.Magic != headerMagic {
		return nil, ErrInvalidFileMagic
	}

	if int(brlyt.Header.FileLength) != len(contents) {
		return nil, ErrFileSizeMismatch
	}

	// Parse the sections. Every BRLYT contains these sections
	txl, err := ParseTXL(contents)
	if err != nil {
		return nil, err
	}

	txt, err := ParseTXT(contents)
	if err != nil {
		return nil, err
	}

	pan, err := ParsePAN(contents)
	if err != nil {
		return nil, err
	}

	// Now that everything has been parsed, append the panes into 1 slice
	for _, pane := range txt {
		pan = append(pan, pane)
	}

	// Finally, build the xml
	return xml.MarshalIndent(Root{
		XMLName: xml.Name{},
		TPLName: txl,
		Panes:   pan,
	}, "", "\t")
}
