package opcodes

import (
	mem "nr-groth16/acir/brillig/memory"
	"nr-groth16/bn254"
	"os"
	"testing"
)

func TestCallDataCopyUnmarshalReader(t *testing.T) {
	file, err := os.Open("../../../binaries/brillig/opcodes/call_data_copy/call_data_copy.bin")
	if err != nil {
		t.Fatalf("failed to open file: %v", err)
	}

	op := BrilligOpcode[bn254.BN254Field]{}
	if err := op.UnmarshalReader(file); err != nil {
		t.Fatalf("failed to unmarshal CallDataCopy: %v", err)
	}

	expected := BrilligOpcode[bn254.BN254Field]{
		OpCode: ACIRBrilligOpcodeCalldataCopy,
		CalldataCopy: &CallDataCopy{
			DestinationAddress: mem.MemoryAddress{
				Kind:  mem.MemoryAddressKindDirect,
				Value: 1234,
			},
			SizeAddress: mem.MemoryAddress{
				Kind:  mem.MemoryAddressKindDirect,
				Value: 5678,
			},
			OffsetAddress: mem.MemoryAddress{
				Kind:  mem.MemoryAddressKindRelative,
				Value: 91011,
			},
		},
	}

	if !op.Equals(expected) {
		t.Errorf("expected %v, got %v", expected, op)
	}

	defer file.Close()
}
