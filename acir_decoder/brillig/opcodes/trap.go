package opcodes

import (
	"io"
	mem "nr-groth16/acir_decoder/brillig/memory"
)

type Trap struct {
	RevertData mem.HeapVector
}

func (t *Trap) UnmarshalReader(r io.Reader) error {
	if err := t.RevertData.UnmarshalReader(r); err != nil {
		return err
	}

	return nil
}
