package expression

import (
	"io"
	shr "nr-groth16/acir/shared"
)

type MulTerm[T shr.ACIRField] struct {
	Term         T
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

func (mt *MulTerm[T]) Equals(other *MulTerm[T]) bool {
	if !mt.Term.Equals(other.Term) {
		return false
	}

	if !mt.WitnessLeft.Equals(&other.WitnessLeft) {
		return false
	}

	if !mt.WitnessRight.Equals(&other.WitnessRight) {
		return false
	}

	return true
}
