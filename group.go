package brlyt

import (
	"bytes"
	"encoding/binary"
	"strings"
)

func (r *Root) ParseGRP(data []byte) (*XMLGRP, error) {
	var grp GRP

	err := binary.Read(bytes.NewReader(data), binary.BigEndian, &grp)
	if err != nil {
		return nil, err
	}

	// Strip the null bytes from the strings
	name := strings.Replace(string(grp.Name[:]), "\x00", "", -1)

	entries := make([]string, grp.NumOfEntries)

	for i := 0; i < int(grp.NumOfEntries); i++ {
		offset := 20 + (i * 16)

		object := strings.Replace(string(data[offset:offset+16]), "\x00", "", -1)
		entries[i] = object
	}

	xmlData := XMLGRP{
		Name:    name,
		Entries: entries,
	}

	if name == "RootGroup" {
		r.RootGroup = xmlData
		if r.HasChildren() {
			r.RootGroup.Children, err = r.ParseChildren()
			if err != nil {
				return nil, err
			}
		}
	} else {
		if r.HasChildren() {
			xmlData.Children, err = r.ParseChildren()
			if err != nil {
				return nil, err
			}
		}
	}

	return &xmlData, nil
}

func (b *BRLYTWriter) WriteGRP(data XMLGRP) error {
	header := SectionHeader{
		Type: SectionTypeGRP,
		Size: uint32(28 + (16 * len(data.Entries))),
	}

	var name [16]byte
	copy(name[:], data.Name)

	grp := GRP{
		Name:         name,
		NumOfEntries: uint16(len(data.Entries)),
	}

	err := write(b, header)
	if err != nil {
		return err
	}

	err = write(b, grp)
	if err != nil {
		return err
	}

	for _, str := range data.Entries {
		var entry [16]byte
		copy(entry[:], str)

		err = write(b, entry)
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *BRLYTWriter) WriteGRS() error {
	header := SectionHeader{
		Type: SectionTypeGRS,
		Size: 8,
	}

	return write(b, header)
}

func (b *BRLYTWriter) WriteGRE() error {
	header := SectionHeader{
		Type: SectionTypeGRE,
		Size: 8,
	}

	return write(b, header)
}
