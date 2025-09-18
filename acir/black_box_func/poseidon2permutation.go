package blackboxfunc

import (
	"encoding/binary"
	"io"
	shr "nr-groth16/acir/shared"

	"github.com/consensys/gnark/frontend"
)

type Poseidon2Permutation[T shr.ACIRField] struct {
	Inputs  []FunctionInput[T]
	Outputs []shr.Witness
	Len     uint32
}

func (a *Poseidon2Permutation[T]) UnmarshalReader(r io.Reader) error {
	var NumInputs uint64
	if err := binary.Read(r, binary.LittleEndian, &NumInputs); err != nil {
		return err
	}
	a.Inputs = make([]FunctionInput[T], NumInputs)
	for i := uint64(0); i < NumInputs; i++ {
		if err := a.Inputs[i].UnmarshalReader(r); err != nil {
			return err
		}
	}

	var NumOutputs uint64
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

func (a *Poseidon2Permutation[T]) Equals(other *Poseidon2Permutation[T]) bool {
	if len(a.Inputs) != len(other.Inputs) || len(a.Outputs) != len(other.Outputs) || a.Len != other.Len {
		return false
	}

	for i := range a.Inputs {
		if !a.Inputs[i].Equals(&other.Inputs[i]) {
			return false
		}
	}

	for i := range a.Outputs {
		if a.Outputs[i] != other.Outputs[i] {
			return false
		}
	}

	return true
}

func (a *Poseidon2Permutation[T]) Define(api frontend.Builder, witnesses map[shr.Witness]frontend.Variable) error {
	return nil
}
