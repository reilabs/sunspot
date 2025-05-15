package blackboxfunc

import (
	"encoding/binary"
	"io"
	shr "nr-groth16/acir_decoder/shared"
)

type Blake3[T shr.ACIRField] struct {
	Inputs []FunctionInput[T]
	output [32]shr.Witness
}

func (a *Blake3[T]) UnmarshalReader(r io.Reader) error {
	NumInputs := uint32(0)
	if err := binary.Read(r, binary.LittleEndian, &NumInputs); err != nil {
		return err
	}

	a.Inputs = make([]FunctionInput[T], NumInputs)
	for i := uint32(0); i < NumInputs; i++ {
		if err := a.Inputs[i].UnmarshalReader(r); err != nil {
			return err
		}
	}

	if err := binary.Read(r, binary.LittleEndian, &a.output); err != nil {
		return err
	}

	return nil
}
