package assertion_payload

import (
	"fmt"
	exp "sunspot/go/acir/expression"
	"sunspot/go/acir/msgpackutil"
	shr "sunspot/go/acir/shared"

	"github.com/consensys/gnark/constraint"
)

// Expression or memory is the basic type used in assertion payloads
type ExpressionOrMemory[T shr.ACIRField, E constraint.Element] struct {
	Expression *exp.Expression[T, E]
	BlockId    *uint32
}

// ExpressionOrMemory: 0=Expression(Expression<F>), 1=Memory(BlockId).
func (e *ExpressionOrMemory[T, E]) UnmarshalReader(r *msgpackutil.Reader) error {
	return msgpackutil.ReadEnum(r, e.decode)
}

func (e *ExpressionOrMemory[T, E]) decode(tag int, r *msgpackutil.Reader) error {
	switch tag {
	case 0:
		e.Expression = new(exp.Expression[T, E])
		return e.Expression.UnmarshalReader(r)
	case 1:
		v, err := r.ReadUint()
		if err != nil {
			return err
		}
		e.BlockId = new(uint32)
		*e.BlockId = uint32(v)
		return nil
	default:
		return fmt.Errorf("unknown ExpressionOrMemoryKind: %d", tag)
	}
}
