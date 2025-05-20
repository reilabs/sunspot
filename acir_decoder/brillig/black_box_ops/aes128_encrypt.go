package black_box_ops

import (
	"io"
	mem "nr-groth16/acir_decoder/brillig/memory"
)

type AES128Encrypt struct {
	Inputs  mem.HeapVector
	IV      mem.HeapArray
	Key     mem.HeapArray
	Outputs mem.HeapVector
}

func (a *AES128Encrypt) UnmarshalReader(r io.Reader) error {
	if err := a.Inputs.UnmarshalReader(r); err != nil {
		return err
	}

	if err := a.IV.UnmarshalReader(r); err != nil {
		return err
	}

	if err := a.Key.UnmarshalReader(r); err != nil {
		return err
	}

	if err := a.Outputs.UnmarshalReader(r); err != nil {
		return err
	}

	return nil
}
