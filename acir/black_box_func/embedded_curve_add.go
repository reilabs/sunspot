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

func (a *EmbeddedCurveAdd[T]) Equals(other *EmbeddedCurveAdd[T]) bool {
	if len(a.Input1) != len(other.Input1) || len(a.Input2) != len(other.Input2) {
		return false
	}

	for i := 0; i < 3; i++ {
		if !a.Input1[i].Equals(&other.Input1[i]) || !a.Input2[i].Equals(&other.Input2[i]) {
			return false
		}
	}

	for i := 0; i < 3; i++ {
		if a.Outputs[i] != other.Outputs[i] {
			return false
		}
	}

	return true
}
