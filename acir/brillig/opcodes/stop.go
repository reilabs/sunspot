package opcodes

import (
	"io"
	mem "sunpot/acir/brillig/memory"
)

type Stop struct {
	ReturnData mem.HeapVector
}

func (s *Stop) UnmarshalReader(r io.Reader) error {
	if err := s.ReturnData.UnmarshalReader(r); err != nil {
		return err
	}

	return nil
}

func (s *Stop) Equals(other Stop) bool {
	return s.ReturnData.Equals(other.ReturnData)
}
