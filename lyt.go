package main

import (
	"bytes"
	"encoding/binary"
	"encoding/xml"
)

// LYT defines the main layout of the BRLYT
type LYT struct {
	Centered [1]byte
	Padding  [3]byte
	Width    float32
	Height   float32
}

func (r *Root) ParseLYT(data []byte) {
	// Parse LYT section
	readable := bytes.NewReader(data)

	var lyt LYT
	err := binary.Read(readable, binary.BigEndian, &lyt)
	if err != nil {
		panic(err)
	}

	r.LYT = LYTNode{
		XMLName:  xml.Name{},
		Centered: uint16(lyt.Centered[0]),
		Width:    lyt.Width,
		Height:   lyt.Height,
	}
}

func (b *BRLYTWriter) WriteLYT(data Root) {
	header := SectionHeader{
		Type: SectionTypeLYT,
		Size: 20,
	}

	lyt := LYT{
		Centered: [1]byte{byte(data.LYT.Centered)},
		Padding:  [3]byte{},
		Width:    data.LYT.Width,
		Height:   data.LYT.Height,
	}

	write(b, header)
	write(b, lyt)
}
