package blackboxfunc

import (
	"os"
	shr "sunspot/go/acir/shared"
	"sunspot/go/bn254"
	"testing"

	"github.com/consensys/gnark/constraint"
)

func TestBlake3UnmarshalReaderEmpty(t *testing.T) {
	type T = *bn254.BN254Field
	type E = constraint.U64

	file, err := os.Open("../../binaries/black_box_func/blake3/blake3_test_empty.bin")
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}

	// read the encoded black box function type before reading the actual content
	kind := shr.ParseThrough32bits(t, file)
	if kind != 5 {
		t.Fatalf("The kind of error code should have been 5, was %d", kind)
	}
	blackBoxFuncCall := BlackBoxFuncCall[T, E]{function: &Blake3[T, E]{}}

	if err := blackBoxFuncCall.UnmarshalReader(file); err != nil {
		t.Fatalf("Failed to unmarshal BlackBoxFuncCall: %v", err)
	}

	expectedFunctionCall := &Blake3[T, E]{
		Inputs:  []FunctionInput[T]{},
		Outputs: [32]shr.Witness{},
	}

	for i := 0; i < 32; i++ {
		expectedFunctionCall.Outputs[i] = shr.Witness(0)
	}

	if !blackBoxFuncCall.Equals(BlackBoxFuncCall[T, E]{function: expectedFunctionCall}) {
		t.Errorf("Expected BlackBoxFuncCall to be %v, got %v", expectedFunctionCall, blackBoxFuncCall)
	}

	defer file.Close()
}

func TestBlake3UnmarshalReaderWithInputs(t *testing.T) {
	type T = *bn254.BN254Field
	type E = constraint.U64
	file, err := os.Open("../../binaries/black_box_func/blake3/blake3_test_with_inputs.bin")
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}

	// read the encoded black box function type before reading the actual content
	kind := shr.ParseThrough32bits(t, file)
	if kind != 5 {
		t.Fatalf("The kind of error code should have been 5, was %d", kind)
	}
	blackBoxFuncCall := BlackBoxFuncCall[T, E]{function: &Blake3[T, E]{}}
	if err := blackBoxFuncCall.UnmarshalReader(file); err != nil {
		t.Fatalf("Failed to unmarshal BlackBoxFuncCall: %v", err)
	}

	expectedWitness1 := shr.Witness(1234)
	expectedWitness2 := shr.Witness(5678)
	expectedFunctionCall := &Blake3[T, E]{
		Inputs: []FunctionInput[T]{
			{
				FunctionInputKind: ACIRFunctionInputKindWitness,
				Witness:           &expectedWitness1,
			},
			{
				FunctionInputKind: ACIRFunctionInputKindWitness,
				Witness:           &expectedWitness2,
			},
		},
		Outputs: [32]shr.Witness{},
	}
	for i := 0; i < 32; i++ {
		expectedFunctionCall.Outputs[i] = shr.Witness(1234)
	}

	if !blackBoxFuncCall.Equals(BlackBoxFuncCall[T, E]{function: expectedFunctionCall}) {
		t.Errorf("Expected BlackBoxFuncCall to be %v, got %v", expectedFunctionCall, blackBoxFuncCall)
	}

	defer file.Close()
}
