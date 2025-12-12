package black_box_ops

import (
	"io"
	mem "sunspot/go/acir/brillig/memory"
)

type Sha256Compression struct {
	Input      mem.HeapArray
	HashValues mem.HeapArray
	Output     mem.HeapArray
}

func (s *Sha256Compression) UnmarshalReader(r io.Reader) error {
	if err := s.Input.UnmarshalReader(r); err != nil {
		return err
	}

	if err := s.HashValues.UnmarshalReader(r); err != nil {
		return err
	}

	if err := s.Output.UnmarshalReader(r); err != nil {
		return err
	}

	return nil
}
