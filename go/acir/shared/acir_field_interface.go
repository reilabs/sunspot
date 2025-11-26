package shared

import (
	"io"
	"math/big"
	"reflect"

	"github.com/consensys/gnark/frontend"
)

func MakeNonNil[T any](v T) T {
	val := reflect.ValueOf(v)

	// Step 1: Is T a pointer type?
	if val.Kind() != reflect.Ptr {
		return v // Not a pointer — leave as is
	}

	// Step 2: Is the pointer nil?
	if val.IsNil() {
		// Step 3: Allocate a new instance of the pointed-to type
		elemType := val.Type().Elem()
		newVal := reflect.New(elemType) // *Elem
		return newVal.Interface().(T)   // Convert to T
	}

	return v // Already non-nil
}

type ACIRField interface {
	UnmarshalReader(r io.Reader) error
	Equals(other ACIRField) bool
	ToFrontendVariable() frontend.Variable
	String() string
	ToBigInt() *big.Int
}
