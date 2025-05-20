package opcodes

import (
	"encoding/binary"
	"io"
	mem "nr-groth16/acir/brillig/memory"
)

type Jump struct {
	location mem.Label
}

func (j *Jump) UnmarshalReader(r io.Reader) error {
	if err := binary.Read(r, binary.LittleEndian, &j.location); err != nil {
		return err
	}

	return nil
}
