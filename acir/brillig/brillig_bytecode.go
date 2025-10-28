package brillig

import (
	"encoding/binary"
	"io"
	ops "sunspot/acir/brillig/opcodes"
	shr "sunspot/acir/shared"
)

type BrilligBytecode[T shr.ACIRField] struct {
	Bytecode []ops.BrilligOpcode[T]
}

func (b *BrilligBytecode[T]) UnmarshalReader(r io.Reader) error {
	var bytecodeSize uint64
	if err := binary.Read(r, binary.LittleEndian, &bytecodeSize); err != nil {
		return err
	}

	b.Bytecode = make([]ops.BrilligOpcode[T], bytecodeSize)
	for i := uint64(0); i < bytecodeSize; i++ {
		if err := b.Bytecode[i].UnmarshalReader(r); err != nil {
			return err
		}
	}

	return nil
}
