package shared

import (
	"io"

	"github.com/consensys/gnark/frontend"
)

type ACIRField interface {
	UnmarshalReader(r io.Reader) error
	Equals(other ACIRField) bool
	ToFrontendVariable() frontend.Variable
}
