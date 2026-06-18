package blackboxfunc

import (
	"fmt"
	"sunspot/go/acir/msgpackutil"
	shr "sunspot/go/acir/shared"

	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/rangecheck"
)

type Range[T shr.ACIRField, E constraint.Element] struct {
	Input FunctionInput[T]
	nBits uint32
}

func (a *Range[T, E]) decode(tag int, r *msgpackutil.Reader) error {
	switch tag {
	case 0:
		return a.Input.UnmarshalReader(r)
	case 1:
		n, err := r.ReadU32()
		if err != nil {
			return err
		}
		a.nBits = n
		return nil
	default:
		return fmt.Errorf("RANGE: unknown field tag %d", tag)
	}
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
