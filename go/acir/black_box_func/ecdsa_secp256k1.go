package blackboxfunc

import (
	"math/big"
	"sunspot/go/acir/msgpackutil"
	shr "sunspot/go/acir/shared"

	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/algebra/emulated/sw_emulated"
	"github.com/consensys/gnark/std/math/bits"
	"github.com/consensys/gnark/std/math/emulated"
	"github.com/consensys/gnark/std/signature/ecdsa"
)

type ECDSASECP256K1[T shr.ACIRField, E constraint.Element] struct {
	PublicKeyX    [32]FunctionInput[T]
	PublicKeyY    [32]FunctionInput[T]
	Signature     [64]FunctionInput[T]
	HashedMessage [32]FunctionInput[T]
	predicate     FunctionInput[T]
	Output        shr.Witness
}

func (a *ECDSASECP256K1[T, E]) UnmarshalReader(r *msgpackutil.Reader) error {
	return msgpackutil.ReadStruct(r, "EcdsaSecp256k1", []msgpackutil.Field{
		{Name: "public_key_x", Decode: func(r *msgpackutil.Reader) error { return msgpackutil.ReadArrayInto(r, a.PublicKeyX[:]) }},
		{Name: "public_key_y", Decode: func(r *msgpackutil.Reader) error { return msgpackutil.ReadArrayInto(r, a.PublicKeyY[:]) }},
		{Name: "signature", Decode: func(r *msgpackutil.Reader) error { return msgpackutil.ReadArrayInto(r, a.Signature[:]) }},
		{Name: "hashed_message", Decode: func(r *msgpackutil.Reader) error { return msgpackutil.ReadArrayInto(r, a.HashedMessage[:]) }},
		{Name: "predicate", Decode: a.predicate.UnmarshalReader},
		{Name: "output", Decode: a.Output.UnmarshalReader},
	})
}

func (a *ECDSASECP256K1[T, E]) Equals(other BlackBoxFunction[E]) bool {
	value, ok := other.(*ECDSASECP256K1[T, E])
	if !ok {
		return false
	}
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

	// PK on-curve validation
	cr, err := sw_emulated.New[emulated.Secp256k1Fp, emulated.Secp256k1Fr](api, sw_emulated.GetSecp256k1Params())
	if err != nil {
		return err
	}
	cr.AssertIsOnCurve(&sw_emulated.AffinePoint[emulated.Secp256k1Fp]{X: Q.X, Y: Q.Y})
	isIdentity := api.And(primeField.IsZero(&Q.X), primeField.IsZero(&Q.Y))
	api.AssertIsEqual(isIdentity, frontend.Variable(0))

	sig := ecdsa.Signature[emulated.Secp256k1Fr]{
		R: *scalarField.NewElement(rValue),
		S: *scalarField.NewElement(sValue),
	}

	msg := scalarField.NewElement(hash_value)
	pred, err := a.predicate.ToVariable(witnesses)
	if err != nil {
		return err
	}

	validSig := Q.IsValid(api, sw_emulated.GetSecp256k1Params(), msg, &sig)

	// Noir's verify_signature rejects signatures with s > n/2 (low-s form) to
	// prevent malleability; gnark's IsValid does not, so enforce it here.
	halfOrder := new(big.Int).Rsh(emulated.Secp256k1Fr{}.Modulus(), 1)
	sBits := scalarField.ToBits(&sig.S)
	isLowS := isLessOrEqualConstant(api, sBits, halfOrder)

	result := api.Mul(validSig, isLowS)

	api.AssertIsEqual(frontend.Variable(0), api.Mul(pred, api.Sub(witnesses[a.Output], result)))
	return nil
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

func (*ECDSASECP256K1[T, E]) SerdeName() string { return "EcdsaSecp256k1" }
