package bn254

import (
	"encoding/binary"
	"io"
	"math/big"
	shr "nr-groth16/acir/shared"

	"github.com/consensys/gnark-crypto/ecc/bn254/fp"
)

const BN254_MODULUS_STRING = "21888242871839275222246405745257275088548364400416034343698204186575808495617"

var Bn254Modulus, _ = new(big.Int).SetString(BN254_MODULUS_STRING, 10)

type BN254Field struct {
	Modulus *big.Int
}

func (b *BN254Field) Zero() *BN254Field {
	return &BN254Field{
		Modulus: new(big.Int).SetInt64(0),
	}
}

func (b *BN254Field) One() *BN254Field {
	return &BN254Field{
		Modulus: new(big.Int).SetInt64(1),
	}
}

func (b *BN254Field) UnmarshalReader(r io.Reader) error {
	// Implement the unmarshalling logic here
	bn254Bytes := make([]byte, 72)
	if err := binary.Read(r, binary.LittleEndian, bn254Bytes); err != nil {
		return err
	}

	//val := new(big.Int).SetBytes(bn254Bytes)
	//b.Modulus = val.Mod(val, Bn254Modulus)

	return nil
}

func (b *BN254Field) Equals(other shr.ACIRField) bool {
	return true // Implement the equality check logic here
}

func (b *BN254Field) ToElement() shr.GenericFPElement {
	var element fp.Element
	element.SetBigInt(b.Modulus)
	return shr.GenericFPElement{
		Kind:           shr.GenericFPElementKindBN254,
		BN254FpElement: &element,
	}
}
