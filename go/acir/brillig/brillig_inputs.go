package brillig

import (
	"encoding/binary"
	"fmt"
	"io"
	exp "sunspot/acir/expression"
	shr "sunspot/acir/shared"

	"github.com/consensys/gnark/constraint"
)

type BrilligInputs[T shr.ACIRField, E constraint.Element] struct {
	Kind    BrilligInputsKind
	Single  *exp.Expression[T, E]
	Array   *[]exp.Expression[T, E]
	BlockID *uint32
}

type BrilligInputsKind uint32

const (
	ACIRBrilligInputsKindSingle BrilligInputsKind = iota
	ACIRBrilligInputsKindArray
	ACIRBrilligInputsKindMemoryArray
)

func (b *BrilligInputs[T, E]) UnmarshalReader(r io.Reader) error {
	if err := binary.Read(r, binary.LittleEndian, &b.Kind); err != nil {
		return err
	}
	switch b.Kind {
	case ACIRBrilligInputsKindSingle:
		b.Single = new(exp.Expression[T, E])
		if err := b.Single.UnmarshalReader(r); err != nil {
			return err
		}
	case ACIRBrilligInputsKindArray:
		var numInputs uint64
		if err := binary.Read(r, binary.LittleEndian, &numInputs); err != nil {
			return err
		}
		b.Array = new([]exp.Expression[T, E])
		*b.Array = make([]exp.Expression[T, E], numInputs)
		for i := uint64(0); i < numInputs; i++ {
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
