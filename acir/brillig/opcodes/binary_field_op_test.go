package opcodes

import (
	mem "nr-groth16/acir/brillig/memory"
	"nr-groth16/bn254"
	"os"
	"testing"
)

func TestBinaryFieldOpUnmarshalReaderAdd(t *testing.T) {
	file, err := os.Open("../../../binaries/brillig/opcodes/binary_field_op/add.bin")
	if err != nil {
		t.Fatalf("failed to open file: %v", err)
	}

	op := BrilligOpcode[bn254.BN254Field]{}
	err = op.UnmarshalReader(file)
	if err != nil {
		t.Fatalf("failed to unmarshal BinaryFieldOp: %v", err)
	}

	expectedOp := BrilligOpcode[bn254.BN254Field]{
		OpCode: ACIRBrilligOpcodeBinaryFieldOp,
		BinaryFieldOp: &BinaryFieldOp{
			Destination: mem.MemoryAddress{
				Kind:  mem.MemoryAddressKindDirect,
				Value: 1234,
			},
			Op: BinaryFieldOpAdd,
			Lhs: mem.MemoryAddress{
				Kind:  mem.MemoryAddressKindDirect,
				Value: 5678,
			},
			Rhs: mem.MemoryAddress{
				Kind:  mem.MemoryAddressKindRelative,
				Value: 91011,
			},
		},
	}

	if !op.Equals(expectedOp) {
		t.Fatalf("expected BinaryFieldOp to be %v, got %v", expectedOp, op)
	}

	defer file.Close()
}

func TestBinaryFieldOpUnmarshalReaderLessThanOrEquals(t *testing.T) {
	file, err := os.Open("../../../binaries/brillig/opcodes/binary_field_op/less_than_equals.bin")
	if err != nil {
		t.Fatalf("failed to open file: %v", err)
	}

	op := BrilligOpcode[bn254.BN254Field]{}
	err = op.UnmarshalReader(file)
	if err != nil {
		t.Fatalf("failed to unmarshal BinaryFieldOp: %v", err)
	}

	expectedOp := BrilligOpcode[bn254.BN254Field]{
		OpCode: ACIRBrilligOpcodeBinaryFieldOp,
		BinaryFieldOp: &BinaryFieldOp{
			Destination: mem.MemoryAddress{
				Kind:  mem.MemoryAddressKindDirect,
				Value: 1234,
			},
			Op: BinaryFieldOpLessThanEquals,
			Lhs: mem.MemoryAddress{
				Kind:  mem.MemoryAddressKindDirect,
				Value: 5678,
			},
			Rhs: mem.MemoryAddress{
				Kind:  mem.MemoryAddressKindRelative,
				Value: 91011,
			},
		},
	}

	if !op.Equals(expectedOp) {
		t.Fatalf("expected BinaryFieldOp to be %v, got %v", expectedOp, op)
	}

	defer file.Close()
}
