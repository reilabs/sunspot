package shared

import (
	"encoding/binary"
	"io"

	"github.com/google/btree"
)

type Witness uint32

func (w *Witness) UnmarshalReader(r io.Reader) error {
	var witness uint32
	if err := binary.Read(r, binary.LittleEndian, &witness); err != nil {
		return err
	}
	*w = Witness(witness)
	return nil
}

func (w Witness) Less(other btree.Item) bool {
	otherWitness, ok := other.(Witness)
	if !ok {
		return false
	}
	return w < otherWitness
}
