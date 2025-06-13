package opcodes

import (
	mem "nr-groth16/acir/brillig/memory"
	"nr-groth16/bn254"
	"os"
	"testing"
)

func TestForeignCallUnmarshalEmpty(t *testing.T) {
	file, err := os.Open("../../../binaries/brillig/opcodes/foreign_call/foreign_call_empty.bin")
	if err != nil {
		t.Fatalf("failed to open file: %v", err)
	}

	op := BrilligOpcode[bn254.BN254Field]{}
	if err := op.UnmarshalReader(file); err != nil {
		t.Fatalf("failed to unmarshal ForeignCall: %v", err)
	}

	expected := BrilligOpcode[bn254.BN254Field]{
		OpCode: ACIRBrilligOpcodeForeignCall,
		ForeignCall: &ForeignCall{
			Function:              "example_function",
			Destinations:          []mem.ValueOrArray{},
			DestinationValueTypes: []mem.HeapValueType{},
			Inputs:                []mem.ValueOrArray{},
			InputValueTypes:       []mem.HeapValueType{},
		},
	}

	if !op.Equals(expected) {
		t.Errorf("expected %v, got %v", expected, op)
	}

	defer file.Close()
}

func TestForeignCallUnmarshalWithInputs(t *testing.T) {
	file, err := os.Open("../../../binaries/brillig/opcodes/foreign_call/foreign_call_with_inputs.bin")
	if err != nil {
		t.Fatalf("failed to open file: %v", err)
	}

	op := BrilligOpcode[bn254.BN254Field]{}
	if err := op.UnmarshalReader(file); err != nil {
		t.Fatalf("failed to unmarshal ForeignCall: %v", err)
	}

	expectedIntegerBitSize := mem.IntegerBitSizeU32
	expected := BrilligOpcode[bn254.BN254Field]{
		OpCode: ACIRBrilligOpcodeForeignCall,
		ForeignCall: &ForeignCall{
			Function: "example_function",
			Destinations: []mem.ValueOrArray{
				mem.ValueOrArray{
					Kind: mem.ACIRBrilligValueOrArrayKindMemoryAddress,
					MemoryAddress: &mem.MemoryAddress{
						Kind:  mem.MemoryAddressKindDirect,
						Value: 1234,
					},
				},
				mem.ValueOrArray{
					Kind: mem.ACIRBrilligValueOrArrayKindHeapArray,
					HeapArray: &mem.HeapArray{
						Pointer: mem.MemoryAddress{
							Kind:  mem.MemoryAddressKindDirect,
							Value: 5678,
						},
						Size: 10,
					},
				},
			},
			DestinationValueTypes: []mem.HeapValueType{
				{
					Kind: mem.ACIRBrilligHeapValueTypeKindSimple,
					Simple: &mem.BitSize{
						Kind: mem.BitSizeKindField,
					},
				},
			},
			Inputs: []mem.ValueOrArray{
				{
					Kind: mem.ACIRBrilligValueOrArrayKindHeapVector,
					HeapVector: &mem.HeapVector{
						Pointer: mem.MemoryAddress{
							Kind:  mem.MemoryAddressKindDirect,
							Value: 91011,
						},
						Size: mem.MemoryAddress{
							Kind:  mem.MemoryAddressKindRelative,
							Value: 1212,
						},
					},
				},
			},
			InputValueTypes: []mem.HeapValueType{
				{
					Kind: mem.ACIRBrilligHeapValueTypeKindSimple,
					Simple: &mem.BitSize{
						Kind:           mem.BitSizeKindInteger,
						IntegerBitSize: &expectedIntegerBitSize,
					},
				},
			},
		},
	}

	if !op.Equals(expected) {
		t.Errorf("expected %v, got %v", expected, op)
	}

	defer file.Close()
}
