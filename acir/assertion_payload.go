package acir

import (
	"encoding/binary"
	"io"
	shr "nr-groth16/acir/shared"
)

type AssertionPayload[T shr.ACIRField] struct {
	ErrorSelector uint64
	Payload       []ExpressionOrMemory[T]
}

func (a *AssertionPayload[T]) UnmarshalReader(r io.Reader) error {
	if err := binary.Read(r, binary.LittleEndian, &a.ErrorSelector); err != nil {
		return err
	}

	var numPayload uint32
	if err := binary.Read(r, binary.LittleEndian, &numPayload); err != nil {
		return err
	}

	a.Payload = make([]ExpressionOrMemory[T], numPayload)
	for i := uint32(0); i < numPayload; i++ {
		if err := a.Payload[i].UnmarshalReader(r); err != nil {
			return err
		}
	}

	return nil
}
