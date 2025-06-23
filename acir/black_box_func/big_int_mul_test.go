package blackboxfunc

import (
	"nr-groth16/bn254"
	"os"
	"testing"
)

func TestBigIntMulUnmarshalReader(t *testing.T) {
	file, err := os.Open("../../binaries/black_box_func/big_int_mul/big_int_mul_test.bin")
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}

	blackBoxFuncCall := BlackBoxFuncCall[*bn254.BN254Field]{}
	if err := blackBoxFuncCall.UnmarshalReader(file); err != nil {
		t.Fatalf("Failed to unmarshal BlackBoxFuncCall: %v", err)
	}

	expectedFunctionCall := BlackBoxFuncCall[*bn254.BN254Field]{
		Kind: ACIRBlackBoxFuncKindBigIntMul,
		BigIntMul: &BigIntMul{
			Lhs:    1234,
			Rhs:    5678,
			Output: 91011,
		},
	}

	if !blackBoxFuncCall.Equals(&expectedFunctionCall) {
		t.Errorf("Expected BlackBoxFuncCall to be %v, got %v", expectedFunctionCall, blackBoxFuncCall)
	}

	defer file.Close()
}
