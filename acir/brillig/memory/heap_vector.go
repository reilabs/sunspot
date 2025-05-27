package memory

import "io"

type HeapVector struct {
	Pointer MemoryAddress
	Size    MemoryAddress
}

func (hv *HeapVector) UnmarshalReader(r io.Reader) error {
	if err := hv.Pointer.UnmarshalReader(r); err != nil {
		return err
	}

	if err := hv.Size.UnmarshalReader(r); err != nil {
		return err
	}

	return nil
}

func (hv *HeapVector) Equals(other HeapVector) bool {
	return hv.Pointer.Equals(other.Pointer) && hv.Size.Equals(other.Size)
}
