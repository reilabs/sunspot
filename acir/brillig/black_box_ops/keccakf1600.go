package black_box_ops

import (
	"io"
	mem "nr-groth16/acir/brillig/memory"
)

type Keccakf1600 struct {
	Input  mem.HeapArray
	Output mem.HeapArray
}

func (k *Keccakf1600) UnmarshalReader(r io.Reader) error {
	if err := k.Input.UnmarshalReader(r); err != nil {
		return err
	}

	if err := k.Output.UnmarshalReader(r); err != nil {
		return err
	}

	return nil
}
