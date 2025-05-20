package opcodes

import (
	"io"
	mem "nr-groth16/acir_decoder/brillig/memory"
)

type Not struct {
	destination mem.MemoryAddress
	source      mem.MemoryAddress
	bit_size    mem.IntegerBitSize
}

func (n *Not) UnmarshalReader(r io.Reader) error {
	if err := n.destination.UnmarshalReader(r); err != nil {
		return err
	}

	if err := n.source.UnmarshalReader(r); err != nil {
		return err
	}

	if err := n.bit_size.UnmarshalReader(r); err != nil {
		return err
	}

	return nil
}
