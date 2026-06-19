package blackboxfunc

import (
	"github.com/reilabs/sunspot/go/acir/msgpackutil"
	shr "github.com/reilabs/sunspot/go/acir/shared"
	grumpkin "github.com/reilabs/sunspot/go/sw-grumpkin"

	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/algebra/algopts"
)

type MultiScalarMul[T shr.ACIRField, E constraint.Element] struct {
	Points    []FunctionInput[T]
	Scalars   []FunctionInput[T]
	predicate FunctionInput[T]
	Outputs   [2]shr.Witness
}

func (a *MultiScalarMul[T, E]) UnmarshalReader(r *msgpackutil.Reader) error {
	return msgpackutil.ReadStruct(r, "MultiScalarMul", []msgpackutil.Field{
		{Name: "points", Decode: func(r *msgpackutil.Reader) error { return msgpackutil.ReadVec(r, &a.Points) }},
		{Name: "scalars", Decode: func(r *msgpackutil.Reader) error { return msgpackutil.ReadVec(r, &a.Scalars) }},
		{Name: "predicate", Decode: a.predicate.UnmarshalReader},
		{Name: "outputs", Decode: func(r *msgpackutil.Reader) error { return msgpackutil.ReadArrayInto(r, a.Outputs[:]) }},
	})
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

	for i := range a.Outputs {
		if a.Outputs[i] != value.Outputs[i] {
			return false
		}
	}

	return true
}

func (a *MultiScalarMul[T, E]) Define(api frontend.Builder[E], witnesses map[shr.Witness]frontend.Variable) error {
	points := make([]*grumpkin.G1Affine, len(a.Points)/2)
	scalars := make([]interface{}, len(a.Scalars)/2)

	pred, err := a.predicate.ToVariable(witnesses)
	if err != nil {
		return err
	}

	for i := 0; i < len(a.Points); i += 2 {
		point, err := EmbeddedPointFromInputs(a.Points[i], a.Points[i+1], witnesses)
		if err != nil {
			return err
		}
		points[i/2] = &point
	}

	for i := 0; i < len(a.Scalars); i += 2 {
		scalar, err := ScalarFromLimbs(api, witnesses, a.Scalars[i], a.Scalars[i+1])
		if err != nil {
			return err
		}
		scalars[i/2] = scalar
	}

	output := grumpkin.G1Affine{
		X: witnesses[a.Outputs[0]],
		Y: witnesses[a.Outputs[1]],
	}

	for i := range points {
		points[i].AssertIsOnCurve(api)
	}
	output.AssertIsOnCurve(api)

	constrained_output := grumpkin.MultiScalarMul(api, points, scalars, algopts.WithCompleteArithmetic())

	// Predicate-gated equality on each coordinate: when pred=0 the
	// constraint is trivially satisfied.
	api.AssertIsEqual(frontend.Variable(0), api.Mul(pred, api.Sub(constrained_output.X, output.X)))
	api.AssertIsEqual(frontend.Variable(0), api.Mul(pred, api.Sub(constrained_output.Y, output.Y)))
	return nil
}

func (*MultiScalarMul[T, E]) SerdeName() string { return "MultiScalarMul" }
