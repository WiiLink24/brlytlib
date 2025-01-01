package brlyt

import (
	"bytes"
	"encoding/binary"
	"encoding/xml"
	"io"
)

// Header represents the header of our BRLYT
type Header struct {
	Magic        [4]byte
	BOM          uint32
	FileSize     uint32
	HeaderLen    uint16
	SectionCount uint16
}

// SectionTypes are known parts of a BRLYT.
type SectionTypes [4]byte

var (
	headerMagic    SectionTypes = [4]byte{'R', 'L', 'Y', 'T'}
	SectionTypeLYT SectionTypes = [4]byte{'l', 'y', 't', '1'}
	SectionTypeTXL SectionTypes = [4]byte{'t', 'x', 'l', '1'}
	SectionTypeFNL SectionTypes = [4]byte{'f', 'n', 'l', '1'}
	SectionTypeMAT SectionTypes = [4]byte{'m', 'a', 't', '1'}
	SectionTypePAN SectionTypes = [4]byte{'p', 'a', 'n', '1'}
	SectionTypePAS SectionTypes = [4]byte{'p', 'a', 's', '1'}
	SectionTypePAE SectionTypes = [4]byte{'p', 'a', 'e', '1'}
	SectionTypeBND SectionTypes = [4]byte{'b', 'n', 'd', '1'}
	SectionTypePIC SectionTypes = [4]byte{'p', 'i', 'c', '1'}
	SectionTypeTXT SectionTypes = [4]byte{'t', 'x', 't', '1'}
	SectionTypeWND SectionTypes = [4]byte{'w', 'n', 'd', '1'}
	SectionTypeGRP SectionTypes = [4]byte{'g', 'r', 'p', '1'}
	SectionTypeGRS SectionTypes = [4]byte{'g', 'r', 's', '1'}
	SectionTypeGRE SectionTypes = [4]byte{'g', 'r', 'e', '1'}

	sectionCount = 0
)

type SectionHeader struct {
	Type SectionTypes
	Size uint32
}

type BRLYTWriter struct {
	*bytes.Buffer
}

func ParseBRLYT(contents []byte) (*Root, error) {
	readable := bytes.NewReader(contents)

	var header Header
	err := binary.Read(readable, binary.BigEndian, &header)
	if err != nil {
		return nil, err
	}

	if !bytes.Equal(headerMagic[:], header.Magic[:]) {
		return nil, ErrInvalidFileMagic
	}

	if readable.Size() != int64(header.FileSize) {
		return nil, ErrFileSizeMismatch
	}

	root := Root{
		XMLName: xml.Name{},
		LYT:     LYTNode{},
		FNL:     nil,
		TXL:     nil,
		reader:  readable,
	}

	for root.count = header.SectionCount; root.count != 0; root.count-- {
		var sectionHeader SectionHeader
		err = binary.Read(readable, binary.BigEndian, &sectionHeader)
		if err != nil {
			return nil, err
		}

		// Subtract the header size
		sectionSize := int(sectionHeader.Size) - 8
		temp := make([]byte, sectionSize)
		_, err = readable.Read(temp)
		if err != nil {
			return nil, err
		}

		switch sectionHeader.Type {
		case SectionTypeLYT:
			err = root.ParseLYT(temp)
			if err != nil {
				return nil, err
			}
		case SectionTypeTXL:
			err = root.ParseTXL(temp, sectionHeader.Size)
			if err != nil {
				return nil, err
			}
		case SectionTypeFNL:
			err = root.ParseFNL(temp, sectionHeader.Size)
			if err != nil {
				return nil, err
			}
		case SectionTypeMAT:
			err = root.ParseMAT(temp, sectionHeader.Size)
			if err != nil {
				return nil, err
			}
		case SectionTypePAN:
			// Root Pane is guaranteed to exist. We will sequentially read from it.
			_, err = root.ParsePAN(temp)
			if err != nil {
				return nil, err
			}
		case SectionTypeGRP:
			_, err = root.ParseGRP(temp)
			if err != nil {
				return nil, err
			}
		}
	}

	return &root, nil
}

func WriteBRLYT(data []byte) ([]byte, error) {
	var root Root
	err := xml.Unmarshal(data, &root)
	if err != nil {
		return nil, err
	}

	writer := BRLYTWriter{bytes.NewBuffer(nil)}

	// First write the header
	header := Header{
		Magic:        headerMagic,
		BOM:          0xFEFF000A,
		FileSize:     0,
		HeaderLen:    16,
		SectionCount: 0,
	}

	err = binary.Write(writer, binary.BigEndian, header)
	if err != nil {
		return nil, err
	}

	sectionCount = 0

	// Write the LYT1 section
	err = writer.WriteLYT(root)
	if err != nil {
		return nil, err
	}
	sectionCount++

	if root.TXL != nil {
		// Write TXL section
		err = writer.WriteTXL(root)
		if err != nil {
			return nil, err
		}

		sectionCount++
	}

	if root.FNL != nil {
		// Write FNL section
		err = writer.WriteFNL(root)
		if err != nil {
			return nil, err
		}

		sectionCount++
	}

	// Write MAT section
	err = writer.WriteMAT(root)
	if err != nil {
		return nil, err
	}

	sectionCount++

	// Write RootPane then children
	err = writer.WritePane(root.RootPane)
	if err != nil {
		return nil, err
	}

	err = writer.WriteChildren(root.RootPane.Children)
	if err != nil {
		return nil, err
	}

	sectionCount++

	// Same with RootGroup.
	err = writer.WriteGRP(root.RootGroup)
	if err != nil {
		return nil, err
	}

	err = writer.WriteGroupChildren(root.RootGroup.Children)
	if err != nil {
		return nil, err
	}

	sectionCount++

	binary.BigEndian.PutUint32(writer.Bytes()[8:12], uint32(writer.Len()))
	binary.BigEndian.PutUint16(writer.Bytes()[14:16], uint16(sectionCount))

	return writer.Bytes(), nil
}

func (b *BRLYTWriter) WriteGroupChildren(children []Children) error {
	if children == nil {
		return nil
	}

	err := b.WriteGRS()
	if err != nil {
		return err
	}

	sectionCount++

	for _, child := range children {
		if child.GRP != nil {
			err = b.WriteGRP(*child.GRP)
			if err != nil {
				return err
			}

			err = b.WriteChildren(child.GRP.Children)
			if err != nil {
				return err
			}
		}

		sectionCount++
	}

	err = b.WriteGRE()
	if err != nil {
		return err
	}

	sectionCount++
	return nil
}

func (b *BRLYTWriter) WriteChildren(children []Children) error {
	if children == nil {
		return nil
	}

	// Write pane start
	err := b.WritePAS()
	if err != nil {
		return err
	}

	sectionCount++

	for _, child := range children {
		if child.Pane != nil {
			err = b.WritePane(*child.Pane)
			if err != nil {
				return err
			}

			err = b.WriteChildren(child.Pane.Children)
			if err != nil {
				return err
			}
		}
		if child.BND != nil {
			err = b.WriteBND(*child.BND)
			if err != nil {
				return err
			}

			err = b.WriteChildren(child.BND.Children)
			if err != nil {
				return err
			}
		}
		if child.PIC != nil {
			err = b.WritePIC(*child.PIC)
			if err != nil {
				return err
			}

			err = b.WriteChildren(child.PIC.Children)
			if err != nil {
				return err
			}
		}
		if child.TXT != nil {
			err = b.WriteTXT(*child.TXT)
			if err != nil {
				return err
			}

			err = b.WriteChildren(child.TXT.Children)
			if err != nil {
				return err
			}
		}
		if child.WND != nil {
			err = b.WriteWND(*child.WND)
			if err != nil {
				return err
			}

			err = b.WriteChildren(child.WND.Children)
			if err != nil {
				return err
			}
		}

		sectionCount++
	}

	// End this pane
	err = b.WritePAE()
	if err != nil {
		return err
	}

	sectionCount++
	return nil
}

func write(writer io.Writer, data any) error {
	return binary.Write(writer, binary.BigEndian, data)
}
