package assertion_payload

import (
	"fmt"
	"sunspot/go/acir/msgpackutil"
	shr "sunspot/go/acir/shared"

	"github.com/consensys/gnark/constraint"
)

type AssertionPayload[T shr.ACIRField, E constraint.Element] struct {
	ErrorSelector uint64
	Payload       []ExpressionOrMemory[T, E]
}

// AssertionPayload fields: 0=error_selector (u64), 1=payload (Vec<ExpressionOrMemory>).
func (a *AssertionPayload[T, E]) UnmarshalReader(r *msgpackutil.Reader) error {
	return msgpackutil.ReadStruct(r, a.decode)
}

func (a *AssertionPayload[T, E]) decode(tag int, r *msgpackutil.Reader) error {
	switch tag {
	case 0:
		v, err := r.ReadUint()
		if err != nil {
			return err
		}
		a.ErrorSelector = v
		return nil
	case 1:
		n, err := r.ReadArrayLen()
		if err != nil {
			return err
		}
		a.Payload = make([]ExpressionOrMemory[T, E], n)
		for i := 0; i < n; i++ {
			if err := a.Payload[i].UnmarshalReader(r); err != nil {
				return err
			}
		}
		return nil
	default:
		return fmt.Errorf("AssertionPayload: unknown field tag %d", tag)
	}
}
