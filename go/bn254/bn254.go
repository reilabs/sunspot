// Package deals with bn254 field elements and utility
package bn254

import (
	"math/big"
	"github.com/reilabs/sunspot/go/acir/msgpackutil"
	shr "github.com/reilabs/sunspot/go/acir/shared"

	"github.com/consensys/gnark-crypto/ecc/bn254/fp"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/consensys/gnark/frontend"
)

type BN254Field struct {
	Value big.Int
}

func Zero() *BN254Field {
	return &BN254Field{
		Value: *new(big.Int).SetUint64(0),
	}
}

func One() *BN254Field {
	return &BN254Field{
		Value: *new(big.Int).SetInt64(1),
	}
}

func (b *BN254Field) UnmarshalReader(r *msgpackutil.Reader) error {
	bytes, err := r.ReadBytes()
	if err != nil {
		return err
	}
	b.Value.SetBytes(bytes)
	return nil
}

func (b BN254Field) Equals(other shr.ACIRField) bool {
	o, ok := other.(*BN254Field)
	if !ok {
		return false
	}
	return b.Value.Cmp(&o.Value) == 0
}

func (b BN254Field) ToElement() shr.GenericFPElement {
	var element fp.Element
	element.SetBigInt(&b.Value)
	return shr.GenericFPElement{
		Kind:           shr.GenericFPElementKindBN254,
		BN254FpElement: &element,
	}
}

func (b BN254Field) ToFrontendVariable() frontend.Variable {
	var element fr.Element
	element.SetBigInt(&b.Value)
	return element
}

func (b BN254Field) String() string {
	return b.Value.String()
}

func (b BN254Field) ToBigInt() *big.Int {
	return new(big.Int).Set(&b.Value)
}
