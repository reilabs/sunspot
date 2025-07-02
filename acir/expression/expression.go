package expression

import (
	"encoding/binary"
	"fmt"
	"io"
	shr "nr-groth16/acir/shared"

	"github.com/consensys/gnark/frontend"
	"github.com/google/btree"
	"github.com/rs/zerolog/log"
)

type Expression[T shr.ACIRField] struct {
	MulTerms           []MulTerm[T]           `json:"mul_terms"`           // Terms that are multiplied together
	LinearCombinations []LinearCombination[T] `json:"linear_combinations"` // Linear combinations of variables
	Constant           T                      `json:"constant"`
	constantWitnessID  shr.Witness            // Constant term in the expression
}

func (e *Expression[T]) UnmarshalReader(r io.Reader) error {
	e.Constant = shr.MakeNonNil(e.Constant) // Ensure Constant is non-nil

	// Read the number of MulTerms
	var numMulTerms uint64
	if err := binary.Read(r, binary.LittleEndian, &numMulTerms); err != nil {
		return err
	}

	log.Trace().Msg("Unmarshalling Expression with " + fmt.Sprint(numMulTerms) + " MulTerms")
	// Initialize the MulTerms slice with the read size
	e.MulTerms = make([]MulTerm[T], numMulTerms)
	// Unmarshal each MulTerm
	for i := uint64(0); i < numMulTerms; i++ {
		log.Trace().Msg("Unmarshalling MulTerm at index: " + fmt.Sprint(i))
		if err := e.MulTerms[i].UnmarshalReader(r); err != nil {
			return err
		}
	}

	// Read the number of LinearCombinations
	var numLinearCombinations uint64
	if err := binary.Read(r, binary.LittleEndian, &numLinearCombinations); err != nil {
		return err
	}

	log.Trace().Msg("Unmarshalling Expression with " + fmt.Sprint(numLinearCombinations) + " LinearCombinations")
	// Initialize the LinearCombinations slice with the read size
	e.LinearCombinations = make([]LinearCombination[T], numLinearCombinations)

	// Unmarshal each LinearCombination
	for i := uint64(0); i < numLinearCombinations; i++ {
		log.Trace().Msg("Unmarshalling LinearCombination at index: " + fmt.Sprint(i))
		if err := e.LinearCombinations[i].UnmarshalReader(r); err != nil {
			return err
		}
	}

	log.Trace().Msg("Unmarshalling Expression with Constant value")
	// Unmarshal the Constant value
	if err := e.Constant.UnmarshalReader(r); err != nil {
		return err
	}

	return nil
}

func (e *Expression[T]) Equals(other *Expression[T]) bool {
	if len(e.MulTerms) != len(other.MulTerms) {
		return false
	}
	for i := range e.MulTerms {
		if !e.MulTerms[i].Equals(&other.MulTerms[i]) {
			return false
		}
	}

	if len(e.LinearCombinations) != len(other.LinearCombinations) {
		return false
	}
	for i := range e.LinearCombinations {
		if !e.LinearCombinations[i].Equals(&other.LinearCombinations[i]) {
			return false
		}
	}

	return e.Constant.Equals(other.Constant)
}

func (e *Expression[T]) Calculate(api frontend.API, witnesses map[shr.Witness]frontend.Variable) frontend.Variable {
	sum := e.Constant.ToFrontendVariable()
	log.Trace().Msg("EXPRESSION: Calculating Expression with " + fmt.Sprint(len(e.MulTerms)) + " MulTerms and " + fmt.Sprint(len(e.LinearCombinations)) + " LinearCombinations")
	for _, term := range e.MulTerms {
		sum = api.Add(sum, term.Calculate(api, witnesses))
	}
	for _, lc := range e.LinearCombinations {
		sum = api.Add(sum, lc.Calculate(api, witnesses))
	}

	log.Trace().Msg("EXPRESSION: Sum after all MulTerms and LinearCombinations: " + fmt.Sprint(sum))
	log.Trace().Msg("EXPRESSION: Adding constant to sum: " + fmt.Sprint(e.Constant.ToBigInt().Uint64()))
	sum = api.Add(sum, witnesses[e.constantWitnessID])
	log.Trace().Msg("EXPRESSION: Final sum after all MulTerms and LinearCombinations and Constant: " + fmt.Sprint(sum))
	return sum
}

func (e *Expression[T]) FillWitnessTree(tree *btree.BTree) bool {
	if tree == nil {
		return false
	}

	for _, term := range e.MulTerms {
		if !term.FillWitnessTree(tree) {
			return false
		}
	}

	for _, lc := range e.LinearCombinations {
		if !lc.FillWitnessTree(tree) {
			return false
		}
	}

	return true
}

func (e *Expression[T]) CollectConstantsAsWitnesses(start uint32, tree *btree.BTree) bool {
	if tree == nil {
		return false
	}

	e.constantWitnessID = shr.Witness(start + uint32(tree.Len()))
	tree.ReplaceOrInsert(e.constantWitnessID)
	log.Trace().Msgf("Collecting constant %s as witness with ID %d", e.Constant.ToBigInt().String(), e.constantWitnessID)

	return true
}
