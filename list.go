package brlyt

import (
	"encoding/binary"
	"errors"
	"io"
)

func (r *Root) ParseChildren() ([]Children, error) {
	var children []Children

	isOver := false
	for {
		var sectionHeader SectionHeader
		err := binary.Read(r.reader, binary.BigEndian, &sectionHeader)
		if err != nil {
			return nil, err
		}

		// Subtract the header size
		sectionSize := int(sectionHeader.Size) - 8
		// We could end off at a pane end, meaning it would try to read EOF
		if sectionSize == 0 {
			r.count--
			break
		}

		temp := make([]byte, sectionSize)
		_, err = r.reader.Read(temp)
		if err != nil {
			return nil, err
		}

		switch sectionHeader.Type {
		case SectionTypeBND:
			bnd, err := r.ParseBND(temp)
			if err != nil {
				return nil, err
			}

			children = append(children, Children{BND: bnd})
		case SectionTypePIC:
			pic, err := r.ParsePIC(temp)
			if err != nil {
				return nil, err
			}

			children = append(children, Children{PIC: pic})
		case SectionTypeTXT:
			txt, err := r.ParseTXT(temp, sectionHeader.Size)
			if err != nil {
				return nil, err
			}

			children = append(children, Children{TXT: txt})
		case SectionTypeWND:
			wnd, err := r.ParseWND(temp)
			if err != nil {
				return nil, err
			}

			children = append(children, Children{WND: wnd})
		case SectionTypePAN:
			pan, err := r.ParsePAN(temp)
			if err != nil {
				return nil, err
			}

			children = append(children, Children{Pane: pan})
		case SectionTypeGRP:
			grp, err := r.ParseGRP(temp)
			if err != nil {
				return nil, err
			}

			children = append(children, Children{GRP: grp})
		case SectionTypePAE:
			isOver = true
		case SectionTypeGRE:
			isOver = true
		}

		// Deincrement the amount of sections left to read.
		r.count--
		if isOver {
			break
		}
	}

	return children, nil
}

func (r *Root) HasChildren() bool {
	var sectionHeader SectionHeader
	err := binary.Read(r.reader, binary.BigEndian, &sectionHeader)
	if err != nil {
		if errors.Is(err, io.EOF) {
			return false
		}

		panic(err)
	}

	if sectionHeader.Type == SectionTypePAS || sectionHeader.Type == SectionTypeGRS {
		// Read the pane start
		r.count--
		return true
	}

	_, err = r.reader.Seek(-8, io.SeekCurrent)
	if err != nil {
		panic(err)
	}
	return false
}
