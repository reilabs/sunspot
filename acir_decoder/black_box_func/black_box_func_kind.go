package blackboxfunc

import (
	"encoding/binary"
	"io"
)

type BlackBoxFuncKind int

const (
	ACIRBlackBoxFuncKindAES128Encrypt BlackBoxFuncKind = iota
	ACIRBlackBoxFuncKindAnd
	ACIRBlackBoxFuncKindXor
	ACIRBlackBoxFuncKindRange
	ACIRBlackBoxFuncKindBlake2s
	ACIRBlackBoxFuncKindBlake3
	ACIRBlackBoxFuncKindEcdsaSecp256k1
	ACIRBlackBoxFuncKindEcdsaSecp256r1
	ACIRBlackBoxFuncKindMultiScalarMul
	ACIRBlackBoxFuncKindEmbeddedCurveAdd
	ACIRBlackBoxFuncKindKeccakf1600
	ACIRBlackBoxFuncKindRecursiveAggregation
	ACIRBlackBoxFuncKindBigIntAdd
	ACIRBlackBoxFuncKindBigIntSub
	ACIRBlackBoxFuncKindBigIntMul
	ACIRBlackBoxFuncKindBigIntDiv
	ACIRBlackBoxFuncKindBigIntFromLeBytes
	ACIRBlackBoxFuncKindBigIntToLeBytes
	ACIRBlackBoxFuncKindPoseidon2Permutation
	ACIRBlackBoxFuncKindSha256Compression
)

func (k *BlackBoxFuncKind) UnmarshalReader(r io.Reader) error {
	if err := binary.Read(r, binary.LittleEndian, k); err != nil {
		return err
	}
	return nil
}
