package opcodes

import (
	"encoding/binary"
	"io"
	mem "nr-groth16/acir/brillig/memory"
)

type Call struct {
	Location mem.Label
}

func (c *Call) UnmarshalReader(r io.Reader) error {
	if err := binary.Read(r, binary.LittleEndian, &c.Location); err != nil {
		return err
	}

	return nil
}
