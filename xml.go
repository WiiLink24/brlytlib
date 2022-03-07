package main

import "encoding/xml"

type XMLPane struct {
	Name       string  `xml:"name,attr"`
	UserData   string  `xml:"user_data,attr"`
	Visible    uint8   `xml:"visible"`
	Widescreen uint8   `xml:"widescreen_affected"`
	Flag       uint8   `xml:"flag"`
	Origin     Coord2D `xml:"origin"`
	Alpha      uint8   `xml:"alpha"`
	Padding    uint8   `xml:"padding"`
	Translate  Coord3D `xml:"translate"`
	Rotate     Coord3D `xml:"rotate"`
	Scale      Coord2D `xml:"scale"`
	Width      float32 `xml:"width"`
	Height     float32 `xml:"height"`
}

type XMLBND struct {
	Name       string  `xml:"name,attr"`
	UserData   string  `xml:"user_data,attr"`
	Visible    uint8   `xml:"visible"`
	Widescreen uint8   `xml:"widescreen_affected"`
	Flag       uint8   `xml:"flag"`
	Origin     Coord2D `xml:"origin"`
	Alpha      uint8   `xml:"alpha"`
	Padding    uint8   `xml:"padding"`
	Translate  Coord3D `xml:"translate"`
	Rotate     Coord3D `xml:"rotate"`
	Scale      Coord2D `xml:"scale"`
	Width      float32 `xml:"width"`
	Height     float32 `xml:"height"`
}

type XMLPIC struct {
	Name             string     `xml:"name,attr"`
	UserData         string     `xml:"user_data,attr"`
	Visible          uint8      `xml:"visible"`
	Widescreen       uint8      `xml:"widescreen_affected"`
	Flag             uint8      `xml:"flag"`
	Origin           Coord2D    `xml:"origin"`
	Alpha            uint8      `xml:"alpha"`
	Padding          uint8      `xml:"padding"`
	Translate        Coord3D    `xml:"translate"`
	Rotate           Coord3D    `xml:"rotate"`
	Scale            Coord2D    `xml:"scale"`
	Width            float32    `xml:"width"`
	Height           float32    `xml:"height"`
	TopLeftColor     Color8     `xml:"topLeftColor"`
	TopRightColor    Color8     `xml:"topRightColor"`
	BottomLeftColor  Color8     `xml:"bottomLeftColor"`
	BottomRightColor Color8     `xml:"bottomRightColor"`
	MatIndex         uint16     `xml:"matIndex"`
	UVSets           *XMLUVSets `xml:"uv_sets"`
}

type XMLWindow struct {
	Name             string         `xml:"name,attr"`
	UserData         string         `xml:"user_data,attr"`
	Visible          uint8          `xml:"visible"`
	Widescreen       uint8          `xml:"widescreen_affected"`
	Flag             uint8          `xml:"flag"`
	Origin           Coord2D        `xml:"origin"`
	Alpha            uint8          `xml:"alpha"`
	Padding          uint8          `xml:"padding"`
	Translate        Coord3D        `xml:"translate"`
	Rotate           Coord3D        `xml:"rotate"`
	Scale            Coord2D        `xml:"scale"`
	Width            float32        `xml:"width"`
	Height           float32        `xml:"height"`
	Coordinate1      float32        `xml:"coordinate_1"`
	Coordinate2      float32        `xml:"coordinate_2"`
	Coordinate3      float32        `xml:"coordinate_3"`
	Coordinate4      float32        `xml:"coordinate_4"`
	TopLeftColor     Color8         `xml:"topLeftColor"`
	TopRightColor    Color8         `xml:"topRightColor"`
	BottomLeftColor  Color8         `xml:"bottomLeftColor"`
	BottomRightColor Color8         `xml:"bottomRightColor"`
	MatIndex         uint16         `xml:"matIndex"`
	UVSets           *XMLUVSets     `xml:"uv_sets"`
	Materials        *XMLWindowMats `xml:"materials"`
}

type XMLGRP struct {
	Name    string   `xml:"name,attr"`
	Entries []string `xml:"entries"`
}

type XMLWindowMat struct {
	MatIndex uint16 `xml:"matIndex"`
	Index    uint8  `xml:"index"`
}

type XMLWindowMats struct {
	Mats []XMLWindowMat `xml:"mats"`
}

type XMLUVSets struct {
	Set []XMLUVSet `xml:"set"`
}

type XMLUVSet struct {
	CoordTL STCoordinates `xml:"coordTL"`
	CoordTR STCoordinates `xml:"coordTR"`
	CoordBL STCoordinates `xml:"coordBL"`
	CoordBR STCoordinates `xml:"coordBR"`
}

type STCoordinates struct {
	S float32 `xml:"s"`
	T float32 `xml:"t"`
}

type XMLTXT struct {
	Name            string  `xml:"name,attr"`
	UserData        string  `xml:"user_data,attr"`
	Visible         uint8   `xml:"visible"`
	Widescreen      uint8   `xml:"widescreen_affected"`
	Flag            uint8   `xml:"flag"`
	Origin          Coord2D `xml:"origin"`
	Alpha           uint8   `xml:"alpha"`
	Padding         uint8   `xml:"padding"`
	Translate       Coord3D `xml:"translate"`
	Rotate          Coord3D `xml:"rotate"`
	Scale           Coord2D `xml:"scale"`
	Width           float32 `xml:"width"`
	Height          float32 `xml:"height"`
	MaxStringLength uint16  `xml:"max_string_length"`
	XSize           float32 `xml:"x_size"`
	YSize           float32 `xml:"y_size"`
	CharSize        float32 `xml:"charsize"`
	LineSize        float32 `xml:"linesize"`
	TopColor        Color8  `xml:"top_color"`
	BottomColor     Color8  `xml:"bottom_color"`
	Text            string  `xml:"text"`
}

type XMLPAS struct{}
type XMLPAE struct{}

type XMLGRS struct{}
type XMLGRE struct{}

// Children contains all the possible children a brlyt can contain.
// This is needed for unmarshalling when we put together a new brlyt.
type Children struct {
	Pane *XMLPane   `xml:"pan1"`
	PAS  *XMLPAS    `xml:"pas1"`
	PAE  *XMLPAE    `xml:"pae1"`
	BND  *XMLBND    `xml:"bnd1"`
	PIC  *XMLPIC    `xml:"pic1"`
	TXT  *XMLTXT    `xml:"txt1"`
	WND  *XMLWindow `xml:"wnd1"`
	GRP  *XMLGRP    `xml:"grp1"`
	GRS  *XMLGRS    `xml:"grs1"`
	GRE  *XMLGRE    `xml:"gre1"`
}

// TPLNames represents the structure of the txl1 section.
type TPLNames struct {
	TPLName []string `xml:"tpl_name"`
}

type FNLNames struct {
	FNLName []string `xml:"font_name"`
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
	XMLName xml.Name   `xml:"root"`
	LYT     LYTNode    `xml:"lyt1"`
	TXL     *TPLNames  `xml:"txl1"`
	FNL     *FNLNames  `xml:"fnt1"`
	MAT     MATNode    `xml:"mat1"`
	Panes   []Children `xml:"children"`
}

type MATNode struct {
	Entries []MATEntries `xml:"entries"`
}

type MATEntries struct {
	Name      string        `xml:"name,attr"`
	ForeColor Color16       `xml:"foreColor"`
	BackColor Color16       `xml:"backColor"`
	ColorReg3 Color16       `xml:"colorReg3"`
	TevColor1 Color8        `xml:"tevColor1"`
	TevColor2 Color8        `xml:"tevColor2"`
	TevColor3 Color8        `xml:"tevColor3"`
	TevColor4 Color8        `xml:"tevColor4"`
	Textures  []MATTexture  `xml:"texture"`
	SRT       []MATSRT      `xml:"textureSRT"`
	CoordGen  []MATCoordGen `xml:"coordGen"`
}

type MATTexture struct {
	Name  string `xml:"name,attr"`
	SWrap uint8
	TWrap uint8
}

type MATSRT struct {
	XTrans   float32 `xml:"XTrans"`
	YTrans   float32 `xml:"YTrans"`
	Rotation float32 `xml:"Rotation"`
	XScale   float32 `xml:"XScale"`
	YScale   float32 `xml:"YScale"`
}

type MATCoordGen struct {
	Type         uint8 `xml:"type"`
	Source       uint8 `xml:"source"`
	MatrixSource uint8 `xml:"matrixSource"`
}

type Color8 struct {
	R uint8
	G uint8
	B uint8
	A uint8
}

type Color16 struct {
	R int16
	G int16
	B int16
	A int16
}

type Coord3D struct {
	X float32 `xml:"x"`
	Y float32 `xml:"y"`
	Z float32 `xml:"z"`
}

type Coord2D struct {
	X float32 `xml:"x"`
	Y float32 `xml:"y"`
}
