package opcodes

import (
	"io"
	mem "nr-groth16/acir/brillig/memory"
)

type ConditionalMov struct {
	Destination mem.MemoryAddress
	SourceA     mem.MemoryAddress
	SourceB     mem.MemoryAddress
	Condition   mem.MemoryAddress
}

func (c *ConditionalMov) UnmarshalReader(r io.Reader) error {
	if err := c.Destination.UnmarshalReader(r); err != nil {
		return err
	}

	if err := c.SourceA.UnmarshalReader(r); err != nil {
		return err
	}

	if err := c.SourceB.UnmarshalReader(r); err != nil {
		return err
	}

	if err := c.Condition.UnmarshalReader(r); err != nil {
		return err
	}

	return nil
}
