package expression

import (
	"fmt"
	"github.com/reilabs/sunspot/go/acir/msgpackutil"
	shr "github.com/reilabs/sunspot/go/acir/shared"

	"github.com/consensys/gnark/frontend"
)

type MulTerm[T shr.ACIRField] struct {
	Term         T           `json:"term"`          // The term that is multiplied with the witnesses
	WitnessLeft  shr.Witness `json:"witness_left"`  // Left witness for multiplication
	WitnessRight shr.Witness `json:"witness_right"` // Right witness for multiplication
}

// On the wire each mul_term is a serde tuple `(F, Witness, Witness)` —
// always a 3-element fixarray, no tagged-map shape.
func (mt *MulTerm[T]) UnmarshalReader(r *msgpackutil.Reader) error {
	n, err := r.ReadArrayLen()
	if err != nil {
		return err
	}
	if n != 3 {
		return fmt.Errorf("mul_term: expected 3-tuple, got %d elements", n)
	}
	mt.Term = shr.MakeNonNil(mt.Term)
	if err := mt.Term.UnmarshalReader(r); err != nil {
		return err
	}
	if err := mt.WitnessLeft.UnmarshalReader(r); err != nil {
		return err
	}
	return mt.WitnessRight.UnmarshalReader(r)
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
