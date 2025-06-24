package opcodes

import (
	"io"
	mem "nr-groth16/acir/brillig/memory"
	shr "nr-groth16/acir/shared"
)

type IndirectConst[T shr.ACIRField] struct {
	DestinationPointer mem.MemoryAddress
	BitSize            mem.BitSize
	Value              T
}

func (c *IndirectConst[T]) UnmarshalReader(r io.Reader) error {
	c.Value = shr.MakeNonNil(c.Value)

	if err := c.DestinationPointer.UnmarshalReader(r); err != nil {
		return err
	}

	if err := c.BitSize.UnmarshalReader(r); err != nil {
		return err
	}

	if err := c.Value.UnmarshalReader(r); err != nil {
		return err
	}

	return nil
}

func (c *IndirectConst[T]) Equals(other IndirectConst[T]) bool {
	return c.DestinationPointer.Equals(other.DestinationPointer) &&
		c.BitSize.Equals(other.BitSize) &&
		c.Value.Equals(other.Value)
}
