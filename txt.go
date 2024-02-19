package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
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
	TopColor        [4]uint8
	BottomColor     [4]uint8
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
	fmt.Println(len(decodedString))

	if decodedString == "あああああああああああああああああああ" {
		decodedString = PlaceHolderString
	}

	txtXML := XMLTXT{
		Name:            name,
		UserData:        userData,
		Visible:         text.Flag & 0x1,
		Widescreen:      (text.Flag & 0x2) >> 1,
		Flag:            text.Flag,
		Origin:          Coord2D{X: float32(text.Origin % 3), Y: float32(text.Origin / 3)},
		Alpha:           text.Alpha,
		Padding:         0,
		Translate:       Coord3D{X: text.XTranslation, Y: text.YTranslation, Z: text.ZTranslation},
		Rotate:          Coord3D{X: text.XRotate, Y: text.YRotate, Z: text.ZRotate},
		Scale:           Coord2D{X: text.XScale, Y: text.YScale},
		Width:           text.Width,
		Height:          text.Height,
		MaxStringLength: text.MaxStringLength,
		MatIndex:        text.MatIndex,
		TextAlignment:   text.Alignment,
		XSize:           text.FontSizeX,
		YSize:           text.FontSizeY,
		CharSize:        text.CharacterSize,
		LineSize:        text.LineSize,
		TopColor: Color8{
			R: text.TopColor[0],
			G: text.TopColor[1],
			B: text.TopColor[2],
			A: text.TopColor[3],
		},
		BottomColor: Color8{
			R: text.BottomColor[0],
			G: text.BottomColor[1],
			B: text.BottomColor[2],
			A: text.BottomColor[3],
		},
		Text: decodedString,
	}

	r.Panes = append(r.Panes, Children{TXT: &txtXML})
}

func (b *BRLYTWriter) WriteTXT(txt XMLTXT) {
	temp := bytes.NewBuffer(nil)

	header := SectionHeader{
		Type: SectionTypeTXT,
		Size: 124,
	}

	var name [16]byte
	copy(name[:], txt.Name)

	var userData [8]byte
	copy(userData[:], txt.UserData)

	text := strings.Replace(txt.Text, "\\n", "\n", -1)
	encodedText := utf16.Encode([]rune(text))

	textLength := len(encodedText)*2 + 2

	pane := TXT{
		Flag:            txt.Flag,
		Origin:          uint8(txt.Origin.X + (txt.Origin.Y * 3)),
		Alpha:           txt.Alpha,
		PaneName:        name,
		UserData:        userData,
		XTranslation:    txt.Translate.X,
		YTranslation:    txt.Translate.Y,
		ZTranslation:    txt.Translate.Z,
		XRotate:         txt.Rotate.X,
		YRotate:         txt.Rotate.Y,
		ZRotate:         txt.Rotate.Z,
		XScale:          txt.Scale.X,
		YScale:          txt.Scale.Y,
		Width:           txt.Width,
		Height:          txt.Height,
		StringLength:    uint16(textLength),
		MaxStringLength: txt.MaxStringLength,
		MatIndex:        txt.MatIndex,
		FontIndex:       0,
		Alignment:       txt.TextAlignment,
		TextOffset:      116,
		TopColor:        [4]uint8{txt.TopColor.R, txt.TopColor.G, txt.TopColor.B, txt.TopColor.A},
		BottomColor:     [4]uint8{txt.BottomColor.R, txt.BottomColor.G, txt.BottomColor.B, txt.BottomColor.A},
		FontSizeX:       txt.XSize,
		FontSizeY:       txt.YSize,
		CharacterSize:   txt.CharSize,
		LineSize:        txt.LineSize,
	}

	write(temp, header)
	write(temp, pane)
	write(temp, encodedText)

	pos := 0
	for (b.Len()+temp.Len())%4 != 0 {
		temp.WriteByte(0)
		pos += 1
	}

	// If there is no modulo padding, pad with an u32
	if pos == 0 {
		write(temp, uint32(0))
	}

	binary.BigEndian.PutUint32(temp.Bytes()[4:8], uint32(temp.Len()))
	write(b, temp.Bytes())
}
