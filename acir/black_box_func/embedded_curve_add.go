package blackboxfunc

import (
	"encoding/binary"
	"fmt"
	"io"
	"math/big"
	shr "nr-groth16/acir/shared"

	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/algebra/emulated/sw_emulated"
	"github.com/consensys/gnark/std/math/bits"
	"github.com/consensys/gnark/std/math/emulated"
	"github.com/consensys/gnark/std/math/emulated/emparams"
	"github.com/google/btree"
)

type EmbeddedCurveAdd[T shr.ACIRField, E constraint.Element] struct {
	Input1  [3]FunctionInput[T]
	Input2  [3]FunctionInput[T]
	Outputs [3]shr.Witness
}

func (a *EmbeddedCurveAdd[T, E]) UnmarshalReader(r io.Reader) error {
	for i := 0; i < 3; i++ {
		if err := a.Input1[i].UnmarshalReader(r); err != nil {
			return err
		}
	}

	for i := 0; i < 3; i++ {
		if err := a.Input2[i].UnmarshalReader(r); err != nil {
			return err
		}
	}

	if err := binary.Read(r, binary.LittleEndian, &a.Outputs); err != nil {
		return err
	}

	return nil
}

func (a *EmbeddedCurveAdd[T, E]) Equals(other BlackBoxFunction[E]) bool {
	value, ok := other.(*EmbeddedCurveAdd[T, E])
	if !ok || len(a.Input1) != len(value.Input1) || len(a.Input2) != len(value.Input2) {
		return false
	}

	for i := 0; i < 3; i++ {
		if !a.Input1[i].Equals(&value.Input1[i]) || !a.Input2[i].Equals(&value.Input2[i]) {
			return false
		}
	}

	for i := 0; i < 3; i++ {
		if a.Outputs[i] != value.Outputs[i] {
			return false
		}
	}

	return true
}

func (a *EmbeddedCurveAdd[T, E]) Define(api frontend.Builder[E], witnesses map[shr.Witness]frontend.Variable) error {
	gy, ok := new(big.Int).SetString("17631683881184975370165255887551781615748388533673675138860", 10)
	if !ok {
		return fmt.Errorf("unable to initialise generator y value for grumpkin curve")
	}
	grumpkin_params := sw_emulated.CurveParams{
		A:  big.NewInt(0),
		B:  big.NewInt(-17),
		Gx: big.NewInt(1),
		Gy: gy,
	}
	// Grumpkin prime field  = BN254Fr
	// Grumpkin scalar field = BN254Fp
	curve, err := sw_emulated.New[emparams.BN254Fr, emparams.BN254Fp](api, grumpkin_params)
	if err != nil {
		return fmt.Errorf("new curve: %w", err)
	}
	primeField, err := emulated.NewField[emparams.BN254Fr](api)
	if err != nil {
		return err
	}

	point1X, err := a.Input1[0].ToVariable(witnesses)
	if err != nil {
		return err
	}

	point1Y, err := a.Input1[1].ToVariable(witnesses)
	if err != nil {
		return err
	}
	point2X, err := a.Input2[0].ToVariable(witnesses)
	if err != nil {
		return err
	}
	point2Y, err := a.Input2[1].ToVariable(witnesses)
	if err != nil {
		return err
	}

	x := sw_emulated.AffinePoint[emparams.BN254Fr]{
		X: *primeField.NewElement(DecomposeTo4x64(api, point1X)),
		Y: *primeField.NewElement(DecomposeTo4x64(api, point1Y)),
	}
	api.Println(x.X.Limbs...)
	api.Println(x.Y.Limbs...)
	y := sw_emulated.AffinePoint[emparams.BN254Fr]{
		X: *primeField.NewElement(DecomposeTo4x64(api, point2X)),
		Y: *primeField.NewElement(DecomposeTo4x64(api, point2Y)),
	}
	api.Println(y.X.Limbs...)
	api.Println(y.Y.Limbs...)
	z := sw_emulated.AffinePoint[emparams.BN254Fr]{
		X: *primeField.NewElement(DecomposeTo4x64(api, witnesses[a.Outputs[0]])),
		Y: *primeField.NewElement(DecomposeTo4x64(api, witnesses[a.Outputs[1]])),
	}
	api.Println(z.X.Limbs...)
	curve.AssertIsOnCurve(&x)
	curve.AssertIsOnCurve(&y)
	curve.AssertIsOnCurve(&z)
	api.Println(z.X.Limbs...)
	curve.AssertIsEqual(&z, curve.AddUnified(&x, &y))

	return nil
}

func (a *EmbeddedCurveAdd[T, E]) FillWitnessTree(tree *btree.BTree) bool {
	if tree == nil {
		return false
	}
	for _, input := range a.Input1 {

		tree.ReplaceOrInsert(*input.Witness)
	}

	for _, input := range a.Input2 {
		tree.ReplaceOrInsert(*input.Witness)
	}

	for _, output := range a.Outputs {
		tree.ReplaceOrInsert(output)
	}
	return true
}

// DecomposeTo4x64 decomposes x into 4 little-endian 64-bit limbs:
func DecomposeTo4x64(api frontend.API, x frontend.Variable) []frontend.Variable {
	const bitsPerLimb = 64
	const nbLimbs = 4
	limbs := make([]frontend.Variable, nbLimbs)
	// Get 256 bits (little-endian) representing x
	allBits := bits.ToBinary(api, x, bits.WithNbDigits(bitsPerLimb*nbLimbs))
	for i := 0; i < nbLimbs; i++ {
		start := i * bitsPerLimb
		end := start + bitsPerLimb
		chunk := allBits[start:end]
		limbs[i] = bits.FromBinary(api, chunk)
	}
	return limbs
}
