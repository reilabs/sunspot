package black_box_ops

import (
	"encoding/binary"
	"fmt"
	"io"
)

type BlackBoxOp struct {
	OpType               BlackBoxOpType
	AES128Encrypt        *AES128Encrypt
	Blake2s              *Blake2s
	Blake3               *Blake3
	Keccakf1600          *Keccakf1600
	EcdsaSecp256k1       *EcdsaSecp256k1
	EcdsaSecp256r1       *EcdsaSecp256r1
	MultiScalarMul       *MultiScalarMul
	EmbeddedCurveAdd     *EmbeddedCurveAdd
	BigIntAdd            *BigIntAdd
	BigIntSub            *BigIntSub
	BigIntMul            *BigIntMul
	BigIntDiv            *BigIntDiv
	BigIntFromLeBytes    *BigIntFromLEBytes
	BigIntToLeBytes      *BigIntToLEBytes
	Poseidon2Permutation *Poseidon2Permutation
	Sha256Compression    *Sha256Compression
	ToRadix              *ToRadix
}

type BlackBoxOpType uint32

const (
	ACIRBlackBoxOpAES128Encrypt BlackBoxOpType = iota
	ACIRBlackBoxOpBlake2s
	ACIRBlackBoxOpBlake3
	ACIRBlackBoxOpKeccakf1600
	ACIRBlackBoxOpEcdsaSecp256k1
	ACIRBlackBoxOpEcdsaSecp256r1
	ACIRBlackBoxOpMultiScalarMul
	ACIRBlackBoxOpEmbeddedCurveAdd
	ACIRBlackBoxOpBigIntAdd
	ACIRBlackBoxOpBigIntSub
	ACIRBlackBoxOpBigIntMul
	ACIRBlackBoxOpBigIntDiv
	ACIRBlackBoxOpBigIntFromLeBytes
	ACIRBlackBoxOpBigIntToLeBytes
	ACIRBlackBoxOpPoseidon2Permutation
	ACIRBlackBoxOpSha256Compression
	ACIRBlackBoxOpToRadix
)

func (b *BlackBoxOpType) UnmarshalReader(r io.Reader) error {
	if err := binary.Read(r, binary.LittleEndian, b); err != nil {
		return err
	}

	if *b > ACIRBlackBoxOpToRadix {
		return fmt.Errorf("invalid BlackBoxOpType: %d", *b)
	}

	return nil
}

func (bbo *BlackBoxOp) UnmarshalReader(r io.Reader) error {
	if err := bbo.OpType.UnmarshalReader(r); err != nil {
		return err
	}

	switch bbo.OpType {
	case ACIRBlackBoxOpAES128Encrypt:
		bbo.AES128Encrypt = &AES128Encrypt{}
		return bbo.AES128Encrypt.UnmarshalReader(r)
	case ACIRBlackBoxOpBlake2s:
		bbo.Blake2s = &Blake2s{}
		return bbo.Blake2s.UnmarshalReader(r)
	case ACIRBlackBoxOpBlake3:
		bbo.Blake3 = &Blake3{}
		return bbo.Blake3.UnmarshalReader(r)
	case ACIRBlackBoxOpKeccakf1600:
		bbo.Keccakf1600 = &Keccakf1600{}
		return bbo.Keccakf1600.UnmarshalReader(r)
	case ACIRBlackBoxOpEcdsaSecp256k1:
		bbo.EcdsaSecp256k1 = &EcdsaSecp256k1{}
		return bbo.EcdsaSecp256k1.UnmarshalReader(r)
	case ACIRBlackBoxOpEcdsaSecp256r1:
		bbo.EcdsaSecp256r1 = &EcdsaSecp256r1{}
		return bbo.EcdsaSecp256r1.UnmarshalReader(r)
	case ACIRBlackBoxOpMultiScalarMul:
		bbo.MultiScalarMul = &MultiScalarMul{}
		return bbo.MultiScalarMul.UnmarshalReader(r)
	case ACIRBlackBoxOpEmbeddedCurveAdd:
		bbo.EmbeddedCurveAdd = &EmbeddedCurveAdd{}
		return bbo.EmbeddedCurveAdd.UnmarshalReader(r)
	case ACIRBlackBoxOpBigIntAdd:
		bbo.BigIntAdd = &BigIntAdd{}
		return bbo.BigIntAdd.UnmarshalReader(r)
	case ACIRBlackBoxOpBigIntSub:
		bbo.BigIntSub = &BigIntSub{}
		return bbo.BigIntSub.UnmarshalReader(r)
	case ACIRBlackBoxOpBigIntMul:
		bbo.BigIntMul = &BigIntMul{}
		return bbo.BigIntMul.UnmarshalReader(r)
	case ACIRBlackBoxOpBigIntDiv:
		bbo.BigIntDiv = &BigIntDiv{}
		return bbo.BigIntDiv.UnmarshalReader(r)
	case ACIRBlackBoxOpBigIntFromLeBytes:
		bbo.BigIntFromLeBytes = &BigIntFromLEBytes{}
		return bbo.BigIntFromLeBytes.UnmarshalReader(r)
	case ACIRBlackBoxOpBigIntToLeBytes:
		bbo.BigIntToLeBytes = &BigIntToLEBytes{}
		return bbo.BigIntToLeBytes.UnmarshalReader(r)
	case ACIRBlackBoxOpPoseidon2Permutation:
		bbo.Poseidon2Permutation = &Poseidon2Permutation{}
		return bbo.Poseidon2Permutation.UnmarshalReader(r)
	case ACIRBlackBoxOpSha256Compression:
		bbo.Sha256Compression = &Sha256Compression{}
		return bbo.Sha256Compression.UnmarshalReader(r)
	case ACIRBlackBoxOpToRadix:
		bbo.ToRadix = &ToRadix{}
		return bbo.ToRadix.UnmarshalReader(r)
	default:
		return fmt.Errorf("unknown BlackBoxOpType: %d", bbo.OpType)
	}
}
