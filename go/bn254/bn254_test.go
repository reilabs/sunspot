package bn254

import (
	"os"
	"testing"
)

func TestBN254FieldUnmarshalReaderZero(t *testing.T) {
	file, err := os.Open("../binaries/acir_field/zero_field.bin")
	if err != nil {
		t.Fatalf("Failed to open test file: %v", err)
	}

	var field BN254Field
	if err := field.UnmarshalReader(file); err != nil {
		t.Fatalf("Failed to unmarshal BN254Field from zero file: %v", err)
	}

	defer file.Close()
}

func TestBN254FieldUnmarshalReader1234(t *testing.T) {
	file, err := os.Open("../binaries/acir_field/field_1234.bin")
	if err != nil {
		t.Fatalf("Failed to open test file: %v", err)
	}

	var field BN254Field
	if err := field.UnmarshalReader(file); err != nil {
		t.Fatalf("Failed to unmarshal BN254Field from 1234 file: %v", err)
	}

	defer file.Close()
}
