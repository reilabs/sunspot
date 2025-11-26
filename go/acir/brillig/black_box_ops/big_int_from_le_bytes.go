package black_box_ops

import (
	"io"
	mem "sunspot/acir/brillig/memory"
)

type BigIntFromLEBytes struct {
	Inputs  mem.HeapVector
	Modulus mem.HeapArray
	Output  mem.MemoryAddress
}

func (b *BigIntFromLEBytes) UnmarshalReader(r io.Reader) error {
	if err := b.Inputs.UnmarshalReader(r); err != nil {
		return err
	}

	if err := b.Modulus.UnmarshalReader(r); err != nil {
		return err
	}

	if err := b.Output.UnmarshalReader(r); err != nil {
		return err
	}

	return nil
}
