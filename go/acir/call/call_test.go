package call

import (
	"os"
	exp "sunspot/go/acir/expression"
	"sunspot/go/acir/msgpackutil"
	shr "sunspot/go/acir/shared"
	"sunspot/go/bn254"
	"testing"

	"github.com/consensys/gnark/constraint"
)

type fixtureT = *bn254.BN254Field
type fixtureE = constraint.U64

func loadCallFixture(t *testing.T, path string) Call[fixtureT, fixtureE] {
	t.Helper()
	file, err := os.Open(path)
	if err != nil {
		t.Fatalf("Failed to open %s: %v", path, err)
	}
	t.Cleanup(func() { file.Close() })

	r := msgpackutil.NewReader(file)
	if tag := shr.ConsumeEnumTag(t, r); tag != 5 {
		t.Fatalf("expected Opcode variant 5 (Call), got %d", tag)
	}
	var opcode Call[fixtureT, fixtureE]
	if err := opcode.UnmarshalReader(r); err != nil {
		t.Fatalf("Failed to unmarshal call: %v", err)
	}
	return opcode
}

func emptyExpr() exp.Expression[fixtureT, fixtureE] {
	return exp.Expression[fixtureT, fixtureE]{
		MulTerms:           []exp.MulTerm[fixtureT]{},
		LinearCombinations: []exp.LinearCombination[fixtureT]{},
		Constant:           bn254.Zero(),
	}
}

func oneConstantExpr() exp.Expression[fixtureT, fixtureE] {
	return exp.Expression[fixtureT, fixtureE]{
		MulTerms:           []exp.MulTerm[fixtureT]{},
		LinearCombinations: []exp.LinearCombination[fixtureT]{},
		Constant:           bn254.One(),
	}
}

func TestCallUnmarshalReaderEmpty(t *testing.T) {
	opcode := loadCallFixture(t, "../../binaries/opcodes/call/call_empty.bin")
	expected := Call[fixtureT, fixtureE]{
		ID:        0,
		Inputs:    []shr.Witness{},
		Outputs:   []shr.Witness{},
		Predicate: emptyExpr(),
	}
	if !opcode.Equals(&expected) {
		t.Errorf("Expected opcode to be %+v, got %+v", expected, opcode)
	}
}

func TestCallUnmarshalReaderWithInputs(t *testing.T) {
	opcode := loadCallFixture(t, "../../binaries/opcodes/call/call_with_inputs.bin")
	expected := Call[fixtureT, fixtureE]{
		ID:        1,
		Inputs:    []shr.Witness{0, 1, 2, 3, 4},
		Outputs:   []shr.Witness{},
		Predicate: emptyExpr(),
	}
	if !opcode.Equals(&expected) {
		t.Errorf("Expected opcode to be %+v, got %+v", expected, opcode)
	}
}

func TestCallUnmarshalReaderWithOutputs(t *testing.T) {
	opcode := loadCallFixture(t, "../../binaries/opcodes/call/call_with_outputs.bin")
	expected := Call[fixtureT, fixtureE]{
		ID:        2,
		Inputs:    []shr.Witness{},
		Outputs:   []shr.Witness{0, 1},
		Predicate: emptyExpr(),
	}
	if !opcode.Equals(&expected) {
		t.Errorf("Expected opcode to be %+v, got %+v", expected, opcode)
	}
}

func TestCallUnmarshalReaderWithPredicate(t *testing.T) {
	opcode := loadCallFixture(t, "../../binaries/opcodes/call/call_with_predicate.bin")
	expected := Call[fixtureT, fixtureE]{
		ID:        3,
		Inputs:    []shr.Witness{},
		Outputs:   []shr.Witness{},
		Predicate: oneConstantExpr(),
	}
	if !opcode.Equals(&expected) {
		t.Errorf("Expected opcode to be %+v, got %+v", expected, opcode)
	}
}

func TestCallUnmarshalReaderWithInputsAndOutputs(t *testing.T) {
	opcode := loadCallFixture(t, "../../binaries/opcodes/call/call_with_inputs_and_outputs.bin")
	expected := Call[fixtureT, fixtureE]{
		ID:        4,
		Inputs:    []shr.Witness{0, 1},
		Outputs:   []shr.Witness{2, 3},
		Predicate: oneConstantExpr(),
	}
	if !opcode.Equals(&expected) {
		t.Errorf("Expected opcode to be %+v, got %+v", expected, opcode)
	}
}
