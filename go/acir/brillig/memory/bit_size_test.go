package memory

import (
	"os"
	"testing"
)

func TestBitSizeField(t *testing.T) {
	file, err := os.Open("../../../binaries/brillig/memory/bitsize/bitsize_field.bin")
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}

	bitSize := BitSize{}
	err = bitSize.UnmarshalReader(file)
	if err != nil {
		t.Fatalf("Failed to unmarshal BitSize: %v", err)
	}

	expectedBitSize := BitSize{
		Kind:           BitSizeKindField,
		IntegerBitSize: nil,
	}
	if !bitSize.Equals(expectedBitSize) {
		t.Fatalf("Expected BitSize to be %v, got %v", expectedBitSize, bitSize)
	}

	defer file.Close()
}

func TestBitSizeIntegerU1(t *testing.T) {
	file, err := os.Open("../../../binaries/brillig/memory/bitsize/bitsize_integer_u1.bin")
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}

	bitSize := BitSize{}
	err = bitSize.UnmarshalReader(file)
	if err != nil {
		t.Fatalf("Failed to unmarshal BitSize: %v", err)
	}

	expectedIntegerBitSize := IntegerBitSizeU1
	expectedBitSize := BitSize{
		Kind:           BitSizeKindInteger,
		IntegerBitSize: &expectedIntegerBitSize,
	}
	if !bitSize.Equals(expectedBitSize) {
		t.Fatalf("Expected BitSize to be %v, got %v", expectedBitSize, bitSize)
	}

	defer file.Close()
}

func TestBitSizeIntegerU8(t *testing.T) {
	file, err := os.Open("../../../binaries/brillig/memory/bitsize/bitsize_integer_u8.bin")
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}

	bitSize := BitSize{}
	err = bitSize.UnmarshalReader(file)
	if err != nil {
		t.Fatalf("Failed to unmarshal BitSize: %v", err)
	}

	expectedIntegerBitSize := IntegerBitSizeU8
	expectedBitSize := BitSize{
		Kind:           BitSizeKindInteger,
		IntegerBitSize: &expectedIntegerBitSize,
	}

	if !bitSize.Equals(expectedBitSize) {
		t.Fatalf("Expected BitSize to be %v, got %v", expectedBitSize, bitSize)
	}

	defer file.Close()
}

func TestBitSizeIntegerU128(t *testing.T) {
	file, err := os.Open("../../../binaries/brillig/memory/bitsize/bitsize_integer_u128.bin")
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}

	bitSize := BitSize{}
	err = bitSize.UnmarshalReader(file)
	if err != nil {
		t.Fatalf("Failed to unmarshal BitSize: %v", err)
	}

	expectedIntegerBitSize := IntegerBitSizeU128
	expectedBitSize := BitSize{
		Kind:           BitSizeKindInteger,
		IntegerBitSize: &expectedIntegerBitSize,
	}

	if !bitSize.Equals(expectedBitSize) {
		t.Fatalf("Expected BitSize to be %v, got %v", expectedBitSize, bitSize)
	}

	defer file.Close()
}
