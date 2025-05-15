package blackboxfunc

import (
	"encoding/binary"
	"io"
	shr "nr-groth16/acir_decoder/shared"
)

type BigIntFromLEBytes[T shr.ACIRField] struct {
	Inputs  []FunctionInput[T]
	Modulus []uint8
	Output  uint32
}

func (a *BigIntFromLEBytes[T]) UnmarshalReader(r io.Reader) error {
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

	NumModulus := uint32(0)
	if err := binary.Read(r, binary.LittleEndian, &NumModulus); err != nil {
		return err
	}

	a.Modulus = make([]uint8, NumModulus)
	if err := binary.Read(r, binary.LittleEndian, &a.Modulus); err != nil {
		return err
	}

	if err := binary.Read(r, binary.LittleEndian, &a.Output); err != nil {
		return err
	}

	return nil
}
