package call

import (
	"encoding/binary"
	"os"
	exp "sunpot/acir/expression"
	shr "sunpot/acir/shared"
	"sunpot/bn254"
	"testing"

	"github.com/consensys/gnark/constraint"
)

func TestCallUnmarshalReaderEmpty(t *testing.T) {
	type T = *bn254.BN254Field
	type E = constraint.U64
	file, err := os.Open("../../binaries/opcodes/call/call_empty.bin")
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}
	shr.ParseThrough32bits(t, file)
	var opcode Call[T, E]
	if err := opcode.UnmarshalReader(file); err != nil {
		t.Fatalf("Failed to unmarshal call: %v", err)
	}

	expectedOpcode := Call[T, E]{
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
	type T = *bn254.BN254Field
	type E = constraint.U64
	file, err := os.Open("../../binaries/opcodes/call/call_with_inputs.bin")
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}
	// read the encoded call type before reading the actual content
	var kind uint32
	if err := binary.Read(file, binary.LittleEndian, &kind); err != nil {
		t.Fatal("was not able to read type")
	}

	var opcode Call[T, E]
	if err := opcode.UnmarshalReader(file); err != nil {
		t.Fatalf("Failed to unmarshal call: %v", err)
	}

	expectedOpcode := Call[T, E]{
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
	type T = *bn254.BN254Field
	type E = constraint.U64
	file, err := os.Open("../../binaries/opcodes/call/call_with_outputs.bin")
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}

	// read the encoded call type before reading the actual content
	shr.ParseThrough32bits(t, file)
	var opcode Call[T, E]
	if err := opcode.UnmarshalReader(file); err != nil {
		t.Fatalf("Failed to unmarshal call: %v", err)
	}

	expectedOpcode := Call[T, E]{
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
	type T = *bn254.BN254Field
	type E = constraint.U64
	file, err := os.Open("../../binaries/opcodes/call/call_with_predicate.bin")
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}

	// read the encoded call type before reading the actual content
	var kind uint32
	if err := binary.Read(file, binary.LittleEndian, &kind); err != nil {
		t.Fatal("was not able to read type")
	}

	var opcode Call[T, E]
	if err := opcode.UnmarshalReader(file); err != nil {
		t.Fatalf("Failed to unmarshal call: %v", err)
	}

	expectedOpcode := Call[T, E]{
		ID:      3,
		Inputs:  []shr.Witness{},
		Outputs: []shr.Witness{},
		Predicate: &exp.Expression[T, E]{
			MulTerms:           []exp.MulTerm[T]{},
			LinearCombinations: []exp.LinearCombination[T]{},
			Constant:           bn254.Zero(),
		}, // Assuming a valid predicate expression
	}

	if !opcode.Equals(&expectedOpcode) {
		t.Errorf("Expected opcode to be %v, got %v", expectedOpcode, opcode)
	}

	defer file.Close()
}

func TestCallUnmarshalReaderWithInputsAndOutputs(t *testing.T) {
	type T = *bn254.BN254Field
	type E = constraint.U64
	file, err := os.Open("../../binaries/opcodes/call/call_with_inputs_and_outputs.bin")
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}

	// read the encoded call type before reading the actual content
	var kind uint32
	if err := binary.Read(file, binary.LittleEndian, &kind); err != nil {
		t.Fatal("was not able to read type")
	}
	var opcode Call[T, E]
	if err := opcode.UnmarshalReader(file); err != nil {
		t.Fatalf("Failed to unmarshal call: %v", err)
	}

	expectedOpcode := Call[T, E]{
		ID:      4,
		Inputs:  []shr.Witness{0, 1},
		Outputs: []shr.Witness{2, 3},
		Predicate: &exp.Expression[T, E]{
			MulTerms:           []exp.MulTerm[T]{},
			LinearCombinations: []exp.LinearCombination[T]{},
			Constant:           bn254.Zero(),
		}, // Assuming a valid predicate expression
	}

	if !opcode.Equals(&expectedOpcode) {
		t.Errorf("Expected opcode to be %v, got %v", expectedOpcode, opcode)
	}

	defer file.Close()
}
