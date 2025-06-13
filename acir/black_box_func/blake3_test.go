package blackboxfunc

import (
	shr "nr-groth16/acir/shared"
	"nr-groth16/bn254"
	"os"
	"testing"
)

func TestBlake3UnmarshalReaderEmpty(t *testing.T) {
	file, err := os.Open("../../binaries/black_box_func/blake3/blake3_test_empty.bin")
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}

	blackBoxFuncCall := BlackBoxFuncCall[bn254.BN254Field]{}
	if err := blackBoxFuncCall.UnmarshalReader(file); err != nil {
		t.Fatalf("Failed to unmarshal BlackBoxFuncCall: %v", err)
	}

	expectedFunctionCall := BlackBoxFuncCall[bn254.BN254Field]{
		Kind: ACIRBlackBoxFuncKindBlake3,
		Blake3: &Blake3[bn254.BN254Field]{
			Inputs:  []FunctionInput[bn254.BN254Field]{},
			Outputs: [32]shr.Witness{},
		},
	}

	for i := 0; i < 32; i++ {
		expectedFunctionCall.Blake3.Outputs[i] = shr.Witness(0)
	}

	if !blackBoxFuncCall.Equals(&expectedFunctionCall) {
		t.Errorf("Expected BlackBoxFuncCall to be %v, got %v", expectedFunctionCall, blackBoxFuncCall)
	}

	defer file.Close()
}

func TestBlake3UnmarshalReaderWithInputs(t *testing.T) {
	file, err := os.Open("../../binaries/black_box_func/blake3/blake3_test_with_inputs.bin")
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}

	blackBoxFuncCall := BlackBoxFuncCall[bn254.BN254Field]{}
	if err := blackBoxFuncCall.UnmarshalReader(file); err != nil {
		t.Fatalf("Failed to unmarshal BlackBoxFuncCall: %v", err)
	}

	expectedWitness1 := shr.Witness(1234)
	expectedWitness2 := shr.Witness(5678)
	expectedFunctionCall := BlackBoxFuncCall[bn254.BN254Field]{
		Kind: ACIRBlackBoxFuncKindBlake3,
		Blake3: &Blake3[bn254.BN254Field]{
			Inputs: []FunctionInput[bn254.BN254Field]{
				{
					FunctionInputKind: ACIRFunctionInputKindWitness,
					Witness:           &expectedWitness1,
					NumberOfBits:      1024,
				},
				{
					FunctionInputKind: ACIRFunctionInputKindWitness,
					Witness:           &expectedWitness2,
					NumberOfBits:      2048,
				},
			},
			Outputs: [32]shr.Witness{},
		},
	}

	for i := 0; i < 32; i++ {
		expectedFunctionCall.Blake3.Outputs[i] = shr.Witness(1234)
	}

	if !blackBoxFuncCall.Equals(&expectedFunctionCall) {
		t.Errorf("Expected BlackBoxFuncCall to be %v, got %v", expectedFunctionCall, blackBoxFuncCall)
	}

	defer file.Close()
}
