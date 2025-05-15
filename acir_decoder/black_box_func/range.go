package blackboxfunc

import (
	"io"
	shr "nr-groth16/acir_decoder/shared"
)

type Range[T shr.ACIRField] struct {
	Input FunctionInput[T]
}

func (a *Range[T]) UnmarshalReader(r io.Reader) error {
	if err := a.Input.UnmarshalReader(r); err != nil {
		return err
	}
	return nil
}
