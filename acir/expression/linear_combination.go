package expression

import (
	"io"
	shr "nr-groth16/acir/shared"

	"github.com/consensys/gnark/frontend"
	"github.com/google/btree"
)

type LinearCombination[T shr.ACIRField] struct {
	Term    T           `json:"term"`    // The term that is multiplied with the witness
	Witness shr.Witness `json:"witness"` // Witness for the linear combination
}

func (lc *LinearCombination[T]) UnmarshalReader(r io.Reader) error {
	lc.Term = shr.MakeNonNil(lc.Term) // Ensure Term is non-nil

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

	left, ok := witnesses[lc.Witness]
	if !ok {
		witnesses[lc.Witness] = api.Compiler().InternalVariable(uint32(lc.Witness))
		left = witnesses[lc.Witness]
	}
	return api.Mul(left, lc.Term.ToFrontendVariable())
}

func (lc *LinearCombination[T]) FillWitnessTree(tree *btree.BTree) bool {
	if tree == nil {
		return false
	}
	tree.ReplaceOrInsert(lc.Witness)
	return true
}
