package memory

import (
	"encoding/binary"
	"fmt"
	"io"
)

type ValueOrArray struct {
	Kind          ValueOrArrayKind
	MemoryAddress *MemoryAddress
	HeapArray     *HeapArray
	HeapVector    *HeapVector
}

type ValueOrArrayKind uint32

const (
	ACIRBrilligValueOrArrayKindMemoryAddress ValueOrArrayKind = iota
	ACIRBrilligValueOrArrayKindHeapArray
	ACIRBrilligValueOrArrayKindHeapVector
)

func (v *ValueOrArrayKind) UnmarshalReader(r io.Reader) error {
	if err := binary.Read(r, binary.LittleEndian, v); err != nil {
		return err
	}

	if *v > ACIRBrilligValueOrArrayKindHeapVector {
		return fmt.Errorf("invalid ValueOrArrayKind: %d", *v)
	}

	return nil
}

func (v *ValueOrArray) UnmarshalReader(r io.Reader) error {
	if err := v.Kind.UnmarshalReader(r); err != nil {
		return err
	}

	switch v.Kind {
	case ACIRBrilligValueOrArrayKindMemoryAddress:
		v.MemoryAddress = &MemoryAddress{}
		return v.MemoryAddress.UnmarshalReader(r)
	case ACIRBrilligValueOrArrayKindHeapArray:
		v.HeapArray = &HeapArray{}
		return v.HeapArray.UnmarshalReader(r)
	case ACIRBrilligValueOrArrayKindHeapVector:
		v.HeapVector = &HeapVector{}
		return v.HeapVector.UnmarshalReader(r)
	default:
		return fmt.Errorf("invalid ValueOrArrayKind: %d", v.Kind)
	}
}
