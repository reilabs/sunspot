package blackboxfunc

import (
	"encoding/binary"
	"fmt"
	"io"
	shr "nr-groth16/acir/shared"

	"github.com/consensys/gnark/frontend"
	"github.com/rs/zerolog/log"
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

	log.Trace().Msgf("Unmarshalling black box function call kind %v", a.Kind)

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
	case ACIRBlackBoxFuncKindAES128Encrypt:
		return a.AES128Encrypt.Define(api, witnesses)
	case ACIRBlackBoxFuncKindAnd:
		return a.And.Define(api, witnesses)
	case ACIRBlackBoxFuncKindXor:
		return a.Xor.Define(api, witnesses)
	case ACIRBlackBoxFuncKindRange:
		return a.Range.Define(api, witnesses)
	case ACIRBlackBoxFuncKindBlake2s:
		return fmt.Errorf("Blake2s is not implemented yet")
	case ACIRBlackBoxFuncKindBlake3:
		return fmt.Errorf("Blake3 is not implemented yet")
	case ACIRBlackBoxFuncKindEcdsaSecp256k1:
		return a.ECDSASECP256K1.Define(api, witnesses)
	case ACIRBlackBoxFuncKindEcdsaSecp256r1:
		return a.ECDSASECP256R1.Define(api, witnesses)
	case ACIRBlackBoxFuncKindMultiScalarMul:
		return fmt.Errorf("MultiScalarMul is not implemented yet")
	case ACIRBlackBoxFuncKindEmbeddedCurveAdd:
		return fmt.Errorf("EmbeddedCurveAdd is not implemented yet")
	case ACIRBlackBoxFuncKindKeccakf1600:
		return fmt.Errorf("Keccakf1600 is not implemented yet")
	case ACIRBlackBoxFuncKindRecursiveAggregation:
		return fmt.Errorf("RecursiveAggregation is not implemented yet")
	case ACIRBlackBoxFuncKindPoseidon2Permutation:
		return a.Poseidon2Permutation.Define(api, witnesses)
	case ACIRBlackBoxFuncKindSha256Compression:
		return a.Sha256Compression.Define(api, witnesses)
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
