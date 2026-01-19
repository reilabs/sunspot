package blackboxfunc

import (
	"os"
	shr "sunspot/go/acir/shared"
	"sunspot/go/bn254"
	"testing"

	"github.com/consensys/gnark/constraint"
)

func TestAES128EncryptUnmarshalReaderEmpty(t *testing.T) {
	type T = *bn254.BN254Field
	type E = constraint.U64
	file, err := os.Open("../../binaries/black_box_func/aes128encrypt/aes128encrypt_empty.bin")
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}

	kind := shr.ParseThrough32bits(t, file)
	if kind != 0 {
		t.Fatalf("The kind of error code should have been 0, was %d", kind)
	}

	blackBoxFuncCall := BlackBoxFuncCall[T, E]{function: &AES128Encrypt[T, E]{}}
	if err := blackBoxFuncCall.UnmarshalReader(file); err != nil {
		t.Fatalf("Failed to unmarshal BlackBoxFuncCall: %v", err)
	}

	expectedIvWitness := shr.Witness(1234)
	expectedKeyWitness := shr.Witness(5678)
	expectedIv := [16]FunctionInput[T]{}
	for i := 0; i < 16; i++ {
		expectedIv[i] = FunctionInput[T]{
			FunctionInputKind: ACIRFunctionInputKindWitness,
			Witness:           &expectedIvWitness,
		}
	}
	expectedKey := [16]FunctionInput[T]{}
	for i := 0; i < 16; i++ {
		expectedKey[i] = FunctionInput[T]{
			FunctionInputKind: ACIRFunctionInputKindWitness,
			Witness:           &expectedKeyWitness,
		}
	}

	expected := BlackBoxFuncCall[T, E]{
		function: &AES128Encrypt[T, E]{
			Inputs:  []FunctionInput[T]{},
			Iv:      expectedIv,
			Key:     expectedKey,
			Outputs: []shr.Witness{},
		},
	}

	if !blackBoxFuncCall.Equals(expected) {
		t.Errorf("Expected BlackBoxFuncCall to be %v, got %v", expected, blackBoxFuncCall)
	}

	defer file.Close()
}

func TestAES128EncryptUnmarshalReaderWithInputsAndOutputs(t *testing.T) {
	type T = *bn254.BN254Field
	type E = constraint.U64
	file, err := os.Open("../../binaries/black_box_func/aes128encrypt/aes128encrypt_with_inputs_and_outputs.bin")
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}
	kind := shr.ParseThrough32bits(t, file)
	if kind != 0 {
		t.Fatalf("The kind of error code should have been 0, was %d", kind)
	}
	blackBoxFuncCall := BlackBoxFuncCall[T, E]{function: &AES128Encrypt[T, E]{}}
	if err := blackBoxFuncCall.UnmarshalReader(file); err != nil {
		t.Fatalf("Failed to unmarshal BlackBoxFuncCall: %v", err)
	}

	expectedIvWitness := shr.Witness(3456)
	expectedKeyWitness := shr.Witness(4567)
	expectedIv := [16]FunctionInput[T]{}
	for i := 0; i < 16; i++ {
		expectedIv[i] = FunctionInput[T]{
			FunctionInputKind: ACIRFunctionInputKindWitness,
			Witness:           &expectedIvWitness,
		}
	}
	expectedKey := [16]FunctionInput[T]{}
	for i := 0; i < 16; i++ {
		expectedKey[i] = FunctionInput[T]{
			FunctionInputKind: ACIRFunctionInputKindWitness,
			Witness:           &expectedKeyWitness,
		}
	}

	expectedWitnessInput1 := shr.Witness(1234)
	expectedWitnessInput2 := shr.Witness(2345)
	expectedInputs := []FunctionInput[T]{
		{
			FunctionInputKind: ACIRFunctionInputKindWitness,
			Witness:           &expectedWitnessInput1,
		},
		{
			FunctionInputKind: ACIRFunctionInputKindWitness,
			Witness:           &expectedWitnessInput2,
		},
	}

	expectedOutputs := []shr.Witness{
		shr.Witness(1234),
		shr.Witness(2345),
		shr.Witness(3456),
	}
	expected := BlackBoxFuncCall[T, E]{
		function: &AES128Encrypt[T, E]{
			Inputs:  expectedInputs,
			Iv:      expectedIv,
			Key:     expectedKey,
			Outputs: expectedOutputs,
		},
	}

	if !blackBoxFuncCall.Equals(expected) {
		t.Errorf("Expected BlackBoxFuncCall to be %v, got %v", expected,
			blackBoxFuncCall)
	}

	defer file.Close()
}
