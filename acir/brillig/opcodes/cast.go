package opcodes

import (
	"io"
	mem "sunpot/acir/brillig/memory"
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

func (c *Cast) Equals(other Cast) bool {
	return c.Destination.Equals(other.Destination) &&
		c.Source.Equals(other.Source) &&
		c.BitSize.Equals(other.BitSize)
}
