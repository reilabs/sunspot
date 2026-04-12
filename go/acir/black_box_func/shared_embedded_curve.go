package blackboxfunc

import (
	"math/big"
	shr "sunspot/go/acir/shared"
	grumpkin "sunspot/go/sw-grumpkin"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
)

var twoTo128 = new(big.Int).Lsh(big.NewInt(1), 128)

// grumpkinScalarModulusHighLimb is floor(p / 2^128) where p is the grumpkin
// scalar modulus — the upper bound for a scalar's high 128-bit limb.
var grumpkinScalarModulusHighLimb = new(big.Int).Rsh(ecc.GRUMPKIN.ScalarField(), 128)

// ScalarFromLimbs recomposes a scalar from its (lo, hi) 128-bit limb
// FunctionInputs into lo + hi * 2^128.
func ScalarFromLimbs[T shr.ACIRField, E constraint.Element](
	api frontend.Builder[E],
	witnesses map[shr.Witness]frontend.Variable,
	lo, hi FunctionInput[T],
) (frontend.Variable, error) {
	scalarLo, err := lo.ToVariable(witnesses)
	if err != nil {
		return nil, err
	}
	scalarHi, err := hi.ToVariable(witnesses)
	if err != nil {
		return nil, err
	}
	api.AssertIsLessOrEqual(scalarHi, grumpkinScalarModulusHighLimb)
	return api.Add(scalarLo, api.Mul(scalarHi, twoTo128)), nil
}

// maskedEmbeddedPoint constrains isInf to a boolean (gated by pred) and returns
// a grumpkin point whose coordinates are masked to (0, 0) when isInf is set.
func maskedEmbeddedPoint[E constraint.Element](
	api frontend.Builder[E],
	pred, x, y, isInf frontend.Variable,
) grumpkin.G1Affine {
	api.AssertIsEqual(frontend.Variable(0), api.Mul(pred, isInf, api.Sub(frontend.Variable(1), isInf)))
	notInf := api.Sub(frontend.Variable(1), isInf)
	return grumpkin.G1Affine{
		X: api.Mul(notInf, x),
		Y: api.Mul(notInf, y),
	}
}

// EmbeddedPointFromInputs resolves an (x, y, is_infinite) triple of
// FunctionInputs into a grumpkin point whose coordinates are masked to (0, 0)
// when is_infinite is set.
func EmbeddedPointFromInputs[T shr.ACIRField, E constraint.Element](
	api frontend.Builder[E],
	witnesses map[shr.Witness]frontend.Variable,
	pred frontend.Variable,
	in [3]FunctionInput[T],
) (grumpkin.G1Affine, error) {
	x, err := in[0].ToVariable(witnesses)
	if err != nil {
		return grumpkin.G1Affine{}, err
	}
	y, err := in[1].ToVariable(witnesses)
	if err != nil {
		return grumpkin.G1Affine{}, err
	}
	isInf, err := in[2].ToVariable(witnesses)
	if err != nil {
		return grumpkin.G1Affine{}, err
	}
	return maskedEmbeddedPoint(api, pred, x, y, isInf), nil
}
