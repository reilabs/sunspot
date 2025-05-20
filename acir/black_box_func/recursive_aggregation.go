package blackboxfunc

import (
	"encoding/binary"
	"io"
	shr "nr-groth16/acir/shared"
)

type RecursiveAggregation[T shr.ACIRField] struct {
	VerificationKey []FunctionInput[T]
	Proof           []FunctionInput[T]
	PublicInputs    []FunctionInput[T]
	KeyHash         FunctionInput[T]
	ProofType       uint32
}

func (a *RecursiveAggregation[T]) UnmarshalReader(r io.Reader) error {
	var VerificationKeyCount uint32
	if err := binary.Read(r, binary.LittleEndian, &VerificationKeyCount); err != nil {
		return err
	}
	a.VerificationKey = make([]FunctionInput[T], VerificationKeyCount)
	for i := uint32(0); i < VerificationKeyCount; i++ {
		if err := a.VerificationKey[i].UnmarshalReader(r); err != nil {
			return err
		}
	}

	var ProofCount uint32
	if err := binary.Read(r, binary.LittleEndian, &ProofCount); err != nil {
		return err
	}
	a.Proof = make([]FunctionInput[T], ProofCount)
	for i := uint32(0); i < ProofCount; i++ {
		if err := a.Proof[i].UnmarshalReader(r); err != nil {
			return err
		}
	}

	var PublicInputsCount uint32
	if err := binary.Read(r, binary.LittleEndian, &PublicInputsCount); err != nil {
		return err
	}
	a.PublicInputs = make([]FunctionInput[T], PublicInputsCount)
	for i := uint32(0); i < PublicInputsCount; i++ {
		if err := a.PublicInputs[i].UnmarshalReader(r); err != nil {
			return err
		}
	}

	if err := a.KeyHash.UnmarshalReader(r); err != nil {
		return err
	}

	if err := binary.Read(r, binary.LittleEndian, &a.ProofType); err != nil {
		return err
	}

	return nil
}
