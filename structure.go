package brlyt

// Pane represents the structure of a pan1 section.
type Pane struct {
	Flag         uint8
	Origin       uint8
	Alpha        uint8
	_            uint8
	PaneName     [16]byte
	UserData     [8]byte
	XTranslation float32
	YTranslation float32
	ZTranslation float32
	XRotate      float32
	YRotate      float32
	ZRotate      float32
	XScale       float32
	YScale       float32
	Width        float32
	Height       float32
}

// PIC defines the image pane in a brlyt
type PIC struct {
	Flag             uint8
	Origin           uint8
	Alpha            uint8
	_                uint8
	PaneName         [16]byte
	UserData         [8]byte
	XTranslation     float32
	YTranslation     float32
	ZTranslation     float32
	XRotate          float32
	YRotate          float32
	ZRotate          float32
	XScale           float32
	YScale           float32
	Width            float32
	Height           float32
	TopLeftColor     [4]uint8
	TopRightColor    [4]uint8
	BottomLeftColor  [4]uint8
	BottomRightColor [4]uint8
	MatIndex         uint16
	NumOfUVSets      uint8
	_                uint8
}

type UVSet struct {
	TopLeftS     float32
	TopLeftT     float32
	TopRightS    float32
	TopRightT    float32
	BottomLeftS  float32
	BottomLeftT  float32
	BottomRightS float32
	BottomRightT float32
}

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
	StringOrigin    uint8
	LineAlignment   uint8
	_               uint16
	TextOffset      uint32
	TopColor        [4]uint8
	BottomColor     [4]uint8
	FontSizeX       float32
	FontSizeY       float32
	CharacterSize   float32
	LineSize        float32
}

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

type GRP struct {
	Name         [16]byte
	NumOfEntries uint16
	_            uint16
}
