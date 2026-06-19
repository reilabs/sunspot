package blackboxfunc

import (
	"fmt"
	"github.com/reilabs/sunspot/go/acir/msgpackutil"
	shr "github.com/reilabs/sunspot/go/acir/shared"

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

func (a *And[T, E]) UnmarshalReader(r *msgpackutil.Reader) error {
	return msgpackutil.ReadStruct(r, "And", []msgpackutil.Field{
		{Name: "lhs", Decode: a.Lhs.UnmarshalReader},
		{Name: "rhs", Decode: a.Rhs.UnmarshalReader},
		{Name: "num_bits", Decode: func(r *msgpackutil.Reader) error {
			n, err := r.ReadU32()
			if err != nil {
				return err
			}
			if n > 128 {
				return fmt.Errorf("AND: num_bits=%d exceeds supported maximum of 128", n)
			}
			a.nBits = n
			return nil
		}},
		{Name: "output", Decode: a.Output.UnmarshalReader},
	})
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

func (*And[T, E]) SerdeName() string { return "AND" }
