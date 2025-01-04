package brlyt

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

type TevSwapModeTable struct {
	B1 uint8
	B2 uint8
	B3 uint8
	B4 uint8
}

type MATIndirectTextureOrderEntry struct {
	TexCoord uint8
	TexMap   uint8
	ScaleS   uint8
	ScaleT   uint8
}

type MATTevStageEntry struct {
	TexCoor uint8
	Color   uint8
	U16     uint16
	B1      uint8
	B2      uint8
	B3      uint8
	B4      uint8
	B5      uint8
	B6      uint8
	B7      uint8
	B8      uint8
	B9      uint8
	B10     uint8
	B11     uint8
	B12     uint8
}

type MatAlphaCompare struct {
	Temp    uint8
	AlphaOP uint8
	Ref0    uint8
	Ref1    uint8
}

func (r *Root) ParseMAT(data []byte, sectionSize uint32) error {
	var mat MAT
	var materialOffsets []uint32
	var matEntries []MATEntries

	err := binary.Read(bytes.NewReader(data), binary.BigEndian, &mat)
	if err != nil {
		return err
	}

	for i := 0; i < int(mat.NumOfMats); i++ {
		var matOffset MATOffset
		offset := 4 + (i * 4)

		err = binary.Read(bytes.NewReader(data[offset:]), binary.BigEndian, &matOffset)
		if err != nil {
			return err
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

		err = binary.Read(bytes.NewReader(data[materialOffsets[i]:materialOffsets[i+1]]), binary.BigEndian, &matMaterials)
		if err != nil {
			return err
		}

		// Read the bitfield
		offset := materialOffsets[i] + 64
		textureEntries := make([]MATTexture, BitExtract(matMaterials.BitFlag, 28, 31))

		for i := 0; i < BitExtract(matMaterials.BitFlag, 28, 31); i++ {
			var texEntry MATTextureEntry

			err = binary.Read(bytes.NewReader(data[offset:]), binary.BigEndian, &texEntry)
			if err != nil {
				return err
			}

			xmlTexture := MATTexture{
				Name:  "nil",
				SWrap: texEntry.SWrap,
				TWrap: texEntry.TWrap,
			}

			if r.TXL.TPLName != nil {
				xmlTexture.Name = r.TXL.TPLName[texEntry.TexIndex]
			}

			offset += 4
			textureEntries[i] = xmlTexture
		}

		textureSRTEntries := make([]MATSRT, BitExtract(matMaterials.BitFlag, 24, 27))
		for i := 0; i < BitExtract(matMaterials.BitFlag, 24, 27); i++ {
			var texSRTEntry MATTextureSRTEntry

			err = binary.Read(bytes.NewReader(data[offset:]), binary.BigEndian, &texSRTEntry)
			if err != nil {
				return err
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

			err = binary.Read(bytes.NewReader(data[offset:]), binary.BigEndian, &texCoorGenEntry)
			if err != nil {
				return err
			}

			xmlCoorGen := MATCoordGen{
				Type:         texCoorGenEntry.Type,
				Source:       texCoorGenEntry.Source,
				MatrixSource: texCoorGenEntry.MatrixSource,
			}

			offset += 4
			texCoorGenEntries[i] = xmlCoorGen
		}

		var chanControlXML *ChanControlXML
		if BitExtract(matMaterials.BitFlag, 6, 100) == 1 {
			var chanControl MATChanControl

			err = binary.Read(bytes.NewReader(data[offset:]), binary.BigEndian, &chanControl)
			if err != nil {
				return err
			}

			chanControlXML = &ChanControlXML{
				ColorMaterialSource: chanControl.ColorMaterialSource,
				AlphaMaterialSource: chanControl.ColorMaterialSource,
			}

			offset += 4
		}

		var matColorXML *Color8
		if BitExtract(matMaterials.BitFlag, 4, 100) == 1 {
			var matColor MATColor

			err = binary.Read(bytes.NewReader(data[offset:]), binary.BigEndian, &matColor)
			if err != nil {
				return err
			}

			matColorXML = &Color8{
				R: matColor.R,
				G: matColor.G,
				B: matColor.B,
				A: matColor.A,
			}

			offset += 4
		}

		var tevSwapModeTableXML *TevSwapModeTableXML
		if BitExtract(matMaterials.BitFlag, 19, 100) == 1 {
			var tevSwapModeTable TevSwapModeTable

			err = binary.Read(bytes.NewReader(data[offset:]), binary.BigEndian, &tevSwapModeTable)
			if err != nil {
				return err
			}

			tevSwapModeTableXML = &TevSwapModeTableXML{
				AR: (tevSwapModeTable.B1 >> 0) & 0x3,
				AG: (tevSwapModeTable.B1 >> 2) & 0x3,
				AB: (tevSwapModeTable.B1 >> 4) & 0x3,
				AA: (tevSwapModeTable.B1 >> 6) & 0x3,
				BR: (tevSwapModeTable.B2 >> 0) & 0x3,
				BG: (tevSwapModeTable.B2 >> 2) & 0x3,
				BB: (tevSwapModeTable.B2 >> 4) & 0x3,
				BA: (tevSwapModeTable.B2 >> 6) & 0x3,
				CR: (tevSwapModeTable.B3 >> 0) & 0x3,
				CG: (tevSwapModeTable.B3 >> 2) & 0x3,
				CB: (tevSwapModeTable.B3 >> 4) & 0x3,
				CA: (tevSwapModeTable.B3 >> 6) & 0x3,
				DR: (tevSwapModeTable.B4 >> 0) & 0x3,
				DG: (tevSwapModeTable.B4 >> 2) & 0x3,
				DB: (tevSwapModeTable.B4 >> 4) & 0x3,
				DA: (tevSwapModeTable.B4 >> 6) & 0x3,
			}

			offset += 4
		}

		indirectTextureSRTEntries := make([]MATSRT, BitExtract(matMaterials.BitFlag, 17, 18))
		for i := 0; i < BitExtract(matMaterials.BitFlag, 17, 18); i++ {
			var texSRTEntry MATTextureSRTEntry

			err = binary.Read(bytes.NewReader(data[offset:]), binary.BigEndian, &texSRTEntry)
			if err != nil {
				return err
			}

			xmlSRT := MATSRT{
				XTrans:   texSRTEntry.XTrans,
				YTrans:   texSRTEntry.YTrans,
				Rotation: texSRTEntry.Rotation,
				XScale:   texSRTEntry.XScale,
				YScale:   texSRTEntry.YScale,
			}

			offset += 20
			indirectTextureSRTEntries[i] = xmlSRT
		}

		indirectTextureOrderEntries := make([]MATIndirectOrderEntryXML, BitExtract(matMaterials.BitFlag, 14, 16))
		for i := 0; i < BitExtract(matMaterials.BitFlag, 14, 16); i++ {
			var texOrderEntry MATIndirectTextureOrderEntry

			err = binary.Read(bytes.NewReader(data[offset:]), binary.BigEndian, &texOrderEntry)
			if err != nil {
				return err
			}

			xmlEntry := MATIndirectOrderEntryXML{
				TexCoord: texOrderEntry.TexCoord,
				TexMap:   texOrderEntry.TexMap,
				ScaleS:   texOrderEntry.ScaleS,
				ScaleT:   texOrderEntry.ScaleT,
			}

			offset += 4
			indirectTextureOrderEntries[i] = xmlEntry
		}

		tevStageEntries := make([]MATTevStageEntryXML, BitExtract(matMaterials.BitFlag, 9, 13))
		for i := 0; i < BitExtract(matMaterials.BitFlag, 9, 13); i++ {
			var temp MATTevStageEntry

			err = binary.Read(bytes.NewReader(data[offset:]), binary.BigEndian, &temp)
			if err != nil {
				return err
			}

			colorClamp := 0
			if temp.B4&0x1 == 1 {
				colorClamp = 1
			}

			alphaClamp := 0
			if temp.B8&0x1 == 1 {
				alphaClamp = 1
			}

			entry := MATTevStageEntryXML{
				TexCoor:          temp.TexCoor,
				Color:            temp.Color,
				TexMap:           temp.U16 & 0x1ff,
				RasSel:           uint8((temp.U16 & 0x7ff) >> 9),
				TexSel:           uint8(temp.U16 >> 11),
				ColorA:           temp.B1 & 0xf,
				ColorB:           temp.B1 >> 4,
				ColorC:           temp.B2 & 0xf,
				ColorD:           temp.B2 >> 4,
				ColorOP:          temp.B3 & 0xf,
				ColorBias:        (temp.B3 & 0x3f) >> 4,
				ColorScale:       temp.B3 >> 6,
				ColorClamp:       uint8(colorClamp),
				ColorRegID:       (temp.B4 & 0x7) >> 1,
				ColorConstantSel: temp.B4 >> 3,
				AlphaA:           temp.B5 & 0xf,
				AlphaB:           temp.B5 >> 4,
				AlphaC:           temp.B6 & 0xf,
				AlphaD:           temp.B6 >> 4,
				AlphaOP:          temp.B7 & 0xf,
				AlphaBias:        (temp.B7 & 0x3f) >> 4,
				AlphaScale:       temp.B7 >> 6,
				AlphaClamp:       uint8(alphaClamp),
				AlphaRegID:       (temp.B8 & 0x7) >> 1,
				AlphaConstantSel: temp.B8 >> 3,
				TexID:            temp.B9 & 0x3,
				Bias:             temp.B10 & 0x7,
				Matrix:           (temp.B10 & 0x7F) >> 3,
				WrapS:            temp.B11 & 0x7,
				WrapT:            (temp.B11 & 0x3F) >> 3,
				Format:           temp.B12 & 0x3,
				AddPrevious:      (temp.B12 & 0x7) >> 2,
				UTCLod:           (temp.B12 & 0xF) >> 3,
				Alpha:            (temp.B12 & 0x3F) >> 4,
			}

			tevStageEntries[i] = entry

			offset += 16
		}

		var alphaCompareXML *MATAlphaCompareXML
		if BitExtract(matMaterials.BitFlag, 8, 8) == 1 {
			var alphaCompare MatAlphaCompare

			err = binary.Read(bytes.NewReader(data[offset:]), binary.BigEndian, &alphaCompare)
			if err != nil {
				return err
			}

			alphaCompareXML = &MATAlphaCompareXML{
				Comp0:   alphaCompare.Temp & 0x7,
				Comp1:   (alphaCompare.Temp >> 4) & 0x7,
				AlphaOP: alphaCompare.AlphaOP,
				Ref0:    alphaCompare.Ref0,
				Ref1:    alphaCompare.Ref1,
			}

			offset += 4
		}

		var blendModeXML *MATBlendMode
		if BitExtract(matMaterials.BitFlag, 7, 7) == 1 {
			var blendMode MATBlendMode

			err = binary.Read(bytes.NewReader(data[offset:]), binary.BigEndian, &blendMode)
			if err != nil {
				return err
			}

			blendModeXML = &blendMode
			offset += 4
		}

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
			BitFlag:              matMaterials.BitFlag,
			Textures:             textureEntries,
			SRT:                  textureSRTEntries,
			CoordGen:             texCoorGenEntries,
			ChanControl:          chanControlXML,
			MatColor:             matColorXML,
			TevSwapMode:          tevSwapModeTableXML,
			IndirectSRT:          indirectTextureSRTEntries,
			IndirectTextureOrder: indirectTextureOrderEntries,
			TevStageEntry:        tevStageEntries,
			AlphaCompare:         alphaCompareXML,
			BlendMode:            blendModeXML,
		}

		matEntries = append(matEntries, xmlData)
	}

	r.MAT = MATNode{Entries: matEntries}
	return nil
}

func (b *BRLYTWriter) WriteMAT(data Root) error {
	temp := bytes.NewBuffer(nil)

	header := SectionHeader{
		Type: SectionTypeMAT,
		Size: 0,
	}

	meta := MAT{NumOfMats: uint16(len(data.MAT.Entries))}

	offsets := make([]MATOffset, len(data.MAT.Entries))
	offsets[0].Offset = uint32(len(data.MAT.Entries)*4 + 12)

	count := 12 + (len(data.MAT.Entries) * 4)
	for i, entry := range data.MAT.Entries {
		if i != 0 {
			offsets[i].Offset = uint32(count)
		}

		var name [20]byte
		copy(name[:], entry.Name)

		material := MATMaterials{
			Name:      name,
			ForeColor: [4]int16{entry.ForeColor.R, entry.ForeColor.G, entry.ForeColor.B, entry.ForeColor.A},
			BackColor: [4]int16{entry.BackColor.R, entry.BackColor.G, entry.BackColor.B, entry.BackColor.A},
			ColorReg3: [4]int16{entry.ColorReg3.R, entry.ColorReg3.G, entry.ColorReg3.B, entry.ColorReg3.A},
			TevColor1: [4]uint8{entry.TevColor1.R, entry.TevColor1.G, entry.TevColor1.B, entry.TevColor1.A},
			TevColor2: [4]uint8{entry.TevColor2.R, entry.TevColor2.G, entry.TevColor2.B, entry.TevColor2.A},
			TevColor3: [4]uint8{entry.TevColor3.R, entry.TevColor3.G, entry.TevColor3.B, entry.TevColor3.A},
			TevColor4: [4]uint8{entry.TevColor4.R, entry.TevColor4.G, entry.TevColor4.B, entry.TevColor4.A},
			BitFlag:   entry.BitFlag,
		}

		err := write(temp, material)
		if err != nil {
			return err
		}

		count += 64

		for _, texture := range entry.Textures {
			for i2, s := range data.TXL.TPLName {
				if texture.Name == s {
					tex := MATTextureEntry{
						TexIndex: uint16(i2),
						SWrap:    texture.SWrap,
						TWrap:    texture.TWrap,
					}

					err = write(temp, tex)
					if err != nil {
						return err
					}

					count += 4
				}
			}
		}

		for _, srt := range entry.SRT {
			srtEntry := MATTextureSRTEntry{
				XTrans:   srt.XTrans,
				YTrans:   srt.YTrans,
				Rotation: srt.Rotation,
				XScale:   srt.XScale,
				YScale:   srt.YScale,
			}

			err = write(temp, srtEntry)
			if err != nil {
				return err
			}

			count += 20
		}

		for _, gen := range entry.CoordGen {
			coorEntry := MATTexCoordGenEntry{
				Type:         gen.Type,
				Source:       gen.Source,
				MatrixSource: gen.MatrixSource,
			}

			err = write(temp, coorEntry)
			if err != nil {
				return err
			}

			count += 4
		}

		if entry.ChanControl != nil {
			chanControl := MATChanControl{
				ColorMaterialSource: entry.ChanControl.ColorMaterialSource,
				AlphaMaterialSource: entry.ChanControl.AlphaMaterialSource,
			}

			err = write(temp, chanControl)
			if err != nil {
				return err
			}

			count += 4
		}

		if entry.MatColor != nil {
			matColor := MATColor{
				R: entry.MatColor.R,
				G: entry.MatColor.B,
				B: entry.MatColor.G,
				A: entry.MatColor.A,
			}

			err = write(temp, matColor)
			if err != nil {
				return err
			}

			count += 4
		}

		if entry.TevSwapMode != nil {
			var b1 uint8 = 0
			b1 |= (entry.TevSwapMode.AA & 0x3) << 6
			b1 |= (entry.TevSwapMode.AB & 0x3) << 4
			b1 |= (entry.TevSwapMode.AG & 0x3) << 2
			b1 |= (entry.TevSwapMode.AR & 0x3) << 0

			var b2 uint8 = 0
			b2 |= (entry.TevSwapMode.BA & 0x3) << 6
			b2 |= (entry.TevSwapMode.BB & 0x3) << 4
			b2 |= (entry.TevSwapMode.BG & 0x3) << 2
			b2 |= (entry.TevSwapMode.BR & 0x3) << 0

			var b3 uint8 = 0
			b3 |= (entry.TevSwapMode.CA & 0x3) << 6
			b3 |= (entry.TevSwapMode.CB & 0x3) << 4
			b3 |= (entry.TevSwapMode.CG & 0x3) << 2
			b3 |= (entry.TevSwapMode.CR & 0x3) << 0

			var b4 uint8 = 0
			b4 |= (entry.TevSwapMode.DA & 0x3) << 6
			b4 |= (entry.TevSwapMode.DB & 0x3) << 4
			b4 |= (entry.TevSwapMode.DG & 0x3) << 2
			b4 |= (entry.TevSwapMode.DR & 0x3) << 0

			tevSwap := TevSwapModeTable{
				B1: b1,
				B2: b2,
				B3: b3,
				B4: b4,
			}

			err = write(temp, tevSwap)
			if err != nil {
				return err
			}

			count += 4
		}

		for _, matsrt := range entry.IndirectSRT {
			srt := MATSRT{
				XTrans:   matsrt.XTrans,
				YTrans:   matsrt.YTrans,
				Rotation: matsrt.Rotation,
				XScale:   matsrt.XScale,
				YScale:   matsrt.YScale,
			}

			err = write(temp, srt)
			if err != nil {
				return err
			}

			count += 20
		}

		for _, tex := range entry.IndirectTextureOrder {
			indirectTex := MATIndirectTextureOrderEntry{
				TexCoord: tex.TexCoord,
				TexMap:   tex.TexMap,
				ScaleS:   tex.ScaleS,
				ScaleT:   tex.ScaleT,
			}

			err = write(temp, indirectTex)
			if err != nil {
				return err
			}

			count += 4
		}

		for _, stageEntry := range entry.TevStageEntry {
			var U16 uint16 = 0
			U16 |= uint16(stageEntry.TexSel&0x3F) << 11
			U16 |= uint16(stageEntry.RasSel&0x7) << 9
			U16 |= (stageEntry.TexMap & 0x1ff) << 0

			var B1 uint8 = 0
			B1 |= (stageEntry.ColorB & 0xf) << 4
			B1 |= (stageEntry.ColorA & 0xf) << 0

			var B2 uint8 = 0
			B2 |= (stageEntry.ColorD & 0xf) << 4
			B2 |= (stageEntry.ColorC & 0xf) << 0

			var B3 uint8 = 0
			B3 |= (stageEntry.ColorScale & 0x3) << 6
			B3 |= (stageEntry.ColorBias & 0x3) << 4
			B3 |= (stageEntry.ColorOP & 0xf) << 0

			var B4 uint8 = 0
			B4 |= (stageEntry.ColorConstantSel & 0x1F) << 3
			B4 |= (stageEntry.ColorRegID & 0x7) << 1
			B4 |= (stageEntry.ColorClamp) << 0

			var B5 uint8 = 0
			B5 |= (stageEntry.AlphaB & 0xf) << 4
			B5 |= (stageEntry.AlphaA & 0xf) << 0

			var B6 uint8 = 0
			B6 |= (stageEntry.AlphaD & 0xf) << 4
			B6 |= (stageEntry.AlphaC & 0xf) << 0

			var B7 uint8 = 0
			B7 |= (stageEntry.AlphaScale & 0x3) << 6
			B7 |= (stageEntry.AlphaBias & 0x3) << 4
			B7 |= (stageEntry.AlphaOP & 0xf) << 0

			var B8 uint8 = 0
			B8 |= (stageEntry.AlphaConstantSel & 0x1F) << 3
			B8 |= (stageEntry.AlphaRegID & 0x7) << 1
			B8 |= (stageEntry.AlphaClamp) << 0

			var B10 uint8 = 0
			B10 |= (stageEntry.Matrix & 0x1F) << 3
			B10 |= (stageEntry.Bias & 0x7) << 0

			var B11 uint8 = 0
			B11 |= (stageEntry.WrapT & 0x7) << 3
			B11 |= (stageEntry.WrapS & 0x7) << 0

			var B12 uint8 = 0
			B12 |= (stageEntry.Alpha & 0xF) << 4
			B12 |= (stageEntry.UTCLod & 0x1) << 3
			B12 |= (stageEntry.AddPrevious & 0x1) << 2
			B12 |= (stageEntry.Format & 0x3) << 0

			entry := MATTevStageEntry{
				TexCoor: stageEntry.TexCoor,
				Color:   stageEntry.Color,
				U16:     U16,
				B1:      B1,
				B2:      B2,
				B3:      B3,
				B4:      B4,
				B5:      B5,
				B6:      B6,
				B7:      B7,
				B8:      B8,
				B9:      stageEntry.TexID & 0x3,
				B10:     B10,
				B11:     B11,
				B12:     B12,
			}

			err = write(temp, entry)
			if err != nil {
				return err
			}

			count += 16
		}

		if entry.AlphaCompare != nil {
			var tempValue uint8 = 0
			tempValue |= (entry.AlphaCompare.Comp1 & 0x7) << 4
			tempValue |= (entry.AlphaCompare.Comp0 & 0x7) << 0

			entry := MatAlphaCompare{
				Temp:    tempValue,
				AlphaOP: entry.AlphaCompare.AlphaOP,
				Ref0:    entry.AlphaCompare.Ref0,
				Ref1:    entry.AlphaCompare.Ref1,
			}

			err = write(temp, entry)
			if err != nil {
				return err
			}

			count += 4
		}

		if entry.BlendMode != nil {
			err = write(temp, entry.BlendMode)
			if err != nil {
				return err
			}

			count += 4
		}
	}

	header.Size = uint32(count)
	err := write(b, header)
	if err != nil {
		return err
	}

	err = write(b, meta)
	if err != nil {
		return err
	}

	err = write(b, offsets)
	if err != nil {
		return err
	}

	return write(b, temp.Bytes())
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
