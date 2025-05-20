package blackboxfunc

import (
	"encoding/binary"
	"io"
	shr "nr-groth16/acir/shared"
)

type EmbeddedCurveAdd[T shr.ACIRField] struct {
	Input1  [3]FunctionInput[T]
	Input2  [3]FunctionInput[T]
	Outputs [3]shr.Witness
}

func (a *EmbeddedCurveAdd[T]) UnmarshalReader(r io.Reader) error {
	for i := 0; i < 3; i++ {
		if err := a.Input1[i].UnmarshalReader(r); err != nil {
			return err
		}
	}

	for i := 0; i < 3; i++ {
		if err := a.Input2[i].UnmarshalReader(r); err != nil {
			return err
		}
	}

	if err := binary.Read(r, binary.LittleEndian, &a.Outputs); err != nil {
		return err
	}

	return nil
}
