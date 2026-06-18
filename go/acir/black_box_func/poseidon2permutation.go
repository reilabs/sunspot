package blackboxfunc

import (
	"fmt"
	"sunspot/go/acir/msgpackutil"
	shr "sunspot/go/acir/shared"
	"sunspot/go/poseidon2"

	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
)

type Poseidon2Permutation[T shr.ACIRField, E constraint.Element] struct {
	Inputs  []FunctionInput[T]
	Outputs []shr.Witness
}

func (a *Poseidon2Permutation[T, E]) decode(tag int, r *msgpackutil.Reader) error {
	switch tag {
	case 0:
		return readFunctionInputVec(r, &a.Inputs)
	case 1:
		return shr.ReadWitnessVec(r, &a.Outputs)
	default:
		return fmt.Errorf("Poseidon2Permutation: unknown field tag %d", tag)
	}
}

func (a *Poseidon2Permutation[T, E]) Equals(other BlackBoxFunction[E]) bool {
	value, ok := other.(*Poseidon2Permutation[T, E])
	if !ok || len(a.Inputs) != len(value.Inputs) || len(a.Outputs) != len(value.Outputs) {
		return false
	}

	for i := range a.Inputs {
		if !a.Inputs[i].Equals(&value.Inputs[i]) {
			return false
		}
	}

	for i := range a.Outputs {
		if a.Outputs[i] != value.Outputs[i] {
			return false
		}
	}

	return true
}

func (a *Poseidon2Permutation[T, E]) Define(api frontend.Builder[E], witnesses map[shr.Witness]frontend.Variable) error {
	inputs := make([]frontend.Variable, 4)

	for i := range a.Inputs {
		input, err := a.Inputs[i].ToVariable(witnesses)
		if err != nil {
			return err
		}
		inputs[i] = input
	}

	poseidon2.Permute(api, inputs)

	for i := range a.Inputs {
		api.AssertIsEqual(inputs[i], witnesses[a.Outputs[i]])
	}
	return nil
}
