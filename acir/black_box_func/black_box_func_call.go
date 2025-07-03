package blackboxfunc

import (
	"encoding/binary"
	"fmt"
	"io"
	shr "nr-groth16/acir/shared"

	"github.com/consensys/gnark/frontend"
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
	if err := binary.Read(r, binary.LittleEndian, &a.Kind); err != nil {
		return err
	}

	if a.Kind > ACIRBlackBoxFuncKindSha256Compression {
		return BlackBoxFuncKindError{
			Code: uint32(a.Kind),
		}
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

func (a BlackBoxFuncCall[T]) Define(api frontend.API, witnesses map[shr.Witness]frontend.Variable) error {
	switch a.Kind {
	/*case ACIRBlackBoxFuncKindAES128Encrypt:
	return a.AES128Encrypt.Define(api, witnesses)*/
	case ACIRBlackBoxFuncKindAnd:
		return a.And.Define(api, witnesses)
	case ACIRBlackBoxFuncKindXor:
		return a.Xor.Define(api, witnesses)
	case ACIRBlackBoxFuncKindRange:
		return a.Range.Define(api, witnesses)
	/*case ACIRBlackBoxFuncKindBlake2s:
		return a.Blake2s.Define(api, witnesses)
	case ACIRBlackBoxFuncKindBlake3:
		return a.Blake3.Define(api, witnesses)
	case ACIRBlackBoxFuncKindEcdsaSecp256k1:
		return a.ECDSASECP256K1.Define(api, witnesses)
	case ACIRBlackBoxFuncKindEcdsaSecp256r1:
		return a.ECDSASECP256R1.Define(api, witnesses)
	case ACIRBlackBoxFuncKindMultiScalarMul:
		return a.MultiScalarMul.Define(api, witnesses)
	case ACIRBlackBoxFuncKindEmbeddedCurveAdd:
		return a.EmbeddedCurveAdd.Define(api, witnesses)
	case ACIRBlackBoxFuncKindKeccakf1600:
		return a.Keccakf1600.Define(api, witnesses)
	case ACIRBlackBoxFuncKindRecursiveAggregation:
		return a.RecursiveAggregation.Define(api, witnesses)
	case ACIRBlackBoxFuncKindBigIntAdd:
		return a.BigIntAdd.Define(api, witnesses)
	case ACIRBlackBoxFuncKindBigIntSub:
		return a.BigIntSub.Define(api, witnesses)
	case ACIRBlackBoxFuncKindBigIntMul:
		return a.BigIntMul.Define(api, witnesses)
	case ACIRBlackBoxFuncKindBigIntDiv:
		return a.BigIntDiv.Define(api, witnesses)
	case ACIRBlackBoxFuncKindBigIntFromLeBytes:
		return a.BigIntFromLEBytes.Define(api, witnesses)
	case ACIRBlackBoxFuncKindBigIntToLeBytes:
		return a.BigIntToLEBytes.Define(api, witnesses)
	case ACIRBlackBoxFuncKindPoseidon2Permutation:
		return a.Poseidon2Permutation.Define(api, witnesses)
	case ACIRBlackBoxFuncKindSha256Compression:
		return a.Sha256Compression.Define(api, witnesses)*/
	default:
		return fmt.Errorf("unknown black box function kind: %d", a.Kind)
	}
}

func (a BlackBoxFuncCall[T]) Equals(other BlackBoxFuncCall[T]) bool {
	if a.Kind != other.Kind {
		fmt.Println("BlackBoxFuncCall: Kind mismatch")
		return false
	}

	switch a.Kind {
	case ACIRBlackBoxFuncKindAES128Encrypt:
		if a.AES128Encrypt == nil || other.AES128Encrypt == nil {
			return a.AES128Encrypt == nil && other.AES128Encrypt == nil
		}
		return a.AES128Encrypt.Equals(other.AES128Encrypt)
	case ACIRBlackBoxFuncKindAnd:
		if a.And == nil || other.And == nil {
			return a.And == nil && other.And == nil
		}
		return a.And.Equals(other.And)
	case ACIRBlackBoxFuncKindXor:
		if a.Xor == nil || other.Xor == nil {
			return a.Xor == nil && other.Xor == nil
		}
		return a.Xor.Equals(other.Xor)
	case ACIRBlackBoxFuncKindRange:
		if a.Range == nil || other.Range == nil {
			return a.Range == nil && other.Range == nil
		}
		return a.Range.Equals(*other.Range)
	case ACIRBlackBoxFuncKindBlake2s:
		if a.Blake2s == nil || other.Blake2s == nil {
			return a.Blake2s == nil && other.Blake2s == nil
		}
		return a.Blake2s.Equals(other.Blake2s)
	case ACIRBlackBoxFuncKindBlake3:
		if a.Blake3 == nil || other.Blake3 == nil {
			return a.Blake3 == nil && other.Blake3 == nil
		}
		return a.Blake3.Equals(other.Blake3)
	case ACIRBlackBoxFuncKindEcdsaSecp256k1:
		if a.ECDSASECP256K1 == nil || other.ECDSASECP256K1 == nil {
			return a.ECDSASECP256K1 == nil && other.ECDSASECP256K1 == nil
		}
		return a.ECDSASECP256K1.Equals(other.ECDSASECP256K1)
	case ACIRBlackBoxFuncKindEcdsaSecp256r1:
		if a.ECDSASECP256R1 == nil || other.ECDSASECP256R1 == nil {
			return a.ECDSASECP256R1 == nil && other.ECDSASECP256R1 == nil
		}
		return a.ECDSASECP256R1.Equals(other.ECDSASECP256R1)
	case ACIRBlackBoxFuncKindMultiScalarMul:
		if a.MultiScalarMul == nil || other.MultiScalarMul == nil {
			return a.MultiScalarMul == nil && other.MultiScalarMul == nil
		}
		return a.MultiScalarMul.Equals(other.MultiScalarMul)
	case ACIRBlackBoxFuncKindEmbeddedCurveAdd:
		if a.EmbeddedCurveAdd == nil || other.EmbeddedCurveAdd == nil {
			return a.EmbeddedCurveAdd == nil && other.EmbeddedCurveAdd == nil
		}
		return a.EmbeddedCurveAdd.Equals(other.EmbeddedCurveAdd)
	case ACIRBlackBoxFuncKindKeccakf1600:
		if a.Keccakf1600 == nil || other.Keccakf1600 == nil {
			return a.Keccakf1600 == nil && other.Keccakf1600 == nil
		}
		return a.Keccakf1600.Equals(other.Keccakf1600)
	case ACIRBlackBoxFuncKindRecursiveAggregation:
		if a.RecursiveAggregation == nil || other.RecursiveAggregation == nil {
			return a.RecursiveAggregation == nil && other.RecursiveAggregation == nil
		}
		return a.RecursiveAggregation.Equals(other.RecursiveAggregation)
	case ACIRBlackBoxFuncKindBigIntAdd:
		if a.BigIntAdd == nil || other.BigIntAdd == nil {
			return a.BigIntAdd == nil && other.BigIntAdd == nil
		}
		return a.BigIntAdd.Equals(other.BigIntAdd)
	case ACIRBlackBoxFuncKindBigIntSub:
		if a.BigIntSub == nil || other.BigIntSub == nil {
			return a.BigIntSub == nil && other.BigIntSub == nil
		}
		return a.BigIntSub.Equals(other.BigIntSub)
	case ACIRBlackBoxFuncKindBigIntMul:
		if a.BigIntMul == nil || other.BigIntMul == nil {
			return a.BigIntMul == nil && other.BigIntMul == nil
		}
		return a.BigIntMul.Equals(other.BigIntMul)
	case ACIRBlackBoxFuncKindBigIntDiv:
		if a.BigIntDiv == nil || other.BigIntDiv == nil {
			return a.BigIntDiv == nil && other.BigIntDiv == nil
		}
		return a.BigIntDiv.Equals(other.BigIntDiv)
	case ACIRBlackBoxFuncKindBigIntFromLeBytes:
		if a.BigIntFromLEBytes == nil || other.BigIntFromLEBytes == nil {
			return a.BigIntFromLEBytes == nil && other.BigIntFromLEBytes == nil
		}
		return a.BigIntFromLEBytes.Equals(other.BigIntFromLEBytes)
	case ACIRBlackBoxFuncKindBigIntToLeBytes:
		if a.BigIntToLEBytes == nil || other.BigIntToLEBytes == nil {
			return a.BigIntToLEBytes == nil && other.BigIntToLEBytes == nil
		}
		return a.BigIntToLEBytes.Equals(other.BigIntToLEBytes)
	case ACIRBlackBoxFuncKindPoseidon2Permutation:
		if a.Poseidon2Permutation == nil || other.Poseidon2Permutation == nil {
			return a.Poseidon2Permutation == nil && other.Poseidon2Permutation == nil
		}
		return a.Poseidon2Permutation.Equals(other.Poseidon2Permutation)
	case ACIRBlackBoxFuncKindSha256Compression:
		if a.Sha256Compression == nil || other.Sha256Compression == nil {
			return a.Sha256Compression == nil && other.Sha256Compression == nil
		}
		return a.Sha256Compression.Equals(other.Sha256Compression)
	default:
		return false
	}
}

type BlackBoxFuncKindError struct {
	Code uint32
}

func (e BlackBoxFuncKindError) Error() string {
	return "unknown black box function kind"
}
