package main

import (
	"bytes"
	"encoding/binary"
	"strings"
)

type MAT struct {
	NumOfMats uint16
	_         uint16
}

type MATOffset struct {
	Offset uint32
}

type MATMaterials struct {
	Name      [20]byte
	ForeColor [4]int16
	BackColor [4]int16
	ColorReg3 [4]int16
	TevColor1 [4]uint8
	TevColor2 [4]uint8
	TevColor3 [4]uint8
	TevColor4 [4]uint8
	BitFlag   uint32
}

type MATTextureEntry struct {
	TexIndex uint16
	SWrap    uint8
	TWrap    uint8
}

type MATTextureSRTEntry struct {
	XTrans   float32
	YTrans   float32
	Rotation float32
	XScale   float32
	YScale   float32
}

type MATTexCoordGenEntry struct {
	Type         uint8
	Source       uint8
	MatrixSource uint8
	_            uint8
}

type MATChanControl struct {
	ColorMaterialSource uint8
	AlphaMaterialSource uint8
	_                   uint8
	_                   uint8
}

type MATColor struct {
	R uint8
	G uint8
	B uint8
	A uint8
}

type MATIndirectTextureOrderEntry struct {
	TexCoord uint8
	TexMap   uint8
	ScaleS   uint8
	ScaleT   uint8
}

func (r *Root) ParseMAT(data []byte, sectionSize uint32) {
	var mat MAT
	var materialOffsets []uint32
	var matEntries []MATEntries

	err := binary.Read(bytes.NewReader(data), binary.BigEndian, &mat)
	if err != nil {
		panic(err)
	}

	for i := 0; i < int(mat.NumOfMats); i++ {
		var matOffset MATOffset
		offset := 4 + (i * 4)

		err := binary.Read(bytes.NewReader(data[offset:]), binary.BigEndian, &matOffset)
		if err != nil {
			panic(err)
		}

		materialOffsets = append(materialOffsets, matOffset.Offset-8)

		// If we have reached the last index, append the section size to the slice.
		if i == int(mat.NumOfMats)-1 {
			materialOffsets = append(materialOffsets, sectionSize-8)
		}
	}

	// Now that we have the offsets, parse the mat section.
	for i := 0; i < int(mat.NumOfMats); i++ {
		var matMaterials MATMaterials

		err := binary.Read(bytes.NewReader(data[materialOffsets[i]:materialOffsets[i+1]]), binary.BigEndian, &matMaterials)
		if err != nil {
			panic(err)
		}

		// Read the bitfield
		offset := materialOffsets[i] + 64
		textureEntries := make([]MATTexture, BitExtract(matMaterials.BitFlag, 28, 31))

		for i := 0; i < BitExtract(matMaterials.BitFlag, 28, 31); i++ {
			var texEntry MATTextureEntry

			err := binary.Read(bytes.NewReader(data[offset:]), binary.BigEndian, &texEntry)
			if err != nil {
				panic(err)
			}

			xmlTexture := MATTexture{
				Name:  "nil",
				SWrap: texEntry.SWrap,
				TWrap: texEntry.TWrap,
			}

			if r.TXL.TPLName != nil {
				xmlTexture.Name = r.TXL.TPLName[i]
			}

			offset += 4
			textureEntries[i] = xmlTexture
		}

		textureSRTEntries := make([]MATSRT, BitExtract(matMaterials.BitFlag, 24, 27))
		for i := 0; i < BitExtract(matMaterials.BitFlag, 24, 27); i++ {
			var texSRTEntry MATTextureSRTEntry

			err := binary.Read(bytes.NewReader(data[offset:]), binary.BigEndian, &texSRTEntry)
			if err != nil {
				panic(err)
			}

			xmlSRT := MATSRT{
				XTrans:   texSRTEntry.XTrans,
				YTrans:   texSRTEntry.YTrans,
				Rotation: texSRTEntry.Rotation,
				XScale:   texSRTEntry.XScale,
				YScale:   texSRTEntry.YScale,
			}

			offset += 20
			textureSRTEntries[i] = xmlSRT
		}

		texCoorGenEntries := make([]MATCoordGen, BitExtract(matMaterials.BitFlag, 20, 23))
		for i := 0; i < BitExtract(matMaterials.BitFlag, 20, 23); i++ {
			var texCoorGenEntry MATTexCoordGenEntry

			err := binary.Read(bytes.NewReader(data[offset:]), binary.BigEndian, &texCoorGenEntry)
			if err != nil {
				panic(err)
			}

			xmlCoorGen := MATCoordGen{
				Type:         texCoorGenEntry.Type,
				Source:       texCoorGenEntry.Source,
				MatrixSource: texCoorGenEntry.MatrixSource,
			}

			offset += 4
			texCoorGenEntries[i] = xmlCoorGen
		}

		// TODO: Implement the below into the XML. I have not seen these types so they are not priority

		var chanControl MATChanControl
		if BitExtract(matMaterials.BitFlag, 6, 100) == 1 {
			err := binary.Read(bytes.NewReader(data[offset:]), binary.BigEndian, &chanControl)
			if err != nil {
				panic(err)
			}

			offset += 4
		}

		var matColor MATColor
		if BitExtract(matMaterials.BitFlag, 4, 100) == 1 {
			err := binary.Read(bytes.NewReader(data[offset:]), binary.BigEndian, &matColor)
			if err != nil {
				panic(err)
			}

			offset += 4
		}

		// TODO: Implement MATTevSwapModeTable
		indirectTextureSRTEntries := make([]MATTextureSRTEntry, BitExtract(matMaterials.BitFlag, 17, 18))
		for i := 0; i < BitExtract(matMaterials.BitFlag, 17, 18); i++ {
			var texSRTEntry MATTextureSRTEntry

			err := binary.Read(bytes.NewReader(data[offset:]), binary.BigEndian, &texSRTEntry)
			if err != nil {
				panic(err)
			}

			offset += 20
			indirectTextureSRTEntries = append(indirectTextureSRTEntries, texSRTEntry)
		}

		indirectTextureOrderEntries := make([]MATIndirectTextureOrderEntry, BitExtract(matMaterials.BitFlag, 14, 16))
		for i := 0; i < BitExtract(matMaterials.BitFlag, 14, 16); i++ {
			var indirectTexOrderEntry MATIndirectTextureOrderEntry

			err := binary.Read(bytes.NewReader(data[offset:]), binary.BigEndian, &indirectTexOrderEntry)
			if err != nil {
				panic(err)
			}

			offset += 4
			indirectTextureOrderEntries = append(indirectTextureOrderEntries, indirectTexOrderEntry)
		}

		// TODO: Implement MATTevStageEntry, MATAlphaCompare and MATBlendMode
		matName := string(matMaterials.Name[:])
		matName = strings.Replace(matName, "\x00", "", -1)

		xmlData := MATEntries{
			Name: matName,
			ForeColor: Color16{
				R: matMaterials.ForeColor[0],
				G: matMaterials.ForeColor[1],
				B: matMaterials.ForeColor[2],
				A: matMaterials.ForeColor[3],
			},
			BackColor: Color16{
				R: matMaterials.BackColor[0],
				G: matMaterials.BackColor[1],
				B: matMaterials.BackColor[2],
				A: matMaterials.BackColor[3],
			},
			ColorReg3: Color16{
				R: matMaterials.ColorReg3[0],
				G: matMaterials.ColorReg3[1],
				B: matMaterials.ColorReg3[2],
				A: matMaterials.ColorReg3[3],
			},
			TevColor1: Color8{
				R: matMaterials.TevColor1[0],
				G: matMaterials.TevColor1[1],
				B: matMaterials.TevColor1[2],
				A: matMaterials.TevColor1[3],
			},
			TevColor2: Color8{
				R: matMaterials.TevColor2[0],
				G: matMaterials.TevColor2[1],
				B: matMaterials.TevColor2[2],
				A: matMaterials.TevColor2[3],
			},
			TevColor3: Color8{
				R: matMaterials.TevColor3[0],
				G: matMaterials.TevColor3[1],
				B: matMaterials.TevColor3[2],
				A: matMaterials.TevColor3[3],
			},
			TevColor4: Color8{
				R: matMaterials.TevColor4[0],
				G: matMaterials.TevColor4[1],
				B: matMaterials.TevColor4[2],
				A: matMaterials.TevColor4[3],
			},
			Textures: textureEntries,
			SRT:      textureSRTEntries,
			CoordGen: texCoorGenEntries,
		}

		matEntries = append(matEntries, xmlData)
	}

	r.MAT = MATNode{Entries: matEntries}
}

func BitExtract(num uint32, start int, end int) int {
	if end == 100 {
		end = start
	}

	firstMask := 1

	for first := 0; first < 31-start+1; first++ {
		firstMask *= 2
	}

	firstMask -= 1
	secondMask := 1

	for first := 0; first < 31-end; first++ {
		secondMask *= 2
	}

	secondMask -= 1

	mask := firstMask - secondMask

	return (int(num) & mask) >> (31 - end)
}
