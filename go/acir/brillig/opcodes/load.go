package opcodes

import (
	"io"
	mem "sunspot/acir/brillig/memory"
)

type Load struct {
	Destination   mem.MemoryAddress
	SourcePointer mem.MemoryAddress
}

func (l *Load) UnmarshalReader(r io.Reader) error {
	if err := l.Destination.UnmarshalReader(r); err != nil {
		return err
	}

	if err := l.SourcePointer.UnmarshalReader(r); err != nil {
		return err
	}

	return nil
}

func (l *Load) Equals(other Load) bool {
	return l.Destination.Equals(other.Destination) &&
		l.SourcePointer.Equals(other.SourcePointer)
}
