package expression

import (
	"io"
	shr "nr-groth16/acir_decoder/shared"
)

type LinearCombination[T shr.ACIRField] struct {
	Term    T
	Witness shr.Witness
}

func (lc *LinearCombination[T]) UnmarshalReader(r io.Reader) error {
	if err := lc.Term.UnmarshalReader(r); err != nil {
		return err
	}

	if err := lc.Witness.UnmarshalReader(r); err != nil {
		return err
	}

	return nil
}
