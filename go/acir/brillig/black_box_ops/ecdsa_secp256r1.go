package black_box_ops

import (
	"io"
	mem "sunspot/go/acir/brillig/memory"
)

type EcdsaSecp256r1 struct {
	HashedMsg  mem.HeapVector
	PublicKeyX mem.HeapArray
	PublicKeyY mem.HeapArray
	Signature  mem.HeapArray
	Result     mem.MemoryAddress
}

func (e *EcdsaSecp256r1) UnmarshalReader(r io.Reader) error {
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
