package opcodes

import (
	"io"
	mem "nr-groth16/acir/brillig/memory"
)

type Cast struct {
	Destination mem.MemoryAddress
	Source      mem.MemoryAddress
	BitSize     mem.BitSize
}

func (c *Cast) UnmarshalReader(r io.Reader) error {
	if err := c.Destination.UnmarshalReader(r); err != nil {
		return err
	}

	if err := c.Source.UnmarshalReader(r); err != nil {
		return err
	}

	if err := c.BitSize.UnmarshalReader(r); err != nil {
		return err
	}

	return nil
}
