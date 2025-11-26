package opcodes

import (
	"os"
	"testing"
)

func TestOpcodeLocationUnmarshalReaderACIR(t *testing.T) {
	file, err := os.Open("../../binaries/opcodes/opcode_location/opcode_location_acir.bin")
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}

	var opcode OpcodeLocation
	if err := opcode.UnmarshalReader(file); err != nil {
		t.Fatalf("Failed to unmarshal opcode location: %v", err)
	}

	expectedOpcode := OpcodeLocation{
		Kind:        ACIROpcodeLocationKindACIR,
		ACIRAddress: new(uint64),
	}

	*expectedOpcode.ACIRAddress = 1234

	if !opcode.Equals(&expectedOpcode) {
		t.Errorf("Expected opcode to be %v, got %v", expectedOpcode, opcode)
	}

	defer file.Close()
}

func TestOpcodeLocationUnmarshalReaderBrillig(t *testing.T) {
	file, err := os.Open("../../binaries/opcodes/opcode_location/opcode_location_brillig.bin")
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}

	var opcode OpcodeLocation
	if err := opcode.UnmarshalReader(file); err != nil {
		t.Fatalf("Failed to unmarshal opcode location: %v", err)
	}

	expectedOpcode := OpcodeLocation{
		Kind:         ACIROpcodeLocationKindBrillig,
		ACIRIndex:    new(uint64),
		BrilligIndex: new(uint64),
	}

	*expectedOpcode.ACIRIndex = 5678
	*expectedOpcode.BrilligIndex = 1234

	if !opcode.Equals(&expectedOpcode) {
		t.Errorf("Expected opcode to be %v, got %v", expectedOpcode, opcode)
	}

	defer file.Close()
}
