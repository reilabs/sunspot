package opcodes

import (
	"io"
	mem "nr-groth16/acir/brillig/memory"
	shr "nr-groth16/acir/shared"
)

type Const[T shr.ACIRField] struct {
	Destination mem.MemoryAddress
	BitSize     mem.BitSize
	Value       T
}

func (c *Const[T]) UnmarshalReader(r io.Reader) error {
	if err := c.Destination.UnmarshalReader(r); err != nil {
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

func (c *Const[T]) Equals(other Const[T]) bool {
	return c.Destination.Equals(other.Destination) &&
		c.BitSize.Equals(other.BitSize) &&
		c.Value.Equals(other.Value)
}
