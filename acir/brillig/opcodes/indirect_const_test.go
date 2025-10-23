package opcodes

import (
	"os"
	mem "sunpot/acir/brillig/memory"
	"sunpot/bn254"
	"testing"
)

func TestIndirectConstUnmarshalReader(t *testing.T) {
	file, err := os.Open("../../../binaries/brillig/opcodes/indirect_const/indirect_const.bin")
	if err != nil {
		t.Fatalf("failed to open file: %v", err)
	}

	op := BrilligOpcode[*bn254.BN254Field]{}
	if err := op.UnmarshalReader(file); err != nil {
		t.Fatalf("failed to unmarshal IndirectConst: %v", err)
	}

	expectedIntegerBitSize := mem.IntegerBitSizeU32
	expected := BrilligOpcode[*bn254.BN254Field]{
		OpCode: ACIRBrilligOpcodeIndirectConst,
		IndirectConst: &IndirectConst[*bn254.BN254Field]{
			DestinationPointer: mem.MemoryAddress{
				Kind:  mem.MemoryAddressKindDirect,
				Value: 1234,
			},
			BitSize: mem.BitSize{
				Kind:           mem.BitSizeKindInteger,
				IntegerBitSize: &expectedIntegerBitSize,
			},
			Value: bn254.Zero(),
		},
	}

	if !op.Equals(expected) {
		t.Errorf("expected %v, got %v", expected, op)
	}

	defer file.Close()
}
