package shared

import "io"

type ACIRField interface {
	UnmarshalReader(r io.Reader) error
	Equals(other ACIRField) bool
}
