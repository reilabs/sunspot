package blackboxfunc

import (
	"math/big"

	"github.com/consensys/gnark/frontend"
)

// isLessOrEqualConstant returns 1 if the unsigned integer encoded by the
// little-endian bits is ≤ c, otherwise 0. bitsLE must be boolean-constrained
// and c must fit in len(bitsLE) bits.
func isLessOrEqualConstant(api frontend.API, bitsLE []frontend.Variable, c *big.Int) frontend.Variable {
	var isLess frontend.Variable = 0
	var isEqual frontend.Variable = 1
	for i := len(bitsLE) - 1; i >= 0; i-- {
		sBit := bitsLE[i]
		if c.Bit(i) == 1 {
			// isEqual*(1-sBit) = isEqual - isEqual*sBit; reuse the product to
			// keep this branch at a single multiplication constraint.
			newIsEqual := api.Mul(isEqual, sBit)
			isLess = api.Add(isLess, api.Sub(isEqual, newIsEqual))
			isEqual = newIsEqual
		} else {
			isEqual = api.Mul(isEqual, api.Sub(1, sBit))
		}
	}
	return api.Add(isLess, isEqual)
}
