package blackboxfunc

import (
	"io"
	shr "nr-groth16/acir/shared"
)

type And[T shr.ACIRField] struct {
	Lhs    FunctionInput[T]
	Rhs    FunctionInput[T]
	Output shr.Witness
}

func (a *And[T]) UnmarshalReader(r io.Reader) error {
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

func (a *And[T]) Equals(other *And[T]) bool {
	return a.Lhs.Equals(&other.Lhs) && a.Rhs.Equals(&other.Rhs) && a.Output.Equals(&other.Output)
}
