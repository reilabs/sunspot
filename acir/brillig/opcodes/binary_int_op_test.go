package opcodes

import (
	mem "nr-groth16/acir/brillig/memory"
	"nr-groth16/bn254"
	"os"
	"testing"
)

func TestBinaryIntOpUnmarshalReaderAdd(t *testing.T) {
	f, err := os.Open("../../../binaries/brillig/opcodes/binary_int_op/add.bin")
	if err != nil {
		t.Fatalf("failed to open test file: %v", err)
	}

	op := BrilligOpcode[*bn254.BN254Field]{}
	if err := op.UnmarshalReader(f); err != nil {
		t.Fatalf("failed to unmarshal BinaryIntOpAdd: %v", err)
	}

	expected := BrilligOpcode[*bn254.BN254Field]{
		OpCode: ACIRBrilligOpcodeBinaryIntOp,
		BinaryIntOp: &BinaryIntOp{
			Destination: mem.MemoryAddress{
				Kind:  mem.MemoryAddressKindDirect,
				Value: 1234,
			},
			BitSize: mem.IntegerBitSizeU32,
			Op:      BinaryIntOpKindAdd,
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

	if !op.Equals(expected) {
		t.Errorf("expected %v, got %v", expected, op)
	}

	defer f.Close()
}

func TestBinaryIntOpUnmarshalReaderLessThanOrEquals(t *testing.T) {
	f, err := os.Open("../../../binaries/brillig/opcodes/binary_int_op/less_than_equals.bin")
	if err != nil {
		t.Fatalf("failed to open test file: %v", err)
	}

	op := BrilligOpcode[*bn254.BN254Field]{}
	if err := op.UnmarshalReader(f); err != nil {
		t.Fatalf("failed to unmarshal BinaryIntOpLessThanOrEquals: %v", err)
	}

	expected := BrilligOpcode[*bn254.BN254Field]{
		OpCode: ACIRBrilligOpcodeBinaryIntOp,
		BinaryIntOp: &BinaryIntOp{
			Destination: mem.MemoryAddress{
				Kind:  mem.MemoryAddressKindDirect,
				Value: 123456,
			},
			BitSize: mem.IntegerBitSizeU32,
			Op:      BinaryIntOpKindLessThanEquals,
			Lhs: mem.MemoryAddress{
				Kind:  mem.MemoryAddressKindDirect,
				Value: 5678910,
			},
			Rhs: mem.MemoryAddress{
				Kind:  mem.MemoryAddressKindRelative,
				Value: 9101112,
			},
		},
	}

	if !op.Equals(expected) {
		t.Errorf("expected %v, got %v", expected, op)
	}

	defer f.Close()
}
