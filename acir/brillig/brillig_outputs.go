package brillig

import (
	"encoding/binary"
	"fmt"
	"io"
	shr "sunpot/acir/shared"
)

type BrilligOutputs struct {
	Kind   BrilligOutputsKind
	Single *shr.Witness
	Array  *[]shr.Witness
}

type BrilligOutputsKind uint32

const (
	ACIRBrilligOutputsKindSimple BrilligOutputsKind = iota
	ACIRBrilligOutputsKindArray
)

func (b *BrilligOutputs) UnmarshalReader(r io.Reader) error {
	if err := binary.Read(r, binary.LittleEndian, &b.Kind); err != nil {
		return err
	}

	switch b.Kind {
	case ACIRBrilligOutputsKindSimple:
		b.Single = new(shr.Witness)
		if err := b.Single.UnmarshalReader(r); err != nil {
			return err
		}
	case ACIRBrilligOutputsKindArray:
		var numOutputs uint64
		if err := binary.Read(r, binary.LittleEndian, &numOutputs); err != nil {
			return err
		}
		b.Array = new([]shr.Witness)
		*b.Array = make([]shr.Witness, numOutputs)
		for i := uint64(0); i < numOutputs; i++ {
			if err := (*b.Array)[i].UnmarshalReader(r); err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("unknown BrilligOutputsKind: %d", b.Kind)
	}

	return nil
}
