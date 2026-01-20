package blackboxfunc

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"sunspot/go/acir/opcodes"
	shr "sunspot/go/acir/shared"

	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/google/btree"
)

type BlackBoxFunction[E constraint.Element] interface {
	UnmarshalReader(r io.Reader) error
	Define(api frontend.Builder[E], witnesses map[shr.Witness]frontend.Variable) error
	Equals(other BlackBoxFunction[E]) bool
	FillWitnessTree(tree *btree.BTree, index uint32) bool
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
	bbf, ok := other.(BlackBoxFuncCall[T, E])
	if !ok {
		return false
	}
	return b.function.Equals(bbf.function)
}

func (b BlackBoxFuncCall[T, E]) FillWitnessTree(tree *btree.BTree, index uint32) bool {
	return b.function.FillWitnessTree(tree, index)
}

func (b BlackBoxFuncCall[T, E]) MarshalJSON() ([]byte, error) {
	stringMap := make(map[string]interface{})
	stringMap["black_box_func_call"] = b
	return json.Marshal(stringMap)
}

func (b BlackBoxFuncCall[T, E]) UnmarshalReader(r io.Reader) error {
	return b.function.UnmarshalReader(r)
}

func NewBlackBoxFunction[T shr.ACIRField, E constraint.Element](r io.Reader) (*BlackBoxFuncCall[T, E], error) {
	var kind uint32
	if err := binary.Read(r, binary.LittleEndian, &kind); err != nil {
		return nil, err
	}
	switch kind {
	case 0:
		return &BlackBoxFuncCall[T, E]{&AES128Encrypt[T, E]{}}, nil
	case 1:
		return &BlackBoxFuncCall[T, E]{&And[T, E]{}}, nil
	case 2:
		return &BlackBoxFuncCall[T, E]{&Xor[T, E]{}}, nil
	case 3:
		return &BlackBoxFuncCall[T, E]{&Range[T, E]{}}, nil
	case 4:
		return &BlackBoxFuncCall[T, E]{&Blake2s[T, E]{}}, nil
	case 5:
		return &BlackBoxFuncCall[T, E]{&Blake3[T, E]{}}, nil
	case 6:
		return &BlackBoxFuncCall[T, E]{&ECDSASECP256K1[T, E]{}}, nil
	case 7:
		return &BlackBoxFuncCall[T, E]{&ECDSASECP256R1[T, E]{}}, nil
	case 8:
		return &BlackBoxFuncCall[T, E]{&MultiScalarMul[T, E]{}}, nil
	case 9:
		return &BlackBoxFuncCall[T, E]{&EmbeddedCurveAdd[T, E]{}}, nil
	case 10:
		return &BlackBoxFuncCall[T, E]{&Keccakf1600[T, E]{}}, nil
	case 11:
		return &BlackBoxFuncCall[T, E]{&RecursiveAggregation[T, E]{}}, nil
	case 12:
		return &BlackBoxFuncCall[T, E]{&Poseidon2Permutation[T, E]{}}, nil
	case 13:
		return &BlackBoxFuncCall[T, E]{&SHA256Compression[T, E]{}}, nil
	default:
		return nil, fmt.Errorf("blackbox opcode %d not yet implemented", kind)
	}
}

type BlackBoxFuncKindError struct {
	Code uint32
}

func (e BlackBoxFuncKindError) Error() string {
	return "unknown black box function kind"
}
