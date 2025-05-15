package blackboxfunc

import (
	"encoding/binary"
	"io"
	shr "nr-groth16/acir_decoder/shared"
)

type SHA256Compression[T shr.ACIRField] struct {
	Inputs     [16]FunctionInput[T]
	HashValues [8]FunctionInput[T]
	Outputs    [8]shr.Witness
}

func (a *SHA256Compression[T]) UnmarshalReader(r io.Reader) error {
	for i := 0; i < 16; i++ {
		if err := a.Inputs[i].UnmarshalReader(r); err != nil {
			return err
		}
	}

	for i := 0; i < 8; i++ {
		if err := a.HashValues[i].UnmarshalReader(r); err != nil {
			return err
		}
	}

	if err := binary.Read(r, binary.LittleEndian, &a.Outputs); err != nil {
		return err
	}

	return nil
}
