package call

import (
	"encoding/binary"
	exp "nr-groth16/acir/expression"
	shr "nr-groth16/acir/shared"
	"nr-groth16/bn254"
	"os"
	"testing"
)

func TestCallUnmarshalReaderEmpty(t *testing.T) {
	file, err := os.Open("../../binaries/opcodes/call/call_empty.bin")
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}
	shr.ParseThrough32bits(t, file)
	var opcode Call[*bn254.BN254Field]
	if err := opcode.UnmarshalReader(file); err != nil {
		t.Fatalf("Failed to unmarshal call: %v", err)
	}

	expectedOpcode := Call[*bn254.BN254Field]{
		ID:        0,
		Inputs:    []shr.Witness{},
		Outputs:   []shr.Witness{},
		Predicate: nil,
	}

	if !opcode.Equals(&expectedOpcode) {
		t.Errorf("Expected opcode to be %v, got %v", expectedOpcode, opcode)
	}

	defer file.Close()
}

func TestCallUnmarshalReaderWithInputs(t *testing.T) {
	file, err := os.Open("../../binaries/opcodes/call/call_with_inputs.bin")
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}
	// read the encoded call type before reading the actual content
	var kind uint32
	if err := binary.Read(file, binary.LittleEndian, &kind); err != nil {
		t.Fatal("was not able to read type")
	}

	var opcode Call[*bn254.BN254Field]
	if err := opcode.UnmarshalReader(file); err != nil {
		t.Fatalf("Failed to unmarshal call: %v", err)
	}

	expectedOpcode := Call[*bn254.BN254Field]{
		ID:        1,
		Inputs:    []shr.Witness{0, 1, 2, 3, 4},
		Outputs:   []shr.Witness{},
		Predicate: nil,
	}

	if !opcode.Equals(&expectedOpcode) {
		t.Errorf("Expected opcode to be %v, got %v", expectedOpcode, opcode)
	}

	defer file.Close()
}

func TestCallUnmarshalReaderWithOutputs(t *testing.T) {
	file, err := os.Open("../../binaries/opcodes/call/call_with_outputs.bin")
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}

	// read the encoded call type before reading the actual content
	shr.ParseThrough32bits(t, file)
	var opcode Call[*bn254.BN254Field]
	if err := opcode.UnmarshalReader(file); err != nil {
		t.Fatalf("Failed to unmarshal call: %v", err)
	}

	expectedOpcode := Call[*bn254.BN254Field]{
		ID:        2,
		Inputs:    []shr.Witness{},
		Outputs:   []shr.Witness{0, 1},
		Predicate: nil,
	}

	if !opcode.Equals(&expectedOpcode) {
		t.Errorf("Expected opcode to be %v, got %v", expectedOpcode, opcode)
	}

	defer file.Close()
}

func TestCallUnmarshalReaderWithPredicate(t *testing.T) {
	file, err := os.Open("../../binaries/opcodes/call/call_with_predicate.bin")
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}

	// read the encoded call type before reading the actual content
	var kind uint32
	if err := binary.Read(file, binary.LittleEndian, &kind); err != nil {
		t.Fatal("was not able to read type")
	}

	var opcode Call[*bn254.BN254Field]
	if err := opcode.UnmarshalReader(file); err != nil {
		t.Fatalf("Failed to unmarshal call: %v", err)
	}

	expectedOpcode := Call[*bn254.BN254Field]{
		ID:      3,
		Inputs:  []shr.Witness{},
		Outputs: []shr.Witness{},
		Predicate: &exp.Expression[*bn254.BN254Field]{
			MulTerms:           []exp.MulTerm[*bn254.BN254Field]{},
			LinearCombinations: []exp.LinearCombination[*bn254.BN254Field]{},
			Constant:           bn254.Zero(),
		}, // Assuming a valid predicate expression
	}

	if !opcode.Equals(&expectedOpcode) {
		t.Errorf("Expected opcode to be %v, got %v", expectedOpcode, opcode)
	}

	defer file.Close()
}

func TestCallUnmarshalReaderWithInputsAndOutputs(t *testing.T) {
	file, err := os.Open("../../binaries/opcodes/call/call_with_inputs_and_outputs.bin")
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}

	// read the encoded call type before reading the actual content
	var kind uint32
	if err := binary.Read(file, binary.LittleEndian, &kind); err != nil {
		t.Fatal("was not able to read type")
	}
	var opcode Call[*bn254.BN254Field]
	if err := opcode.UnmarshalReader(file); err != nil {
		t.Fatalf("Failed to unmarshal call: %v", err)
	}

	expectedOpcode := Call[*bn254.BN254Field]{
		ID:      4,
		Inputs:  []shr.Witness{0, 1},
		Outputs: []shr.Witness{2, 3},
		Predicate: &exp.Expression[*bn254.BN254Field]{
			MulTerms:           []exp.MulTerm[*bn254.BN254Field]{},
			LinearCombinations: []exp.LinearCombination[*bn254.BN254Field]{},
			Constant:           bn254.Zero(),
		}, // Assuming a valid predicate expression
	}

	if !opcode.Equals(&expectedOpcode) {
		t.Errorf("Expected opcode to be %v, got %v", expectedOpcode, opcode)
	}

	defer file.Close()
}
