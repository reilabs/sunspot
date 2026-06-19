package blackboxfunc

import (
	"fmt"
	"github.com/reilabs/sunspot/go/acir/msgpackutil"
	shr "github.com/reilabs/sunspot/go/acir/shared"

	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/rangecheck"
)

type Range[T shr.ACIRField, E constraint.Element] struct {
	Input FunctionInput[T]
	nBits uint32
}

func (a *Range[T, E]) UnmarshalReader(r *msgpackutil.Reader) error {
	return msgpackutil.ReadStruct(r, "Range", []msgpackutil.Field{
		{Name: "input", Decode: a.Input.UnmarshalReader},
		{Name: "num_bits", Decode: func(r *msgpackutil.Reader) error {
			n, err := r.ReadU32()
			if err != nil {
				return err
			}
			a.nBits = n
			return nil
		}},
	})
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

func (*Range[T, E]) SerdeName() string { return "RANGE" }
