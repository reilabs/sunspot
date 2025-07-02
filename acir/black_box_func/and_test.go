package blackboxfunc

import (
	shr "nr-groth16/acir/shared"
	"nr-groth16/bn254"
	"os"
	"testing"
)

func TestAndUnmarshalReader(t *testing.T) {
	file, err := os.Open("../../binaries/black_box_func/and/and_test.bin")
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}

	blackBoxFuncCall := BlackBoxFuncCall[*bn254.BN254Field]{}
	if err := blackBoxFuncCall.UnmarshalReader(file); err != nil {
		t.Fatalf("Failed to unmarshal BlackBoxFuncCall: %v", err)
	}

	expectedWitnessLhs := shr.Witness(1234)
	expectedWitnessRhs := shr.Witness(2345)
	expectedFunctionCall := BlackBoxFuncCall[*bn254.BN254Field]{
		Kind: ACIRBlackBoxFuncKindAnd,
		And: &And[*bn254.BN254Field]{
			Lhs: FunctionInput[*bn254.BN254Field]{
				FunctionInputKind: ACIRFunctionInputKindWitness,
				Witness:           &expectedWitnessLhs,
				NumberOfBits:      5678,
			},
			Rhs: FunctionInput[*bn254.BN254Field]{
				FunctionInputKind: ACIRFunctionInputKindWitness,
				Witness:           &expectedWitnessRhs,
				NumberOfBits:      6789,
			},
			Output: shr.Witness(3456),
		},
	}

	if !blackBoxFuncCall.Equals(expectedFunctionCall) {
		t.Errorf("Expected BlackBoxFuncCall to be %v, got %v", expectedFunctionCall, blackBoxFuncCall)
	}

	defer file.Close()
}
