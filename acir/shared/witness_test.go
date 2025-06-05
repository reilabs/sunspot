package shared

import (
	"os"
	"testing"
)

func TestWitnessZero(t *testing.T) {
	file, err := os.Open("../../binaries/shared/witness/witness_zero.bin")
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}

	witness := Witness(0)
	err = witness.UnmarshalReader(file)
	if err != nil {
		t.Fatalf("Failed to unmarshal witness: %v", err)
	}

	if witness != 0 {
		t.Errorf("Expected witness to be 0, got %d", witness)
	}
	defer file.Close()
}

func TestWitness1234(t *testing.T) {
	file, err := os.Open("../../binaries/shared/witness/witness_1234.bin")
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}

	witness := Witness(0)
	err = witness.UnmarshalReader(file)
	if err != nil {
		t.Fatalf("Failed to unmarshal witness: %v", err)
	}

	if witness != 0x1234 {
		t.Errorf("Expected witness to be 0x1234, got %x", witness)
	}
	defer file.Close()
}
