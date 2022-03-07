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
