package brillig

import (
	"encoding/binary"
	"io"
	mem "nr-groth16/acir/brillig/memory"
)

type JumpIfNot struct {
	condition mem.MemoryAddress
	location  mem.Label
}

func (j *JumpIfNot) UnmarshalReader(r io.Reader) error {
	if err := j.condition.UnmarshalReader(r); err != nil {
		return err
	}

	if err := binary.Read(r, binary.LittleEndian, &j.location); err != nil {
		return err
	}

	return nil
}
