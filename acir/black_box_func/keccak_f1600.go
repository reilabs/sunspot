package blackboxfunc

import (
	"encoding/binary"
	"io"
	shr "nr-groth16/acir/shared"
)

type Keccakf1600[T shr.ACIRField] struct {
	Inputs  [25]FunctionInput[T]
	Outputs [25]shr.Witness
}

func (a *Keccakf1600[T]) UnmarshalReader(r io.Reader) error {
	for i := 0; i < 25; i++ {
		if err := a.Inputs[i].UnmarshalReader(r); err != nil {
			return err
		}
	}

	if err := binary.Read(r, binary.LittleEndian, &a.Outputs); err != nil {
		return err
	}

	return nil
}

func (a *Keccakf1600[T]) Equals(other *Keccakf1600[T]) bool {
	if len(a.Inputs) != len(other.Inputs) {
		return false
	}

	for i := 0; i < 25; i++ {
		if !a.Inputs[i].Equals(&other.Inputs[i]) {
			return false
		}
	}

	for i := 0; i < 25; i++ {
		if a.Outputs[i] != other.Outputs[i] {
			return false
		}
	}

	return true
}
