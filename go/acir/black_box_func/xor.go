package blackboxfunc

import (
	"fmt"
	"sunspot/go/acir/msgpackutil"
	shr "sunspot/go/acir/shared"

	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/math/uints"
)

type Xor[T shr.ACIRField, E constraint.Element] struct {
	Lhs    FunctionInput[T]
	Rhs    FunctionInput[T]
	Output shr.Witness
	nBits  uint32
}

func (a *Xor[T, E]) decode(f msgpackutil.Field, r *msgpackutil.Reader) error {
	switch f.Tag {
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
			return fmt.Errorf("XOR: num_bits=%d exceeds supported maximum of 128", n)
		}
		a.nBits = n
		return nil
	case 3:
		return a.Output.UnmarshalReader(r)
	default:
		return fmt.Errorf("Xor: unknown field %s", f)
	}
}

func (a *Xor[T, E]) Equals(other BlackBoxFunction[E]) bool {
	value, ok := other.(*Xor[T, E])

	if !ok || !a.Lhs.Equals(&value.Lhs) || !a.Rhs.Equals(&value.Rhs) || a.nBits != value.nBits {
		return false
	}
	return a.Output == value.Output
}

func (a *Xor[T, E]) Define(api frontend.Builder[E], witnesses map[shr.Witness]frontend.Variable) error {
	uapi, err := uints.New[uints.U64](api)
	if err != nil {
		return err
	}
	return defineBitwise(api, uapi, witnesses, a.Lhs, a.Rhs, a.Output, int(a.nBits), uapi.Xor)
}

func (*Xor[T, E]) SerdeName() string { return "XOR" }
