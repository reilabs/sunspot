package opcodes

import (
	"io"
	mem "nr-groth16/acir/brillig/memory"
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
