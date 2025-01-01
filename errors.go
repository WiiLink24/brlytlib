package brlyt

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidFileMagic         = errors.New("file is not a BRLYT")
	ErrFileSizeMismatch         = errors.New("file size is mismatched")
	ErrInvalidTXLHeader         = errors.New("txl1 header magic is invalid")
	ErrMisMatchedTXT1StringSize = func(stringSize int, correctSize uint16) error {
		return fmt.Errorf("string Size (%d) does not match the size found (%d)", stringSize, correctSize)
	}
)
