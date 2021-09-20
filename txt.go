package brlytlib

import (
	"bytes"
	"encoding/binary"
	"strings"
	"unicode/utf16"
)

// TextChunk represents the text data of the txt1 section
type TextChunk struct {
	StringLength    uint16
	MaxStringLength uint16
	MatIndex        uint16
	FontIndex       uint16
	Alignment       uint8
	_               uint8
	_               uint16
	_               uint32
	Color1          uint32
	Color2          uint32
	FontSizeX       float32
	FontSizeY       float32
	CharacterSize   float32
	LineSize        float32
}

// TXT represents the structure of the txt1 section
type TXT struct {
	Pane      Pane
	TextChunk TextChunk
}

func ParseTXT(contents []byte) ([]PaneValues, error) {
	txtOffsets := findAllOccurrences(contents, []string{"txt1"})
	var xmlNode []PaneValues

	for _, txtOffset := range txtOffsets {
		var text TXT
		err := binary.Read(bytes.NewReader(contents[txtOffset:]), binary.BigEndian, &text)
		if err != nil {
			return nil, err
		}

		// Strip the null bytes from the strings
		name := strings.Replace(string(text.Pane.PaneName[:]), "\x00", "", -1)
		userData := strings.Replace(string(text.Pane.UserData[:]), "\x00", "", -1)

		// The text is always at offset 116 relative to the start of the section.
		utf16String := contents[txtOffset+116 : txtOffset+int(text.TextChunk.StringLength)+116]
		if len(utf16String) != int(text.TextChunk.StringLength) {
			// The string size found in the text chunk does not match the size of the string we pulled.
			return nil, ErrMisMatchedTXT1StringSize(len(utf16String), text.TextChunk.StringLength)
		}

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

		txtXML := PaneValues{
			Type:       "txt1",
			Name:       name,
			UserData:   userData,
			Visible:    text.Pane.Flag1 & 0x1,
			Widescreen: (text.Pane.Flag1 & 0x2) >> 1,
			Flag:       (text.Pane.Flag1 & 0x4) >> 2,
			Origin: struct {
				X string `xml:"x"`
				Y string `xml:"y"`
			}{X: OriginX[text.Pane.Origin%3], Y: OriginY[text.Pane.Origin/3]},
			Alpha:   text.Pane.Alpha,
			Padding: 0,
			Translate: struct {
				X float32 `xml:"x"`
				Y float32 `xml:"y"`
				Z float32 `xml:"z"`
			}{X: text.Pane.XTranslation, Y: text.Pane.YTranslation, Z: text.Pane.ZTranslation},
			Rotate: struct {
				X float32 `xml:"x"`
				Y float32 `xml:"y"`
				Z float32 `xml:"z"`
			}{X: text.Pane.XRotate, Y: text.Pane.YRotate, Z: text.Pane.ZRotate},
			Scale: struct {
				X float32 `xml:"x"`
				Y float32 `xml:"y"`
			}{X: text.Pane.XScale, Y: text.Pane.YScale},
			Width:  text.Pane.Width,
			Height: text.Pane.Height,
			TXT1Text: []TXTStrings{{
				StringLength:    text.TextChunk.StringLength,
				MaxStringLength: text.TextChunk.MaxStringLength,
				XSize:           text.TextChunk.FontSizeX,
				YSize:           text.TextChunk.FontSizeY,
				CharSize:        text.TextChunk.CharacterSize,
				LineSize:        text.TextChunk.LineSize,
				TopColor: struct {
					R uint32 `xml:"r"`
					G uint32 `xml:"g"`
					B uint32 `xml:"b"`
					A uint32 `xml:"a"`
				}{R: (text.TextChunk.Color1 >> 24) & 0xff,
					G: (text.TextChunk.Color1 >> 16) & 0xff,
					B: (text.TextChunk.Color1 >> 8) & 0xff,
					A: (text.TextChunk.Color1 >> 0) & 0xff},
				BottomColor: struct {
					R uint32 `xml:"r"`
					G uint32 `xml:"g"`
					B uint32 `xml:"b"`
					A uint32 `xml:"a"`
				}{R: (text.TextChunk.Color2 >> 24) & 0xff,
					G: (text.TextChunk.Color2 >> 16) & 0xff,
					B: (text.TextChunk.Color2 >> 8) & 0xff,
					A: (text.TextChunk.Color2 >> 0) & 0xff,
				},
				Alignment: struct {
					X string `xml:"x"`
					Y string `xml:"y"`
				}{X: OriginX[text.Pane.Origin%3], Y: OriginY[text.Pane.Origin/3]},
				Text: decodedString,
			}},
		}

		xmlNode = append(xmlNode, txtXML)
	}

	return xmlNode, nil
}
