package blackboxfunc

import (
	"encoding/binary"
	"fmt"
	"io"
	shr "sunspot/go/acir/shared"

	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/math/uints"
	"github.com/google/btree"
)

type Xor[T shr.ACIRField, E constraint.Element] struct {
	Lhs    FunctionInput[T]
	Rhs    FunctionInput[T]
	Output shr.Witness
	nBits  uint32
}

func (a *Xor[T, E]) UnmarshalReader(r io.Reader) error {
	if err := a.Lhs.UnmarshalReader(r); err != nil {
		return err
	}
	if err := a.Rhs.UnmarshalReader(r); err != nil {
		return err
	}
	if err := binary.Read(r, binary.LittleEndian, &a.nBits); err != nil {
		return err
	}
	if err := a.Output.UnmarshalReader(r); err != nil {
		return err
	}
	return nil
}

func (a *Xor[T, E]) Equals(other BlackBoxFunction[E]) bool {
	value, ok := other.(*Xor[T, E])

	if !ok || !a.Lhs.Equals(&value.Lhs) || !a.Rhs.Equals(&value.Rhs) {
		return false
	}
	return a.Output == value.Output
}

func (a *Xor[T, E]) Define(api frontend.Builder[E], witnesses map[shr.Witness]frontend.Variable) error {
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

	uapi.AssertEq(output_b, uapi.Xor(lhs_b, rhs_b))

	return nil
}

func (a *Xor[T, E]) FillWitnessTree(tree *btree.BTree, index uint32) bool {
	if tree == nil {
		return false
	}

	if a.Lhs.IsWitness() {
		tree.ReplaceOrInsert(*a.Lhs.Witness + shr.Witness(index))
	}
	if a.Rhs.IsWitness() {
		tree.ReplaceOrInsert(*a.Rhs.Witness + shr.Witness(index))
	}
	tree.ReplaceOrInsert(a.Output + shr.Witness(index))

	return true
}
