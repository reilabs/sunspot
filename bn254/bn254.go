package bn254

import (
	"encoding/binary"
	"fmt"
	"io"
	"math/big"
	shr "nr-groth16/acir/shared"
	"strings"

	"github.com/consensys/gnark-crypto/ecc/bn254/fp"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/consensys/gnark/frontend"
)

const BN254_MODULUS_STRING = "21888242871839275222246405745257275088548364400416034343698204186575808495617"

var Bn254Modulus, _ = new(big.Int).SetString(BN254_MODULUS_STRING, 10)

type BN254Field struct {
	Modulus big.Int
}

func Zero() BN254Field {
	return BN254Field{
		Modulus: *new(big.Int).SetUint64(0),
	}
}

func One() BN254Field {
	return BN254Field{
		Modulus: *new(big.Int).SetInt64(1),
	}
}

func (b BN254Field) UnmarshalReader(r io.Reader) error {
	// Implement the unmarshalling logic here

	var bn254len uint64
	if err := binary.Read(r, binary.LittleEndian, &bn254len); err != nil {
		return err
	}

	fmt.Println("Read BN254 field bytes:", bn254len)
	bn254Bytes := make([]byte, bn254len)
	if _, err := io.ReadFull(r, bn254Bytes); err != nil {
		return fmt.Errorf("failed to read BN254 field bytes: %w", err)
	}
	fmt.Println("BN254 field bytes read successfully:", bn254Bytes)

	fmt.Println("Creating new BN254 field modulus...")

	fmt.Println("Setting BN254 field modulus from bytes...", bn254Bytes)
	str := string(bn254Bytes)
	fmt.Println("Setting BN254 field modulus from bytes as string...", str)

	// if b == nil {
	// 	fmt.Println("BN254Field is nil, cannot set modulus")
	// 	return fmt.Errorf("BN254Field is nil, cannot set modulus")
	// }
	if len(str) >= 2 && strings.HasPrefix(str, "0x") {
		fmt.Println("Setting BN254 field modulus from hex string:", str[2:])
		if _, ok := b.Modulus.SetString(str[2:], 16); !ok {
			return fmt.Errorf("failed to set BN254 field modulus from hex string: %s", str)
		}
	} else if strings.ContainsAny(str, "abcdefABCDEF") {
		fmt.Println("Setting BN254 field modulus from hex string:", str)
		if _, ok := b.Modulus.SetString(str, 16); !ok {
			return fmt.Errorf("failed to set BN254 field modulus from hex string: %s", str)
		}
	} else {
		fmt.Println("Setting BN254 field modulus from decimal string:", str)
		if _, ok := b.Modulus.SetString(str, 10); !ok {
			return fmt.Errorf("failed to set BN254 field modulus from string: %s", str)
		}
	}
	fmt.Println("BN254 field modulus set to:", b.Modulus)

	//val := new(big.Int).SetBytes(bn254Bytes)
	//b.Modulus = val.Mod(val, Bn254Modulus)

	return nil
}

func (b BN254Field) Equals(other shr.ACIRField) bool {
	return true // Implement the equality check logic here
}

func (b BN254Field) ToElement() shr.GenericFPElement {
	var element fp.Element
	element.SetBigInt(&b.Modulus)
	return shr.GenericFPElement{
		Kind:           shr.GenericFPElementKindBN254,
		BN254FpElement: &element,
	}
}

func (b BN254Field) ToFrontendVariable() frontend.Variable {
	var element fr.Element
	element.SetBigInt(&b.Modulus)
	return element
}
