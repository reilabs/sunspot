package blackboxfunc

import (
	"fmt"
	"io"
	shr "nr-groth16/acir/shared"

	"github.com/consensys/gnark/frontend"
	"github.com/google/btree"
)

type Xor[T shr.ACIRField] struct {
	Lhs    FunctionInput[T]
	Rhs    FunctionInput[T]
	Output shr.Witness
}

func (a *Xor[T]) UnmarshalReader(r io.Reader) error {
	if err := a.Lhs.UnmarshalReader(r); err != nil {
		return err
	}
	if err := a.Rhs.UnmarshalReader(r); err != nil {
		return err
	}
	if err := a.Output.UnmarshalReader(r); err != nil {
		return err
	}
	return nil
}

func (a *Xor[T]) Equals(other BlackBoxFunction) bool {
	value, ok := other.(*Xor[T])

	if !ok || !a.Lhs.Equals(&value.Lhs) || !a.Rhs.Equals(&value.Rhs) {
		return false
	}
	return a.Output == value.Output
}

func (a *Xor[T]) Define(api frontend.API, witnesses map[shr.Witness]frontend.Variable) error {
	lhs, err := a.Lhs.ToVariable(witnesses)
	if err != nil {
		return err
	}
	lhs_binary := api.ToBinary(lhs)
	rhs, err := a.Rhs.ToVariable(witnesses)
	if err != nil {
		return err
	}
	rhs_binary := api.ToBinary(rhs)
	output, ok := witnesses[a.Output]
	if !ok {
		return fmt.Errorf("witness %d not found in witnesses map", a.Output)
	}
	output_binary := api.ToBinary(output)
	verifiable_len := min(len(lhs_binary), len(rhs_binary), len(output_binary))
	for i := 0; i < verifiable_len; i++ {
		lhs_bit := lhs_binary[i]
		rhs_bit := rhs_binary[i]
		output_bit := output_binary[i]

		api.AssertIsEqual(api.Xor(lhs_bit, rhs_bit), output_bit)
	}

	return nil
}

func (a *Xor[T]) FillWitnessTree(tree *btree.BTree) bool {
	if tree == nil {
		return false
	}

	tree.ReplaceOrInsert(*a.Lhs.Witness)
	tree.ReplaceOrInsert(*a.Rhs.Witness)

	return true
}
