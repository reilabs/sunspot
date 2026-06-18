package expression

import (
	"fmt"
	"sunspot/go/acir/msgpackutil"
	shr "sunspot/go/acir/shared"

	"github.com/consensys/gnark/frontend"
)

type LinearCombination[T shr.ACIRField] struct {
	Term    T           `json:"term"`    // The term that is multiplied with the witness
	Witness shr.Witness `json:"witness"` // Witness for the linear combination
}

// On the wire each linear_combination is a serde tuple `(F, Witness)` —
// a 2-element fixarray.
func (lc *LinearCombination[T]) UnmarshalReader(r *msgpackutil.Reader) error {
	n, err := r.ReadArrayLen()
	if err != nil {
		return err
	}
	if n != 2 {
		return fmt.Errorf("linear_combination: expected 2-tuple, got %d elements", n)
	}
	lc.Term = shr.MakeNonNil(lc.Term)
	if err := lc.Term.UnmarshalReader(r); err != nil {
		return err
	}
	return lc.Witness.UnmarshalReader(r)
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
