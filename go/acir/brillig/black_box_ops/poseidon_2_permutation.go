package black_box_ops

import (
	"io"
	mem "sunspot/go/acir/brillig/memory"
)

type Poseidon2Permutation struct {
	Message mem.HeapVector
	Output  mem.HeapArray
}

func (p *Poseidon2Permutation) UnmarshalReader(r io.Reader) error {
	if err := p.Message.UnmarshalReader(r); err != nil {
		return err
	}

	if err := p.Output.UnmarshalReader(r); err != nil {
		return err
	}

	return nil
}
