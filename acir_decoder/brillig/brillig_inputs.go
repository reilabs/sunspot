package brillig

import (
	"encoding/binary"
	"fmt"
	"io"
	exp "nr-groth16/acir_decoder/expression"
	shr "nr-groth16/acir_decoder/shared"
)

type BrilligInputs[T shr.ACIRField] struct {
	Kind    BrilligInputsKind
	Single  *exp.Expression[T]
	Array   *[]exp.Expression[T]
	BlockID *uint32
}

type BrilligInputsKind uint32

const (
	ACIRBrilligInputsKindSingle BrilligInputsKind = iota
	ACIRBrilligInputsKindArray
	ACIRBrilligInputsKindMemoryArray
)

func (b *BrilligInputs[T]) UnmarshalReader(r io.Reader) error {
	if err := binary.Read(r, binary.LittleEndian, &b.Kind); err != nil {
		return err
	}
	switch b.Kind {
	case ACIRBrilligInputsKindSingle:
		b.Single = new(exp.Expression[T])
		if err := b.Single.UnmarshalReader(r); err != nil {
			return err
		}
	case ACIRBrilligInputsKindArray:
		var numInputs uint32
		if err := binary.Read(r, binary.LittleEndian, &numInputs); err != nil {
			return err
		}
		b.Array = new([]exp.Expression[T])
		*b.Array = make([]exp.Expression[T], numInputs)
		for i := uint32(0); i < numInputs; i++ {
			if err := (*b.Array)[i].UnmarshalReader(r); err != nil {
				return err
			}
		}
	case ACIRBrilligInputsKindMemoryArray:
		var blockID uint32
		if err := binary.Read(r, binary.LittleEndian, &blockID); err != nil {
			return err
		}
		b.BlockID = new(uint32)
		*b.BlockID = blockID
	default:
		return fmt.Errorf("unknown BrilligInputsKind: %d", b.Kind)
	}
	return nil
}
