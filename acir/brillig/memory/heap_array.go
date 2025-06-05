package memory

import (
	"encoding/binary"
	"io"
)

type HeapArray struct {
	Pointer MemoryAddress
	Size    uint64
}

func (h *HeapArray) UnmarshalReader(r io.Reader) error {
	if err := h.Pointer.UnmarshalReader(r); err != nil {
		return err
	}

	if err := binary.Read(r, binary.LittleEndian, &h.Size); err != nil {
		return err
	}

	return nil
}

func (h HeapArray) Equals(other HeapArray) bool {
	if !h.Pointer.Equals(other.Pointer) {
		return false
	}
	return h.Size == other.Size
}
