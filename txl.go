package brlyt

import (
	"bytes"
	"encoding/binary"
	"strings"
)

type TXL struct {
	NumOfTPL uint16
	Unknown  uint16
}

type TPLOffSet struct {
	// Offset is relative to the beginning of the txl1 section
	Offset  uint32
	Padding uint32
}

func (r *Root) ParseTXL(data []byte, sectionSize uint32) error {
	var tplOffsets []uint32
	var tplNames []string

	var txl TXL
	err := binary.Read(bytes.NewReader(data), binary.BigEndian, &txl)
	if err != nil {
		return err
	}

	for i := 0; i < int(txl.NumOfTPL); i++ {
		// By now we have only read the header.
		// We will read the TPLOffset table in order to get our names.
		var tplTable TPLOffSet
		offset := 4 + (i * 8)

		err = binary.Read(bytes.NewReader(data[offset:]), binary.BigEndian, &tplTable)
		if err != nil {
			return err
		}

		tplOffsets = append(tplOffsets, tplTable.Offset+4)

		// If we have reached the last index, append the section size to the slice.
		if i == int(txl.NumOfTPL)-1 {
			tplOffsets = append(tplOffsets, sectionSize-8)
		}
	}

	// Now that we have the offsets, retrieve the TPL names.
	for i := 0; i < int(txl.NumOfTPL); i++ {
		tplName := string(data[tplOffsets[i]:tplOffsets[i+1]])

		// Strip the null terminator
		tplName = strings.Replace(tplName, "\x00", "", -1)

		tplNames = append(tplNames, tplName)
	}

	r.TXL = &TPLNames{TPLName: tplNames}
	return nil
}

func (b *BRLYTWriter) WriteTXL(data Root) error {
	sectionWriter := bytes.NewBuffer(nil)

	header := SectionHeader{
		Type: SectionTypeTXL,
		Size: 0,
	}

	txl := TXL{
		NumOfTPL: uint16(len(data.TXL.TPLName)),
		Unknown:  0,
	}

	offset := len(data.TXL.TPLName) * 8
	offsets := make([]TPLOffSet, len(data.TXL.TPLName))
	for i, _ := range data.TXL.TPLName {
		if i != 0 {
			offset += len(data.TXL.TPLName[i-1]) + 1
		}

		tplOffset := TPLOffSet{
			Offset:  uint32(offset),
			Padding: 0,
		}

		offsets[i] = tplOffset
	}

	for _, s := range data.TXL.TPLName {
		_, err := sectionWriter.WriteString(s)
		if err != nil {
			return err
		}

		// Write null terminator
		_, _ = sectionWriter.Write([]byte{0})
	}

	for (b.Len()+sectionWriter.Len())%4 != 0 {
		_, _ = sectionWriter.Write([]byte{0})
	}

	header.Size = uint32(12 + (8 * len(data.TXL.TPLName)) + sectionWriter.Len())

	err := write(b, header)
	if err != nil {
		return err
	}

	err = write(b, txl)
	if err != nil {
		return err
	}

	err = write(b, offsets)
	if err != nil {
		return err
	}

	return write(b, sectionWriter.Bytes())
}
