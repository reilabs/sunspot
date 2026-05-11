package blackboxfunc

import (
	"encoding/binary"
	"io"
	shr "sunspot/go/acir/shared"
	grumpkin "sunspot/go/sw-grumpkin"

	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
)

type MultiScalarMul[T shr.ACIRField, E constraint.Element] struct {
	Points    []FunctionInput[T]
	Scalars   []FunctionInput[T]
	predicate FunctionInput[T]
	Outputs   [3]shr.Witness
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
	if err := a.predicate.UnmarshalReader(r); err != nil {
		return err
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

	pred, err := a.predicate.ToVariable(witnesses)
	if err != nil {
		return err
	}

	for i := 0; i < len(a.Points); i += 3 {
		point, err := EmbeddedPointFromInputs(api, witnesses, pred,
			[3]FunctionInput[T]{a.Points[i], a.Points[i+1], a.Points[i+2]})
		if err != nil {
			return err
		}
		points[i/3] = &point
	}

	for i := 0; i < len(a.Scalars); i += 2 {
		scalar, err := ScalarFromLimbs(api, witnesses, a.Scalars[i], a.Scalars[i+1])
		if err != nil {
			return err
		}
		scalars[i/2] = scalar
	}

	output := maskedEmbeddedPoint(api, pred,
		witnesses[a.Outputs[0]], witnesses[a.Outputs[1]], witnesses[a.Outputs[2]])

	constrained_output := grumpkin.MultiScalarMul(api, points, scalars)

	// To assert the two points are the same (and ignore if predicate is zero), we have to split into
	// its X and Y coordinates
	api.AssertIsEqual(frontend.Variable(0), api.Mul(pred, api.Sub(constrained_output.X, output.X)))
	api.AssertIsEqual(frontend.Variable(0), api.Mul(pred, api.Sub(constrained_output.Y, output.Y)))
	return nil
}

