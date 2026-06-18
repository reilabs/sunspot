package blackboxfunc

import (
	"math/big"
	"os"
	"sunspot/go/acir/msgpackutil"
	shr "sunspot/go/acir/shared"
	"sunspot/go/bn254"
	"testing"
)

func TestFunctionInputUnmarshalReaderConstant(t *testing.T) {
	file, err := os.Open("../../binaries/black_box_func/function_input/function_input_constant.bin")
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}

	var input FunctionInput[*bn254.BN254Field]
	if err := input.UnmarshalReader(msgpackutil.NewReader(file)); err != nil {
		t.Fatalf("Failed to unmarshal FunctionInput: %v", err)
	}

	expectedField := &bn254.BN254Field{Value: *big.NewInt(1234)}
	expected := FunctionInput[*bn254.BN254Field]{
		ConstantInput: &expectedField,
		Witness:       nil,
	}

	if !input.Equals(&expected) {
		t.Errorf("Expected FunctionInput to be %v, got %v", expected, input)
	}

	defer file.Close()
}

func TestFunctionInputUnmarshalReaderWitness(t *testing.T) {
	file, err := os.Open("../../binaries/black_box_func/function_input/function_input_witness.bin")
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}

	var input FunctionInput[*bn254.BN254Field]
	if err := input.UnmarshalReader(msgpackutil.NewReader(file)); err != nil {
		t.Fatalf("Failed to unmarshal FunctionInput: %v", err)
	}

	expectedWitness := shr.Witness(1234)
	expected := FunctionInput[*bn254.BN254Field]{
		ConstantInput: nil,
		Witness:       &expectedWitness,
	}

	if !input.Equals(&expected) {
		t.Errorf("Expected FunctionInput to be %v, got %v", expected, input)
	}

	defer file.Close()
}
