package main

import (
	"bytes"
	"encoding/binary"
	"strings"
)

type GRP struct {
	Name         [16]byte
	NumOfEntries uint16
	_            uint16
}

func (r *Root) ParseGRP(data []byte) {
	var grp GRP

	err := binary.Read(bytes.NewReader(data), binary.BigEndian, &grp)
	if err != nil {
		panic(err)
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

	r.Panes = append(r.Panes, Children{GRP: &xmlData})
}

func (r *Root) ParseGRS() {
	r.Panes = append(r.Panes, Children{
		GRS: &XMLGRS{},
	})
}

func (r *Root) ParseGRE() {
	r.Panes = append(r.Panes, Children{
		GRE: &XMLGRE{},
	})
}

func (b *BRLYTWriter) WriteGRP(data XMLGRP) {
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

	write(b, header)
	write(b, grp)

	for _, str := range data.Entries {
		var entry [16]byte
		copy(entry[:], str)

		write(b, entry)
	}
}

func (b *BRLYTWriter) WriteGRS() {
	header := SectionHeader{
		Type: SectionTypeGRS,
		Size: 8,
	}

	write(b, header)
}

func (b *BRLYTWriter) WriteGRE() {
	header := SectionHeader{
		Type: SectionTypeGRE,
		Size: 8,
	}

	write(b, header)
}
