package blackboxfunc

import (
	"encoding/binary"
	"io"
	shr "nr-groth16/acir/shared"
)

type ECDSASECP256R1[T shr.ACIRField] struct {
	PublicKeyX    [32]FunctionInput[T]
	PublicKeyY    [32]FunctionInput[T]
	Signature     [64]FunctionInput[T]
	HashedMessage [32]FunctionInput[T]
	Output        shr.Witness
}

func (a *ECDSASECP256R1[T]) UnmarshalReader(r io.Reader) error {
	for i := 0; i < 32; i++ {
		if err := a.PublicKeyX[i].UnmarshalReader(r); err != nil {
			return err
		}
	}

	for i := 0; i < 32; i++ {
		if err := a.PublicKeyY[i].UnmarshalReader(r); err != nil {
			return err
		}
	}

	for i := 0; i < 64; i++ {
		if err := a.Signature[i].UnmarshalReader(r); err != nil {
			return err
		}
	}

	for i := 0; i < 32; i++ {
		if err := a.HashedMessage[i].UnmarshalReader(r); err != nil {
			return err
		}
	}

	if err := binary.Read(r, binary.LittleEndian, &a.Output); err != nil {
		return err
	}
	return nil
}
