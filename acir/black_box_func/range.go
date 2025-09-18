package blackboxfunc

import (
	"fmt"
	"io"
	shr "nr-groth16/acir/shared"

	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/rangecheck"
	"github.com/google/btree"
)

type Range[T shr.ACIRField] struct {
	Input FunctionInput[T]
}

func (a *Range[T]) UnmarshalReader(r io.Reader) error {
	if err := a.Input.UnmarshalReader(r); err != nil {
		return err
	}
	return nil
}

func (a Range[T]) Equals(other BlackBoxFunction) bool {
	value, ok := other.(*Range[T])
	return ok && a.Input.Equals(&value.Input)
}

func (a Range[T]) Define(api frontend.Builder, witnesses map[shr.Witness]frontend.Variable) error {
	if a.Input.FunctionInputKind == ACIRFunctionInputKindConstant {
		return nil
	}

	witness := a.Input.Witness
	if witness == nil {
		return fmt.Errorf("witness is nil for Range function input")
	}

	w, ok := witnesses[*witness]
	if !ok {
		return fmt.Errorf("witness %v not found in witnesses map", *witness)
	}

	rangechecker := rangecheck.New(api)
	rangechecker.Check(w, int(a.Input.NumberOfBits))
	return nil
}

func (a *Range[T]) FillWitnessTree(tree *btree.BTree) bool {
	if tree == nil {
		return false
	}
	tree.ReplaceOrInsert(*a.Input.Witness)
	return true
}
