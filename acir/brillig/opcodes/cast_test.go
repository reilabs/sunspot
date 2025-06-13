package opcodes

import (
	mem "nr-groth16/acir/brillig/memory"
	"nr-groth16/bn254"
	"os"
	"testing"
)

func TestCastUnmarshalReader(t *testing.T) {
	file, err := os.Open("../../../binaries/brillig/opcodes/cast/cast.bin")
	if err != nil {
		t.Fatalf("failed to open file: %v", err)
	}

	op := BrilligOpcode[bn254.BN254Field]{}
	if err := op.UnmarshalReader(file); err != nil {
		t.Fatalf("failed to unmarshal Cast: %v", err)
	}

	expectedIntegerBitSize := mem.IntegerBitSizeU32
	expected := BrilligOpcode[bn254.BN254Field]{
		OpCode: ACIRBrilligOpcodeCast,
		Cast: &Cast{
			Destination: mem.MemoryAddress{
				Kind:  mem.MemoryAddressKindDirect,
				Value: 1234,
			},
			Source: mem.MemoryAddress{
				Kind:  mem.MemoryAddressKindDirect,
				Value: 5678,
			},
			BitSize: mem.BitSize{
				Kind:           mem.BitSizeKindInteger,
				IntegerBitSize: &expectedIntegerBitSize,
			},
		},
	}

	if !op.Equals(expected) {
		t.Errorf("expected %v, got %v", expected, op)
	}

	defer file.Close()
}
