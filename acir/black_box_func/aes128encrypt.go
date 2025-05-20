package blackboxfunc

import (
	"encoding/binary"
	"io"
	shr "nr-groth16/acir/shared"
)

type AES128Encrypt[T shr.ACIRField] struct {
	Inputs  []FunctionInput[T]
	Iv      []FunctionInput[T]
	Key     []FunctionInput[T]
	Outputs []shr.Witness
}

func (a *AES128Encrypt[T]) UnmarshalReader(r io.Reader) error {
	InputsNum := uint32(0)
	if err := binary.Read(r, binary.LittleEndian, &InputsNum); err != nil {
		return err
	}
	for i := uint32(0); i < InputsNum; i++ {
		var input FunctionInput[T]
		if err := input.UnmarshalReader(r); err != nil {
			return err
		}
		a.Inputs = append(a.Inputs, input)
	}

	for i := uint32(0); i < 16; i++ {
		var iv FunctionInput[T]
		if err := iv.UnmarshalReader(r); err != nil {
			return err
		}
		a.Iv = append(a.Iv, iv)
	}

	for i := uint32(0); i < 16; i++ {
		var key FunctionInput[T]
		if err := key.UnmarshalReader(r); err != nil {
			return err
		}
		a.Key = append(a.Key, key)
	}

	OutputsNum := uint32(0)
	if err := binary.Read(r, binary.LittleEndian, &OutputsNum); err != nil {
		return err
	}
	for i := uint32(0); i < OutputsNum; i++ {
		var witness shr.Witness
		if err := witness.UnmarshalReader(r); err != nil {
			return err
		}
		a.Outputs = append(a.Outputs, witness)
	}

	return nil
}
