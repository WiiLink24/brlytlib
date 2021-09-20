package brlytlib

import "encoding/xml"

// PaneValues is represents all the keys that a Pane can have.
type PaneValues struct {
	Type       string `xml:"type,attr"`
	Name       string `xml:"name,attr"`
	UserData   string `xml:"user_data,attr"`
	Visible    uint8  `xml:"visible"`
	Widescreen uint8  `xml:"widescreen_affected"`
	Flag       uint8  `xml:"flag"`
	Origin     struct {
		X string `xml:"x"`
		Y string `xml:"y"`
	} `xml:"origin"`
	Alpha     uint8 `xml:"alpha"`
	Padding   uint8 `xml:"padding"`
	Translate struct {
		X float32 `xml:"x"`
		Y float32 `xml:"y"`
		Z float32 `xml:"z"`
	} `xml:"translate"`
	Rotate struct {
		X float32 `xml:"x"`
		Y float32 `xml:"y"`
		Z float32 `xml:"z"`
	} `xml:"rotate"`
	Scale struct {
		X float32 `xml:"x"`
		Y float32 `xml:"y"`
	} `xml:"scale"`
	Width    float32      `xml:"width"`
	Height   float32      `xml:"height"`
	TXT1Text []TXTStrings `xml:"text"`
}

// TPLNames represents the structure of the txl1 section.
type TPLNames struct {
	TPLName []TPLNamesFormat `xml:"tpl_name"`
}

// TPLNamesFormat specifies the index and string of TPL names in the txl1 section.
type TPLNamesFormat struct {
	Index  int    `xml:"key,attr"`
	String string `xml:",innerxml"`
}

type FNLNames struct {
	FNLName []FNLNamesFormat `xml:"font_name"`
}

type FNLNamesFormat struct {
	Index  int    `xml:"key,attr"`
	String string `xml:",innerxml"`
}

// TXTStrings specifies the values that TextChunk contains
type TXTStrings struct {
	StringLength    uint16  `xml:"string_length"`
	MaxStringLength uint16  `xml:"max_string_length"`
	XSize           float32 `xml:"x_size"`
	YSize           float32 `xml:"y_size"`
	CharSize        float32 `xml:"charsize"`
	LineSize        float32 `xml:"linesize"`
	Alignment       struct {
		X string `xml:"x"`
		Y string `xml:"y"`
	} `xml:"alignment"`
	TopColor struct {
		R uint32 `xml:"r"`
		G uint32 `xml:"g"`
		B uint32 `xml:"b"`
		A uint32 `xml:"a"`
	} `xml:"top_color"`
	BottomColor struct {
		R uint32 `xml:"r"`
		G uint32 `xml:"g"`
		B uint32 `xml:"b"`
		A uint32 `xml:"a"`
	} `xml:"bottom_color"`
	Text string `xml:"text"`
}

// LYTNode specifies the values that LYT contains
type LYTNode struct {
	XMLName  xml.Name `xml:"lyt1"`
	Centered uint16   `xml:"is_centered"`
	Width    float32  `xml:"width"`
	Height   float32  `xml:"height"`
}

// Root is the main structure of our XML
type Root struct {
	XMLName xml.Name     `xml:"root"`
	LYT     LYTNode      `xml:"lyt1"`
	FNL     []FNLNames   `xml:"fnt1"`
	TPLName []TPLNames   `xml:"txl1"`
	Panes   []PaneValues `xml:"tag"`
}
