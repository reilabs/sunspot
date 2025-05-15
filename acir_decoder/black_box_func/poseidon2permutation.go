package blackboxfunc

import (
	"encoding/binary"
	"io"
	shr "nr-groth16/acir_decoder/shared"
)

type Poseidon2Permutation[T shr.ACIRField] struct {
	Inputs  []FunctionInput[T]
	Outputs []shr.Witness
	Len     uint32
}

func (a *Poseidon2Permutation[T]) UnmarshalReader(r io.Reader) error {
	var NumInputs uint32
	if err := binary.Read(r, binary.LittleEndian, &NumInputs); err != nil {
		return err
	}
	a.Inputs = make([]FunctionInput[T], NumInputs)
	for i := uint32(0); i < NumInputs; i++ {
		if err := a.Inputs[i].UnmarshalReader(r); err != nil {
			return err
		}
	}

	var NumOutputs uint32
	if err := binary.Read(r, binary.LittleEndian, &NumOutputs); err != nil {
		return err
	}

	a.Outputs = make([]shr.Witness, NumOutputs)
	if err := binary.Read(r, binary.LittleEndian, &a.Outputs); err != nil {
		return err
	}

	if err := binary.Read(r, binary.LittleEndian, &a.Len); err != nil {
		return err
	}

	return nil
}
