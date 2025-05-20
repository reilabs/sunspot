package acir

import (
	"encoding/binary"
	"fmt"
	"io"
	exp "nr-groth16/acir/expression"
	shr "nr-groth16/acir/shared"
)

type ExpressionOrMemory[T shr.ACIRField] struct {
	Kind       ExpressionOrMemoryKind
	Expression *exp.Expression[T]
	BlockId    *uint32
}

type ExpressionOrMemoryKind uint32

const (
	ACIRExpressionOrMemoryKindExpression ExpressionOrMemoryKind = iota
	ACIRExpressionOrMemoryKindMemory
)

func (e *ExpressionOrMemory[T]) UnmarshalReader(r io.Reader) error {
	if err := binary.Read(r, binary.LittleEndian, &e.Kind); err != nil {
		return err
	}

	switch e.Kind {
	case ACIRExpressionOrMemoryKindExpression:
		e.Expression = new(exp.Expression[T])
		if err := e.Expression.UnmarshalReader(r); err != nil {
			return err
		}
	case ACIRExpressionOrMemoryKindMemory:
		var blockID uint32
		if err := binary.Read(r, binary.LittleEndian, &blockID); err != nil {
			return err
		}
		e.BlockId = new(uint32)
		*e.BlockId = blockID
	default:
		return fmt.Errorf("unknown ExpressionOrMemoryKind: %d", e.Kind)
	}
	return nil
}
