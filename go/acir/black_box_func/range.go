package blackboxfunc

import (
	"encoding/binary"
	"fmt"
	"io"
	shr "sunspot/go/acir/shared"

	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/rangecheck"
	"github.com/google/btree"
)

type Range[T shr.ACIRField, E constraint.Element] struct {
	Input FunctionInput[T]
	nBits uint32
}

func (a *Range[T, E]) UnmarshalReader(r io.Reader) error {
	if err := a.Input.UnmarshalReader(r); err != nil {
		return err
	}
	if err := binary.Read(r, binary.LittleEndian, &a.nBits); err != nil {
		return err
	}
	return nil
}

func (a Range[T, E]) Equals(other BlackBoxFunction[E]) bool {
	value, ok := other.(*Range[T, E])
	return ok && a.Input.Equals(&value.Input)
}

func (a Range[T, E]) Define(api frontend.Builder[E], witnesses map[shr.Witness]frontend.Variable) error {
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
	rangechecker.Check(w, int(a.nBits))
	return nil
}

func (a *Range[T, E]) FillWitnessTree(tree *btree.BTree, index uint32) bool {
	if tree == nil {
		return false
	}
	if a.Input.IsWitness() {
		tree.ReplaceOrInsert(*a.Input.Witness + shr.Witness(index))
	}
	return true
}
