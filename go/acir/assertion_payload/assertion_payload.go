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
	return msgpackutil.ReadStruct(r, assertionPayloadSchema, a.decode)
}

func (a *AssertionPayload[T, E]) decode(f msgpackutil.Field, r *msgpackutil.Reader) error {
	switch f.Tag {
	case 0:
		v, err := r.ReadUint()
		if err != nil {
			return err
		}
		a.ErrorSelector = v
		return nil
	case 1:
		return msgpackutil.ReadVec(r, &a.Payload)
	default:
		return fmt.Errorf("AssertionPayload: unknown field %s", f)
	}
}

var assertionPayloadSchema = msgpackutil.NewSchema(map[string]int{
	"error_selector": 0, "payload": 1,
})
