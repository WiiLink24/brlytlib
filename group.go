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
