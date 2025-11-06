package blackboxfunc

import (
	"encoding/binary"
	"io"
	shr "sunspot/acir/shared"
	grumpkin "sunspot/sw-grumpkin"

	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/google/btree"
)

type MultiScalarMul[T shr.ACIRField, E constraint.Element] struct {
	Points  []FunctionInput[T]
	Scalars []FunctionInput[T]
	Outputs [3]shr.Witness
}

func (a *MultiScalarMul[T, E]) UnmarshalReader(r io.Reader) error {

	var numPoints uint64
	if err := binary.Read(r, binary.LittleEndian, &numPoints); err != nil {
		return err
	}

	a.Points = make([]FunctionInput[T], numPoints)
	for i := uint64(0); i < numPoints; i++ {
		if err := a.Points[i].UnmarshalReader(r); err != nil {
			return err
		}
	}

	var numScalars uint64
	if err := binary.Read(r, binary.LittleEndian, &numScalars); err != nil {
		return err
	}

	a.Scalars = make([]FunctionInput[T], numScalars)
	for i := uint64(0); i < numScalars; i++ {
		if err := a.Scalars[i].UnmarshalReader(r); err != nil {
			return err
		}
	}

	for i := 0; i < 3; i++ {
		if err := a.Outputs[i].UnmarshalReader(r); err != nil {
			return err
		}
	}

	return nil
}

func (a *MultiScalarMul[T, E]) Equals(other BlackBoxFunction[E]) bool {
	value, ok := other.(*MultiScalarMul[T, E])

	if !ok || len(a.Points) != len(value.Points) || len(a.Scalars) != len(value.Scalars) {
		return false
	}

	for i := range a.Points {
		if !a.Points[i].Equals(&value.Points[i]) {
			return false
		}
	}

	for i := range a.Scalars {
		if !a.Scalars[i].Equals(&value.Scalars[i]) {
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

func (a *MultiScalarMul[T, E]) Define(api frontend.Builder[E], witnesses map[shr.Witness]frontend.Variable) error {
	points := make([]*grumpkin.G1Affine, len(a.Points)/3)

	scalars := make([]interface{}, len(a.Scalars)/2)

	for i := 0; i < len(a.Points); i += 3 {
		pointX, err := a.Points[i].ToVariable(witnesses)
		if err != nil {
			return err
		}

		pointY, err := a.Points[i+1].ToVariable(witnesses)
		if err != nil {
			return err
		}

		point := grumpkin.G1Affine{
			X: pointX,
			Y: pointY,
		}
		points[i/3] = &point
	}

	for i := 0; i < len(a.Scalars); i += 2 {
		scalar, err := a.Scalars[i].ToVariable(witnesses)
		if err != nil {
			return err
		}
		scalars[i/2] = scalar
	}

	output := grumpkin.G1Affine{
		X: witnesses[a.Outputs[0]],
		Y: witnesses[a.Outputs[1]],
	}

	constrained_output := grumpkin.MultiScalarMul(api, points, scalars)
	output.AssertIsEqual(api, *constrained_output)
	return nil
}

func (a *MultiScalarMul[T, E]) FillWitnessTree(tree *btree.BTree, index uint32) bool {
	if tree == nil {
		return false
	}
	for _, point := range a.Points {
		if point.IsWitness() {
			tree.ReplaceOrInsert(*point.Witness + shr.Witness(index))
		}
	}

	for _, scalar := range a.Scalars {
		if scalar.IsWitness() {
			tree.ReplaceOrInsert(*scalar.Witness + shr.Witness(index))
		}
	}
	for _, output := range a.Outputs {
		tree.ReplaceOrInsert(output + shr.Witness(index))
	}
	return true
}
