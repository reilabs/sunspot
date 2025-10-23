package memory_op

import (
	"os"
	exp "sunpot/acir/expression"
	shr "sunpot/acir/shared"
	"sunpot/bn254"
	"testing"

	"github.com/consensys/gnark/constraint"
)

func TestMemoryOpWithoutPredicate(t *testing.T) {
	type E = constraint.U64
	type T = *bn254.BN254Field
	file, err := os.Open("../../binaries/opcodes/memory_op/memory_op_without_predicate.bin")
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}

	kind := shr.ParseThrough32bits(t, file)
	if kind != 2 {
		t.Fatal("Failed: mem op code should be 2")
	}

	var opcode MemoryOp[T, E]

	if err := opcode.UnmarshalReader(file); err != nil {
		t.Fatalf("Failed to unmarshal memory operation: %v", err)
	}

	expectedOpcode := MemoryOp[T, E]{
		BlockID: 0,
		Operation: exp.Expression[T, E]{
			MulTerms:           []exp.MulTerm[*bn254.BN254Field]{},
			LinearCombinations: []exp.LinearCombination[*bn254.BN254Field]{},
			Constant:           bn254.Zero(),
		},
		Index: exp.Expression[T, E]{
			MulTerms:           []exp.MulTerm[T]{},
			LinearCombinations: []exp.LinearCombination[T]{},
			Constant:           bn254.Zero(),
		},
		Value: exp.Expression[T, E]{
			MulTerms:           []exp.MulTerm[T]{},
			LinearCombinations: []exp.LinearCombination[T]{},
			Constant:           bn254.Zero(),
		},
		Predicate: nil,
	}

	if !opcode.Equals(&expectedOpcode) {
		t.Errorf("Expected opcode to be %v, got %v", expectedOpcode, opcode)
	}

	defer file.Close()
}

func TestMemoryOpWithPredicate(t *testing.T) {
	type E = constraint.U64
	type T = *bn254.BN254Field
	file, err := os.Open("../../binaries/opcodes/memory_op/memory_op_with_predicate.bin")
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}

	kind := shr.ParseThrough32bits(t, file)
	if kind != 2 {
		t.Fatal("Failed: mem op code should be 2")
	}

	var opcode MemoryOp[T, E]
	if err := opcode.UnmarshalReader(file); err != nil {
		t.Fatalf("Failed to unmarshal memory operation: %v", err)
	}

	expectedOpcode := MemoryOp[T, E]{
		BlockID: 1,
		Operation: exp.Expression[T, E]{
			MulTerms:           []exp.MulTerm[T]{},
			LinearCombinations: []exp.LinearCombination[T]{},
			Constant:           bn254.Zero(),
		},
		Index: exp.Expression[T, E]{
			MulTerms:           []exp.MulTerm[T]{},
			LinearCombinations: []exp.LinearCombination[T]{},
			Constant:           bn254.Zero(),
		},
		Value: exp.Expression[T, E]{
			MulTerms:           []exp.MulTerm[T]{},
			LinearCombinations: []exp.LinearCombination[T]{},
			Constant:           bn254.Zero(),
		},
		Predicate: &exp.Expression[T, E]{
			MulTerms:           []exp.MulTerm[T]{},
			LinearCombinations: []exp.LinearCombination[T]{},
			Constant:           bn254.Zero(),
		},
	}

	if !opcode.Equals(&expectedOpcode) {
		t.Errorf("Expected opcode to be %v, got %v", expectedOpcode, opcode)
	}

	defer file.Close()
}
