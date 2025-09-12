package memory_init

import (
	"encoding/binary"
	shr "nr-groth16/acir/shared"
	"nr-groth16/bn254"
	"os"
	"testing"
)

func TestMemoryInitUnmarshalReaderBlockTest(t *testing.T) {
	file, err := os.Open("../../binaries/opcodes/memory_init/memory_init_memory_block.bin")
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}
	// read the encoded call type before reading the actual content
	var kind uint32
	if err := binary.Read(file, binary.LittleEndian, &kind); err != nil {
		t.Fatal("was not able to read type")
	}
	var opcode MemoryInit[*bn254.BN254Field]
	if err := opcode.UnmarshalReader(file); err != nil {
		t.Fatalf("Failed to unmarshal memory init: %v", err)
	}

	expectedOpcode := MemoryInit[*bn254.BN254Field]{
		BlockType: ACIRMemoryBlockMemory,
		BlockID:   0,
		Init:      []shr.Witness{},
	}

	if !opcode.Equals(&expectedOpcode) {
		t.Errorf("Expected opcode to be %v, got %v", expectedOpcode, opcode)
	}

	defer file.Close()
}

func TestMemoryInitUnmarshalReaderCallDataTest(t *testing.T) {
	file, err := os.Open("../../binaries/opcodes/memory_init/memory_init_calldata.bin")
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}
	// read the encoded call type before reading the actual content
	var kind uint32
	if err := binary.Read(file, binary.LittleEndian, &kind); err != nil {
		t.Fatal("was not able to read type")
	}
	var opcode MemoryInit[*bn254.BN254Field]
	if err := opcode.UnmarshalReader(file); err != nil {
		t.Fatalf("Failed to unmarshal memory init: %v", err)
	}

	expectedOpcode := MemoryInit[*bn254.BN254Field]{
		BlockType: ACIRMemoryBlockCallData,
		BlockID:   1,
		Init:      []shr.Witness{0, 1, 2},
		CallData:  new(uint32),
	}
	*expectedOpcode.CallData = 1234

	if !opcode.Equals(&expectedOpcode) {
		t.Errorf("Expected opcode to be %v, got %v", expectedOpcode, opcode)
	}

	defer file.Close()
}

func TestMemoryInitUnmarshalReaderReturnDataTest(t *testing.T) {
	file, err := os.Open("../../binaries/opcodes/memory_init/memory_init_return_data.bin")
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}
	// read the encoded call type before reading the actual content
	var kind uint32
	if err := binary.Read(file, binary.LittleEndian, &kind); err != nil {
		t.Fatal("was not able to read type")
	}
	var opcode MemoryInit[*bn254.BN254Field]
	if err := opcode.UnmarshalReader(file); err != nil {
		t.Fatalf("Failed to unmarshal memory init: %v", err)
	}

	expectedOpcode := MemoryInit[*bn254.BN254Field]{
		BlockType: ACIRMemoryBlockReturnData,
		BlockID:   2,
		Init:      []shr.Witness{0, 1, 2},
		CallData:  nil,
	}

	if !opcode.Equals(&expectedOpcode) {
		t.Errorf("Expected opcode to be %v, got %v", expectedOpcode, opcode)
	}

	defer file.Close()
}
