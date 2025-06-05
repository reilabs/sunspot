package memory

import (
	"os"
	"testing"
)

func TestHeapValueTypeUnmarshalReaderSimple(t *testing.T) {
	file, err := os.Open("../../../binaries/brillig/memory/heap_value_type/heap_value_type_simple.bin")
	if err != nil {
		t.Fatalf("failed to open file: %v", err)
	}

	heapValueType := HeapValueType{}
	err = heapValueType.UnmarshalReader(file)
	if err != nil {
		t.Fatalf("failed to unmarshal HeapValueType: %v", err)
	}

	expectedIntegerBitSizeU16 := IntegerBitSizeU16
	expectedHeapValueType := HeapValueType{
		Kind: ACIRBrilligHeapValueTypeKindSimple,
		Simple: &BitSize{
			Kind:           BitSizeKindInteger,
			IntegerBitSize: &expectedIntegerBitSizeU16,
		},
	}

	if !heapValueType.Equals(expectedHeapValueType) {
		t.Fatalf("expected HeapValueType to be %v, got %v", expectedHeapValueType, heapValueType)
	}

	defer file.Close()
}

func TestHeapValueTypeUnmarshalReaderArrayEmpty(t *testing.T) {
	file, err := os.Open("../../../binaries/brillig/memory/heap_value_type/heap_value_type_array_empty.bin")
	if err != nil {
		t.Fatalf("failed to open file: %v", err)
	}

	heapValueType := HeapValueType{}
	err = heapValueType.UnmarshalReader(file)
	if err != nil {
		t.Fatalf("failed to unmarshal HeapValueType: %v", err)
	}

	var expectedArraySize uint64 = 10
	expectedHeapValueType := HeapValueType{
		Kind:            ACIRBrilligHeapValueTypeKindArray,
		ArraySize:       &expectedArraySize,
		ArrayValueTypes: &[]HeapValueType{},
	}

	if !heapValueType.Equals(expectedHeapValueType) {
		t.Fatalf("expected HeapValueType to be %v, got %v", expectedHeapValueType, heapValueType)
	}

	defer file.Close()
}

func TestHeapValueTypeUnmarshalReaderArray(t *testing.T) {
	file, err := os.Open("../../../binaries/brillig/memory/heap_value_type/heap_value_type_array.bin")
	if err != nil {
		t.Fatalf("failed to open file: %v", err)
	}

	heapValueType := HeapValueType{}
	err = heapValueType.UnmarshalReader(file)
	if err != nil {
		t.Fatalf("failed to unmarshal HeapValueType: %v", err)
	}

	var expectedIntegerBitSizeU8 = IntegerBitSizeU8
	var expectedIntegerBitSizeU16 = IntegerBitSizeU16
	var expectedIntegerBitSizeU32 = IntegerBitSizeU32
	var expectedArraySize uint64 = 1234
	var expectedInnerArraySize uint64 = 123
	expectedHeapValueType := HeapValueType{
		Kind:      ACIRBrilligHeapValueTypeKindArray,
		ArraySize: &expectedArraySize,
		ArrayValueTypes: &[]HeapValueType{
			{
				Kind: ACIRBrilligHeapValueTypeKindSimple,
				Simple: &BitSize{
					Kind:           BitSizeKindInteger,
					IntegerBitSize: &expectedIntegerBitSizeU16,
				},
			},
			{
				Kind: ACIRBrilligHeapValueTypeKindSimple,
				Simple: &BitSize{
					Kind:           BitSizeKindInteger,
					IntegerBitSize: &expectedIntegerBitSizeU32,
				},
			},
			{
				Kind:      ACIRBrilligHeapValueTypeKindArray,
				ArraySize: &expectedInnerArraySize,
				ArrayValueTypes: &[]HeapValueType{
					{
						Kind: ACIRBrilligHeapValueTypeKindSimple,
						Simple: &BitSize{
							Kind:           BitSizeKindInteger,
							IntegerBitSize: &expectedIntegerBitSizeU8,
						},
					},
					{
						Kind: ACIRBrilligHeapValueTypeKindSimple,
						Simple: &BitSize{
							Kind:           BitSizeKindInteger,
							IntegerBitSize: &expectedIntegerBitSizeU16,
						},
					},
				},
			},
		},
	}

	if !heapValueType.Equals(expectedHeapValueType) {
		t.Fatalf("expected HeapValueType to be %v, got %v", expectedHeapValueType, heapValueType)
	}

	defer file.Close()
}

func TestHeapValueTypeUnmarshalReaderVectorEmpty(t *testing.T) {
	file, err := os.Open("../../../binaries/brillig/memory/heap_value_type/heap_value_type_vector_empty.bin")
	if err != nil {
		t.Fatalf("failed to open file: %v", err)
	}

	heapValueType := HeapValueType{}
	err = heapValueType.UnmarshalReader(file)
	if err != nil {
		t.Fatalf("failed to unmarshal HeapValueType: %v", err)
	}

	expectedHeapValueType := HeapValueType{
		Kind:             ACIRBrilligHeapValueTypeKindVector,
		VectorValueTypes: &[]HeapValueType{},
	}

	if !heapValueType.Equals(expectedHeapValueType) {
		t.Fatalf("expected HeapValueType to be %v, got %v", expectedHeapValueType, heapValueType)
	}

	defer file.Close()
}

func TestHeapValueTypeUnmarshalReaderVector(t *testing.T) {
	file, err := os.Open("../../../binaries/brillig/memory/heap_value_type/heap_value_type_vector.bin")
	if err != nil {
		t.Fatalf("failed to open file: %v", err)
	}

	heapValueType := HeapValueType{}
	err = heapValueType.UnmarshalReader(file)
	if err != nil {
		t.Fatalf("failed to unmarshal HeapValueType: %v", err)
	}

	var expectedIntegerBitSizeU16 = IntegerBitSizeU16
	var expectedIntegerBitSizeU32 = IntegerBitSizeU32
	expectedHeapValueType := HeapValueType{
		Kind: ACIRBrilligHeapValueTypeKindVector,
		VectorValueTypes: &[]HeapValueType{
			{
				Kind: ACIRBrilligHeapValueTypeKindSimple,
				Simple: &BitSize{
					Kind:           BitSizeKindInteger,
					IntegerBitSize: &expectedIntegerBitSizeU16,
				},
			},
			{
				Kind: ACIRBrilligHeapValueTypeKindSimple,
				Simple: &BitSize{
					Kind:           BitSizeKindInteger,
					IntegerBitSize: &expectedIntegerBitSizeU32,
				},
			},
		},
	}

	if !heapValueType.Equals(expectedHeapValueType) {
		t.Fatalf("expected HeapValueType to be %v, got %v", expectedHeapValueType, heapValueType)
	}

	defer file.Close()
}
