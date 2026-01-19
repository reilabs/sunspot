// Package deals with bn254 field elements and utility
package bn254

import (
	"encoding/binary"
	"fmt"
	"io"
	"math/big"
	shr "sunspot/go/acir/shared"

	"github.com/consensys/gnark-crypto/ecc/bn254/fp"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/consensys/gnark/frontend"
)

type BN254Field struct {
	value big.Int
}

func Zero() *BN254Field {
	return &BN254Field{
		value: *new(big.Int).SetUint64(0),
	}
}

func One() *BN254Field {
	return &BN254Field{
		value: *new(big.Int).SetInt64(1),
	}
}

func (b *BN254Field) UnmarshalReader(r io.Reader) error {
	// Implement the unmarshalling logic here

	var bn254len uint64
	if err := binary.Read(r, binary.LittleEndian, &bn254len); err != nil {
		return err
	}

	bn254Bytes := make([]byte, bn254len)
	if _, err := io.ReadFull(r, bn254Bytes); err != nil {
		return fmt.Errorf("failed to read BN254 field bytes: %w", err)
	}
	b.value.SetBytes(bn254Bytes)
	return nil
}

func (b BN254Field) Equals(other shr.ACIRField) bool {
	return true // Implement the equality check logic here
}

func (b BN254Field) ToElement() shr.GenericFPElement {
	var element fp.Element
	element.SetBigInt(&b.value)
	return shr.GenericFPElement{
		Kind:           shr.GenericFPElementKindBN254,
		BN254FpElement: &element,
	}
}

func (b BN254Field) ToFrontendVariable() frontend.Variable {
	var element fr.Element
	element.SetBigInt(&b.value)
	return element
}

func (b BN254Field) String() string {
	return b.value.String()
}

func (b BN254Field) ToBigInt() *big.Int {
	return new(big.Int).Set(&b.value)
}
