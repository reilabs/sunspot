package opcodes

import (
	"encoding/binary"
	"fmt"
	"io"
	mem "sunspot/go/acir/brillig/memory"
)

type BinaryIntOp struct {
	Destination mem.MemoryAddress
	Op          BinaryIntOpKind
	BitSize     mem.IntegerBitSize
	Lhs         mem.MemoryAddress
	Rhs         mem.MemoryAddress
}

type BinaryIntOpKind uint32

const (
	BinaryIntOpKindAdd BinaryIntOpKind = iota
	BinaryIntOpKindSub
	BinaryIntOpKindMul
	BinaryIntOpKindDiv
	BinaryIntOpKindEquals
	BinaryIntOpKindLessThan
	BinaryIntOpKindLessThanEquals
	BinaryIntOpKindAnd
	BinaryIntOpKindOr
	BinaryIntOpKindXor
	BinaryIntOpKindShl
	BinaryIntOpKindShr
)

func (b *BinaryIntOp) UnmarshalReader(r io.Reader) error {
	if err := b.Destination.UnmarshalReader(r); err != nil {
		return err
	}

	if err := binary.Read(r, binary.LittleEndian, &b.Op); err != nil {
		return err
	}

	if b.Op > BinaryIntOpKindShr {
		return fmt.Errorf("invalid BinaryIntOpKind: %d", b.Op)
	}

	if err := b.BitSize.UnmarshalReader(r); err != nil {
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

func (b *BinaryIntOp) Equals(other BinaryIntOp) bool {
	return b.Destination.Equals(other.Destination) &&
		b.Op == other.Op &&
		b.BitSize == other.BitSize &&
		b.Lhs.Equals(other.Lhs) &&
		b.Rhs.Equals(other.Rhs)
}
