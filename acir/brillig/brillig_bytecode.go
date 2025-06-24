package brillig

import (
	"encoding/binary"
	"io"
	ops "nr-groth16/acir/brillig/opcodes"
	shr "nr-groth16/acir/shared"

	"fmt"

	"github.com/rs/zerolog/log"
)

type BrilligBytecode[T shr.ACIRField] struct {
	Bytecode []ops.BrilligOpcode[T]
}

func (b *BrilligBytecode[T]) UnmarshalReader(r io.Reader) error {
	var bytecodeSize uint64
	if err := binary.Read(r, binary.LittleEndian, &bytecodeSize); err != nil {
		return err
	}
	log.Trace().Msg("Unmarshalling BrilligBytecode with size: " + fmt.Sprintf("%x", bytecodeSize))

	b.Bytecode = make([]ops.BrilligOpcode[T], bytecodeSize)
	for i := uint64(0); i < bytecodeSize; i++ {
		log.Trace().Msg("Unmarshalling BrilligOpcode at index: " + fmt.Sprint(i))
		if err := b.Bytecode[i].UnmarshalReader(r); err != nil {
			return err
		}
	}

	return nil
}
