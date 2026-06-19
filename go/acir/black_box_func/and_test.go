package blackboxfunc

import (
	"os"
	"github.com/reilabs/sunspot/go/acir/msgpackutil"
	shr "github.com/reilabs/sunspot/go/acir/shared"
	"github.com/reilabs/sunspot/go/bn254"
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

	var bbf BlackBoxFuncCall[T, E]
	if err := bbf.UnmarshalReader(msgpackutil.NewReader(file)); err != nil {
		t.Fatalf("Failed to unmarshal BlackBoxFuncCall: %v", err)
	}

	expectedWitnessLhs := shr.Witness(1234)
	expectedWitnessRhs := shr.Witness(2345)
	expected := BlackBoxFuncCall[T, E]{
		function: &And[T, E]{
			Lhs:    FunctionInput[T]{Witness: &expectedWitnessLhs},
			Rhs:    FunctionInput[T]{Witness: &expectedWitnessRhs},
			Output: shr.Witness(3456),
			nBits:  64,
		},
	}
	if !bbf.Equals(&expected) {
		t.Errorf("Expected BlackBoxFuncCall to be %+v, got %+v", expected, bbf)
	}

	defer file.Close()
}
