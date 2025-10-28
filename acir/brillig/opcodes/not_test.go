package opcodes

import (
	"os"
	mem "sunspot/acir/brillig/memory"
	"sunspot/bn254"
	"testing"
)

func TestNotUnmarshalReader(t *testing.T) {
	file, err := os.Open("../../../binaries/brillig/opcodes/not/not.bin")
	if err != nil {
		t.Fatalf("failed to open file: %v", err)
	}

	op := BrilligOpcode[*bn254.BN254Field]{}
	if err := op.UnmarshalReader(file); err != nil {
		t.Fatalf("failed to unmarshal Not: %v", err)
	}

	expected := BrilligOpcode[*bn254.BN254Field]{
		OpCode: ACIRBrilligOpcodeNot,
		Not: &Not{
			Destination: mem.MemoryAddress{
				Kind:  mem.MemoryAddressKindDirect,
				Value: 1234,
			},
			Source: mem.MemoryAddress{
				Kind:  mem.MemoryAddressKindRelative,
				Value: 5678,
			},
			BitSize: mem.IntegerBitSizeU32,
		},
	}

	if !op.Equals(expected) {
		t.Errorf("expected %v, got %v", expected, op)
	}

	defer file.Close()
}
