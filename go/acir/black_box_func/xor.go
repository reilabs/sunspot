package blackboxfunc

import (
	"encoding/binary"
	"fmt"
	"io"
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

func (a *Xor[T, E]) UnmarshalReader(r io.Reader) error {
	if err := a.Lhs.UnmarshalReader(r); err != nil {
		return err
	}
	if err := a.Rhs.UnmarshalReader(r); err != nil {
		return err
	}
	if err := binary.Read(r, binary.LittleEndian, &a.nBits); err != nil {
		return err
	}
	if a.nBits > 128 {
		panic(fmt.Sprintf("XOR: nBits=%d exceeds supported maximum of 128", a.nBits))
	}
	if err := a.Output.UnmarshalReader(r); err != nil {
		return err
	}
	return nil
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

