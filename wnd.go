package main

import (
	"bytes"
	"encoding/binary"
	"strings"
)

type Window struct {
	Flag              uint8
	Origin            uint8
	Alpha             uint8
	_                 uint8
	PaneName          [16]byte
	UserData          [8]byte
	XTranslation      float32
	YTranslation      float32
	ZTranslation      float32
	XRotate           float32
	YRotate           float32
	ZRotate           float32
	XScale            float32
	YScale            float32
	Width             float32
	Height            float32
	Coordinate1       float32
	Coordinate2       float32
	Coordinate3       float32
	Coordinate4       float32
	FrameCount        uint8
	_                 [3]byte
	WindowOffset      uint32
	WindowFrameOffset uint32
	TopLeftColor      [4]uint8
	TopRightColor     [4]uint8
	BottomLeftColor   [4]uint8
	BottomRightColor  [4]uint8
	MatIndex          uint16
	NumOfUVSets       uint8
	_                 uint8
}

type WindowMat struct {
	MatIndex uint16
	Index    uint8
	_        uint8
}

func (r *Root) ParseWND(data []byte) {
	var wnd Window
	err := binary.Read(bytes.NewReader(data), binary.BigEndian, &wnd)
	if err != nil {
		panic(err)
	}

	// Strip the null bytes from the strings
	name := strings.Replace(string(wnd.PaneName[:]), "\x00", "", -1)
	userData := strings.Replace(string(wnd.UserData[:]), "\x00", "", -1)

	// Get the UVSets
	uvSets := make([]XMLUVSet, wnd.NumOfUVSets)
	for i := 0; i < int(wnd.NumOfUVSets); i++ {
		offset := 116 + (i * 32)

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

	// Parse Window Mat table
	mats := make([]XMLWindowMat, wnd.FrameCount)
	for i := 0; i < int(wnd.FrameCount); i++ {
		offset := 116 + (int(wnd.NumOfUVSets) * 32) + (i * 4)

		var actualOffset uint32
		err := binary.Read(bytes.NewReader(data[offset:]), binary.BigEndian, &actualOffset)
		if err != nil {
			panic(err)
		}

		var mat WindowMat
		err = binary.Read(bytes.NewReader(data[actualOffset-8:]), binary.BigEndian, &mat)
		if err != nil {
			panic(err)
		}

		material := XMLWindowMat{
			MatIndex: mat.MatIndex,
			Index:    mat.Index,
		}

		mats[i] = material
	}

	xmlData := XMLWindow{
		Name:        name,
		UserData:    userData,
		Visible:     wnd.Flag & 0x1,
		Widescreen:  (wnd.Flag & 0x2) >> 1,
		Flag:        (wnd.Flag & 0x4) >> 2,
		Origin:      Coord2D{X: float32(wnd.Origin % 3), Y: float32(wnd.Origin / 3)},
		Alpha:       wnd.Alpha,
		Padding:     0,
		Translate:   Coord3D{X: wnd.XTranslation, Y: wnd.YTranslation, Z: wnd.ZTranslation},
		Rotate:      Coord3D{X: wnd.XRotate, Y: wnd.YRotate, Z: wnd.ZRotate},
		Scale:       Coord2D{X: wnd.XScale, Y: wnd.YScale},
		Width:       wnd.Width,
		Height:      wnd.Height,
		Coordinate1: wnd.Coordinate1,
		Coordinate2: wnd.Coordinate2,
		Coordinate3: wnd.Coordinate3,
		Coordinate4: wnd.Coordinate4,
		TopLeftColor: Color8{
			R: wnd.TopLeftColor[0],
			G: wnd.TopLeftColor[1],
			B: wnd.TopLeftColor[2],
			A: wnd.TopLeftColor[3],
		},
		TopRightColor: Color8{
			R: wnd.TopRightColor[0],
			G: wnd.TopRightColor[1],
			B: wnd.TopRightColor[2],
			A: wnd.TopRightColor[3],
		},
		BottomLeftColor: Color8{
			R: wnd.BottomLeftColor[0],
			G: wnd.BottomLeftColor[1],
			B: wnd.BottomLeftColor[2],
			A: wnd.BottomLeftColor[3],
		},
		BottomRightColor: Color8{
			R: wnd.BottomRightColor[0],
			G: wnd.BottomRightColor[1],
			B: wnd.BottomRightColor[2],
			A: wnd.BottomRightColor[3],
		},
		MatIndex:  wnd.MatIndex,
		UVSets:    &XMLUVSets{Set: uvSets},
		Materials: &XMLWindowMats{Mats: mats},
	}

	r.Panes = append(r.Panes, Children{WND: &xmlData})
}
