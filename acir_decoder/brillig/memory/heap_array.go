package memory

import (
	"encoding/binary"
	"io"
)

type HeapArray struct {
	Pointer MemoryAddress
	Length  uint64
}

func (h *HeapArray) UnmarshalReader(r io.Reader) error {
	if err := h.Pointer.UnmarshalReader(r); err != nil {
		return err
	}

	if err := binary.Read(r, binary.LittleEndian, &h.Length); err != nil {
		return err
	}

	return nil
}
