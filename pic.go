package main

import (
	"bytes"
	"encoding/binary"
	"strings"
)

// PIC defines the image pane in a brlyt
type PIC struct {
	Flag             uint8
	Origin           uint8
	Alpha            uint8
	_                uint8
	PaneName         [16]byte
	UserData         [8]byte
	XTranslation     float32
	YTranslation     float32
	ZTranslation     float32
	XRotate          float32
	YRotate          float32
	ZRotate          float32
	XScale           float32
	YScale           float32
	Width            float32
	Height           float32
	TopLeftColor     [4]uint8
	TopRightColor    [4]uint8
	BottomLeftColor  [4]uint8
	BottomRightColor [4]uint8
	MatIndex         uint16
	NumOfUVSets      uint8
	_                uint8
}

type UVSet struct {
	TopLeftS     float32
	TopLeftT     float32
	TopRightS    float32
	TopRightT    float32
	BottomLeftS  float32
	BottomLeftT  float32
	BottomRightS float32
	BottomRightT float32
}

func (r *Root) ParsePIC(data []byte) {

	var pic PIC
	err := binary.Read(bytes.NewReader(data), binary.BigEndian, &pic)
	if err != nil {
		panic(err)
	}

	// Strip the null bytes from the strings
	name := strings.Replace(string(pic.PaneName[:]), "\x00", "", -1)
	userData := strings.Replace(string(pic.UserData[:]), "\x00", "", -1)

	// Get the UVSets
	uvSets := make([]XMLUVSet, pic.NumOfUVSets)
	for i := 0; i < int(pic.NumOfUVSets); i++ {
		offset := 88 + (i * 32)

		var uv UVSet
		err := binary.Read(bytes.NewReader(data[offset:]), binary.BigEndian, &uv)
		if err != nil {
			panic(err)
		}

		set := XMLUVSet{
			CoordTL: STCoordinates{
				T: uv.TopLeftT,
				S: uv.TopLeftS,
			},
			CoordTR: STCoordinates{
				T: uv.TopRightT,
				S: uv.TopRightS,
			},
			CoordBL: STCoordinates{
				T: uv.BottomLeftT,
				S: uv.BottomLeftS,
			},
			CoordBR: STCoordinates{
				T: uv.BottomRightT,
				S: uv.BottomRightS,
			},
		}

		uvSets[i] = set
	}

	xmlData := XMLPIC{
		Name:       name,
		UserData:   userData,
		Visible:    pic.Flag & 0x1,
		Widescreen: (pic.Flag & 0x2) >> 1,
		Flag:       pic.Flag,
		Origin:     Coord2D{X: float32(pic.Origin % 3), Y: float32(pic.Origin / 3)},
		Alpha:      pic.Alpha,
		Padding:    0,
		Translate:  Coord3D{X: pic.XTranslation, Y: pic.YTranslation, Z: pic.ZTranslation},
		Rotate:     Coord3D{X: pic.XRotate, Y: pic.YRotate, Z: pic.ZRotate},
		Scale:      Coord2D{X: pic.XScale, Y: pic.YScale},
		Width:      pic.Width,
		Height:     pic.Height,
		TopLeftColor: Color8{
			R: pic.TopLeftColor[0],
			G: pic.TopLeftColor[1],
			B: pic.TopLeftColor[2],
			A: pic.TopLeftColor[3],
		},
		TopRightColor: Color8{
			R: pic.TopRightColor[0],
			G: pic.TopRightColor[1],
			B: pic.TopRightColor[2],
			A: pic.TopRightColor[3],
		},
		BottomLeftColor: Color8{
			R: pic.BottomLeftColor[0],
			G: pic.BottomLeftColor[1],
			B: pic.BottomLeftColor[2],
			A: pic.BottomLeftColor[3],
		},
		BottomRightColor: Color8{
			R: pic.BottomRightColor[0],
			G: pic.BottomRightColor[1],
			B: pic.BottomRightColor[2],
			A: pic.BottomRightColor[3],
		},
		MatIndex: pic.MatIndex,
		UVSets:   &XMLUVSets{Set: uvSets},
	}

	r.Panes = append(r.Panes, Children{PIC: &xmlData})
}

func (b *BRLYTWriter) WritePIC(pic XMLPIC) {
	header := SectionHeader{
		Type: SectionTypePIC,
		Size: uint32(96 + (32 * len(pic.UVSets.Set))),
	}
	var name [16]byte
	copy(name[:], pic.Name)

	var userData [8]byte
	copy(userData[:], pic.UserData)

	pane := PIC{
		Flag:             pic.Flag,
		Origin:           uint8(pic.Origin.X + (pic.Origin.Y * 3)),
		Alpha:            pic.Alpha,
		PaneName:         name,
		UserData:         userData,
		XTranslation:     pic.Translate.X,
		YTranslation:     pic.Translate.Y,
		ZTranslation:     pic.Translate.Z,
		XRotate:          pic.Rotate.X,
		YRotate:          pic.Rotate.Y,
		ZRotate:          pic.Rotate.Z,
		XScale:           pic.Scale.X,
		YScale:           pic.Scale.Y,
		Width:            pic.Width,
		Height:           pic.Height,
		TopLeftColor:     [4]uint8{pic.TopLeftColor.R, pic.TopLeftColor.G, pic.TopLeftColor.B, pic.TopLeftColor.A},
		TopRightColor:    [4]uint8{pic.TopRightColor.R, pic.TopRightColor.G, pic.TopRightColor.B, pic.TopRightColor.A},
		BottomLeftColor:  [4]uint8{pic.BottomLeftColor.R, pic.BottomLeftColor.G, pic.BottomLeftColor.B, pic.BottomLeftColor.A},
		BottomRightColor: [4]uint8{pic.BottomRightColor.R, pic.BottomRightColor.G, pic.BottomRightColor.B, pic.BottomRightColor.A},
		MatIndex:         pic.MatIndex,
		NumOfUVSets:      uint8(len(pic.UVSets.Set)),
	}

	write(b, header)
	write(b, pane)

	// Write the UV Sets
	for _, set := range pic.UVSets.Set {
		uvSet := UVSet{
			TopLeftS:     set.CoordTL.S,
			TopLeftT:     set.CoordTL.T,
			TopRightS:    set.CoordTR.S,
			TopRightT:    set.CoordTR.T,
			BottomLeftS:  set.CoordBL.S,
			BottomLeftT:  set.CoordBL.T,
			BottomRightS: set.CoordBR.S,
			BottomRightT: set.CoordBR.T,
		}

		write(b, uvSet)
	}
}
