package opcodes

import (
	"os"
	mem "sunpot/acir/brillig/memory"
	"sunpot/bn254"
	"testing"
)

func TestJumpIfNotUnmarshalReader(t *testing.T) {
	file, err := os.Open("../../../binaries/brillig/opcodes/jump_if_not/jump_if_not.bin")
	if err != nil {
		t.Fatalf("failed to open file: %v", err)
	}

	op := BrilligOpcode[*bn254.BN254Field]{}
	if err := op.UnmarshalReader(file); err != nil {
		t.Fatalf("failed to unmarshal JumpIfNot: %v", err)
	}

	expected := BrilligOpcode[*bn254.BN254Field]{
		OpCode: ACIRBrilligOpcodeJumpIfNot,
		JumpIfNot: &JumpIfNot{
			Condition: mem.MemoryAddress{
				Kind:  mem.MemoryAddressKindDirect,
				Value: 1234,
			},
			Location: 5678,
		},
	}

	if !op.Equals(expected) {
		t.Errorf("expected %v, got %v", expected, op)
	}

	defer file.Close()
}
