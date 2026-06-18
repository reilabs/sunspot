package blackboxfunc

import (
	"encoding/json"
	"fmt"
	"sunspot/go/acir/msgpackutil"
	"sunspot/go/acir/opcodes"
	shr "sunspot/go/acir/shared"

	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
)

type BlackBoxFunction[E constraint.Element] interface {
	Define(api frontend.Builder[E], witnesses map[shr.Witness]frontend.Variable) error
	Equals(other BlackBoxFunction[E]) bool
	decode(tag int, r *msgpackutil.Reader) error
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

// UnmarshalReader reads the enum: dispatches on variant tag, allocates the
// concrete payload type, and delegates payload decoding.
func (b *BlackBoxFuncCall[T, E]) UnmarshalReader(r *msgpackutil.Reader) error {
	return msgpackutil.ReadEnum(r, b.decodeBlackBoxFunction)
}

func (b *BlackBoxFuncCall[T, E]) decodeBlackBoxFunction(tag int, r *msgpackutil.Reader) error {
	fn, err := newBlackBoxFunction[T, E](tag)
	if err != nil {
		return err
	}
	if err := msgpackutil.ReadStruct(r, fn.decode); err != nil {
		return fmt.Errorf("black-box variant %d: %w", tag, err)
	}
	b.function = fn
	return nil
}

func newBlackBoxFunction[T shr.ACIRField, E constraint.Element](tag int) (BlackBoxFunction[E], error) {
	switch tag {
	case 0:
		return &AES128Encrypt[T, E]{}, nil
	case 1:
		return &And[T, E]{}, nil
	case 2:
		return &Xor[T, E]{}, nil
	case 3:
		return &Range[T, E]{}, nil
	case 4:
		return &Blake2s[T, E]{}, nil
	case 5:
		return &Blake3[T, E]{}, nil
	case 6:
		return &ECDSASECP256K1[T, E]{}, nil
	case 7:
		return &ECDSASECP256R1[T, E]{}, nil
	case 8:
		return &MultiScalarMul[T, E]{}, nil
	case 9:
		return &EmbeddedCurveAdd[T, E]{}, nil
	case 10:
		return &Keccakf1600[T, E]{}, nil
	case 11:
		return &RecursiveAggregation[T, E]{}, nil
	case 12:
		return &Poseidon2Permutation[T, E]{}, nil
	case 13:
		return &SHA256Compression[T, E]{}, nil
	default:
		return nil, fmt.Errorf("blackbox variant %d not implemented", tag)
	}
}
