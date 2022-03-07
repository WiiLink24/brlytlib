package main

import (
	"bytes"
	"encoding/binary"
	"strings"
	"unicode/utf16"
)

const PlaceHolderString = "==== THIS IS PLACEHOLDER TEXT PLEASE DO NOT TRANSLATE ===="

// TXT represents the text data of the txt1 section
type TXT struct {
	Flag            uint8
	Origin          uint8
	Alpha           uint8
	_               uint8
	PaneName        [16]byte
	UserData        [8]byte
	XTranslation    float32
	YTranslation    float32
	ZTranslation    float32
	XRotate         float32
	YRotate         float32
	ZRotate         float32
	XScale          float32
	YScale          float32
	Width           float32
	Height          float32
	StringLength    uint16
	MaxStringLength uint16
	MatIndex        uint16
	FontIndex       uint16
	Alignment       uint8
	_               uint8
	_               uint16
	TextOffset      uint32
	Color1          [4]uint8
	Color2          [4]uint8
	FontSizeX       float32
	FontSizeY       float32
	CharacterSize   float32
	LineSize        float32
}

func (r *Root) ParseTXT(data []byte, sectionSize uint32) {
	var text TXT
	err := binary.Read(bytes.NewReader(data), binary.BigEndian, &text)
	if err != nil {
		panic(err)
	}

	// Strip the null bytes from the strings
	name := strings.Replace(string(text.PaneName[:]), "\x00", "", -1)
	userData := strings.Replace(string(text.UserData[:]), "\x00", "", -1)

	utf16String := data[text.TextOffset-8 : sectionSize-8]

	// Convert the UTF-16 string to UTF-8
	var full []uint16
	for i := 0; i < len(utf16String); i += 2 {
		current := binary.BigEndian.Uint16([]byte{utf16String[i], utf16String[i+1]})
		if current == 0 {
			// Our string was terminated
			break
		}
		full = append(full, current)
	}

	// Strip null bytes
	decodedString := strings.Replace(string(utf16.Decode(full)), "\x00", "", -1)
	// Replace newlines with \n
	decodedString = strings.Replace(decodedString, "\r\n", "\\n", -1)
	// Same with Unix newlines
	decodedString = strings.Replace(decodedString, "\n", "\\n", -1)

	if decodedString == "あああああああああああああああああああ" {
		decodedString = PlaceHolderString
	}

	txtXML := XMLTXT{
		Name:            name,
		UserData:        userData,
		Visible:         text.Flag & 0x1,
		Widescreen:      (text.Flag & 0x2) >> 1,
		Flag:            (text.Flag & 0x4) >> 2,
		Origin:          Coord2D{X: float32(text.Origin % 3), Y: float32(text.Origin / 3)},
		Alpha:           text.Alpha,
		Padding:         0,
		Translate:       Coord3D{X: text.XTranslation, Y: text.YTranslation, Z: text.ZTranslation},
		Rotate:          Coord3D{X: text.XRotate, Y: text.YRotate, Z: text.ZRotate},
		Scale:           Coord2D{X: text.XScale, Y: text.YScale},
		Width:           text.Width,
		Height:          text.Height,
		MaxStringLength: text.MaxStringLength,
		XSize:           text.FontSizeX,
		YSize:           text.FontSizeY,
		CharSize:        text.CharacterSize,
		LineSize:        text.LineSize,
		TopColor: Color8{
			R: text.Color1[0],
			G: text.Color1[1],
			B: text.Color1[2],
			A: text.Color1[3],
		},
		BottomColor: Color8{
			R: text.Color2[0],
			G: text.Color2[1],
			B: text.Color2[2],
			A: text.Color2[3],
		},
		Text: decodedString,
	}

	r.Panes = append(r.Panes, Children{TXT: &txtXML})
}
