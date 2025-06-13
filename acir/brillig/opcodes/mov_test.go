package opcodes

import (
	mem "nr-groth16/acir/brillig/memory"
	"nr-groth16/bn254"
	"os"
	"testing"
)

func TestMovUnmarshalReader(t *testing.T) {
	file, err := os.Open("../../../binaries/brillig/opcodes/mov/mov.bin")
	if err != nil {
		t.Fatalf("failed to open file: %v", err)
	}

	op := BrilligOpcode[bn254.BN254Field]{}
	if err := op.UnmarshalReader(file); err != nil {
		t.Fatalf("failed to unmarshal Mov: %v", err)
	}

	expected := BrilligOpcode[bn254.BN254Field]{
		OpCode: ACIRBrilligOpcodeMov,
		Mov: &Mov{
			Destination: mem.MemoryAddress{
				Kind:  mem.MemoryAddressKindDirect,
				Value: 1234,
			},
			Source: mem.MemoryAddress{
				Kind:  mem.MemoryAddressKindRelative,
				Value: 5678,
			},
		},
	}

	if !op.Equals(expected) {
		t.Errorf("expected %v, got %v", expected, op)
	}

	defer file.Close()
}
