package expression

import (
	"io"
	shr "nr-groth16/acir/shared"

	"github.com/consensys/gnark/frontend"
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

func (lc *LinearCombination[T]) Equals(other *LinearCombination[T]) bool {
	if !lc.Term.Equals(other.Term) {
		return false
	}

	if !lc.Witness.Equals(&other.Witness) {
		return false
	}

	return true
}

func (lc *LinearCombination[T]) Calculate(api frontend.API, witnesses map[shr.Witness]frontend.Variable) frontend.Variable {
	left := witnesses[lc.Witness]
	return api.Mul(left, lc.Term.ToFrontendVariable())
}
