package blackboxfunc

import (
	"os"
	shr "sunspot/go/acir/shared"
	"sunspot/go/bn254"
	"testing"

	"github.com/consensys/gnark/constraint"
)

func TestAndUnmarshalReader(t *testing.T) {
	type T = *bn254.BN254Field
	type E = constraint.U64
	file, err := os.Open("../../binaries/black_box_func/and/and_test.bin")
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}

	kind := shr.ParseThrough32bits(t, file)
	if kind != 1 {
		t.Fatalf("The kind of error code should have been 0, was %d", kind)
	}
	blackBoxFuncCall := BlackBoxFuncCall[T, E]{function: &And[T, E]{}}
	if err := blackBoxFuncCall.UnmarshalReader(file); err != nil {
		t.Fatalf("Failed to unmarshal BlackBoxFuncCall: %v", err)
	}

	expectedWitnessLhs := shr.Witness(1234)
	expectedWitnessRhs := shr.Witness(2345)
	expectedFunctionCall := BlackBoxFuncCall[T, E]{
		function: &And[T, E]{
			Lhs: FunctionInput[T]{
				FunctionInputKind: ACIRFunctionInputKindWitness,
				Witness:           &expectedWitnessLhs,
			},
			Rhs: FunctionInput[T]{
				FunctionInputKind: ACIRFunctionInputKindWitness,
				Witness:           &expectedWitnessRhs,
			},
			Output: shr.Witness(3456),
		},
	}

	if !blackBoxFuncCall.Equals(expectedFunctionCall) {
		t.Errorf("Expected BlackBoxFuncCall to be %v, got %v", expectedFunctionCall, blackBoxFuncCall)
	}

	defer file.Close()
}
