package shared

import (
	"encoding/binary"
	"io"
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
