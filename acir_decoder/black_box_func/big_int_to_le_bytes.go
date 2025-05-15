package blackboxfunc

import (
	"encoding/binary"
	"io"
	shr "nr-groth16/acir_decoder/shared"
)

type BigIntToLEBytes struct {
	Input  uint32
	Output []shr.Witness
}

func (a *BigIntToLEBytes) UnmarshalReader(r io.Reader) error {
	if err := binary.Read(r, binary.LittleEndian, &a.Input); err != nil {
		return err
	}

	NumOutputs := uint32(0)
	if err := binary.Read(r, binary.LittleEndian, &NumOutputs); err != nil {
		return err
	}
	a.Output = make([]shr.Witness, NumOutputs)
	for i := uint32(0); i < NumOutputs; i++ {
		if err := binary.Read(r, binary.LittleEndian, &a.Output[i]); err != nil {
			return err
		}
	}

	return nil
}
