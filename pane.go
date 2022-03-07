package main

import (
	"bytes"
	"encoding/binary"
	"strings"
)

// Pane represents the structure of a pan1 section.
type Pane struct {
	Flag         uint8
	Origin       uint8
	Alpha        uint8
	_            uint8
	PaneName     [16]byte
	UserData     [8]byte
	XTranslation float32
	YTranslation float32
	ZTranslation float32
	XRotate      float32
	YRotate      float32
	ZRotate      float32
	XScale       float32
	YScale       float32
	Width        float32
	Height       float32
}

func (r *Root) ParsePAN(data []byte) {
	var pane Pane
	err := binary.Read(bytes.NewReader(data), binary.BigEndian, &pane)
	if err != nil {
		panic(err)
	}

	// Strip the null bytes from the strings
	name := strings.Replace(string(pane.PaneName[:]), "\x00", "", -1)
	userData := strings.Replace(string(pane.UserData[:]), "\x00", "", -1)

	xmlData := XMLPane{
		Name:       name,
		UserData:   userData,
		Visible:    pane.Flag & 0x1,
		Widescreen: (pane.Flag & 0x2) >> 1,
		Flag:       (pane.Flag & 0x4) >> 2,
		Origin:     Coord2D{X: float32(pane.Origin % 3), Y: float32(pane.Origin / 3)},
		Alpha:      pane.Alpha,
		Padding:    0,
		Translate:  Coord3D{X: pane.XTranslation, Y: pane.YTranslation, Z: pane.ZTranslation},
		Rotate:     Coord3D{X: pane.XRotate, Y: pane.YRotate, Z: pane.ZRotate},
		Scale:      Coord2D{X: pane.XScale, Y: pane.YScale},
		Width:      pane.Width,
		Height:     pane.Height,
	}

	r.Panes = append(r.Panes, Children{Pane: &xmlData})
}

// ParseBND is almost 1:1 with ParsePAN as PAN and BND are the same types.
func (r *Root) ParseBND(data []byte) {
	var pane Pane
	err := binary.Read(bytes.NewReader(data), binary.BigEndian, &pane)
	if err != nil {
		panic(err)
	}

	// Strip the null bytes from the strings
	name := strings.Replace(string(pane.PaneName[:]), "\x00", "", -1)
	userData := strings.Replace(string(pane.UserData[:]), "\x00", "", -1)

	xmlData := XMLBND{
		Name:       name,
		UserData:   userData,
		Visible:    pane.Flag & 0x1,
		Widescreen: (pane.Flag & 0x2) >> 1,
		Flag:       (pane.Flag & 0x4) >> 2,
		Origin:     Coord2D{X: float32(pane.Origin % 3), Y: float32(pane.Origin / 3)},
		Alpha:      pane.Alpha,
		Padding:    0,
		Translate:  Coord3D{X: pane.XTranslation, Y: pane.YTranslation, Z: pane.ZTranslation},
		Rotate:     Coord3D{X: pane.XRotate, Y: pane.YRotate, Z: pane.ZRotate},
		Scale:      Coord2D{X: pane.XScale, Y: pane.YScale},
		Width:      pane.Width,
		Height:     pane.Height,
	}

	r.Panes = append(r.Panes, Children{BND: &xmlData})
}

func (r *Root) ParsePAS() {
	r.Panes = append(r.Panes, Children{
		PAS: &XMLPAS{},
	})
}

func (r *Root) ParsePAE() {
	r.Panes = append(r.Panes, Children{
		PAE: &XMLPAE{},
	})
}
