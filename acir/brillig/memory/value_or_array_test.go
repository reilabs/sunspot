package memory

import (
	"os"
	"testing"
)

func TestValueOrArrayUnmarshalReaderMemoryAddress(t *testing.T) {
	file, err := os.Open("../../../binaries/brillig/memory/value_or_array/value_or_array_memory_address.bin")
	if err != nil {
		t.Fatalf("failed to open file: %v", err)
	}

	valueOrArray := ValueOrArray{}
	err = valueOrArray.UnmarshalReader(file)
	if err != nil {
		t.Fatalf("failed to unmarshal ValueOrArray: %v", err)
	}

	expectedValueOrArray := ValueOrArray{
		Kind: ACIRBrilligValueOrArrayKindMemoryAddress,
		MemoryAddress: &MemoryAddress{
			Kind:  MemoryAddressKindDirect,
			Value: 1234,
		},
	}

	if !valueOrArray.Equals(expectedValueOrArray) {
		t.Fatalf("expected ValueOrArray to be %v, got %v", expectedValueOrArray, valueOrArray)
	}

	defer file.Close()
}

func TestValueOrArrayUnmarshalReaderHeapArray(t *testing.T) {
	file, err := os.Open("../../../binaries/brillig/memory/value_or_array/value_or_array_heap_array.bin")
	if err != nil {
		t.Fatalf("failed to open file: %v", err)
	}

	valueOrArray := ValueOrArray{}
	err = valueOrArray.UnmarshalReader(file)
	if err != nil {
		t.Fatalf("failed to unmarshal ValueOrArray: %v", err)
	}

	expectedValueOrArray := ValueOrArray{
		Kind: ACIRBrilligValueOrArrayKindHeapArray,
		HeapArray: &HeapArray{
			Pointer: MemoryAddress{Kind: MemoryAddressKindDirect, Value: 1234},
			Size:    5678,
		},
	}

	if !valueOrArray.Equals(expectedValueOrArray) {
		t.Fatalf("expected ValueOrArray to be %v, got %v", expectedValueOrArray, valueOrArray)
	}

	defer file.Close()
}

func TestValueOrArrayUnmarshalReaderHeapVector(t *testing.T) {
	file, err := os.Open("../../../binaries/brillig/memory/value_or_array/value_or_array_heap_vector.bin")
	if err != nil {
		t.Fatalf("failed to open file: %v", err)
	}

	valueOrArray := ValueOrArray{}
	err = valueOrArray.UnmarshalReader(file)
	if err != nil {
		t.Fatalf("failed to unmarshal ValueOrArray: %v", err)
	}

	expectedValueOrArray := ValueOrArray{
		Kind: ACIRBrilligValueOrArrayKindHeapVector,
		HeapVector: &HeapVector{
			Pointer: MemoryAddress{Kind: MemoryAddressKindDirect, Value: 1234},
			Size:    MemoryAddress{Kind: MemoryAddressKindRelative, Value: 5678},
		},
	}

	if !valueOrArray.Equals(expectedValueOrArray) {
		t.Fatalf("expected ValueOrArray to be %v, got %v", expectedValueOrArray, valueOrArray)
	}

	defer file.Close()
}
