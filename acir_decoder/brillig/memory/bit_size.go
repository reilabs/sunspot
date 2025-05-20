package memory

import (
	"encoding/binary"
	"fmt"
	"io"
)

type BitSize struct {
	Kind           BitSizeKind
	IntegerBitSize *IntegerBitSize
}

type BitSizeKind uint32

const (
	BitSizeKindField BitSizeKind = iota
	BitSizeKindInteger
)

func (b *BitSize) UnmarshalReader(r io.Reader) error {
	if err := binary.Read(r, binary.LittleEndian, &b.Kind); err != nil {
		return err
	}

	if b.Kind > BitSizeKindInteger {
		return fmt.Errorf("invalid BitSizeKind: %d", b.Kind)
	}

	if b.Kind == BitSizeKindInteger {
		var integerBitSize IntegerBitSize
		b.IntegerBitSize = &integerBitSize
		if err := b.IntegerBitSize.UnmarshalReader(r); err != nil {
			return err
		}
	}

	return nil
}
