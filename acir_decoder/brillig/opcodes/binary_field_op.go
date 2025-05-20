package opcodes

import (
	"encoding/binary"
	"io"
	mem "nr-groth16/acir_decoder/brillig/memory"
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
