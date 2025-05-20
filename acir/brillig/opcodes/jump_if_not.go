package opcodes

import (
	"encoding/binary"
	"io"
	mem "nr-groth16/acir/brillig/memory"
)

type Label uint64

type JumpIfNot struct {
	condition mem.MemoryAddress
	location  Label
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
