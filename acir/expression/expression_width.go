package expression

import (
	"encoding/binary"
	"fmt"
	"io"
)

type ExpressionWidth struct {
	Kind  ExpressionWidthKind
	Width *uint64
}

type ExpressionWidthKind uint32

const (
	ACIRExpressionWidthUnbounded ExpressionWidthKind = iota
	ACIRExpressionWidthBounded
)

func (e *ExpressionWidth) UnmarshalReader(r io.Reader) error {
	if err := binary.Read(r, binary.LittleEndian, &e.Kind); err != nil {
		return err
	}

	switch e.Kind {
	case ACIRExpressionWidthUnbounded:
		e.Width = nil
	case ACIRExpressionWidthBounded:
		e.Width = new(uint64)
		if err := binary.Read(r, binary.LittleEndian, e.Width); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown ExpressionWidth Kind: %d", e.Kind)
	}

	return nil
}
