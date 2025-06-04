package shared

import (
	"io"

	"github.com/consensys/gnark/frontend"
)

type ACIRField interface {
	UnmarshalReader(r io.Reader) error
	Equals(other ACIRField) bool
	Mul(api frontend.API, other ACIRField) ACIRField
}
