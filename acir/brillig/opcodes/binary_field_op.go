package opcodes

import (
	"encoding/binary"
	"io"
	mem "nr-groth16/acir/brillig/memory"
)

type BinaryFieldOp struct {
	Destination mem.MemoryAddress
	Op          BinaryFieldOpKind
	Lhs         mem.MemoryAddress
	Rhs         mem.MemoryAddress
}

type BinaryFieldOpKind uint32

const (
	BinaryFieldOpAdd BinaryFieldOpKind = iota
	BinaryFieldOpSub
	BinaryFieldOpMul
	BinaryFieldOpDiv
	BinaryFieldOpIntegerDiv
	BinaryFieldOpEquals
	BinaryFieldOpLessThan
	BinaryFieldOpLessThanEquals
)

func (b *BinaryFieldOp) UnmarshalReader(r io.Reader) error {
	if err := b.Destination.UnmarshalReader(r); err != nil {
		return err
	}

	if err := binary.Read(r, binary.LittleEndian, &b.Op); err != nil {
		return err
	}

	if err := b.Lhs.UnmarshalReader(r); err != nil {
		return err
	}

	if err := b.Rhs.UnmarshalReader(r); err != nil {
		return err
	}

	return nil
}

func (b *BinaryFieldOp) Equals(other BinaryFieldOp) bool {
	return b.Destination.Equals(other.Destination) &&
		b.Op == other.Op &&
		b.Lhs.Equals(other.Lhs) &&
		b.Rhs.Equals(other.Rhs)
}
