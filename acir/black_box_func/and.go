package blackboxfunc

import (
	"fmt"
	"io"
	shr "nr-groth16/acir/shared"

	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/math/uints"
	"github.com/google/btree"
)

type And[T shr.ACIRField] struct {
	Lhs    FunctionInput[T]
	Rhs    FunctionInput[T]
	Output shr.Witness
}

func (a *And[T]) UnmarshalReader(r io.Reader) error {
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

func (a *And[T]) Equals(other BlackBoxFunction) bool {
	value, ok := other.(*And[T])
	return ok && a.Lhs.Equals(&value.Lhs) && a.Rhs.Equals(&value.Rhs) && a.Output.Equals(&value.Output)
}

func (a *And[T]) Define(api frontend.API, witnesses map[shr.Witness]frontend.Variable) error {
	uapi, err := uints.New[uints.U64](api)
	if err != nil {
		return err
	}
	lhs, err := a.Lhs.ToVariable(witnesses)
	if err != nil {
		return err
	}
	lhs_b := uapi.ValueOf(lhs)

	rhs, err := a.Rhs.ToVariable(witnesses)
	if err != nil {
		return err
	}
	rhs_b := uapi.ValueOf(rhs)
	output, ok := witnesses[a.Output]
	if !ok {
		return fmt.Errorf("witness %d not found in witnesses map", a.Output)
	}
	output_b := uapi.ValueOf(output)

	uapi.AssertEq(output_b, uapi.And(lhs_b, rhs_b))

	return nil
}

func (a *And[T]) FillWitnessTree(tree *btree.BTree) bool {
	if tree == nil {
		return false
	}

	if a.Lhs.FunctionInputKind == 1 {
		tree.ReplaceOrInsert(*a.Lhs.Witness)
	}
	if a.Rhs.FunctionInputKind == 1 {
		tree.ReplaceOrInsert(*a.Rhs.Witness)
	}
	tree.ReplaceOrInsert(a.Output)

	return true
}
