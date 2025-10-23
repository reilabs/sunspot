package blackboxfunc

import (
	"encoding/binary"
	"io"
	shr "sunpot/acir/shared"

	grumpkin "sunpot/sw-grumpkin"

	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/google/btree"
)

type EmbeddedCurveAdd[T shr.ACIRField, E constraint.Element] struct {
	Input1  [3]FunctionInput[T]
	Input2  [3]FunctionInput[T]
	Outputs [3]shr.Witness
}

func (a *EmbeddedCurveAdd[T, E]) UnmarshalReader(r io.Reader) error {
	for i := 0; i < 3; i++ {
		if err := a.Input1[i].UnmarshalReader(r); err != nil {
			return err
		}
	}

	for i := 0; i < 3; i++ {
		if err := a.Input2[i].UnmarshalReader(r); err != nil {
			return err
		}
	}

	if err := binary.Read(r, binary.LittleEndian, &a.Outputs); err != nil {
		return err
	}

	return nil
}

func (a *EmbeddedCurveAdd[T, E]) Equals(other BlackBoxFunction[E]) bool {
	value, ok := other.(*EmbeddedCurveAdd[T, E])
	if !ok || len(a.Input1) != len(value.Input1) || len(a.Input2) != len(value.Input2) {
		return false
	}

	for i := 0; i < 3; i++ {
		if !a.Input1[i].Equals(&value.Input1[i]) || !a.Input2[i].Equals(&value.Input2[i]) {
			return false
		}
	}

	for i := 0; i < 3; i++ {
		if a.Outputs[i] != value.Outputs[i] {
			return false
		}
	}

	return true
}

func (a *EmbeddedCurveAdd[T, E]) Define(api frontend.Builder[E], witnesses map[shr.Witness]frontend.Variable) error {
	// Initialise points and pairs
	point1X, err := a.Input1[0].ToVariable(witnesses)
	if err != nil {
		return err
	}

	point1Y, err := a.Input1[1].ToVariable(witnesses)
	if err != nil {
		return err
	}
	point2X, err := a.Input2[0].ToVariable(witnesses)
	if err != nil {
		return err
	}
	point2Y, err := a.Input2[1].ToVariable(witnesses)
	if err != nil {
		return err
	}

	x := grumpkin.G1Affine{
		X: point1X,
		Y: point1Y,
	}

	y := grumpkin.G1Affine{
		X: point2X,
		Y: point2Y,
	}

	z := grumpkin.G1Affine{
		X: witnesses[a.Outputs[0]],
		Y: witnesses[a.Outputs[1]],
	}

	// Assert that the addition is correct
	z.AssertIsEqual(api, *x.AddUnified(api, y))

	return nil
}

func (a *EmbeddedCurveAdd[T, E]) FillWitnessTree(tree *btree.BTree, index uint32) bool {
	if tree == nil {
		return false
	}
	for _, input := range a.Input1 {

		tree.ReplaceOrInsert(*input.Witness + shr.Witness(index))
	}

	for _, input := range a.Input2 {
		tree.ReplaceOrInsert(*input.Witness + shr.Witness(index))
	}

	for _, output := range a.Outputs {
		tree.ReplaceOrInsert(output + shr.Witness(index))
	}
	return true
}
