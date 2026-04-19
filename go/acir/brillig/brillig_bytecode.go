package brillig

import (
	"encoding/binary"
	"io"
	ops "sunspot/go/acir/brillig/opcodes"
	shr "sunspot/go/acir/shared"
)

type BrilligBytecode[T shr.ACIRField] struct {
	FunctionName string
	Bytecode     []ops.BrilligOpcode[T]
}

func (b *BrilligBytecode[T]) UnmarshalReader(r io.Reader) error {
	var nameLen uint64
	if err := binary.Read(r, binary.LittleEndian, &nameLen); err != nil {
		return err
	}
	nameData := make([]byte, nameLen)
	if _, err := io.ReadFull(r, nameData); err != nil {
		return err
	}
	b.FunctionName = string(nameData)

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
