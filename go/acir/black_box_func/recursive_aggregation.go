package blackboxfunc

import (
	"encoding/binary"
	"fmt"
	"io"
	shr "sunspot/go/acir/shared"

	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	sw_bn254 "github.com/consensys/gnark/std/algebra/emulated/sw_bn254"
	"github.com/consensys/gnark/std/commitments/pedersen"
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
	predicate       FunctionInput[T]
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

	if err := a.predicate.UnmarshalReader(r); err != nil {
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
	switch a.ProofType {
	case 0:
		return a.AggregateGroth16(api, witnesses)
	default:
		return fmt.Errorf("proof type %d not supported in recursive aggregation", a.ProofType)
	}
}

func (a *RecursiveAggregation[T, E]) AggregateGroth16(api frontend.Builder[E], witnesses map[shr.Witness]frontend.Variable) error {
	proof, err := newProof(api, a.Proof, witnesses)
	if err != nil {
		return err
	}

	vk, err := newVK(api, a.VerificationKey, witnesses, len(a.PublicInputs)+1)
	if err != nil {
		return err
	}

	witness, err := newWitness(api, a.PublicInputs, witnesses)
	if err != nil {
		return err
	}

	v, err := groth16.NewVerifier[emulated.BN254Fr, sw_bn254.G1Affine, sw_bn254.G2Affine, sw_bn254.GTEl](api)
	if err != nil {
		return err
	}

	// TODO Find a way to make this dependent on the predicate input
	// ATM we attempt to verify all proofs and don't check the predicate

	v.AssertProof(vk, proof, witness)

	return nil
}

func (a *RecursiveAggregation[T, E]) FillWitnessTree(tree *btree.BTree, index uint32) bool {
	for i := range a.VerificationKey {
		if a.VerificationKey[i].IsWitness() {
			tree.ReplaceOrInsert(*a.VerificationKey[i].Witness + shr.Witness(index))
		}
	}

	for i := range a.Proof {
		if a.Proof[i].IsWitness() {
			tree.ReplaceOrInsert(*a.Proof[i].Witness + shr.Witness(index))
		}
	}
	for i := range a.PublicInputs {
		if a.PublicInputs[i].IsWitness() {
			tree.ReplaceOrInsert(*a.PublicInputs[i].Witness + shr.Witness(index))
		}
	}

	if a.KeyHash.IsWitness() {
		tree.ReplaceOrInsert(*a.KeyHash.Witness + shr.Witness(index))
	}
	return tree != nil
}

func newVK[T shr.ACIRField](api frontend.API, vars []FunctionInput[T], witnesses map[shr.Witness]frontend.Variable, kLen int) (groth16.VerifyingKey[sw_bn254.G1Affine, sw_bn254.G2Affine, sw_bn254.GTEl], error) {
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
	k := make([]sw_bn254.G1Affine, kLen)
	for i := range k {
		k[i], err = newG1(api, vars[14+i*2:14+i*2+2], witnesses)
		if err != nil {
			return vk, err
		}
	}
	vk.G1.K = k

	idx := 14 + (kLen)*2
	commitments := make([]pedersen.VerifyingKey[sw_bn254.G2Affine], (len(vars)-idx)/8)

	for i := range commitments {
		g, err := newG2(api, vars[idx:idx+4], witnesses)
		if err != nil {
			return vk, err
		}
		gSigmaNeg, err := newG2(api, vars[idx+4:idx+8], witnesses)
		if err != nil {
			return vk, err
		}
		commitments[i].G = g
		commitments[i].GSigmaNeg = gSigmaNeg
		idx += 8
	}
	vk.CommitmentKeys = commitments

	return vk, nil
}

func newWitness[T shr.ACIRField](api frontend.API, vars []FunctionInput[T], witnesses map[shr.Witness]frontend.Variable) (groth16.Witness[emulated.BN254Fr], error) {
	var witness groth16.Witness[sw_bn254.ScalarField]
	scalarField, err := emulated.NewField[sw_bn254.ScalarField](api)
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
	proof.Krs, err = newG1(api, vars[2:4], witnesses)
	if err != nil {
		return proof, err
	}
	proof.Bs, err = newG2(api, vars[4:8], witnesses)
	if err != nil {
		return proof, err
	}

	pok, err := newG1(api, vars[8:10], witnesses)
	if err != nil {
		return proof, err
	}
	proof.CommitmentPok.G1El = pok

	idx := 10

	commitments := make([]pedersen.Commitment[sw_bn254.G1Affine], (len(vars)-10)/2)

	for i := range commitments {
		commitment, err := newG1(api, vars[idx:idx+2], witnesses)
		if err != nil {
			return proof, err
		}
		commitments[i].G1El = commitment
	}
	proof.Commitments = commitments
	return proof, nil
}

func newG1[T shr.ACIRField](api frontend.API, vars []FunctionInput[T], witnesses map[shr.Witness]frontend.Variable) (sw_bn254.G1Affine, error) {
	var ret sw_bn254.G1Affine
	primeField, err := emulated.NewField[sw_bn254.BaseField](api)

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
	ret := sw_bn254.G2Affine{}
	primeField, err := emulated.NewField[sw_bn254.BaseField](api)
	if err != nil {
		return ret, err
	}
	XA0, err := VariableTo64BitLimbs(api, vars[0], witnesses)
	if err != nil {
		return ret, err
	}

	ret.P.X.A0 = emulated.Element[sw_bn254.BaseField](*primeField.NewElement(XA0))
	XA1, err := VariableTo64BitLimbs(api, vars[1], witnesses)
	if err != nil {
		return ret, err
	}

	ret.P.X.A1 = emulated.Element[sw_bn254.BaseField](*primeField.NewElement(XA1))

	YA0, err := VariableTo64BitLimbs(api, vars[2], witnesses)
	if err != nil {
		return ret, err
	}

	YA1, err := VariableTo64BitLimbs(api, vars[3], witnesses)
	if err != nil {
		return ret, err
	}
	ret.P.Y.A1 = emulated.Element[sw_bn254.BaseField](*primeField.NewElement(YA1))
	ret.P.Y.A0 = emulated.Element[sw_bn254.BaseField](*primeField.NewElement(YA0))

	return ret, nil

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

	bit_array := api.ToBinary(variable, 256)

	for i := 0; i < nbLimbs; i++ {
		start := i * bitsPerLimb
		end := start + bitsPerLimb
		chunk := bit_array[start:end]
		out[i] = bits.FromBinary(api, chunk)
	}
	return out, nil
}
