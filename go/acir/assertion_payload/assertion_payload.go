package assertion_payload

import (
	"encoding/binary"
	"io"
	shr "sunspot/go/acir/shared"

	"github.com/consensys/gnark/constraint"
)

type AssertionPayload[T shr.ACIRField, E constraint.Element] struct {
	ErrorSelector uint64
	Payload       []ExpressionOrMemory[T, E]
}

func (a *AssertionPayload[T, E]) UnmarshalReader(r io.Reader) error {
	if err := binary.Read(r, binary.LittleEndian, &a.ErrorSelector); err != nil {
		return err
	}

	var numPayload uint64
	if err := binary.Read(r, binary.LittleEndian, &numPayload); err != nil {
		return err
	}

	a.Payload = make([]ExpressionOrMemory[T, E], numPayload)
	for i := uint64(0); i < numPayload; i++ {
		if err := a.Payload[i].UnmarshalReader(r); err != nil {
			return err
		}
	}

	return nil
}
