package opcodes

import (
	"os"
	"sunspot/go/acir/msgpackutil"
	"testing"
)

func TestOpcodeLocationUnmarshalReaderACIR(t *testing.T) {
	file, err := os.Open("../../binaries/opcodes/opcode_location/opcode_location_acir.bin")
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	var opcode OpcodeLocation
	if err := opcode.UnmarshalReader(msgpackutil.NewReader(file)); err != nil {
		t.Fatalf("Failed to unmarshal opcode location: %v", err)
	}

	expectedOpcode := OpcodeLocation{
		Kind:        OpcodeLocationAcir,
		ACIRAddress: 1234,
	}

	if opcode != expectedOpcode {
		t.Errorf("Expected opcode to be %v, got %v", expectedOpcode, opcode)
	}
}

func TestOpcodeLocationUnmarshalReaderBrillig(t *testing.T) {
	file, err := os.Open("../../binaries/opcodes/opcode_location/opcode_location_brillig.bin")
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	var opcode OpcodeLocation
	if err := opcode.UnmarshalReader(msgpackutil.NewReader(file)); err != nil {
		t.Fatalf("Failed to unmarshal opcode location: %v", err)
	}

	expectedOpcode := OpcodeLocation{
		Kind:         OpcodeLocationBrillig,
		ACIRIndex:    5678,
		BrilligIndex: 1234,
	}

	if opcode != expectedOpcode {
		t.Errorf("Expected opcode to be %v, got %v", expectedOpcode, opcode)
	}
}
