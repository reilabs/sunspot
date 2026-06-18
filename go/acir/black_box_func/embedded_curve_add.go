package blackboxfunc

import (
	"fmt"
	"sunspot/go/acir/msgpackutil"
	shr "sunspot/go/acir/shared"
	grumpkin "sunspot/go/sw-grumpkin"

	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
)

type EmbeddedCurveAdd[T shr.ACIRField, E constraint.Element] struct {
	Input1    [2]FunctionInput[T]
	Input2    [2]FunctionInput[T]
	predicate FunctionInput[T]
	Outputs   [2]shr.Witness
}

func (a *EmbeddedCurveAdd[T, E]) decode(f msgpackutil.Field, r *msgpackutil.Reader) error {
	switch f.Tag {
	case 0:
		return msgpackutil.ReadArrayInto(r, a.Input1[:])
	case 1:
		return msgpackutil.ReadArrayInto(r, a.Input2[:])
	case 2:
		return a.predicate.UnmarshalReader(r)
	case 3:
		n, err := r.ReadArrayLen()
		if err != nil {
			return err
		}
		if n != 2 {
			return fmt.Errorf("EmbeddedCurveAdd.outputs: expected 2-tuple, got %d", n)
		}
		if err := a.Outputs[0].UnmarshalReader(r); err != nil {
			return err
		}
		return a.Outputs[1].UnmarshalReader(r)
	default:
		return fmt.Errorf("EmbeddedCurveAdd: unknown field %s", f)
	}
}

func (a *EmbeddedCurveAdd[T, E]) Equals(other BlackBoxFunction[E]) bool {
	value, ok := other.(*EmbeddedCurveAdd[T, E])
	if !ok {
		return false
	}
	for i := 0; i < 2; i++ {
		if !a.Input1[i].Equals(&value.Input1[i]) || !a.Input2[i].Equals(&value.Input2[i]) {
			return false
		}
		if a.Outputs[i] != value.Outputs[i] {
			return false
		}
	}
	return true
}

func (a *EmbeddedCurveAdd[T, E]) Define(api frontend.Builder[E], witnesses map[shr.Witness]frontend.Variable) error {
	pred, err := a.predicate.ToVariable(witnesses)
	if err != nil {
		return err
	}

	point1, err := EmbeddedPointFromInputs(a.Input1[0], a.Input1[1], witnesses)
	if err != nil {
		return err
	}
	point2, err := EmbeddedPointFromInputs(a.Input2[0], a.Input2[1], witnesses)
	if err != nil {
		return err
	}

	output := grumpkin.G1Affine{
		X: witnesses[a.Outputs[0]],
		Y: witnesses[a.Outputs[1]],
	}

	point1.AssertIsOnCurve(api)
	point2.AssertIsOnCurve(api)
	output.AssertIsOnCurve(api)

	constrained_output := point1.AddUnified(api, point2)
	api.AssertIsEqual(frontend.Variable(0), api.Mul(pred, api.Sub(constrained_output.X, output.X)))
	api.AssertIsEqual(frontend.Variable(0), api.Mul(pred, api.Sub(constrained_output.Y, output.Y)))
	return nil
}

func (*EmbeddedCurveAdd[T, E]) SerdeName() string { return "EmbeddedCurveAdd" }
