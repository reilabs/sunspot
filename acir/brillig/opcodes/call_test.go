package opcodes

import (
	"os"
	mem "sunpot/acir/brillig/memory"
	"sunpot/bn254"
	"testing"
)

func TestCallUnmarshalReader(t *testing.T) {
	file, err := os.Open("../../../binaries/brillig/opcodes/call/call.bin")
	if err != nil {
		t.Fatalf("failed to open file: %v", err)
	}

	op := BrilligOpcode[*bn254.BN254Field]{}
	if err := op.UnmarshalReader(file); err != nil {
		t.Fatalf("failed to unmarshal Call: %v", err)
	}

	expected := BrilligOpcode[*bn254.BN254Field]{
		OpCode: ACIRBrilligOpcodeCall,
		Call: &Call{
			Location: mem.Label(1234),
		},
	}

	if !op.Equals(expected) {
		t.Errorf("expected %v, got %v", expected, op)
	}

	defer file.Close()
}
