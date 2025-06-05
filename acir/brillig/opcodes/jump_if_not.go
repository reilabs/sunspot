package opcodes

import (
	"encoding/binary"
	"io"
	mem "nr-groth16/acir/brillig/memory"
)

type Label uint64

type JumpIfNot struct {
	Condition mem.MemoryAddress
	Location  Label
}

func (j *JumpIfNot) UnmarshalReader(r io.Reader) error {
	if err := j.Condition.UnmarshalReader(r); err != nil {
		return err
	}

	if err := binary.Read(r, binary.LittleEndian, &j.Location); err != nil {
		return err
	}

	return nil
}

func (j *JumpIfNot) Equals(other JumpIfNot) bool {
	return j.Condition.Equals(other.Condition) &&
		j.Location == other.Location
}
