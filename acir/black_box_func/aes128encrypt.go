package blackboxfunc

import (
	"encoding/binary"
	"io"
	shr "nr-groth16/acir/shared"
)

type AES128Encrypt[T shr.ACIRField] struct {
	Inputs  []FunctionInput[T]
	Iv      [16]FunctionInput[T]
	Key     [16]FunctionInput[T]
	Outputs []shr.Witness
}

func (a *AES128Encrypt[T]) UnmarshalReader(r io.Reader) error {
	InputsNum := uint64(0)
	if err := binary.Read(r, binary.LittleEndian, &InputsNum); err != nil {
		return err
	}
	for i := uint64(0); i < InputsNum; i++ {
		var input FunctionInput[T]
		if err := input.UnmarshalReader(r); err != nil {
			return err
		}
		a.Inputs = append(a.Inputs, input)
	}

	for i := uint32(0); i < 16; i++ {
		if err := a.Iv[i].UnmarshalReader(r); err != nil {
			return err
		}
	}

	for i := uint32(0); i < 16; i++ {
		if err := a.Key[i].UnmarshalReader(r); err != nil {
			return err
		}
	}

	OutputsNum := uint64(0)
	if err := binary.Read(r, binary.LittleEndian, &OutputsNum); err != nil {
		return err
	}
	for i := uint64(0); i < OutputsNum; i++ {
		var witness shr.Witness
		if err := witness.UnmarshalReader(r); err != nil {
			return err
		}
		a.Outputs = append(a.Outputs, witness)
	}

	return nil
}

func (a *AES128Encrypt[T]) Equals(other *AES128Encrypt[T]) bool {
	if len(a.Inputs) != len(other.Inputs) {
		return false
	}
	for i := range a.Inputs {
		if !a.Inputs[i].Equals(&other.Inputs[i]) {
			return false
		}
	}

	for i := uint32(0); i < 16; i++ {
		if !a.Iv[i].Equals(&other.Iv[i]) {
			return false
		}
	}

	for i := uint32(0); i < 16; i++ {
		if !a.Key[i].Equals(&other.Key[i]) {
			return false
		}
	}

	if len(a.Outputs) != len(other.Outputs) {
		return false
	}
	for i := range a.Outputs {
		if !a.Outputs[i].Equals(&other.Outputs[i]) {
			return false
		}
	}

	return true
}
