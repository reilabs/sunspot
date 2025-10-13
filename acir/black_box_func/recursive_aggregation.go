package blackboxfunc

import (
	"encoding/binary"
	"io"
	shr "nr-groth16/acir/shared"

	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	sw_bn254 "github.com/consensys/gnark/std/algebra/emulated/sw_bn254"
	"github.com/consensys/gnark/std/algebra/emulated/sw_emulated"
	"github.com/consensys/gnark/std/math/bits"
	"github.com/consensys/gnark/std/math/emulated"
	"github.com/consensys/gnark/std/recursion/groth16"
	"github.com/google/btree"
)

type RecursiveAggregation[T shr.ACIRField, E constraint.Element] struct {
	VerificationKey []FunctionInput[T]
	Proof           []FunctionInput[T]
	PublicInputs    []FunctionInput[T]
	KeyHash         FunctionInput[T]
	ProofType       uint32
}

func (a *RecursiveAggregation[T, E]) UnmarshalReader(r io.Reader) error {
	var VerificationKeyCount uint64
	if err := binary.Read(r, binary.LittleEndian, &VerificationKeyCount); err != nil {
		return err
	}
	a.VerificationKey = make([]FunctionInput[T], VerificationKeyCount)
	for i := uint64(0); i < VerificationKeyCount; i++ {
		if err := a.VerificationKey[i].UnmarshalReader(r); err != nil {
			return err
		}
	}

	var ProofCount uint64
	if err := binary.Read(r, binary.LittleEndian, &ProofCount); err != nil {
		return err
	}
	a.Proof = make([]FunctionInput[T], ProofCount)
	for i := uint64(0); i < ProofCount; i++ {
		if err := a.Proof[i].UnmarshalReader(r); err != nil {
			return err
		}
	}

	var PublicInputsCount uint64
	if err := binary.Read(r, binary.LittleEndian, &PublicInputsCount); err != nil {
		return err
	}
	a.PublicInputs = make([]FunctionInput[T], PublicInputsCount)
	for i := uint64(0); i < PublicInputsCount; i++ {
		if err := a.PublicInputs[i].UnmarshalReader(r); err != nil {
			return err
		}
	}

	if err := a.KeyHash.UnmarshalReader(r); err != nil {
		return err
	}

	if err := binary.Read(r, binary.LittleEndian, &a.ProofType); err != nil {
		return err
	}

	return nil
}

func (a *RecursiveAggregation[T, E]) Equals(other BlackBoxFunction[E]) bool {
	value, ok := other.(*RecursiveAggregation[T, E])
	if !ok || len(a.VerificationKey) != len(value.VerificationKey) ||
		len(a.Proof) != len(value.Proof) ||
		len(a.PublicInputs) != len(value.PublicInputs) ||
		a.ProofType != value.ProofType {
		return false
	}

	for i := range a.VerificationKey {
		if !a.VerificationKey[i].Equals(&value.VerificationKey[i]) {
			return false
		}
	}

	for i := range a.Proof {
		if !a.Proof[i].Equals(&value.Proof[i]) {
			return false
		}
	}

	for i := range a.PublicInputs {
		if !a.PublicInputs[i].Equals(&value.PublicInputs[i]) {
			return false
		}
	}

	return a.KeyHash.Equals(&value.KeyHash)
}

func (a *RecursiveAggregation[T, E]) Define(api frontend.Builder[E], witnesses map[shr.Witness]frontend.Variable) error {

	e, err := sw_bn254.NewPairing(api)
	if err != nil {
		return err
	}
	curve, err := sw_emulated.New[emulated.BN254Fp, emulated.BN254Fr](api, sw_emulated.GetCurveParams[emulated.BN254Fp]())
	if err != nil {
		panic("initialize new curve")
	}

	proof, err := newProof(api, a.Proof, witnesses)
	if err != nil {
		return err
	}

	vk, err := newVK(api, a.VerificationKey, witnesses)
	if err != nil {
		return err
	}

	witness, err := newWitness(api, a.PublicInputs, witnesses)
	if err != nil {
		return err
	}

	gIC := vk.G1.K[0]

	for i, input := range witness.Public {
		curve.AddUnified(&gIC, curve.ScalarMul(&vk.G1.K[i+1], &input))
	}

	g1_vec := []*sw_bn254.G1Affine{&proof.Ar, &gIC, &proof.Krs}
	g2_vec := []*sw_bn254.G2Affine{&proof.Bs, &vk.G2.GammaNeg, &vk.G2.DeltaNeg}
	qap, err := e.MillerLoop(g1_vec, g2_vec)

	if err != nil {
		return err
	}

	exponent := e.FinalExponentiation(qap)
	e.AssertIsEqual(exponent, &vk.E)

	return nil
}

func (a *RecursiveAggregation[T, E]) FillWitnessTree(tree *btree.BTree) bool {
	for i := range a.VerificationKey {
		if a.VerificationKey[i].IsWitness() {
			tree.ReplaceOrInsert(*a.VerificationKey[i].Witness)
		}
	}

	for i := range a.Proof {
		if a.Proof[i].IsWitness() {
			tree.ReplaceOrInsert(*a.Proof[i].Witness)
		}
	}
	for i := range a.PublicInputs {
		if a.PublicInputs[i].IsWitness() {
			tree.ReplaceOrInsert(*a.PublicInputs[i].Witness)
		}
	}

	if a.KeyHash.IsWitness() {
		tree.ReplaceOrInsert(*a.KeyHash.Witness)
	}
	return tree != nil
}

func VariableTo64BitLimbs[T shr.ACIRField](
	api frontend.API,
	fi FunctionInput[T],
	witnesses map[shr.Witness]frontend.Variable,
) ([]frontend.Variable, error) {
	const bitsPerLimb = 64
	const nbLimbs = 4

	variable, err := fi.ToVariable(witnesses)
	if err != nil {
		return nil, err
	}
	out := make([]frontend.Variable, nbLimbs)

	bit_array := api.ToBinary(variable, 254)

	for i := 0; i < nbLimbs; i++ {
		start := i * bitsPerLimb
		end := start + bitsPerLimb
		chunk := bit_array[start:end]
		out[i] = bits.FromBinary(api, chunk)
	}
	return out, nil
}

func newVK[T shr.ACIRField](api frontend.API, vars []FunctionInput[T], witnesses map[shr.Witness]frontend.Variable) (groth16.VerifyingKey[sw_bn254.G1Affine, sw_bn254.G2Affine, sw_bn254.GTEl], error) {
	vk := groth16.VerifyingKey[sw_bn254.G1Affine, sw_bn254.G2Affine, sw_bn254.GTEl]{}
	g2, err := sw_bn254.NewG2(api)
	if err != nil {
		return vk, err
	}
	e, err := sw_bn254.NewPairing(api)
	if err != nil {
		return vk, err
	}

	alpha, err := newG1(api, vars[0:2], witnesses)
	if err != nil {
		return vk, err
	}
	g2Beta, err := newG2(api, vars[2:6], witnesses)
	if err != nil {
		return vk, err
	}
	pair, err := e.Pair([]*sw_bn254.G1Affine{&alpha}, []*sw_bn254.G2Affine{&g2Beta})
	if err != nil {
		return vk, err
	}
	vk.E = *pair

	g2Gamma, err := newG2(api, vars[6:10], witnesses)
	if err != nil {
		return vk, err
	}
	g2Gamma.P.Y = *g2.Neg(&g2Gamma.P.Y)
	vk.G2.GammaNeg = g2Gamma

	g2Delta, err := newG2(api, vars[10:14], witnesses)
	if err != nil {
		return vk, err
	}
	g2Delta.P.Y = *g2.Neg(&g2Delta.P.Y)
	vk.G2.DeltaNeg = g2Delta

	k := make([]sw_bn254.G1Affine, (len(vars)-14)/2)
	for i := range k {
		k[i], err = newG1(api, vars[i*2:i*2+1], witnesses)
		if err != nil {
			return vk, err
		}
	}
	vk.G1.K = k
	return vk, nil
}

func newG1[T shr.ACIRField](api frontend.API, vars []FunctionInput[T], witnesses map[shr.Witness]frontend.Variable) (sw_bn254.G1Affine, error) {
	var ret sw_bn254.G1Affine
	primeField, err := emulated.NewField[emulated.BN254Fp](api)

	if err != nil {
		return ret, err
	}
	alphaX, err := VariableTo64BitLimbs(api, vars[0], witnesses)
	if err != nil {
		return ret, err
	}
	alphaY, err := VariableTo64BitLimbs(api, vars[1], witnesses)
	if err != nil {
		return ret, err
	}

	ret.X = emulated.Element[sw_bn254.BaseField](*primeField.NewElement(alphaX))
	ret.Y = emulated.Element[sw_bn254.BaseField](*primeField.NewElement(alphaY))

	return ret, err

}

func newG2[T shr.ACIRField](api frontend.API, vars []FunctionInput[T], witnesses map[shr.Witness]frontend.Variable) (sw_bn254.G2Affine, error) {
	ret := sw_bn254.NewG2AffineFixedPlaceholder()
	primeField, err := emulated.NewField[emulated.BN254Fp](api)
	if err != nil {
		return ret, err
	}
	g2BetaXA0, err := VariableTo64BitLimbs(api, vars[0], witnesses)
	if err != nil {
		return ret, err
	}

	ret.P.X.A0 = emulated.Element[sw_bn254.BaseField](*primeField.NewElement(g2BetaXA0))
	g2BetaXA1, err := VariableTo64BitLimbs(api, vars[1], witnesses)
	if err != nil {
		return ret, err
	}

	ret.P.X.A0 = emulated.Element[sw_bn254.BaseField](*primeField.NewElement(g2BetaXA1))

	g2BetaYA0, err := VariableTo64BitLimbs(api, vars[2], witnesses)
	if err != nil {
		return ret, err
	}
	ret.P.Y.A0 = emulated.Element[sw_bn254.BaseField](*primeField.NewElement(g2BetaYA0))

	g2BetaYA1, err := VariableTo64BitLimbs(api, vars[3], witnesses)
	if err != nil {
		return ret, err
	}
	ret.P.Y.A0 = emulated.Element[sw_bn254.BaseField](*primeField.NewElement(g2BetaYA1))

	return ret, nil

}

func newWitness[T shr.ACIRField](api frontend.API, vars []FunctionInput[T], witnesses map[shr.Witness]frontend.Variable) (groth16.Witness[emulated.BN254Fr], error) {
	var witness groth16.Witness[emulated.BN254Fr]
	scalarField, err := emulated.NewField[emulated.BN254Fr](api)
	if err != nil {
		return witness, err
	}
	witnessVector := make([]emulated.Element[sw_bn254.ScalarField], len(vars))

	for i := range vars {
		value, err := VariableTo64BitLimbs(api, vars[i], witnesses)
		if err != nil {
			return witness, err
		}
		witnessVector[i] = emulated.Element[sw_bn254.ScalarField](*scalarField.NewElement(value))

	}
	witness.Public = witnessVector
	return witness, err
}

func newProof[T shr.ACIRField](api frontend.API, vars []FunctionInput[T], witnesses map[shr.Witness]frontend.Variable) (groth16.Proof[sw_bn254.G1Affine, sw_bn254.G2Affine], error) {
	var proof groth16.Proof[sw_bn254.G1Affine, sw_bn254.G2Affine]
	var err error
	proof.Ar, err = newG1(api, vars[0:2], witnesses)
	if err != nil {
		return proof, err
	}
	proof.Bs, err = newG2(api, vars[2:6], witnesses)
	if err != nil {
		return proof, err
	}
	proof.Krs, err = newG1(api, vars[6:8], witnesses)
	if err != nil {
		return proof, err
	}
	return proof, nil
}
