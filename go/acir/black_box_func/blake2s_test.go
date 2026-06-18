package blackboxfunc

import (
	"os"
	"sunspot/go/acir/msgpackutil"
	shr "sunspot/go/acir/shared"
	"sunspot/go/bn254"
	"testing"

	"github.com/consensys/gnark/constraint"
)

func loadBlake2s(t *testing.T, path string) BlackBoxFuncCall[*bn254.BN254Field, constraint.U64] {
	t.Helper()
	file, err := os.Open(path)
	if err != nil {
		t.Fatalf("Failed to open %s: %v", path, err)
	}
	t.Cleanup(func() { file.Close() })
	var b BlackBoxFuncCall[*bn254.BN254Field, constraint.U64]
	if err := b.UnmarshalReader(msgpackutil.NewReader(file)); err != nil {
		t.Fatalf("Failed to unmarshal BlackBoxFuncCall: %v", err)
	}
	return b
}

func TestBlake2sUnmarshalReaderEmpty(t *testing.T) {
	type T = *bn254.BN254Field
	type E = constraint.U64
	bbf := loadBlake2s(t, "../../binaries/black_box_func/blake2s/blake2s_test_empty.bin")

	expected := BlackBoxFuncCall[T, E]{function: &Blake2s[T, E]{
		Inputs:  []FunctionInput[T]{},
		Outputs: [32]shr.Witness{},
	}}
	for i := 0; i < 32; i++ {
		expected.function.(*Blake2s[T, E]).Outputs[i] = shr.Witness(0)
	}
	if !bbf.Equals(&expected) {
		t.Errorf("Expected BlackBoxFuncCall to be %+v, got %+v", expected, bbf)
	}
}

func TestBlake2sUnmarshalReaderWithInputs(t *testing.T) {
	type T = *bn254.BN254Field
	type E = constraint.U64
	bbf := loadBlake2s(t, "../../binaries/black_box_func/blake2s/blake2s_test_with_inputs.bin")

	expectedWitness1 := shr.Witness(1234)
	expectedWitness2 := shr.Witness(5678)
	expected := BlackBoxFuncCall[T, E]{function: &Blake2s[T, E]{
		Inputs: []FunctionInput[T]{
			{Witness: &expectedWitness1},
			{Witness: &expectedWitness2},
		},
		Outputs: [32]shr.Witness{},
	}}
	for i := 0; i < 32; i++ {
		expected.function.(*Blake2s[T, E]).Outputs[i] = shr.Witness(1234)
	}
	if !bbf.Equals(&expected) {
		t.Errorf("Expected BlackBoxFuncCall to be %+v, got %+v", expected, bbf)
	}
}
