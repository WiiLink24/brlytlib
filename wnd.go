package brlyt

import (
	"bytes"
	"encoding/binary"
	"strings"
)

func (r *Root) ParseWND(data []byte) (*XMLWND, error) {
	var wnd Window
	err := binary.Read(bytes.NewReader(data), binary.BigEndian, &wnd)
	if err != nil {
		return nil, err
	}

	// Strip the null bytes from the strings
	name := strings.Replace(string(wnd.PaneName[:]), "\x00", "", -1)
	userData := strings.Replace(string(wnd.UserData[:]), "\x00", "", -1)

	// Get the UVSets
	uvSets := make([]XMLUVSet, wnd.NumOfUVSets)
	for i := 0; i < int(wnd.NumOfUVSets); i++ {
		offset := 116 + (i * 32)

		var uv UVSet
		err = binary.Read(bytes.NewReader(data[offset:]), binary.BigEndian, &uv)
		if err != nil {
			return nil, err
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
		err = binary.Read(bytes.NewReader(data[offset:]), binary.BigEndian, &actualOffset)
		if err != nil {
			return nil, err
		}

		var mat WindowMat
		err = binary.Read(bytes.NewReader(data[actualOffset-8:]), binary.BigEndian, &mat)
		if err != nil {
			return nil, err
		}

		material := XMLWindowMat{
			MatIndex: mat.MatIndex,
			Index:    mat.Index,
		}

		mats[i] = material
	}

	xmlData := XMLWND{
		Name:        name,
		UserData:    userData,
		Visible:     wnd.Flag & 0x1,
		Widescreen:  (wnd.Flag & 0x2) >> 1,
		Flag:        wnd.Flag,
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

	if r.HasChildren() {
		xmlData.Children, err = r.ParseChildren()
		if err != nil {
			return nil, err
		}
	}

	return &xmlData, nil
}

func (b *BRLYTWriter) WriteWND(data XMLWND) error {
	temp := bytes.NewBuffer(nil)

	header := SectionHeader{
		Type: SectionTypeWND,
		Size: 76,
	}

	var name [16]byte
	copy(name[:], data.Name)

	var userData [8]byte
	copy(userData[:], data.UserData)

	wnd := Window{
		Flag:              data.Flag,
		Origin:            uint8(data.Origin.X + (data.Origin.Y * 3)),
		Alpha:             data.Alpha,
		PaneName:          name,
		UserData:          userData,
		XTranslation:      data.Translate.X,
		YTranslation:      data.Translate.Y,
		ZTranslation:      data.Translate.Z,
		XRotate:           data.Rotate.X,
		YRotate:           data.Rotate.Y,
		ZRotate:           data.Rotate.Z,
		XScale:            data.Scale.X,
		YScale:            data.Scale.Y,
		Width:             data.Width,
		Height:            data.Height,
		Coordinate1:       data.Coordinate1,
		Coordinate2:       data.Coordinate2,
		Coordinate3:       data.Coordinate3,
		Coordinate4:       data.Coordinate4,
		FrameCount:        uint8(len(data.Materials.Mats)),
		WindowOffset:      104,
		WindowFrameOffset: uint32(124 + len(data.UVSets.Set)*32),
		TopLeftColor:      [4]uint8{data.TopLeftColor.R, data.TopLeftColor.G, data.TopLeftColor.B, data.TopLeftColor.A},
		TopRightColor:     [4]uint8{data.TopRightColor.R, data.TopRightColor.G, data.TopRightColor.B, data.TopRightColor.A},
		BottomLeftColor:   [4]uint8{data.BottomLeftColor.R, data.BottomLeftColor.G, data.BottomLeftColor.B, data.BottomLeftColor.A},
		BottomRightColor:  [4]uint8{data.BottomRightColor.R, data.BottomRightColor.G, data.BottomRightColor.B, data.BottomRightColor.A},
		MatIndex:          data.MatIndex,
		NumOfUVSets:       uint8(len(data.UVSets.Set)),
	}

	err := write(temp, header)
	if err != nil {
		return err
	}

	err = write(temp, wnd)
	if err != nil {
		return err
	}

	// Write the UV Sets
	for _, set := range data.UVSets.Set {
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

		err = write(temp, uvSet)
		if err != nil {
			return err
		}
	}

	// Write the offsets to the Window Mats
	for i := 0; i < len(data.Materials.Mats); i++ {
		var offset uint32

		offset = uint32(temp.Len() + (len(data.Materials.Mats) * 4))

		err = write(temp, offset)
		if err != nil {
			return err
		}
	}

	// Write Window Mats
	for _, mat := range data.Materials.Mats {
		windowMat := WindowMat{
			MatIndex: mat.MatIndex,
			Index:    mat.Index,
		}

		err = write(temp, windowMat)
		if err != nil {
			return err
		}
	}

	binary.BigEndian.PutUint32(temp.Bytes()[4:8], uint32(temp.Len()))

	return write(b, temp.Bytes())
}
