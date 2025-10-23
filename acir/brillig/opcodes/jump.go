package opcodes

import (
	"encoding/binary"
	"io"
	mem "sunpot/acir/brillig/memory"
)

type Jump struct {
	Location mem.Label
}

func (j *Jump) UnmarshalReader(r io.Reader) error {
	if err := binary.Read(r, binary.LittleEndian, &j.Location); err != nil {
		return err
	}

	return nil
}

func (j *Jump) Equals(other Jump) bool {
	return j.Location == other.Location
}
