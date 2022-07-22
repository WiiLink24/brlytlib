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
		Name:      name,
		UserData:  userData,
		Flag:      pane.Flag,
		Origin:    Coord2D{X: float32(pane.Origin % 3), Y: float32(pane.Origin / 3)},
		Alpha:     pane.Alpha,
		Padding:   0,
		Translate: Coord3D{X: pane.XTranslation, Y: pane.YTranslation, Z: pane.ZTranslation},
		Rotate:    Coord3D{X: pane.XRotate, Y: pane.YRotate, Z: pane.ZRotate},
		Scale:     Coord2D{X: pane.XScale, Y: pane.YScale},
		Width:     pane.Width,
		Height:    pane.Height,
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
		Name:      name,
		UserData:  userData,
		Flag:      pane.Flag,
		Origin:    Coord2D{X: float32(pane.Origin % 3), Y: float32(pane.Origin / 3)},
		Alpha:     pane.Alpha,
		Padding:   0,
		Translate: Coord3D{X: pane.XTranslation, Y: pane.YTranslation, Z: pane.ZTranslation},
		Rotate:    Coord3D{X: pane.XRotate, Y: pane.YRotate, Z: pane.ZRotate},
		Scale:     Coord2D{X: pane.XScale, Y: pane.YScale},
		Width:     pane.Width,
		Height:    pane.Height,
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

func (b *BRLYTWriter) WritePane(pan XMLPane) {
	header := SectionHeader{
		Type: SectionTypePAN,
		Size: 76,
	}
	var name [16]byte
	copy(name[:], pan.Name)

	var userData [8]byte
	copy(userData[:], pan.UserData)

	pane := Pane{
		Flag:         pan.Flag,
		Origin:       uint8(pan.Origin.X + (pan.Origin.Y * 3)),
		Alpha:        pan.Alpha,
		PaneName:     name,
		UserData:     userData,
		XTranslation: pan.Translate.X,
		YTranslation: pan.Translate.Y,
		ZTranslation: pan.Translate.Z,
		XRotate:      pan.Rotate.X,
		YRotate:      pan.Rotate.Y,
		ZRotate:      pan.Rotate.Z,
		XScale:       pan.Scale.X,
		YScale:       pan.Scale.Y,
		Width:        pan.Width,
		Height:       pan.Height,
	}

	write(b, header)
	write(b, pane)
}

func (b *BRLYTWriter) WriteBND(pan XMLBND) {
	header := SectionHeader{
		Type: SectionTypeBND,
		Size: 76,
	}
	var name [16]byte
	copy(name[:], pan.Name)

	var userData [8]byte
	copy(userData[:], pan.UserData)

	pane := Pane{
		Flag:         pan.Flag,
		Origin:       uint8(pan.Origin.X + (pan.Origin.Y * 3)),
		Alpha:        pan.Alpha,
		PaneName:     name,
		UserData:     userData,
		XTranslation: pan.Translate.X,
		YTranslation: pan.Translate.Y,
		ZTranslation: pan.Translate.Z,
		XRotate:      pan.Rotate.X,
		YRotate:      pan.Rotate.Y,
		ZRotate:      pan.Rotate.Z,
		XScale:       pan.Scale.X,
		YScale:       pan.Scale.Y,
		Width:        pan.Width,
		Height:       pan.Height,
	}

	write(b, header)
	write(b, pane)
}

func (b *BRLYTWriter) WritePAS() {
	header := SectionHeader{
		Type: SectionTypePAS,
		Size: 8,
	}

	write(b, header)
}

func (b *BRLYTWriter) WritePAE() {
	header := SectionHeader{
		Type: SectionTypePAE,
		Size: 8,
	}

	write(b, header)
}
