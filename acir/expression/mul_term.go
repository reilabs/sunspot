package expression

import (
	"io"
	shr "nr-groth16/acir/shared"

	"github.com/consensys/gnark/frontend"
	"github.com/google/btree"
)

type MulTerm[T shr.ACIRField] struct {
	Term         T           `json:"term"`          // The term that is multiplied with the witnesses
	WitnessLeft  shr.Witness `json:"witness_left"`  // Left witness for multiplication
	WitnessRight shr.Witness `json:"witness_right"` // Right witness for multiplication
}

func (mt *MulTerm[T]) UnmarshalReader(r io.Reader) error {
	mt.Term = shr.MakeNonNil(mt.Term) // Ensure Term is non-nil

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
	left, ok := witnesses[Mt.WitnessLeft]
	if !ok {
		witnesses[Mt.WitnessLeft] = api.Compiler().InternalVariable(uint32(Mt.WitnessLeft))
		left = witnesses[Mt.WitnessLeft]
	}
	right, ok := witnesses[Mt.WitnessRight]
	if !ok {
		witnesses[Mt.WitnessRight] = api.Compiler().InternalVariable(uint32(Mt.WitnessRight))
		right = witnesses[Mt.WitnessRight]
	}

	return api.Mul(left, right, Mt.Term.ToFrontendVariable())
}

func (mt *MulTerm[T]) FillWitnessTree(tree *btree.BTree) bool {
	if tree == nil {
		return false
	}

	tree.ReplaceOrInsert(mt.WitnessLeft)
	tree.ReplaceOrInsert(mt.WitnessRight)

	return true
}
