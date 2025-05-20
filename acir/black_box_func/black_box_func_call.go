package blackboxfunc

import (
	"io"
	shr "nr-groth16/acir/shared"
)

type BlackBoxFuncCall[T shr.ACIRField] struct {
	Kind                 BlackBoxFuncKind
	AES128Encrypt        *AES128Encrypt[T]
	And                  *And[T]
	Xor                  *Xor[T]
	Range                *Range[T]
	Blake2s              *Blake2s[T]
	Blake3               *Blake3[T]
	ECDSASECP256K1       *ECDSASECP256K1[T]
	ECDSASECP256R1       *ECDSASECP256R1[T]
	MultiScalarMul       *MultiScalarMul[T]
	EmbeddedCurveAdd     *EmbeddedCurveAdd[T]
	Keccakf1600          *Keccakf1600[T]
	RecursiveAggregation *RecursiveAggregation[T]
	BigIntAdd            *BigIntAdd
	BigIntSub            *BigIntSub
	BigIntMul            *BigIntMul
	BigIntDiv            *BigIntDiv
	BigIntFromLEBytes    *BigIntFromLEBytes[T]
	BigIntToLEBytes      *BigIntToLEBytes
	Poseidon2Permutation *Poseidon2Permutation[T]
	Sha256Compression    *SHA256Compression[T]
}

func (a *BlackBoxFuncCall[T]) UnmarshalReader(r io.Reader) error {
	if err := a.Kind.UnmarshalReader(r); err != nil {
		return err
	}

	switch a.Kind {
	case ACIRBlackBoxFuncKindAES128Encrypt:
		a.AES128Encrypt = &AES128Encrypt[T]{}
		if err := a.AES128Encrypt.UnmarshalReader(r); err != nil {
			return err
		}
	case ACIRBlackBoxFuncKindAnd:
		a.And = &And[T]{}
		if err := a.And.UnmarshalReader(r); err != nil {
			return err
		}
	case ACIRBlackBoxFuncKindXor:
		a.Xor = &Xor[T]{}
		if err := a.Xor.UnmarshalReader(r); err != nil {
			return err
		}
	case ACIRBlackBoxFuncKindRange:
		a.Range = &Range[T]{}
		if err := a.Range.UnmarshalReader(r); err != nil {
			return err
		}
	case ACIRBlackBoxFuncKindBlake2s:
		a.Blake2s = &Blake2s[T]{}
		if err := a.Blake2s.UnmarshalReader(r); err != nil {
			return err
		}
	case ACIRBlackBoxFuncKindBlake3:
		a.Blake3 = &Blake3[T]{}
		if err := a.Blake3.UnmarshalReader(r); err != nil {
			return err
		}
	case ACIRBlackBoxFuncKindEcdsaSecp256k1:
		a.ECDSASECP256K1 = &ECDSASECP256K1[T]{}
		if err := a.ECDSASECP256K1.UnmarshalReader(r); err != nil {
			return err
		}
	case ACIRBlackBoxFuncKindEcdsaSecp256r1:
		a.ECDSASECP256R1 = &ECDSASECP256R1[T]{}
		if err := a.ECDSASECP256R1.UnmarshalReader(r); err != nil {
			return err
		}
	case ACIRBlackBoxFuncKindMultiScalarMul:
		a.MultiScalarMul = &MultiScalarMul[T]{}
		if err := a.MultiScalarMul.UnmarshalReader(r); err != nil {
			return err
		}
	case ACIRBlackBoxFuncKindEmbeddedCurveAdd:
		a.EmbeddedCurveAdd = &EmbeddedCurveAdd[T]{}
		if err := a.EmbeddedCurveAdd.UnmarshalReader(r); err != nil {
			return err
		}
	case ACIRBlackBoxFuncKindKeccakf1600:
		a.Keccakf1600 = &Keccakf1600[T]{}
		if err := a.Keccakf1600.UnmarshalReader(r); err != nil {
			return err
		}
	case ACIRBlackBoxFuncKindRecursiveAggregation:
		a.RecursiveAggregation = &RecursiveAggregation[T]{}
		if err := a.RecursiveAggregation.UnmarshalReader(r); err != nil {
			return err
		}
	case ACIRBlackBoxFuncKindBigIntAdd:
		a.BigIntAdd = &BigIntAdd{}
		if err := a.BigIntAdd.UnmarshalReader(r); err != nil {
			return err
		}
	case ACIRBlackBoxFuncKindBigIntSub:
		a.BigIntSub = &BigIntSub{}
		if err := a.BigIntSub.UnmarshalReader(r); err != nil {
			return err
		}
	case ACIRBlackBoxFuncKindBigIntMul:
		a.BigIntMul = &BigIntMul{}
		if err := a.BigIntMul.UnmarshalReader(r); err != nil {
			return err
		}
	case ACIRBlackBoxFuncKindBigIntDiv:
		a.BigIntDiv = &BigIntDiv{}
		if err := a.BigIntDiv.UnmarshalReader(r); err != nil {
			return err
		}
	case ACIRBlackBoxFuncKindBigIntFromLeBytes:
		a.BigIntFromLEBytes = &BigIntFromLEBytes[T]{}
		if err := a.BigIntFromLEBytes.UnmarshalReader(r); err != nil {
			return err
		}
	case ACIRBlackBoxFuncKindBigIntToLeBytes:
		a.BigIntToLEBytes = &BigIntToLEBytes{}
		if err := a.BigIntToLEBytes.UnmarshalReader(r); err != nil {
			return err
		}
	case ACIRBlackBoxFuncKindPoseidon2Permutation:
		a.Poseidon2Permutation = &Poseidon2Permutation[T]{}
		if err := a.Poseidon2Permutation.UnmarshalReader(r); err != nil {
			return err
		}
	case ACIRBlackBoxFuncKindSha256Compression:
		a.Sha256Compression = &SHA256Compression[T]{}
		if err := a.Sha256Compression.UnmarshalReader(r); err != nil {
			return err
		}
	default:
		return BlackBoxFuncKindError{
			Code: uint32(a.Kind),
		}
	}

	return nil
}

type BlackBoxFuncKindError struct {
	Code uint32
}

func (e BlackBoxFuncKindError) Error() string {
	return "unknown black box function kind"
}
