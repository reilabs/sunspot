package bn254

import (
	"encoding/binary"
	"io"
	"math/big"
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
	bn254Bytes := make([]byte, 32)
	if err := binary.Read(r, binary.LittleEndian, bn254Bytes); err != nil {
		return err
	}

	val := new(big.Int).SetBytes(bn254Bytes)
	b.Modulus = val.Mod(val, Bn254Modulus)

	return nil
}
