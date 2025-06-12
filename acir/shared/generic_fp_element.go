package shared

import (
	"fmt"

	bls12_377_fp "github.com/consensys/gnark-crypto/ecc/bls12-377/fp"
	bls12_381_fp "github.com/consensys/gnark-crypto/ecc/bls12-381/fp"
	bls24_315_fp "github.com/consensys/gnark-crypto/ecc/bls24-315/fp"
	bls24_317_fp "github.com/consensys/gnark-crypto/ecc/bls24-317/fp"
	bn254_fp "github.com/consensys/gnark-crypto/ecc/bn254/fp"
	bw6_633_fp "github.com/consensys/gnark-crypto/ecc/bw6-633/fp"
	bw6_761_fp "github.com/consensys/gnark-crypto/ecc/bw6-761/fp"
	grumpkin_fp "github.com/consensys/gnark-crypto/ecc/grumpkin/fp"
	secp256k1_fp "github.com/consensys/gnark-crypto/ecc/secp256k1/fp"
	stark_curve_fp "github.com/consensys/gnark-crypto/ecc/stark-curve/fp"
	"github.com/consensys/gnark/frontend"
)

type GenericFPElement struct {
	Kind                GenericFPElementKind
	BLS12_377FpElement  *bls12_377_fp.Element
	BLS12_381FpElement  *bls12_381_fp.Element
	BLS24_315FpElement  *bls24_315_fp.Element
	BLS24_317FpElement  *bls24_317_fp.Element
	BN254FpElement      *bn254_fp.Element
	BW6_633FpElement    *bw6_633_fp.Element
	BW6_761FpElement    *bw6_761_fp.Element
	GrumpkinFpElement   *grumpkin_fp.Element
	Secp256k1FpElement  *secp256k1_fp.Element
	StarkCurveFpElement *stark_curve_fp.Element
}

type GenericFPElementKind uint32

const (
	GenericFPElementKindBLS12_377 GenericFPElementKind = iota
	GenericFPElementKindBLS12_381
	GenericFPElementKindBLS24_315
	GenericFPElementKindBLS24_317
	GenericFPElementKindBN254
	GenericFPElementKindBW6_633
	GenericFPElementKindBW6_761
	GenericFPElementKindGrumpkin
	GenericFPElementKindSecp256k1
	GenericFPElementKindStarkCurve
)

func (g *GenericFPElement) Equals(other GenericFPElement) bool {
	if g.Kind != other.Kind {
		return false
	}

	switch g.Kind {
	case GenericFPElementKindBLS12_377:
		return g.BLS12_377FpElement.Equal(other.BLS12_377FpElement)
	case GenericFPElementKindBLS12_381:
		return g.BLS12_381FpElement.Equal(other.BLS12_381FpElement)
	case GenericFPElementKindBLS24_315:
		return g.BLS24_315FpElement.Equal(other.BLS24_315FpElement)
	case GenericFPElementKindBLS24_317:
		return g.BLS24_317FpElement.Equal(other.BLS24_317FpElement)
	case GenericFPElementKindBN254:
		return g.BN254FpElement.Equal(other.BN254FpElement)
	case GenericFPElementKindBW6_633:
		return g.BW6_633FpElement.Equal(other.BW6_633FpElement)
	case GenericFPElementKindBW6_761:
		return g.BW6_761FpElement.Equal(other.BW6_761FpElement)
	case GenericFPElementKindGrumpkin:
		return g.GrumpkinFpElement.Equal(other.GrumpkinFpElement)
	case GenericFPElementKindSecp256k1:
		return g.Secp256k1FpElement.Equal(other.Secp256k1FpElement)
	case GenericFPElementKindStarkCurve:
		return g.StarkCurveFpElement.Equal(other.StarkCurveFpElement)
	default:
		return false
	}
}

func (g GenericFPElement) ToFrontendVariable() frontend.Variable {
	switch g.Kind {
	case GenericFPElementKindBLS12_377:
		return g.BLS12_377FpElement
	case GenericFPElementKindBLS12_381:
		return g.BLS12_381FpElement
	case GenericFPElementKindBLS24_315:
		return g.BLS24_315FpElement
	case GenericFPElementKindBLS24_317:
		return g.BLS24_317FpElement
	case GenericFPElementKindBN254:
		return g.BN254FpElement
	case GenericFPElementKindBW6_633:
		return g.BW6_633FpElement
	case GenericFPElementKindBW6_761:
		return g.BW6_761FpElement
	case GenericFPElementKindGrumpkin:
		return g.GrumpkinFpElement
	case GenericFPElementKindSecp256k1:
		return g.Secp256k1FpElement
	case GenericFPElementKindStarkCurve:
		return g.StarkCurveFpElement
	default:
		panic(fmt.Sprintf("Unknown GenericFPElementKind: %d", g.Kind))
	}
}
