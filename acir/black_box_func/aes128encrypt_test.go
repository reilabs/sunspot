package blackboxfunc

import (
	shr "nr-groth16/acir/shared"
	"nr-groth16/bn254"
	"os"
	"testing"
)

func TestAES128EncryptUnmarshalReaderEmpty(t *testing.T) {
	file, err := os.Open("../../binaries/black_box_func/aes128encrypt/aes128encrypt_empty.bin")
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}

	blackBoxFuncCall := BlackBoxFuncCall[*bn254.BN254Field]{}
	if err := blackBoxFuncCall.UnmarshalReader(file); err != nil {
		t.Fatalf("Failed to unmarshal BlackBoxFuncCall: %v", err)
	}

	expectedIvWitness := shr.Witness(1234)
	expectedKeyWitness := shr.Witness(5678)
	expectedIv := [16]FunctionInput[*bn254.BN254Field]{}
	for i := 0; i < 16; i++ {
		expectedIv[i] = FunctionInput[*bn254.BN254Field]{
			FunctionInputKind: ACIRFunctionInputKindWitness,
			Witness:           &expectedIvWitness,
			NumberOfBits:      5678,
		}
	}
	expectedKey := [16]FunctionInput[*bn254.BN254Field]{}
	for i := 0; i < 16; i++ {
		expectedKey[i] = FunctionInput[*bn254.BN254Field]{
			FunctionInputKind: ACIRFunctionInputKindWitness,
			Witness:           &expectedKeyWitness,
			NumberOfBits:      91011,
		}
	}

	expected := BlackBoxFuncCall[*bn254.BN254Field]{
		Kind: ACIRBlackBoxFuncKindAES128Encrypt,
		AES128Encrypt: &AES128Encrypt[*bn254.BN254Field]{
			Inputs:  []FunctionInput[*bn254.BN254Field]{},
			Iv:      expectedIv,
			Key:     expectedKey,
			Outputs: []shr.Witness{},
		},
	}

	if !blackBoxFuncCall.Equals(&expected) {
		t.Errorf("Expected BlackBoxFuncCall to be %v, got %v", expected, blackBoxFuncCall)
	}

	defer file.Close()
}

func TestAES128EncryptUnmarshalReaderWithInputsAndOutputs(t *testing.T) {
	file, err := os.Open("../../binaries/black_box_func/aes128encrypt/aes128encrypt_with_inputs_and_outputs.bin")
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}

	blackBoxFuncCall := BlackBoxFuncCall[*bn254.BN254Field]{}
	if err := blackBoxFuncCall.UnmarshalReader(file); err != nil {
		t.Fatalf("Failed to unmarshal BlackBoxFuncCall: %v", err)
	}

	expectedIvWitness := shr.Witness(3456)
	expectedKeyWitness := shr.Witness(4567)
	expectedIv := [16]FunctionInput[*bn254.BN254Field]{}
	for i := 0; i < 16; i++ {
		expectedIv[i] = FunctionInput[*bn254.BN254Field]{
			FunctionInputKind: ACIRFunctionInputKindWitness,
			Witness:           &expectedIvWitness,
			NumberOfBits:      5678,
		}
	}
	expectedKey := [16]FunctionInput[*bn254.BN254Field]{}
	for i := 0; i < 16; i++ {
		expectedKey[i] = FunctionInput[*bn254.BN254Field]{
			FunctionInputKind: ACIRFunctionInputKindWitness,
			Witness:           &expectedKeyWitness,
			NumberOfBits:      6789,
		}
	}

	expectedWitnessInput1 := shr.Witness(1234)
	expectedWitnessInput2 := shr.Witness(2345)
	expectedInputs := []FunctionInput[*bn254.BN254Field]{
		{
			FunctionInputKind: ACIRFunctionInputKindWitness,
			Witness:           &expectedWitnessInput1,
			NumberOfBits:      5678,
		},
		{
			FunctionInputKind: ACIRFunctionInputKindWitness,
			Witness:           &expectedWitnessInput2,
			NumberOfBits:      6789,
		},
	}

	expectedOutputs := []shr.Witness{
		shr.Witness(1234),
		shr.Witness(2345),
		shr.Witness(3456),
	}
	expected := BlackBoxFuncCall[*bn254.BN254Field]{
		Kind: ACIRBlackBoxFuncKindAES128Encrypt,
		AES128Encrypt: &AES128Encrypt[*bn254.BN254Field]{
			Inputs:  expectedInputs,
			Iv:      expectedIv,
			Key:     expectedKey,
			Outputs: expectedOutputs,
		},
	}

	if !blackBoxFuncCall.Equals(&expected) {
		t.Errorf("Expected BlackBoxFuncCall to be %v, got %v", expected,
			blackBoxFuncCall)
	}

	defer file.Close()
}
