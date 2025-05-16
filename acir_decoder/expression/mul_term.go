package expression

import (
	"io"
	shr "nr-groth16/acir_decoder/shared"
)

type MulTerm[T shr.ACIRField] struct {
	Term         shr.ACIRField
	WitnessLeft  shr.Witness
	WitnessRight shr.Witness
}

func (mt *MulTerm[T]) UnmarshalReader(r io.Reader) error {
	if err := mt.Term.UnmarshalReader(r); err != nil {
		return err
	}

	if err := mt.WitnessLeft.UnmarshalReader(r); err != nil {
		return err
	}

	if err := mt.WitnessRight.UnmarshalReader(r); err != nil {
		return err
	}

	return nil
}
