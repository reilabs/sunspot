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
	var VerificationKeyCount uint64
	if err := binary.Read(r, binary.LittleEndian, &VerificationKeyCount); err != nil {
		return err
	}
	a.VerificationKey = make([]FunctionInput[T], VerificationKeyCount)
	for i := uint64(0); i < VerificationKeyCount; i++ {
		if err := a.VerificationKey[i].UnmarshalReader(r); err != nil {
			return err
		}
	}

	var ProofCount uint64
	if err := binary.Read(r, binary.LittleEndian, &ProofCount); err != nil {
		return err
	}
	a.Proof = make([]FunctionInput[T], ProofCount)
	for i := uint64(0); i < ProofCount; i++ {
		if err := a.Proof[i].UnmarshalReader(r); err != nil {
			return err
		}
	}

	var PublicInputsCount uint64
	if err := binary.Read(r, binary.LittleEndian, &PublicInputsCount); err != nil {
		return err
	}
	a.PublicInputs = make([]FunctionInput[T], PublicInputsCount)
	for i := uint64(0); i < PublicInputsCount; i++ {
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

func (a *RecursiveAggregation[T]) Equals(other *RecursiveAggregation[T]) bool {
	if len(a.VerificationKey) != len(other.VerificationKey) ||
		len(a.Proof) != len(other.Proof) ||
		len(a.PublicInputs) != len(other.PublicInputs) ||
		a.ProofType != other.ProofType {
		return false
	}

	for i := range a.VerificationKey {
		if !a.VerificationKey[i].Equals(&other.VerificationKey[i]) {
			return false
		}
	}

	for i := range a.Proof {
		if !a.Proof[i].Equals(&other.Proof[i]) {
			return false
		}
	}

	for i := range a.PublicInputs {
		if !a.PublicInputs[i].Equals(&other.PublicInputs[i]) {
			return false
		}
	}

	return a.KeyHash.Equals(&other.KeyHash)
}
