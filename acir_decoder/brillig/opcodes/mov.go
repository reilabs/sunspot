package opcodes

import (
	"io"
	mem "nr-groth16/acir_decoder/brillig/memory"
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
