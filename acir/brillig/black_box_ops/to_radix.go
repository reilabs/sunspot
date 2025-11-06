package black_box_ops

import (
	"io"
	mem "sunspot/acir/brillig/memory"
)

type ToRadix struct {
	Input         mem.MemoryAddress
	Radix         mem.MemoryAddress
	OutputPointer mem.MemoryAddress
	NumLimbs      mem.MemoryAddress
	OutputBits    mem.MemoryAddress
}

func (t *ToRadix) UnmarshalReader(r io.Reader) error {
	if err := t.Input.UnmarshalReader(r); err != nil {
		return err
	}

	if err := t.Radix.UnmarshalReader(r); err != nil {
		return err
	}

	if err := t.OutputPointer.UnmarshalReader(r); err != nil {
		return err
	}

	if err := t.NumLimbs.UnmarshalReader(r); err != nil {
		return err
	}

	if err := t.OutputBits.UnmarshalReader(r); err != nil {
		return err
	}

	return nil
}
