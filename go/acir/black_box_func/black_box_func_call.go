package blackboxfunc

import (
	"encoding/json"
	"github.com/reilabs/sunspot/go/acir/msgpackutil"
	"github.com/reilabs/sunspot/go/acir/opcodes"
	shr "github.com/reilabs/sunspot/go/acir/shared"

	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
)

type BlackBoxFunction[E constraint.Element] interface {
	msgpackutil.EnumVariant
	Define(api frontend.Builder[E], witnesses map[shr.Witness]frontend.Variable) error
	Equals(other BlackBoxFunction[E]) bool
}

// Struct that implements the Opcode interface
// Allows us to create generic behaviour for all black box functions
type BlackBoxFuncCall[T shr.ACIRField, E constraint.Element] struct {
	function BlackBoxFunction[E]
}

func (b BlackBoxFuncCall[T, E]) Define(api frontend.Builder[E], witnesses map[shr.Witness]frontend.Variable) error {
	return b.function.Define(api, witnesses)
}

func (b BlackBoxFuncCall[T, E]) Equals(other opcodes.Opcode[E]) bool {
	bbf, ok := other.(*BlackBoxFuncCall[T, E])
	if !ok {
		return false
	}
	return b.function.Equals(bbf.function)
}

func (b BlackBoxFuncCall[T, E]) MarshalJSON() ([]byte, error) {
	stringMap := make(map[string]interface{})
	stringMap["black_box_func_call"] = b
	return json.Marshal(stringMap)
}

// UnmarshalReader dispatches on the enum variant name, allocates the concrete
// payload type, and delegates payload decoding to the variant.
func (b *BlackBoxFuncCall[T, E]) UnmarshalReader(r *msgpackutil.Reader) error {
	return msgpackutil.ReadDispatchedEnum(r, "BlackBoxFuncCall", []BlackBoxFunction[E]{
		&AES128Encrypt[T, E]{},
		&And[T, E]{},
		&Xor[T, E]{},
		&Range[T, E]{},
		&Blake2s[T, E]{},
		&Blake3[T, E]{},
		&ECDSASECP256K1[T, E]{},
		&ECDSASECP256R1[T, E]{},
		&MultiScalarMul[T, E]{},
		&EmbeddedCurveAdd[T, E]{},
		&Keccakf1600[T, E]{},
		&RecursiveAggregation[T, E]{},
		&Poseidon2Permutation[T, E]{},
		&SHA256Compression[T, E]{},
	}, func(v BlackBoxFunction[E]) { b.function = v })
}

func (*BlackBoxFuncCall[T, E]) SerdeName() string { return "BlackBoxFuncCall" }
