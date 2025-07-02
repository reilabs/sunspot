package blackboxfunc

import (
	shr "nr-groth16/acir/shared"
	"nr-groth16/bn254"
	"os"
	"testing"
)

func TestBigIntFromLEBytesUnmarshalReader(t *testing.T) {
	file, err := os.Open("../../binaries/black_box_func/big_int_from_le_bytes/big_int_from_le_bytes_test.bin")
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}

	blackBoxFuncCall := BlackBoxFuncCall[*bn254.BN254Field]{}
	if err := blackBoxFuncCall.UnmarshalReader(file); err != nil {
		t.Fatalf("Failed to unmarshal BlackBoxFuncCall: %v", err)
	}

	expectedFunctionCall := BlackBoxFuncCall[*bn254.BN254Field]{
		Kind: ACIRBlackBoxFuncKindBigIntFromLeBytes,
		BigIntFromLEBytes: &BigIntFromLEBytes[*bn254.BN254Field]{
			Inputs:  []FunctionInput[*bn254.BN254Field]{},
			Modulus: []uint8{},
			Output:  0,
		},
	}

	if !blackBoxFuncCall.Equals(expectedFunctionCall) {
		t.Errorf("Expected BlackBoxFuncCall to be %v, got %v", expectedFunctionCall, blackBoxFuncCall)
	}

	defer file.Close()
}

func TestBigIntFromLEBytesUnmarshalReaderWithInputsAndModulus(t *testing.T) {
	file, err := os.Open("../../binaries/black_box_func/big_int_from_le_bytes/big_int_from_le_bytes_test_with_inputs_and_modulus.bin")
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}

	blackBoxFuncCall := BlackBoxFuncCall[*bn254.BN254Field]{}
	if err := blackBoxFuncCall.UnmarshalReader(file); err != nil {
		t.Fatalf("Failed to unmarshal BlackBoxFuncCall: %v", err)
	}

	expectedWitness := shr.Witness(1)
	expectedFunctionCall := BlackBoxFuncCall[*bn254.BN254Field]{
		Kind: ACIRBlackBoxFuncKindBigIntFromLeBytes,
		BigIntFromLEBytes: &BigIntFromLEBytes[*bn254.BN254Field]{
			Inputs: []FunctionInput[*bn254.BN254Field]{
				{
					FunctionInputKind: ACIRFunctionInputKindWitness,
					Witness:           &expectedWitness,
					NumberOfBits:      5678,
				},
			},
			Modulus: []uint8{1, 2, 3, 4},
			Output:  5678,
		},
	}

	if !blackBoxFuncCall.Equals(expectedFunctionCall) {
		t.Errorf("Expected BlackBoxFuncCall to be %v, got %v", expectedFunctionCall, blackBoxFuncCall)
	}

	defer file.Close()
}
