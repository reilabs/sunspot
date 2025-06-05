package blackboxfunc

import (
	"io"
	shr "nr-groth16/acir/shared"
)

type Xor[T shr.ACIRField] struct {
	Lhs    FunctionInput[T]
	Rhs    FunctionInput[T]
	Output shr.Witness
}

func (a *Xor[T]) UnmarshalReader(r io.Reader) error {
	if err := a.Lhs.UnmarshalReader(r); err != nil {
		return err
	}
	if err := a.Rhs.UnmarshalReader(r); err != nil {
		return err
	}
	if err := a.Output.UnmarshalReader(r); err != nil {
		return err
	}
	return nil
}

func (a *Xor[T]) Equals(other *Xor[T]) bool {
	if !a.Lhs.Equals(&other.Lhs) || !a.Rhs.Equals(&other.Rhs) {
		return false
	}
	return a.Output == other.Output
}
