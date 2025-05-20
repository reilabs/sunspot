package opcodes

import (
	"io"
	mem "nr-groth16/acir_decoder/brillig/memory"
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
