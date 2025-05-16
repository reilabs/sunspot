package expression

import (
	"encoding/binary"
	"io"
	shr "nr-groth16/acir_decoder/shared"
)

type Expression[T shr.ACIRField] struct {
	MulTerms           []MulTerm[T]
	LinearCombinations []LinearCombination[T]
	Constant           T
}

func (e *Expression[T]) UnmarshalReader(r io.Reader) error {
	// Read the number of MulTerms
	var numMulTerms uint32
	if err := binary.Read(r, binary.LittleEndian, &numMulTerms); err != nil {
		return err
	}

	// Initialize the MulTerms slice with the read size
	e.MulTerms = make([]MulTerm[T], numMulTerms)
	// Unmarshal each MulTerm
	for i := uint32(0); i < numMulTerms; i++ {
		if err := e.MulTerms[i].UnmarshalReader(r); err != nil {
			return err
		}
	}

	// Read the number of LinearCombinations
	var numLinearCombinations uint32
	if err := binary.Read(r, binary.LittleEndian, &numLinearCombinations); err != nil {
		return err
	}

	// Initialize the LinearCombinations slice with the read size
	e.LinearCombinations = make([]LinearCombination[T], numLinearCombinations)

	// Unmarshal each LinearCombination
	for i := uint32(0); i < numLinearCombinations; i++ {
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
