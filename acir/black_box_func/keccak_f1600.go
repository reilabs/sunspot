package blackboxfunc

import (
	"encoding/binary"
	"fmt"
	"io"
	shr "nr-groth16/acir/shared"

	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/math/uints"
	"github.com/consensys/gnark/std/permutation/keccakf"
	"github.com/google/btree"
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

func (a *Keccakf1600[T]) Define(api frontend.Builder, witnesses map[shr.Witness]frontend.Variable) error {
	uapi, err := uints.New[uints.U64](api)
	if err != nil {
		return err
	}
	var keccak_inputs [25]uints.U64
	for i, input := range a.Inputs {
		v, err := input.ToVariable(witnesses)
		if err != nil {
			return fmt.Errorf("unable to get input as variable, index %d", i)
		}
		keccak_inputs[i] = uapi.ValueOf(v)
	}

	var keccak_outputs [25]uints.U64
	for i, output := range a.Outputs {
		v := witnesses[output]
		keccak_outputs[i] = uapi.ValueOf(v)
	}

	constrained_outputs := keccakf.Permute(uapi, keccak_inputs)

	for i := 0; i < 25; i++ {
		uapi.AssertEq(constrained_outputs[i], keccak_outputs[i])
	}
	return nil
}

func (a *Keccakf1600[T]) FillWitnessTree(tree *btree.BTree) bool {
	if tree == nil {
		return false
	}
	for _, input := range a.Inputs {

		tree.ReplaceOrInsert(*input.Witness)
	}
	return true
}
