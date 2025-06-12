package expression

import (
	"encoding/binary"
	"io"
	shr "nr-groth16/acir/shared"

	"github.com/consensys/gnark/frontend"
)

type Expression[T shr.ACIRField] struct {
	MulTerms           []MulTerm[T]
	LinearCombinations []LinearCombination[T]
	Constant           T
}

func (e *Expression[T]) UnmarshalReader(r io.Reader) error {
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
	for _, term := range e.MulTerms {
		sum = api.Add(sum, term.Calculate(api, witnesses))
	}
	for _, lc := range e.LinearCombinations {
		sum = api.Add(sum, lc.Calculate(api, witnesses))
	}
	return api.Add(sum, e.Constant.ToFrontendVariable())
}
