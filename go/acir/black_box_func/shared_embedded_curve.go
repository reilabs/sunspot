package blackboxfunc

import (
	"math/big"
	shr "github.com/reilabs/sunspot/go/acir/shared"
	grumpkin "github.com/reilabs/sunspot/go/sw-grumpkin"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
)

var twoTo128 = new(big.Int).Lsh(big.NewInt(1), 128)

// grumpkinScalarModulusHighLimb is floor(p / 2^128) where p is the grumpkin
// scalar modulus — the upper bound for a scalar's high 128-bit limb.
var grumpkinScalarModulusHighLimb = new(big.Int).Rsh(ecc.GRUMPKIN.ScalarField(), 128)

var grumpkinScalarModulusLowLimbMax = new(big.Int).Sub(
	new(big.Int).Mod(ecc.GRUMPKIN.ScalarField(), twoTo128),
	big.NewInt(1),
)

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

	// Assert point is less than scalar field modulus
	hiAtMax := api.IsZero(api.Sub(scalarHi, grumpkinScalarModulusHighLimb))
	api.AssertIsLessOrEqual(api.Mul(hiAtMax, scalarLo), grumpkinScalarModulusLowLimbMax)
	return api.Add(scalarLo, api.Mul(scalarHi, twoTo128)), nil
}

// EmbeddedPointFromInputs resolves an (x, y) pair of FunctionInputs into a
// grumpkin point.
func EmbeddedPointFromInputs[T shr.ACIRField](
	x FunctionInput[T],
	y FunctionInput[T],
	witnesses map[shr.Witness]frontend.Variable,
) (grumpkin.G1Affine, error) {
	xVar, err := x.ToVariable(witnesses)
	if err != nil {
		return grumpkin.G1Affine{}, err
	}
	yVar, err := y.ToVariable(witnesses)
	if err != nil {
		return grumpkin.G1Affine{}, err
	}
	return grumpkin.G1Affine{X: xVar, Y: yVar}, nil
}
