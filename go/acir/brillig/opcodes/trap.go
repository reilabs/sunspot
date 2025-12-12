package opcodes

import (
	"io"
	mem "sunspot/go/acir/brillig/memory"
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

func (t *Trap) Equals(other Trap) bool {
	return t.RevertData.Equals(other.RevertData)
}
