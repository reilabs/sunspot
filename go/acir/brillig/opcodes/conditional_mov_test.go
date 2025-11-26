package opcodes

import (
	"os"
	mem "sunspot/acir/brillig/memory"
	"sunspot/bn254"
	"testing"
)

func TestConditionalMovUnmarshalReader(t *testing.T) {
	file, err := os.Open("../../../binaries/brillig/opcodes/conditional_mov/conditional_mov.bin")
	if err != nil {
		t.Fatalf("failed to open file: %v", err)
	}

	op := BrilligOpcode[*bn254.BN254Field]{}
	if err := op.UnmarshalReader(file); err != nil {
		t.Fatalf("failed to unmarshal ConditionalMov: %v", err)
	}

	expected := BrilligOpcode[*bn254.BN254Field]{
		OpCode: ACIRBrilligOpcodeConditionalMov,
		ConditionalMov: &ConditionalMov{
			Destination: mem.MemoryAddress{
				Kind:  mem.MemoryAddressKindDirect,
				Value: 1234,
			},
			SourceA: mem.MemoryAddress{
				Kind:  mem.MemoryAddressKindDirect,
				Value: 5678,
			},
			SourceB: mem.MemoryAddress{
				Kind:  mem.MemoryAddressKindRelative,
				Value: 91011,
			},
			Condition: mem.MemoryAddress{
				Kind:  mem.MemoryAddressKindDirect,
				Value: 121314,
			},
		},
	}

	if !op.Equals(expected) {
		t.Errorf("expected %v, got %v", expected, op)
	}

	defer file.Close()
}
