package memory

import (
	"os"
	"testing"
)

func TestMemoryAddressUnmarshalReaderDirectZero(t *testing.T) {
	file, err := os.Open("../../../binaries/brillig/memory/memory_address/memory_address_direct_zero.bin")
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}

	address := MemoryAddress{}
	err = address.UnmarshalReader(file)
	if err != nil {
		t.Fatalf("Failed to unmarshal MemoryAddress: %v", err)
	}

	expectedAddress := MemoryAddress{
		Kind:  MemoryAddressKindDirect,
		Value: 0,
	}
	if !address.Equals(expectedAddress) {
		t.Fatalf("Expected MemoryAddress to be %v, got %v", expectedAddress, address)
	}

	defer file.Close()
}

func TestMemoryAddressUnmarshalReaderDirect0x1234(t *testing.T) {
	file, err := os.Open("../../../binaries/brillig/memory/memory_address/memory_address_direct_0x1234.bin")
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}

	address := MemoryAddress{}
	err = address.UnmarshalReader(file)
	if err != nil {
		t.Fatalf("Failed to unmarshal MemoryAddress: %v", err)
	}

	expectedAddress := MemoryAddress{
		Kind:  MemoryAddressKindDirect,
		Value: 0x1234,
	}
	if !address.Equals(expectedAddress) {
		t.Fatalf("Expected MemoryAddress to be %v, got %v", expectedAddress, address)
	}

	defer file.Close()
}

func TestMemoryAddressUnmarshalReaderRelativeZero(t *testing.T) {
	file, err := os.Open("../../../binaries/brillig/memory/memory_address/memory_address_relative_zero.bin")
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}

	address := MemoryAddress{}
	err = address.UnmarshalReader(file)
	if err != nil {
		t.Fatalf("Failed to unmarshal MemoryAddress: %v", err)
	}

	expectedAddress := MemoryAddress{
		Kind:  MemoryAddressKindRelative,
		Value: 0,
	}
	if !address.Equals(expectedAddress) {
		t.Fatalf("Expected MemoryAddress to be %v, got %v", expectedAddress, address)
	}

	defer file.Close()
}

func TestMemoryAddressUnmarshalReaderRelative0x1234(t *testing.T) {
	file, err := os.Open("../../../binaries/brillig/memory/memory_address/memory_address_relative_0x1234.bin")
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}

	address := MemoryAddress{}
	err = address.UnmarshalReader(file)
	if err != nil {
		t.Fatalf("Failed to unmarshal MemoryAddress: %v", err)
	}

	expectedAddress := MemoryAddress{
		Kind:  MemoryAddressKindRelative,
		Value: 0x1234,
	}
	if !address.Equals(expectedAddress) {
		t.Fatalf("Expected MemoryAddress to be %v, got %v", expectedAddress, address)
	}

	defer file.Close()
}
