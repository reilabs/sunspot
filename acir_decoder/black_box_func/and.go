package blackboxfunc

import (
	"io"
	shr "nr-groth16/acir_decoder/shared"
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
