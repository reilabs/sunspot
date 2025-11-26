package memory

import (
	"os"
	"testing"
)

func TestHeapVectorUnmarshalReaderZero(t *testing.T) {
	file, err := os.Open("../../../binaries/brillig/memory/heap_vector/heap_vector_zero.bin")
	if err != nil {
		t.Fatalf("failed to open file: %v", err)
	}

	heapVector := HeapVector{}
	err = heapVector.UnmarshalReader(file)
	if err != nil {
		t.Fatalf("failed to unmarshal HeapVector: %v", err)
	}

	expectedHeapVector := HeapVector{
		Pointer: MemoryAddress{Kind: MemoryAddressKindDirect, Value: 0},
		Size:    MemoryAddress{Kind: MemoryAddressKindRelative, Value: 0},
	}

	if !heapVector.Equals(expectedHeapVector) {
		t.Fatalf("expected HeapVector to be %v, got %v", expectedHeapVector, heapVector)
	}

	defer file.Close()
}

func TestHeapVectorUnmarshalReader1234(t *testing.T) {
	file, err := os.Open("../../../binaries/brillig/memory/heap_vector/heap_vector_1234.bin")
	if err != nil {
		t.Fatalf("failed to open file: %v", err)
	}

	heapVector := HeapVector{}
	err = heapVector.UnmarshalReader(file)
	if err != nil {
		t.Fatalf("failed to unmarshal HeapVector: %v", err)
	}

	expectedHeapVector := HeapVector{
		Pointer: MemoryAddress{Kind: MemoryAddressKindDirect, Value: 1234},
		Size:    MemoryAddress{Kind: MemoryAddressKindRelative, Value: 5678},
	}

	if !heapVector.Equals(expectedHeapVector) {
		t.Fatalf("expected HeapVector to be %v, got %v", expectedHeapVector, heapVector)
	}

	defer file.Close()
}
