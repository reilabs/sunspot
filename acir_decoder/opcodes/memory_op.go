package opcodes

import (
	"encoding/binary"
	"io"
	exp "nr-groth16/acir_decoder/expression"
	shr "nr-groth16/acir_decoder/shared"
)

type MemoryOp[T shr.ACIRField] struct {
	BlockID   uint32
	Operation exp.Expression[T]
	Index     exp.Expression[T]
	Value     exp.Expression[T]
	Predicate *exp.Expression[T]
}

func (m *MemoryOp[T]) UnmarshalReader(r io.Reader) error {
	if err := binary.Read(r, binary.LittleEndian, &m.BlockID); err != nil {
		return err
	}

	if err := m.Operation.UnmarshalReader(r); err != nil {
		return err
	}

	if err := m.Index.UnmarshalReader(r); err != nil {
		return err
	}

	if err := m.Value.UnmarshalReader(r); err != nil {
		return err
	}

	var predicateExists uint8
	if err := binary.Read(r, binary.LittleEndian, &predicateExists); err != nil {
		return err
	}
	if predicateExists == 1 {
		m.Predicate = new(exp.Expression[T])
		if err := m.Predicate.UnmarshalReader(r); err != nil {
			return err
		}
	}
	return nil
}
