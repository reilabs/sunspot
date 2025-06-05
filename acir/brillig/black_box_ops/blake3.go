package black_box_ops

import (
	"io"
	mem "nr-groth16/acir/brillig/memory"
)

type Blake3 struct {
	Message mem.HeapVector
	Output  mem.HeapArray
}

func (b *Blake3) UnmarshalReader(r io.Reader) error {
	if err := b.Message.UnmarshalReader(r); err != nil {
		return err
	}

	if err := b.Output.UnmarshalReader(r); err != nil {
		return err
	}

	return nil
}
