package brlyt

import (
	"bytes"
	"encoding/binary"
	"io"
	"strings"
)

func (r *Root) ParsePAN(data []byte) (*XMLPane, error) {
	var pane Pane
	err := binary.Read(bytes.NewReader(data), binary.BigEndian, &pane)
	if err != nil {
		return nil, err
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

	if name == "RootPane" {
		r.RootPane = xmlData
		r.count--

		_, err = r.reader.Seek(8, io.SeekCurrent)
		if err != nil {
			return nil, err
		}

		r.RootPane.Children, err = r.ParseChildren()
		if err != nil {
			return nil, err
		}
	} else {
		// Only the Root Pane is guaranteed to have children, peek to see if this pane does.
		if r.HasChildren() {
			xmlData.Children, err = r.ParseChildren()
			if err != nil {
				return nil, err
			}
		}
	}

	return &xmlData, nil
}

func (r *Root) ParseBND(data []byte) (*XMLPane, error) {
	var pane Pane
	err := binary.Read(bytes.NewReader(data), binary.BigEndian, &pane)
	if err != nil {
		return nil, err
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

	if r.HasChildren() {
		xmlData.Children, err = r.ParseChildren()
		if err != nil {
			return nil, err
		}
	}

	return &xmlData, nil
}

func (b *BRLYTWriter) WritePane(pan XMLPane) error {
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

	err := write(b, header)
	if err != nil {
		return err
	}

	return write(b, pane)
}

func (b *BRLYTWriter) WriteBND(pan XMLPane) error {
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

	err := write(b, header)
	if err != nil {
		return err
	}

	return write(b, pane)
}

func (b *BRLYTWriter) WritePAS() error {
	header := SectionHeader{
		Type: SectionTypePAS,
		Size: 8,
	}

	return write(b, header)
}

func (b *BRLYTWriter) WritePAE() error {
	header := SectionHeader{
		Type: SectionTypePAE,
		Size: 8,
	}

	return write(b, header)
}
