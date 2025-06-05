package opcodes

import (
	"encoding/binary"
	"io"
	mem "nr-groth16/acir/brillig/memory"
)

type JumpIf struct {
	Condition mem.MemoryAddress
	Location  Label
}

func (j *JumpIf) UnmarshalReader(r io.Reader) error {
	if err := j.Condition.UnmarshalReader(r); err != nil {
		return err
	}

	if err := binary.Read(r, binary.LittleEndian, &j.Location); err != nil {
		return err
	}

	return nil
}

func (j *JumpIf) Equals(other JumpIf) bool {
	return j.Condition.Equals(other.Condition) &&
		j.Location == other.Location
}
