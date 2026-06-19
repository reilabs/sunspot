package memory_init

import (
	"os"
	"github.com/reilabs/sunspot/go/acir/msgpackutil"
	shr "github.com/reilabs/sunspot/go/acir/shared"
	"github.com/reilabs/sunspot/go/bn254"
	"testing"

	"github.com/consensys/gnark/constraint"
)

func loadMemoryInit(t *testing.T, path string) MemoryInit[*bn254.BN254Field, constraint.U64] {
	t.Helper()
	file, err := os.Open(path)
	if err != nil {
		t.Fatalf("Failed to open %s: %v", path, err)
	}
	t.Cleanup(func() { file.Close() })

	r := msgpackutil.NewReader(file)
	if tag := shr.ConsumeEnumTag(t, r); tag != 3 {
		t.Fatalf("expected Opcode variant 3 (MemoryInit), got %d", tag)
	}
	var opcode MemoryInit[*bn254.BN254Field, constraint.U64]
	if err := opcode.UnmarshalReader(r); err != nil {
		t.Fatalf("Failed to unmarshal memory init: %v", err)
	}
	return opcode
}

func TestMemoryInitUnmarshalReaderBlockTest(t *testing.T) {
	opcode := loadMemoryInit(t, "../../binaries/opcodes/memory_init/memory_init_memory_block.bin")
	expected := MemoryInit[*bn254.BN254Field, constraint.U64]{
		BlockID: 0,
		Init:    []shr.Witness{},
	}
	if !opcode.Equals(&expected) {
		t.Errorf("Expected opcode to be %+v, got %+v", expected, opcode)
	}
}

func TestMemoryInitUnmarshalReaderCallDataTest(t *testing.T) {
	opcode := loadMemoryInit(t, "../../binaries/opcodes/memory_init/memory_init_calldata.bin")
	expected := MemoryInit[*bn254.BN254Field, constraint.U64]{
		BlockID: 1,
		Init:    []shr.Witness{0, 1, 2},
	}
	if !opcode.Equals(&expected) {
		t.Errorf("Expected opcode to be %+v, got %+v", expected, opcode)
	}
}

func TestMemoryInitUnmarshalReaderReturnDataTest(t *testing.T) {
	opcode := loadMemoryInit(t, "../../binaries/opcodes/memory_init/memory_init_return_data.bin")
	expected := MemoryInit[*bn254.BN254Field, constraint.U64]{
		BlockID: 2,
		Init:    []shr.Witness{0, 1, 2},
	}
	if !opcode.Equals(&expected) {
		t.Errorf("Expected opcode to be %+v, got %+v", expected, opcode)
	}
}
