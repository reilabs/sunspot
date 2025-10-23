package black_box_ops

import (
	"io"
	mem "sunpot/acir/brillig/memory"
)

type Blake2s struct {
	Message mem.HeapVector
	Output  mem.HeapArray
}

func (b *Blake2s) UnmarshalReader(r io.Reader) error {
	if err := b.Message.UnmarshalReader(r); err != nil {
		return err
	}

	if err := b.Output.UnmarshalReader(r); err != nil {
		return err
	}

	return nil
}
