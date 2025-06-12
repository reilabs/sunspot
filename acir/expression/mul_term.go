package expression

import (
	"io"
	shr "nr-groth16/acir/shared"

	"github.com/consensys/gnark/frontend"
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

func (Mt *MulTerm[T]) Calculate(api frontend.API, witnesses map[shr.Witness]frontend.Variable) frontend.Variable {
	left := witnesses[Mt.WitnessLeft]
	right := witnesses[Mt.WitnessRight]
	api.Println("Calculating MulTerm with left witness:", left, "and right witness:", right, "with term:", Mt.Term.ToFrontendVariable())
	return api.Mul(left, right, Mt.Term.ToFrontendVariable())
}
