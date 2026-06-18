package memory_op

import (
	"os"
	"sunspot/go/acir/msgpackutil"
	shr "sunspot/go/acir/shared"
	"sunspot/go/bn254"
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
	defer file.Close()

	r := msgpackutil.NewReader(file)
	if tag := shr.ConsumeEnumTag(t, r); tag != 2 {
		t.Fatalf("expected Opcode variant 2 (MemoryOp), got %d", tag)
	}

	var opcode MemoryOp[T, E]
	if err := opcode.UnmarshalReader(r); err != nil {
		t.Fatalf("Failed to unmarshal memory operation: %v", err)
	}

	expectedOpcode := MemoryOp[T, E]{
		BlockID: 0,
		IsWrite: false,
		Index:   shr.Witness(2),
		Value:   shr.Witness(3),
	}
	if !opcode.Equals(&expectedOpcode) {
		t.Errorf("Expected opcode to be %+v, got %+v", expectedOpcode, opcode)
	}
}

func TestMemoryOpWithPredicate(t *testing.T) {
	type E = constraint.U64
	type T = *bn254.BN254Field
	file, err := os.Open("../../binaries/opcodes/memory_op/memory_op_with_predicate.bin")
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	r := msgpackutil.NewReader(file)
	if tag := shr.ConsumeEnumTag(t, r); tag != 2 {
		t.Fatalf("expected Opcode variant 2 (MemoryOp), got %d", tag)
	}

	var opcode MemoryOp[T, E]
	if err := opcode.UnmarshalReader(r); err != nil {
		t.Fatalf("Failed to unmarshal memory operation: %v", err)
	}

	expectedOpcode := MemoryOp[T, E]{
		BlockID: 1,
		IsWrite: true,
		Index:   shr.Witness(5),
		Value:   shr.Witness(6),
	}
	if !opcode.Equals(&expectedOpcode) {
		t.Errorf("Expected opcode to be %+v, got %+v", expectedOpcode, opcode)
	}
}
