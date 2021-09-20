package brlytlib

import (
	"bytes"
	"encoding/binary"
	"strings"
)

var (
	OriginX = []string{"Left", "Center", "Right"}
	OriginY = []string{"Top", "Center", "Bottom"}
)

// Pane represents the structure of a pane. All sections follow this structure
type Pane struct {
	Magic			[4]byte
	SectionLength	uint32
	Flag1			uint8
	Origin			uint8
	Alpha			uint8
	_				uint8
	PaneName		[0x10]byte
	UserData		[0x08]byte
	XTranslation	float32
	YTranslation	float32
	ZTranslation	float32
	XRotate			float32
	YRotate			float32
	ZRotate			float32
	XScale			float32
	YScale			float32
	Width			float32
	Height			float32
}

func ParsePAN(contents []byte) ([]PaneValues, error) {
	panOffsets := findAllOccurrences(contents, []string{"pan1"})
	var xmlNode	[]PaneValues

	for _, panOffset := range panOffsets {
		var pane Pane
		err := binary.Read(bytes.NewReader(contents[panOffset:]), binary.BigEndian, &pane)
		if err != nil {
			return nil, err
		}

		// Strip the null bytes from the strings
		name := strings.Replace(string(pane.PaneName[:]), "\x00", "", -1)
		userData := strings.Replace(string(pane.UserData[:]), "\x00", "", -1)

		panXML := PaneValues{
			Type:       "pan1",
			Name:       name,
			UserData:   userData,
			Visible:    pane.Flag1 & 0x1,
			Widescreen: (pane.Flag1 & 0x2) >> 1,
			Flag:       (pane.Flag1 & 0x4) >> 2,
			Origin: struct {
				X string `xml:"x"`
				Y string `xml:"y"`
			}{X: OriginX[pane.Origin%3], Y: OriginY[pane.Origin/3]},
			Alpha:   pane.Alpha,
			Padding: 0,
			Translate: struct {
				X float32 `xml:"x"`
				Y float32 `xml:"y"`
				Z float32 `xml:"z"`
			}{X: pane.XTranslation, Y: pane.YTranslation, Z: pane.ZTranslation},
			Rotate: struct {
				X float32 `xml:"x"`
				Y float32 `xml:"y"`
				Z float32 `xml:"z"`
			}{X: pane.XRotate, Y: pane.YRotate, Z: pane.ZRotate},
			Scale: struct {
				X float32 `xml:"x"`
				Y float32 `xml:"y"`
			}{X: pane.XScale, Y: pane.YScale},
			Width:  pane.Width,
			Height: pane.Height,
			TXT1Text: nil,
		}

		xmlNode = append(xmlNode, panXML)
	}

	return xmlNode, nil
}

// findAllOccurrences finds the offsets of the specified string.
func findAllOccurrences(data []byte, searches []string) []int {
	var results []int
	for _, search := range searches {
		searchData := data
		term := []byte(search)
		for x, d := bytes.Index(searchData, term), 0; x > -1; x, d = bytes.Index(searchData, term), d+x+1 {
			results = append(results, x+d)
			searchData = searchData[x+1:]
		}
	}
	return results
}