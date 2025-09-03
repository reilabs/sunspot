package expression

import (
	"fmt"
	"io"
	shr "nr-groth16/acir/shared"

	"github.com/consensys/gnark/frontend"
	"github.com/google/btree"
	"github.com/rs/zerolog/log"
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
	log.Trace().Msg("Unmarshalling  MulTerm with term: " + mt.Term.String())

	if err := mt.WitnessLeft.UnmarshalReader(r); err != nil {
		return err
	}
	log.Trace().Msg("Unmarshalling  MulTerm with left witness: " + fmt.Sprint(mt.WitnessLeft))

	if err := mt.WitnessRight.UnmarshalReader(r); err != nil {
		return err
	}
	log.Trace().Msg("Unmarshalling  MulTerm with right witness: " + fmt.Sprint(mt.WitnessRight))

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
		log.Trace().Msg("EXPRESSION: MULTERM: Left witness not found, creating internal variable for witness: " + fmt.Sprint(Mt.WitnessLeft))
	}
	right, ok := witnesses[Mt.WitnessRight]
	if !ok {
		witnesses[Mt.WitnessRight] = api.Compiler().InternalVariable(uint32(Mt.WitnessRight))
		right = witnesses[Mt.WitnessRight]
		log.Trace().Msg("EXPRESSION: MULTERM: Right witness not found, creating internal variable for witness: " + fmt.Sprint(Mt.WitnessRight))
	}
	log.Trace().Msg("EXPRESSION: MULTERM: Calculating MulTerm with left witness: " + fmt.Sprint(Mt.WitnessLeft) + " and right witness: " + fmt.Sprint(Mt.WitnessRight))
	//log.Trace().Msg("Witnesses: " + fmt.Sprint(witnesses))
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
