package blackboxfunc

import (
	"fmt"
	"sunspot/go/acir/msgpackutil"
	shr "sunspot/go/acir/shared"

	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/math/uints"
)

type And[T shr.ACIRField, E constraint.Element] struct {
	Lhs    FunctionInput[T]
	Rhs    FunctionInput[T]
	nBits  uint32
	Output shr.Witness
}

func (a *And[T, E]) decode(tag int, r *msgpackutil.Reader) error {
	switch tag {
	case 0:
		return a.Lhs.UnmarshalReader(r)
	case 1:
		return a.Rhs.UnmarshalReader(r)
	case 2:
		n, err := r.ReadU32()
		if err != nil {
			return err
		}
		if n > 128 {
			return fmt.Errorf("AND: num_bits=%d exceeds supported maximum of 128", n)
		}
		a.nBits = n
		return nil
	case 3:
		return a.Output.UnmarshalReader(r)
	default:
		return fmt.Errorf("AND: unknown field tag %d", tag)
	}
}

func (a *And[T, E]) Equals(other BlackBoxFunction[E]) bool {
	value, ok := other.(*And[T, E])
	return ok && a.Lhs.Equals(&value.Lhs) && a.Rhs.Equals(&value.Rhs) && a.Output.Equals(&value.Output) && a.nBits == value.nBits
}

func (a *And[T, E]) Define(api frontend.Builder[E], witnesses map[shr.Witness]frontend.Variable) error {
	uapi, err := uints.New[uints.U64](api)
	if err != nil {
		return err
	}
	return defineBitwise(api, uapi, witnesses, a.Lhs, a.Rhs, a.Output, int(a.nBits), uapi.And)
}
