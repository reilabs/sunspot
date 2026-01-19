package blackboxfunc

import (
	"encoding/binary"
	"io"
	shr "sunspot/go/acir/shared"

	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/algebra/emulated/sw_emulated"
	"github.com/consensys/gnark/std/math/bits"
	"github.com/consensys/gnark/std/math/emulated"
	"github.com/consensys/gnark/std/signature/ecdsa"
	"github.com/google/btree"
)

type ECDSASECP256K1[T shr.ACIRField, E constraint.Element] struct {
	PublicKeyX    [32]FunctionInput[T]
	PublicKeyY    [32]FunctionInput[T]
	Signature     [64]FunctionInput[T]
	HashedMessage [32]FunctionInput[T]
	predicate     FunctionInput[T]
	Output        shr.Witness
}

func (a *ECDSASECP256K1[T, Equals]) UnmarshalReader(r io.Reader) error {
	for i := 0; i < 32; i++ {
		if err := a.PublicKeyX[i].UnmarshalReader(r); err != nil {
			return err
		}
	}

	for i := 0; i < 32; i++ {
		if err := a.PublicKeyY[i].UnmarshalReader(r); err != nil {
			return err
		}
	}

	for i := 0; i < 64; i++ {
		if err := a.Signature[i].UnmarshalReader(r); err != nil {
			return err
		}
	}

	for i := 0; i < 32; i++ {
		if err := a.HashedMessage[i].UnmarshalReader(r); err != nil {
			return err
		}
	}

	if err := a.predicate.UnmarshalReader(r); err != nil {
		return err
	}

	if err := binary.Read(r, binary.LittleEndian, &a.Output); err != nil {
		return err
	}
	return nil
}

func (a *ECDSASECP256K1[T, E]) Equals(other BlackBoxFunction[E]) bool {
	value := other.(*ECDSASECP256K1[T, E])
	if len(a.PublicKeyX) != len(value.PublicKeyX) ||
		len(a.PublicKeyY) != len(value.PublicKeyY) ||
		len(a.Signature) != len(value.Signature) ||
		len(a.HashedMessage) != len(value.HashedMessage) {
		return false
	}

	for i := 0; i < 32; i++ {
		if !a.PublicKeyX[i].Equals(&value.PublicKeyX[i]) ||
			!a.PublicKeyY[i].Equals(&value.PublicKeyY[i]) ||
			!a.HashedMessage[i].Equals(&value.HashedMessage[i]) {
			return false
		}
	}

	for i := 0; i < 64; i++ {
		if !a.Signature[i].Equals(&value.Signature[i]) {
			return false
		}
	}

	return a.Output == value.Output
}

func (a *ECDSASECP256K1[T, E]) Define(api frontend.Builder[E], witnesses map[shr.Witness]frontend.Variable) error {
	primeField, err := emulated.NewField[emulated.Secp256k1Fp](api)
	if err != nil {
		return err
	}
	scalarField, err := emulated.NewField[emulated.Secp256k1Fr](api)
	if err != nil {
		return err
	}

	qXValue, err := BytesTo64BitLimbs(api, a.PublicKeyX[:], witnesses)
	if err != nil {
		return err
	}

	qYValue, err := BytesTo64BitLimbs(api, a.PublicKeyY[:], witnesses)
	if err != nil {
		return err
	}

	rValue, err := BytesTo64BitLimbs(api, a.Signature[0:32], witnesses)
	if err != nil {
		return err
	}

	sValue, err := BytesTo64BitLimbs(api, a.Signature[32:64], witnesses)
	if err != nil {
		return err
	}

	hash_value, err := BytesTo64BitLimbs(api, a.HashedMessage[:], witnesses)
	if err != nil {
		return err
	}

	Q := ecdsa.PublicKey[emulated.Secp256k1Fp, emulated.Secp256k1Fr]{
		X: *primeField.NewElement(qXValue),
		Y: *primeField.NewElement(qYValue),
	}

	sig := ecdsa.Signature[emulated.Secp256k1Fr]{
		R: *scalarField.NewElement(rValue),
		S: *scalarField.NewElement(sValue),
	}

	msg := scalarField.NewElement(hash_value)
	pred, err := a.predicate.ToVariable(witnesses)
	if err != nil {
		return err
	}
	api.AssertIsEqual(frontend.Variable(0), api.Mul(pred, api.Sub(witnesses[a.Output], Q.IsValid(api, sw_emulated.GetSecp256k1Params(), msg, &sig))))
	return nil
}

func (a *ECDSASECP256K1[T, E]) FillWitnessTree(tree *btree.BTree, index uint32) bool {
	if tree == nil {
		return false
	}

	for _, input := range a.PublicKeyX {
		if input.IsWitness() {
			tree.ReplaceOrInsert(*input.Witness + shr.Witness(index))
		}
	}
	for _, input := range a.PublicKeyY {
		if input.IsWitness() {
			tree.ReplaceOrInsert(*input.Witness + shr.Witness(index))
		}
	}
	for _, input := range a.HashedMessage {
		if input.IsWitness() {
			tree.ReplaceOrInsert(*input.Witness + shr.Witness(index))
		}
	}
	for _, input := range a.Signature {
		if input.IsWitness() {
			tree.ReplaceOrInsert(*input.Witness + shr.Witness(index))
		}
	}

	tree.ReplaceOrInsert(a.Output + shr.Witness(index))

	return true
}

// ACIR has signature variables as big endian bytes
// but gnark wants them as 4 * 64 but limbs.
// See https://pkg.go.dev/github.com/consensys/gnark/std/math/emulated@v0.14.0#hdr-Element_representation

func BytesTo64BitLimbs[T shr.ACIRField](
	api frontend.API,
	vars []FunctionInput[T],
	witnesses map[shr.Witness]frontend.Variable,
) ([]frontend.Variable, error) {
	const bitsPerLimb = 64
	const nbLimbs = 4
	if len(vars) != 32 {
		panic("expected 32 variables")
	}

	bit_array := make([]frontend.Variable, 256)
	out := make([]frontend.Variable, nbLimbs)

	// Reverse the byte order: vars[0] is MSB → goes to highest slot
	for i := 0; i < 32; i++ {
		variable, err := vars[31-i].ToVariable(witnesses)
		if err != nil {
			return nil, err
		}
		start := 8 * i
		copy(bit_array[start:start+8], bits.ToBinary(api, variable, bits.WithNbDigits(8)))
	}

	// Now pack into 64-bit limbs (little endian order)
	for i := 0; i < nbLimbs; i++ {
		start := i * bitsPerLimb
		end := start + bitsPerLimb
		chunk := bit_array[start:end]
		out[i] = bits.FromBinary(api, chunk)
	}
	return out, nil
}
