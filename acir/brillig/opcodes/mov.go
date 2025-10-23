package opcodes

import (
	"io"
	mem "sunpot/acir/brillig/memory"
)

type Mov struct {
	Destination mem.MemoryAddress
	Source      mem.MemoryAddress
}

func (m *Mov) UnmarshalReader(r io.Reader) error {
	if err := m.Destination.UnmarshalReader(r); err != nil {
		return err
	}

	if err := m.Source.UnmarshalReader(r); err != nil {
		return err
	}

	return nil
}

func (m *Mov) Equals(other Mov) bool {
	return m.Destination.Equals(other.Destination) &&
		m.Source.Equals(other.Source)
}
