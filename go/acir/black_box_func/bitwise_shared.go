package blackboxfunc

import (
	"errors"
	"fmt"
	"math/big"
	shr "sunspot/go/acir/shared"

	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/constraint/solver"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/math/uints"
	"github.com/consensys/gnark/std/rangecheck"
	"github.com/google/btree"
)

// defineBitwise emits the constraints for a bitwise op (AND, XOR, ...) between
// lhs and rhs with result output, where each operand fits in nBits bits. For
// nBits <= 64 the operands are range-checked and the op is asserted directly;
// for 64 < nBits <= 128 the values are split into two 64-bit limbs and the op
// is asserted limb-wise.
func defineBitwise[T shr.ACIRField, E constraint.Element](
	api frontend.Builder[E],
	uapi *uints.BinaryField[uints.U64],
	witnesses map[shr.Witness]frontend.Variable,
	lhs, rhs FunctionInput[T],
	output shr.Witness,
	nBits int,
	op func(a ...uints.U64) uints.U64,
) error {
	lhsV, err := lhs.ToVariable(witnesses)
	if err != nil {
		return err
	}
	rhsV, err := rhs.ToVariable(witnesses)
	if err != nil {
		return err
	}
	outV, ok := witnesses[output]
	if !ok {
		return fmt.Errorf("witness %d not found in witnesses map", output)
	}

	if nBits <= 64 {
		rc := rangecheck.New(api)
		rc.Check(lhsV, nBits)
		rc.Check(rhsV, nBits)
		uapi.AssertEq(uapi.ValueOf(outV), op(uapi.ValueOf(lhsV), uapi.ValueOf(rhsV)))
		return nil
	}

	lhsLimbs, err := splitToU64Limbs(api, uapi, lhsV, nBits)
	if err != nil {
		return err
	}
	rhsLimbs, err := splitToU64Limbs(api, uapi, rhsV, nBits)
	if err != nil {
		return err
	}
	outLimbs, err := splitToU64Limbs(api, uapi, outV, nBits)
	if err != nil {
		return err
	}
	for i := range lhsLimbs {
		uapi.AssertEq(outLimbs[i], op(lhsLimbs[i], rhsLimbs[i]))
	}
	return nil
}

// fillBitwiseWitnessTree inserts the witness indices referenced by a bitwise
// black-box call (lhs, rhs if they are witnesses, and output) into tree.
func fillBitwiseWitnessTree[T shr.ACIRField](
	tree *btree.BTree,
	index uint32,
	lhs, rhs FunctionInput[T],
	output shr.Witness,
) bool {
	if tree == nil {
		return false
	}
	if lhs.IsWitness() {
		tree.ReplaceOrInsert(*lhs.Witness + shr.Witness(index))
	}
	if rhs.IsWitness() {
		tree.ReplaceOrInsert(*rhs.Witness + shr.Witness(index))
	}
	tree.ReplaceOrInsert(output + shr.Witness(index))
	return true
}

var twoTo64 = new(big.Int).Lsh(big.NewInt(1), 64)

func init() {
	solver.RegisterHint(splitInto64BitLimbsHint)
}

// splitInto64BitLimbsHint outputs [v mod 2^64, v >> 64].
func splitInto64BitLimbsHint(_ *big.Int, inputs, outputs []*big.Int) error {
	if len(inputs) != 1 || len(outputs) != 2 {
		return errors.New("splitInto64BitLimbsHint: expected 1 input and 2 outputs")
	}
	mask := new(big.Int).SetUint64(^uint64(0))
	outputs[0].And(inputs[0], mask)
	outputs[1].Rsh(inputs[0], 64)
	return nil
}

// splitToU64Limbs splits v into [lo, hi] U64 limbs such that v == lo + hi * 2^64
// and v fits in nBits (64 < nBits <= 128). The hint supplies the split; ValueOf
// range-checks each limb to 64 bits, and a tighter check on hi enforces nBits
// when nBits < 128.
func splitToU64Limbs[E constraint.Element](
	api frontend.Builder[E],
	uapi *uints.BinaryField[uints.U64],
	v frontend.Variable,
	nBits int,
) ([2]uints.U64, error) {
	limbs, err := api.NewHint(splitInto64BitLimbsHint, 2, v)
	if err != nil {
		return [2]uints.U64{}, err
	}
	lo, hi := limbs[0], limbs[1]
	if nBits < 128 {
		rangecheck.New(api).Check(hi, nBits-64)
	}
	loU64 := uapi.ValueOf(lo)
	hiU64 := uapi.ValueOf(hi)
	api.AssertIsEqual(v, api.Add(lo, api.Mul(hi, twoTo64)))
	return [2]uints.U64{loU64, hiU64}, nil
}
