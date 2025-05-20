package black_box_ops

import (
	"io"
	mem "nr-groth16/acir/brillig/memory"
)

type BigIntDiv struct {
	Lhs    mem.MemoryAddress
	Rhs    mem.MemoryAddress
	Output mem.MemoryAddress
}

func (b *BigIntDiv) UnmarshalReader(r io.Reader) error {
	if err := b.Lhs.UnmarshalReader(r); err != nil {
		return err
	}

	if err := b.Rhs.UnmarshalReader(r); err != nil {
		return err
	}

	if err := b.Output.UnmarshalReader(r); err != nil {
		return err
	}

	return nil
}
