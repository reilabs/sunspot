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
	return ok && a.Input.Equals(&value.Input) && a.nBits == value.nBits
}

func (a Range[T, E]) Define(api frontend.Builder[E], witnesses map[shr.Witness]frontend.Variable) error {
	input, err := a.Input.ToVariable(witnesses)
	if err != nil {
		return fmt.Errorf("failed to resolve Range function input: %w", err)
	}

	rangechecker := rangecheck.New(api)
	rangechecker.Check(input, int(a.nBits))
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
