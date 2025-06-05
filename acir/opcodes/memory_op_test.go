package opcodes

import (
	exp "nr-groth16/acir/expression"
	"nr-groth16/bn254"
	"os"
	"testing"
)

func TestMemoryOpWithoutPredicate(t *testing.T) {
	file, err := os.Open("../../binaries/opcodes/memory_op/memory_op_without_predicate.bin")
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}

	var opcode Opcode[*bn254.BN254Field]
	if err := opcode.UnmarshalReader(file); err != nil {
		t.Fatalf("Failed to unmarshal memory operation: %v", err)
	}

	expectedOpcode := Opcode[*bn254.BN254Field]{
		Kind: ACIROpcodeMemoryOp,
		MemoryOp: &MemoryOp[*bn254.BN254Field]{
			BlockID: 0,
			Operation: exp.Expression[*bn254.BN254Field]{
				MulTerms:           []exp.MulTerm[*bn254.BN254Field]{},
				LinearCombinations: []exp.LinearCombination[*bn254.BN254Field]{},
				Constant:           &bn254.BN254Field{},
			},
			Index: exp.Expression[*bn254.BN254Field]{
				MulTerms:           []exp.MulTerm[*bn254.BN254Field]{},
				LinearCombinations: []exp.LinearCombination[*bn254.BN254Field]{},
				Constant:           &bn254.BN254Field{},
			},
			Value: exp.Expression[*bn254.BN254Field]{
				MulTerms:           []exp.MulTerm[*bn254.BN254Field]{},
				LinearCombinations: []exp.LinearCombination[*bn254.BN254Field]{},
				Constant:           &bn254.BN254Field{},
			},
			Predicate: nil,
		},
	}

	if !opcode.Equals(&expectedOpcode) {
		t.Errorf("Expected opcode to be %v, got %v", expectedOpcode, opcode)
	}

	defer file.Close()
}

func TestMemoryOpWithPredicate(t *testing.T) {
	file, err := os.Open("../../binaries/opcodes/memory_op/memory_op_with_predicate.bin")
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}

	var opcode Opcode[*bn254.BN254Field]
	if err := opcode.UnmarshalReader(file); err != nil {
		t.Fatalf("Failed to unmarshal memory operation: %v", err)
	}

	expectedOpcode := Opcode[*bn254.BN254Field]{
		Kind: ACIROpcodeMemoryOp,
		MemoryOp: &MemoryOp[*bn254.BN254Field]{
			BlockID: 1,
			Operation: exp.Expression[*bn254.BN254Field]{
				MulTerms:           []exp.MulTerm[*bn254.BN254Field]{},
				LinearCombinations: []exp.LinearCombination[*bn254.BN254Field]{},
				Constant:           &bn254.BN254Field{},
			},
			Index: exp.Expression[*bn254.BN254Field]{
				MulTerms:           []exp.MulTerm[*bn254.BN254Field]{},
				LinearCombinations: []exp.LinearCombination[*bn254.BN254Field]{},
				Constant:           &bn254.BN254Field{},
			},
			Value: exp.Expression[*bn254.BN254Field]{
				MulTerms:           []exp.MulTerm[*bn254.BN254Field]{},
				LinearCombinations: []exp.LinearCombination[*bn254.BN254Field]{},
				Constant:           &bn254.BN254Field{},
			},
			Predicate: &exp.Expression[*bn254.BN254Field]{
				MulTerms:           []exp.MulTerm[*bn254.BN254Field]{},
				LinearCombinations: []exp.LinearCombination[*bn254.BN254Field]{},
				Constant:           &bn254.BN254Field{},
			},
		},
	}

	if !opcode.Equals(&expectedOpcode) {
		t.Errorf("Expected opcode to be %v, got %v", expectedOpcode, opcode)
	}

	defer file.Close()
}
