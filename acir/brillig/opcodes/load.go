package opcodes

import (
	"io"
	mem "nr-groth16/acir/brillig/memory"
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
