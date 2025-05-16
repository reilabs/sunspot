package brillig

import (
	"encoding/binary"
	"io"
	exp "nr-groth16/acir_decoder/expression"
	shr "nr-groth16/acir_decoder/shared"
)

type BrilligCall[T shr.ACIRField] struct {
	ID        uint32
	Inputs    []BrilligInputs[T]
	Outputs   []BrilligOutputs
	Predicate *exp.Expression[T]
}

func (b *BrilligCall[T]) UnmarshalReader(r io.Reader) error {
	if err := binary.Read(r, binary.LittleEndian, &b.ID); err != nil {
		return err
	}

	var numInputs uint32
	if err := binary.Read(r, binary.LittleEndian, &numInputs); err != nil {
		return err
	}
	b.Inputs = make([]BrilligInputs[T], numInputs)
	for i := uint32(0); i < numInputs; i++ {
		if err := b.Inputs[i].UnmarshalReader(r); err != nil {
			return err
		}
	}

	var numOutputs uint32
	if err := binary.Read(r, binary.LittleEndian, &numOutputs); err != nil {
		return err
	}
	b.Outputs = make([]BrilligOutputs, numOutputs)
	for i := uint32(0); i < numOutputs; i++ {
		if err := b.Outputs[i].UnmarshalReader(r); err != nil {
			return err
		}
	}

	var predicateExists uint8
	if err := binary.Read(r, binary.LittleEndian, &predicateExists); err != nil {
		return err
	}
	if predicateExists == 1 {
		b.Predicate = new(exp.Expression[T])
		if err := b.Predicate.UnmarshalReader(r); err != nil {
			return err
		}
	} else {
		b.Predicate = nil
	}

	return nil
}
