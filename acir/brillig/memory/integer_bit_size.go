package memory

import (
	"encoding/binary"
	"fmt"
	"io"
)

type IntegerBitSize uint32

const (
	IntegerBitSizeU1 IntegerBitSize = iota
	IntegerBitSizeU8
	IntegerBitSizeU16
	IntegerBitSizeU32
	IntegerBitSizeU64
	IntegerBitSizeU128
)

func (b *IntegerBitSize) UnmarshalReader(r io.Reader) error {
	if err := binary.Read(r, binary.LittleEndian, b); err != nil {
		return err
	}

	if *b > IntegerBitSizeU128 {
		return fmt.Errorf("invalid IntegerBitSize: %d", *b)
	}

	return nil
}
