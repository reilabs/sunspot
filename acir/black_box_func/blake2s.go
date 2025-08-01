package blackboxfunc

import (
	"encoding/binary"
	"io"
	shr "nr-groth16/acir/shared"

	"github.com/consensys/gnark/frontend"
)

type Blake2s[T shr.ACIRField] struct {
	Inputs  []FunctionInput[T]
	Outputs [32]shr.Witness
}

func (a *Blake2s[T]) UnmarshalReader(r io.Reader) error {
	NumInputs := uint64(0)
	if err := binary.Read(r, binary.LittleEndian, &NumInputs); err != nil {
		return err
	}

	a.Inputs = make([]FunctionInput[T], NumInputs)
	for i := uint64(0); i < NumInputs; i++ {
		if err := a.Inputs[i].UnmarshalReader(r); err != nil {
			return err
		}
	}

	if err := binary.Read(r, binary.LittleEndian, &a.Outputs); err != nil {
		return err
	}

	return nil
}

func (a *Blake2s[T]) Equals(other *Blake2s[T]) bool {
	if len(a.Inputs) != len(other.Inputs) {
		return false
	}
	for i := range a.Inputs {
		if !a.Inputs[i].Equals(&other.Inputs[i]) {
			return false
		}
	}

	for i := 0; i < 32; i++ {
		if a.Outputs[i] != other.Outputs[i] {
			return false
		}
	}

	return true
}

func (a *Blake2s[T]) Define(api frontend.API, witnesses map[shr.Witness]frontend.Variable) error {

	for i := 0; i < 32; i++ {
		api.AssertIsEqual(a.Outputs[i], witnesses[a.Outputs[i]])
	}

	// Here you would implement the actual Blake2s hashing logic using the inputs.
	// This is a placeholder for the actual implementation.
	// api.Blake2s(a.Inputs, a.Outputs[:])
	return nil
}
