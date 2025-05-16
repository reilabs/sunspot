package shared

import "io"

type ACIRField interface {
	UnmarshalReader(r io.Reader) error
}
