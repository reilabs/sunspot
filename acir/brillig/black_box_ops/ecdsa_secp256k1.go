package black_box_ops

import (
	"io"
	mem "nr-groth16/acir/brillig/memory"
)

type EcdsaSecp256k1 struct {
	HashedMsg  mem.HeapVector
	PublicKeyX mem.HeapArray
	PublicKeyY mem.HeapArray
	Signature  mem.HeapArray
	Result     mem.MemoryAddress
}

func (e *EcdsaSecp256k1) UnmarshalReader(r io.Reader) error {
	if err := e.HashedMsg.UnmarshalReader(r); err != nil {
		return err
	}

	if err := e.PublicKeyX.UnmarshalReader(r); err != nil {
		return err
	}

	if err := e.PublicKeyY.UnmarshalReader(r); err != nil {
		return err
	}

	if err := e.Signature.UnmarshalReader(r); err != nil {
		return err
	}

	if err := e.Result.UnmarshalReader(r); err != nil {
		return err
	}

	return nil
}
