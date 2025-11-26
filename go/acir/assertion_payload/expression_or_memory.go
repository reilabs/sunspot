package assertion_payload

import (
	"encoding/binary"
	"fmt"
	"io"
	exp "sunspot/acir/expression"
	shr "sunspot/acir/shared"

	"github.com/consensys/gnark/constraint"
)

// Expression or memory is the basic type used in assertion payloads
type ExpressionOrMemory[T shr.ACIRField, E constraint.Element] struct {
	Kind       ExpressionOrMemoryKind
	Expression *exp.Expression[T, E]
	BlockId    *uint32
}

type ExpressionOrMemoryKind uint32

const (
	ACIRExpressionOrMemoryKindExpression ExpressionOrMemoryKind = iota
	ACIRExpressionOrMemoryKindMemory
)

func (e *ExpressionOrMemory[T, E]) UnmarshalReader(r io.Reader) error {
	if err := binary.Read(r, binary.LittleEndian, &e.Kind); err != nil {
		return err
	}

	switch e.Kind {
	case ACIRExpressionOrMemoryKindExpression:
		e.Expression = new(exp.Expression[T, E])
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
