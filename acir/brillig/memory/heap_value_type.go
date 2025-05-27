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
		var arraySize uint64
		if err := binary.Read(r, binary.LittleEndian, &arraySize); err != nil {
			return err
		}
		fmt.Printf("Array size: %d\n", arraySize)

		h.ArrayValueTypes = &[]HeapValueType{}
		for i := uint64(0); i < arraySize; i++ {
			fmt.Printf("Reading value type %d\n", i)
			var valueType HeapValueType
			if err := valueType.UnmarshalReader(r); err != nil {
				return err
			}
			*h.ArrayValueTypes = append(*h.ArrayValueTypes, valueType)
		}

		var arraySize64 uint64
		if err := binary.Read(r, binary.LittleEndian, &arraySize64); err != nil {
			return err
		}
		fmt.Printf("Array size (64-bit): %d\n", arraySize64)

		h.ArraySize = &arraySize64
	case ACIRBrilligHeapValueTypeKindVector:
		var arraySize uint64
		if err := binary.Read(r, binary.LittleEndian, &arraySize); err != nil {
			return err
		}

		h.VectorValueTypes = &[]HeapValueType{}
		for i := uint64(0); i < arraySize; i++ {
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

func (h HeapValueType) Equals(other HeapValueType) bool {
	if h.Kind != other.Kind {
		fmt.Printf("HeapValueTypeKind does not match: %d != %d\n", h.Kind, other.Kind)
		return false
	}

	switch h.Kind {
	case ACIRBrilligHeapValueTypeKindSimple:
		if h.Simple == nil || other.Simple == nil {
			return h.Simple == other.Simple
		}
		return h.Simple.Equals(*other.Simple)
	case ACIRBrilligHeapValueTypeKindArray:
		if h.ArraySize == nil || other.ArraySize == nil || *h.ArraySize != *other.ArraySize {
			fmt.Println("Array sizes do not match: ", *h.ArraySize, *other.ArraySize)
			return false
		}
		if len(*h.ArrayValueTypes) != len(*other.ArrayValueTypes) {
			fmt.Println("Array value types lengths do not match")
			return false
		}
		for i := range *h.ArrayValueTypes {
			if !(*h.ArrayValueTypes)[i].Equals((*other.ArrayValueTypes)[i]) {
				fmt.Printf("Array value types at index %d do not match\n", i)
				return false
			}
		}
	case ACIRBrilligHeapValueTypeKindVector:
		if len(*h.VectorValueTypes) != len(*other.VectorValueTypes) {
			return false
		}
		for i := range *h.VectorValueTypes {
			if !(*h.VectorValueTypes)[i].Equals((*other.VectorValueTypes)[i]) {
				return false
			}
		}
	default:
		return false
	}

	return true
}
