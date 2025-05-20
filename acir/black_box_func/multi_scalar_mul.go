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
	var numPoints uint32
	if err := binary.Read(r, binary.LittleEndian, &numPoints); err != nil {
		return err
	}

	a.Points = make([]FunctionInput[T], numPoints)
	for i := uint32(0); i < numPoints; i++ {
		if err := a.Points[i].UnmarshalReader(r); err != nil {
			return err
		}
	}

	var numScalars uint32
	if err := binary.Read(r, binary.LittleEndian, &numScalars); err != nil {
		return err
	}

	if numScalars != numPoints {
		panic("numScalars != numPoints")
	}

	a.Scalars = make([]FunctionInput[T], numScalars)
	for i := uint32(0); i < numScalars; i++ {
		if err := a.Scalars[i].UnmarshalReader(r); err != nil {
			return err
		}
	}

	return nil
}
