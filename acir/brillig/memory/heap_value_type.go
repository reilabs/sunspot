package memory

import (
	"encoding/binary"
	"fmt"
	"io"
)

type HeapValueType struct {
	Kind             HeapValueTypeKind
	Simple           *BitSize
	ArrayValueTypes  *[]HeapValueType
	ArraySize        *uint64
	VectorValueTypes *[]HeapValueType
}

type HeapValueTypeKind uint32

const (
	ACIRBrilligHeapValueTypeKindSimple HeapValueTypeKind = iota
	ACIRBrilligHeapValueTypeKindArray
	ACIRBrilligHeapValueTypeKindVector
)

func (h *HeapValueTypeKind) UnmarshalReader(r io.Reader) error {
	if err := binary.Read(r, binary.LittleEndian, h); err != nil {
		return err
	}

	if *h > ACIRBrilligHeapValueTypeKindVector {
		return fmt.Errorf("invalid HeapValueTypeKind: %d", *h)
	}

	return nil
}

func (h *HeapValueType) UnmarshalReader(r io.Reader) error {
	if err := h.Kind.UnmarshalReader(r); err != nil {
		return err
	}

	switch h.Kind {
	case ACIRBrilligHeapValueTypeKindSimple:
		h.Simple = &BitSize{}
		return h.Simple.UnmarshalReader(r)
	case ACIRBrilligHeapValueTypeKindArray:
		var arraySize uint32
		if err := binary.Read(r, binary.LittleEndian, &arraySize); err != nil {
			return err
		}

		h.ArrayValueTypes = &[]HeapValueType{}
		for i := uint32(0); i < arraySize; i++ {
			var valueType HeapValueType
			if err := valueType.UnmarshalReader(r); err != nil {
				return err
			}
			*h.ArrayValueTypes = append(*h.ArrayValueTypes, valueType)
		}

		if err := binary.Read(r, binary.LittleEndian, &h.ArraySize); err != nil {
			return err
		}
	case ACIRBrilligHeapValueTypeKindVector:
		var arraySize uint32
		if err := binary.Read(r, binary.LittleEndian, &arraySize); err != nil {
			return err
		}

		h.VectorValueTypes = &[]HeapValueType{}
		for i := uint32(0); i < arraySize; i++ {
			var valueType HeapValueType
			if err := valueType.UnmarshalReader(r); err != nil {
				return err
			}
			*h.VectorValueTypes = append(*h.VectorValueTypes, valueType)
		}
	default:
		return fmt.Errorf("invalid HeapValueTypeKind: %d", h.Kind)
	}

	return nil
}
