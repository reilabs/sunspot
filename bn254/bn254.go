// Package deals with bn254 field elements and utility
package bn254

import (
	"encoding/binary"
	"fmt"
	"io"
	"math/big"
	"strings"
	shr "sunpot/acir/shared"

	"github.com/consensys/gnark-crypto/ecc/bn254/fp"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/consensys/gnark/frontend"
)

const BN254_MODULUS_STRING = "21888242871839275222246405745257275088548364400416034343698204186575808495617"

var Bn254Modulus, _ = new(big.Int).SetString(BN254_MODULUS_STRING, 10)

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
	str := string(bn254Bytes)

	if len(str) >= 2 && strings.HasPrefix(str, "0x") {
		if _, ok := b.value.SetString(str[2:], 16); !ok {
			return fmt.Errorf("failed to set BN254 field element from hex string: %s", str)
		}
	} else if strings.ContainsAny(str, "abcdefABCDEF") {
		if _, ok := b.value.SetString(str, 16); !ok {
			return fmt.Errorf("failed to set BN254 element from hex string: %s", str)
		}
	} else {
		if _, ok := b.value.SetString(str, 16); !ok {
			return fmt.Errorf("failed to set BN254 element value from string: %s", str)
		}
	}
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
