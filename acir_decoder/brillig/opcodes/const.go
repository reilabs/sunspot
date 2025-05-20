package opcodes

import (
	"io"
	mem "nr-groth16/acir_decoder/brillig/memory"
	shr "nr-groth16/acir_decoder/shared"
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
