package blackboxfunc

import (
	"encoding/binary"
	"io"
	shr "nr-groth16/acir/shared"
)

type BigIntToLEBytes struct {
	Input   uint32
	Outputs []shr.Witness
}

func (a *BigIntToLEBytes) UnmarshalReader(r io.Reader) error {
	if err := binary.Read(r, binary.LittleEndian, &a.Input); err != nil {
		return err
	}

	NumOutputs := uint32(0)
	if err := binary.Read(r, binary.LittleEndian, &NumOutputs); err != nil {
		return err
	}
	a.Outputs = make([]shr.Witness, NumOutputs)
	for i := uint32(0); i < NumOutputs; i++ {
		if err := binary.Read(r, binary.LittleEndian, &a.Outputs[i]); err != nil {
			return err
		}
	}

	return nil
}

func (a *BigIntToLEBytes) Equals(other *BigIntToLEBytes) bool {
	if a.Input != other.Input {
		return false
	}
	if len(a.Outputs) != len(other.Outputs) {
		return false
	}
	for i := range a.Outputs {
		if a.Outputs[i] != other.Outputs[i] {
			return false
		}
	}
	return true
}
