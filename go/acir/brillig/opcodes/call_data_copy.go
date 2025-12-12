package opcodes

import (
	"io"
	mem "sunspot/go/acir/brillig/memory"
)

type CallDataCopy struct {
	DestinationAddress mem.MemoryAddress
	SizeAddress        mem.MemoryAddress
	OffsetAddress      mem.MemoryAddress
}

func (c *CallDataCopy) UnmarshalReader(r io.Reader) error {
	if err := c.DestinationAddress.UnmarshalReader(r); err != nil {
		return err
	}

	if err := c.SizeAddress.UnmarshalReader(r); err != nil {
		return err
	}

	if err := c.OffsetAddress.UnmarshalReader(r); err != nil {
		return err
	}

	return nil
}

func (c *CallDataCopy) Equals(other CallDataCopy) bool {
	return c.DestinationAddress.Equals(other.DestinationAddress) &&
		c.SizeAddress.Equals(other.SizeAddress) &&
		c.OffsetAddress.Equals(other.OffsetAddress)
}
