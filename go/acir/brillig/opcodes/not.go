package opcodes

import (
	"io"
	mem "sunspot/acir/brillig/memory"
)

type Not struct {
	Destination mem.MemoryAddress
	Source      mem.MemoryAddress
	BitSize     mem.IntegerBitSize
}

func (n *Not) UnmarshalReader(r io.Reader) error {
	if err := n.Destination.UnmarshalReader(r); err != nil {
		return err
	}

	if err := n.Source.UnmarshalReader(r); err != nil {
		return err
	}

	if err := n.BitSize.UnmarshalReader(r); err != nil {
		return err
	}

	return nil
}

func (n *Not) Equals(other Not) bool {
	return n.Destination.Equals(other.Destination) &&
		n.Source.Equals(other.Source) &&
		n.BitSize == other.BitSize
}
