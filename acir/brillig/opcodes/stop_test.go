package opcodes

import (
	"os"
	mem "sunpot/acir/brillig/memory"
	"sunpot/bn254"
	"testing"
)

func TestStopUnmarshalReader(t *testing.T) {
	file, err := os.Open("../../../binaries/brillig/opcodes/stop/stop.bin")
	if err != nil {
		t.Fatalf("failed to open file: %v", err)
	}

	op := BrilligOpcode[*bn254.BN254Field]{}
	if err := op.UnmarshalReader(file); err != nil {
		t.Fatalf("failed to unmarshal Stop: %v", err)
	}

	expected := BrilligOpcode[*bn254.BN254Field]{
		OpCode: ACIRBrilligOpcodeStop,
		Stop: &Stop{
			ReturnData: mem.HeapVector{
				Pointer: mem.MemoryAddress{
					Kind:  mem.MemoryAddressKindDirect,
					Value: 1234,
				},
				Size: mem.MemoryAddress{
					Kind:  mem.MemoryAddressKindRelative,
					Value: 5678,
				},
			},
		},
	}

	if !op.Equals(expected) {
		t.Errorf("expected %v, got %v", expected, op)
	}

	defer file.Close()
}
