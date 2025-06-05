package blackboxfunc

import (
	shr "nr-groth16/acir/shared"
	"nr-groth16/bn254"
	"os"
	"testing"
)

func TestBigIntToLEBytesUnmarshalReader(t *testing.T) {
	file, err := os.Open("../../binaries/black_box_func/big_int_to_le_bytes/big_int_to_le_bytes_test.bin")
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}

	blackBoxFuncCall := BlackBoxFuncCall[*bn254.BN254Field]{}
	if err := blackBoxFuncCall.UnmarshalReader(file); err != nil {
		t.Fatalf("Failed to unmarshal BlackBoxFuncCall: %v", err)
	}

	expectedFunctionCall := BlackBoxFuncCall[*bn254.BN254Field]{
		Kind: ACIRBlackBoxFuncKindBigIntToLeBytes,
		BigIntToLEBytes: &BigIntToLEBytes{
			Input:   1234,
			Outputs: []shr.Witness{},
		},
	}

	if !blackBoxFuncCall.Equals(&expectedFunctionCall) {
		t.Errorf("Expected BlackBoxFuncCall to be %v, got %v", expectedFunctionCall, blackBoxFuncCall)
	}

	defer file.Close()
}
