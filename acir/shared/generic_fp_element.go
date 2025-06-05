package shared

import (
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
