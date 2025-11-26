package expression

import (
	"encoding/binary"
	"encoding/json"
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

func (e *ExpressionWidth) Equals(other *ExpressionWidth) bool {
	if e.Kind != other.Kind {
		return false
	}

	if e.Width == nil && other.Width == nil {
		return true
	}

	if e.Width == nil || other.Width == nil {
		return false
	}

	return *e.Width == *other.Width
}

func (e *ExpressionWidth) MarshalJSON() ([]byte, error) {
	fieldsMap := make(map[string]interface{})
	switch e.Kind {
	case ACIRExpressionWidthUnbounded:
		fieldsMap["unbounded"] = nil
	case ACIRExpressionWidthBounded:
		fieldsMap["bounded"] = *e.Width
	}
	return json.Marshal(fieldsMap)
}
