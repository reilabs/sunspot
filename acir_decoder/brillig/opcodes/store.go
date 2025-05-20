package opcodes

import (
	"io"
	mem "nr-groth16/acir_decoder/brillig/memory"
)

type Store struct {
	DestinationPointer mem.MemoryAddress
	Source             mem.MemoryAddress
}

func (s *Store) UnmarshalReader(r io.Reader) error {
	if err := s.DestinationPointer.UnmarshalReader(r); err != nil {
		return err
	}

	if err := s.Source.UnmarshalReader(r); err != nil {
		return err
	}

	return nil
}
