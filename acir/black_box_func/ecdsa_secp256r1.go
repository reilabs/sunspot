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

func (a *ECDSASECP256R1[T]) Equals(other *ECDSASECP256R1[T]) bool {
	if len(a.PublicKeyX) != len(other.PublicKeyX) ||
		len(a.PublicKeyY) != len(other.PublicKeyY) ||
		len(a.Signature) != len(other.Signature) ||
		len(a.HashedMessage) != len(other.HashedMessage) {
		return false
	}

	for i := 0; i < 32; i++ {
		if !a.PublicKeyX[i].Equals(&other.PublicKeyX[i]) ||
			!a.PublicKeyY[i].Equals(&other.PublicKeyY[i]) ||
			!a.HashedMessage[i].Equals(&other.HashedMessage[i]) {
			return false
		}
	}

	for i := 0; i < 64; i++ {
		if !a.Signature[i].Equals(&other.Signature[i]) {
			return false
		}
	}

	return a.Output == other.Output
}
