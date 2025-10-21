package expression

import (
	"encoding/binary"
	"encoding/json"
	"io"
	"nr-groth16/acir/opcodes"
	shr "nr-groth16/acir/shared"

	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/google/btree"
)

type Expression[T shr.ACIRField, E constraint.Element] struct {
	MulTerms           []MulTerm[T]           `json:"mul_terms"`           // Terms that are multiplied together
	LinearCombinations []LinearCombination[T] `json:"linear_combinations"` // Linear combinations of variables
	Constant           T                      `json:"constant"`
	constantWitnessID  shr.Witness            // Constant term in the expression
}

func (e *Expression[T, E]) Define(
	api frontend.Builder[E],
	witnesses map[shr.Witness]frontend.Variable,
) error {
	api.AssertIsEqual(e.Calculate(api, witnesses), 0)
	return nil
}

func (e *Expression[T, E]) UnmarshalReader(r io.Reader) error {
	e.Constant = shr.MakeNonNil(e.Constant) // Ensure Constant is non-nil

	// Read the number of MulTerms
	var numMulTerms uint64
	if err := binary.Read(r, binary.LittleEndian, &numMulTerms); err != nil {
		return err
	}

	// Initialize the MulTerms slice with the read size
	e.MulTerms = make([]MulTerm[T], numMulTerms)
	// Unmarshal each MulTerm
	for i := uint64(0); i < numMulTerms; i++ {
		if err := e.MulTerms[i].UnmarshalReader(r); err != nil {
			return err
		}
	}

	// Read the number of LinearCombinations
	var numLinearCombinations uint64
	if err := binary.Read(r, binary.LittleEndian, &numLinearCombinations); err != nil {
		return err
	}
	// Initialize the LinearCombinations slice with the read size
	e.LinearCombinations = make([]LinearCombination[T], numLinearCombinations)

	// Unmarshal each LinearCombination
	for i := uint64(0); i < numLinearCombinations; i++ {
		if err := e.LinearCombinations[i].UnmarshalReader(r); err != nil {
			return err
		}
	}

	// Unmarshal the Constant value
	if err := e.Constant.UnmarshalReader(r); err != nil {
		return err
	}

	return nil
}

func (e *Expression[T, E]) Equals(other opcodes.Opcode[E]) bool {
	value, ok := other.(*Expression[T, E])
	if !ok {
		return false
	}

	if len(e.MulTerms) != len(value.MulTerms) {
		return false
	}
	for i := range e.MulTerms {
		if !e.MulTerms[i].Equals(&value.MulTerms[i]) {
			return false
		}
	}

	if len(e.LinearCombinations) != len(value.LinearCombinations) {
		return false
	}
	for i := range e.LinearCombinations {
		if !e.LinearCombinations[i].Equals(&value.LinearCombinations[i]) {
			return false
		}
	}

	return e.Constant.Equals(value.Constant)
}

func (e *Expression[T, E]) Calculate(api frontend.API, witnesses map[shr.Witness]frontend.Variable) frontend.Variable {
	sum := e.Constant.ToFrontendVariable()
	for _, term := range e.MulTerms {
		sum = api.Add(sum, term.Calculate(api, witnesses))
	}
	for _, lc := range e.LinearCombinations {
		sum = api.Add(sum, lc.Calculate(api, witnesses))
	}

	return sum
}

func (e *Expression[T, E]) FillWitnessTree(tree *btree.BTree, index uint32) bool {
	if tree == nil {
		return false
	}

	for _, term := range e.MulTerms {
		if !term.FillWitnessTree(tree, index) {
			return false
		}
	}

	for _, lc := range e.LinearCombinations {
		if !lc.FillWitnessTree(tree, index) {
			return false
		}
	}

	return true
}

func (e *Expression[T, E]) MarshalJSON() ([]byte, error) {
	stringMap := make(map[string]interface{})
	stringMap["assert_zero"] = e
	return json.Marshal(stringMap)
}
