package blackboxfunc

import (
	"encoding/binary"
	"io"
	shr "nr-groth16/acir/shared"
)

type MultiScalarMul[T shr.ACIRField] struct {
	Points  []FunctionInput[T]
	Scalars []FunctionInput[T]
	Outputs [3]shr.Witness
}

func (a *MultiScalarMul[T]) UnmarshalReader(r io.Reader) error {
	var numPoints uint64
	if err := binary.Read(r, binary.LittleEndian, &numPoints); err != nil {
		return err
	}

	a.Points = make([]FunctionInput[T], numPoints)
	for i := uint64(0); i < numPoints; i++ {
		if err := a.Points[i].UnmarshalReader(r); err != nil {
			return err
		}
	}

	var numScalars uint64
	if err := binary.Read(r, binary.LittleEndian, &numScalars); err != nil {
		return err
	}

	a.Scalars = make([]FunctionInput[T], numScalars)
	for i := uint64(0); i < numScalars; i++ {
		if err := a.Scalars[i].UnmarshalReader(r); err != nil {
			return err
		}
	}

	return nil
}

func (a *MultiScalarMul[T]) Equals(other *MultiScalarMul[T]) bool {
	if len(a.Points) != len(other.Points) || len(a.Scalars) != len(other.Scalars) {
		return false
	}

	for i := range a.Points {
		if !a.Points[i].Equals(&other.Points[i]) {
			return false
		}
	}

	for i := range a.Scalars {
		if !a.Scalars[i].Equals(&other.Scalars[i]) {
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
