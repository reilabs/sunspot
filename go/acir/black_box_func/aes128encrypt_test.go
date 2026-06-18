package blackboxfunc

import (
	"os"
	"sunspot/go/acir/msgpackutil"
	shr "sunspot/go/acir/shared"
	"sunspot/go/bn254"
	"testing"

	"github.com/consensys/gnark/constraint"
)

func loadAES128(t *testing.T, path string) BlackBoxFuncCall[*bn254.BN254Field, constraint.U64] {
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

func TestAES128EncryptUnmarshalReaderEmpty(t *testing.T) {
	type T = *bn254.BN254Field
	type E = constraint.U64
	bbf := loadAES128(t, "../../binaries/black_box_func/aes128encrypt/aes128encrypt_empty.bin")

	expectedIvWitness := shr.Witness(1234)
	expectedKeyWitness := shr.Witness(5678)
	expectedIv := [16]FunctionInput[T]{}
	for i := 0; i < 16; i++ {
		expectedIv[i] = FunctionInput[T]{
			Witness: &expectedIvWitness,
		}
	}
	expectedKey := [16]FunctionInput[T]{}
	for i := 0; i < 16; i++ {
		expectedKey[i] = FunctionInput[T]{
			Witness: &expectedKeyWitness,
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
	if !bbf.Equals(&expected) {
		t.Errorf("Expected BlackBoxFuncCall to be %+v, got %+v", expected, bbf)
	}
}

func TestAES128EncryptUnmarshalReaderWithInputsAndOutputs(t *testing.T) {
	type T = *bn254.BN254Field
	type E = constraint.U64
	bbf := loadAES128(t, "../../binaries/black_box_func/aes128encrypt/aes128encrypt_with_inputs_and_outputs.bin")

	expectedIvWitness := shr.Witness(3456)
	expectedKeyWitness := shr.Witness(4567)
	expectedIv := [16]FunctionInput[T]{}
	for i := 0; i < 16; i++ {
		expectedIv[i] = FunctionInput[T]{
			Witness: &expectedIvWitness,
		}
	}
	expectedKey := [16]FunctionInput[T]{}
	for i := 0; i < 16; i++ {
		expectedKey[i] = FunctionInput[T]{
			Witness: &expectedKeyWitness,
		}
	}

	expectedWitnessInput1 := shr.Witness(1234)
	expectedWitnessInput2 := shr.Witness(2345)
	expectedInputs := []FunctionInput[T]{
		{Witness: &expectedWitnessInput1},
		{Witness: &expectedWitnessInput2},
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
	if !bbf.Equals(&expected) {
		t.Errorf("Expected BlackBoxFuncCall to be %+v, got %+v", expected, bbf)
	}
}
