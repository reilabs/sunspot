package blackboxfunc

import (
	"encoding/binary"
	"io"
	shr "sunspot/go/acir/shared"
	"sunspot/go/poseidon2"

	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/google/btree"
)

type Poseidon2Permutation[T shr.ACIRField, E constraint.Element] struct {
	Inputs  []FunctionInput[T]
	Outputs []shr.Witness
}

func (a *Poseidon2Permutation[T, E]) UnmarshalReader(r io.Reader) error {
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

	return nil
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

func (a *Poseidon2Permutation[T, E]) FillWitnessTree(tree *btree.BTree, index uint32) bool {
	if tree == nil {
		return false
	}

	for _, input := range a.Inputs {
		if input.IsWitness() {
			tree.ReplaceOrInsert(*input.Witness + shr.Witness(index))
		}
	}

	for _, output := range a.Outputs {
		tree.ReplaceOrInsert(output + shr.Witness(index))
	}

	return true
}
