package black_box_ops

import (
	"io"
	mem "nr-groth16/acir/brillig/memory"
)

type BigIntToLEBytes struct {
	Input  mem.MemoryAddress
	Output mem.HeapVector
}

func (b *BigIntToLEBytes) UnmarshalReader(r io.Reader) error {
	if err := b.Input.UnmarshalReader(r); err != nil {
		return err
	}

	if err := b.Output.UnmarshalReader(r); err != nil {
		return err
	}

	return nil
}
