package blackboxfunc

import (
	"encoding/binary"
	"io"
	shr "nr-groth16/acir/shared"

	"github.com/consensys/gnark/frontend"
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

func (a *Keccakf1600[T]) Equals(other BlackBoxFunction) bool {
	value, ok := other.(*Keccakf1600[T])
	if !ok || len(a.Inputs) != len(value.Inputs) {
		return false
	}

	for i := 0; i < 25; i++ {
		if !a.Inputs[i].Equals(&value.Inputs[i]) {
			return false
		}
	}

	for i := 0; i < 25; i++ {
		if a.Outputs[i] != value.Outputs[i] {
			return false
		}
	}

	return true
}

func (a *Keccakf1600[T]) Define(api frontend.API, witnesses map[shr.Witness]frontend.Variable) error {

	return nil
}
