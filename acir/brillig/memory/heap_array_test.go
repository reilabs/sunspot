package memory

import (
	"os"
	"testing"
)

func TestHeapArrayUnmarshalReaderDirect0x1234(t *testing.T) {
	file, err := os.Open("../../../binaries/brillig/memory/heap_array/heap_array_direct_0x1234.bin")
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}

	heapArray := HeapArray{}
	err = heapArray.UnmarshalReader(file)
	if err != nil {
		t.Fatalf("Failed to unmarshal HeapArray: %v", err)
	}

	expectedHeapArray := HeapArray{
		Pointer: MemoryAddress{
			Kind:  MemoryAddressKindDirect,
			Value: 0x1234,
		},
		Size: 1234,
	}
	if !heapArray.Equals(expectedHeapArray) {
		t.Fatalf("Expected HeapArray to be %v, got %v", expectedHeapArray, heapArray)
	}

	defer file.Close()
}

func TestHeapArrayUnmarshalReaderRelative0x1234(t *testing.T) {
	file, err := os.Open("../../../binaries/brillig/memory/heap_array/heap_array_relative_0x1234.bin")
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}

	heapArray := HeapArray{}
	err = heapArray.UnmarshalReader(file)
	if err != nil {
		t.Fatalf("Failed to unmarshal HeapArray: %v", err)
	}

	expectedHeapArray := HeapArray{
		Pointer: MemoryAddress{
			Kind:  MemoryAddressKindRelative,
			Value: 0x1234,
		},
		Size: 1234,
	}
	if !heapArray.Equals(expectedHeapArray) {
		t.Fatalf("Expected HeapArray to be %v, got %v", expectedHeapArray, heapArray)
	}

	defer file.Close()
}
