package blackboxfunc

import (
	"encoding/binary"
	"io"
	shr "nr-groth16/acir/shared"

	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/google/btree"
)

type RecursiveAggregation[T shr.ACIRField, E constraint.Element] struct {
	VerificationKey []FunctionInput[T]
	Proof           []FunctionInput[T]
	PublicInputs    []FunctionInput[T]
	KeyHash         FunctionInput[T]
	ProofType       uint32
}

func (a *RecursiveAggregation[T, E]) UnmarshalReader(r io.Reader) error {
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

func (a *RecursiveAggregation[T, E]) Equals(other BlackBoxFunction[E]) bool {
	value, ok := other.(*RecursiveAggregation[T, E])
	if !ok || len(a.VerificationKey) != len(value.VerificationKey) ||
		len(a.Proof) != len(value.Proof) ||
		len(a.PublicInputs) != len(value.PublicInputs) ||
		a.ProofType != value.ProofType {
		return false
	}

	for i := range a.VerificationKey {
		if !a.VerificationKey[i].Equals(&value.VerificationKey[i]) {
			return false
		}
	}

	for i := range a.Proof {
		if !a.Proof[i].Equals(&value.Proof[i]) {
			return false
		}
	}

	for i := range a.PublicInputs {
		if !a.PublicInputs[i].Equals(&value.PublicInputs[i]) {
			return false
		}
	}

	return a.KeyHash.Equals(&value.KeyHash)
}

func (a *RecursiveAggregation[T, E]) Define(api frontend.Builder[E], witnesses map[shr.Witness]frontend.Variable) error {
	return nil
}

func (a *RecursiveAggregation[T, E]) FillWitnessTree(tree *btree.BTree) bool {
	for i := range a.VerificationKey {
		if a.VerificationKey[i].IsWitness() {
			tree.ReplaceOrInsert(*a.VerificationKey[i].Witness)
		}
	}

	for i := range a.Proof {
		if a.Proof[i].IsWitness() {
			tree.ReplaceOrInsert(*a.Proof[i].Witness)
		}
	}
	for i := range a.PublicInputs {
		if a.PublicInputs[i].IsWitness() {
			tree.ReplaceOrInsert(*a.PublicInputs[i].Witness)
		}
	}

	if a.KeyHash.IsWitness() {
		tree.ReplaceOrInsert(*a.KeyHash.Witness)
	}
	return tree != nil
}
