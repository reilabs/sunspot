package blackboxfunc

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"nr-groth16/acir/opcodes"
	shr "nr-groth16/acir/shared"

	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/google/btree"
)

type BlackBoxFunction[E constraint.Element] interface {
	UnmarshalReader(r io.Reader) error
	Define(api frontend.API, witnesses map[shr.Witness]frontend.Variable) error
	Equals(other BlackBoxFunction[E]) bool
	FillWitnessTree(tree *btree.BTree) bool
}

// Struct that implements the Opcode interface
// Allows us to create generic behaviour for all black box functions
// TODO: revisit and think about this design
type BlackBoxFuncCall[T shr.ACIRField, E constraint.Element] struct {
	function BlackBoxFunction[E]
}

func (b BlackBoxFuncCall[T, E]) CollectConstantsAsWitnesses(start uint32, tree *btree.BTree) bool {
	return true
}

func (b BlackBoxFuncCall[T, E]) Define(api frontend.API, witnesses map[shr.Witness]frontend.Variable) error {
	return b.function.Define(api, witnesses)
}

func (b BlackBoxFuncCall[T, E]) Equals(other opcodes.Opcode[E]) bool {
	bbf, ok := other.(BlackBoxFuncCall[T, E])
	if !ok {
		return false
	}
	return b.function.Equals(bbf.function)
}

func (b BlackBoxFuncCall[T, E]) FeedConstantsAsWitnesses() []*big.Int {
	values := make([]*big.Int, 0)
	return values
}

func (b BlackBoxFuncCall[T, E]) FillWitnessTree(tree *btree.BTree) bool {
	return b.function.FillWitnessTree(tree)
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
	case 1:
		return &BlackBoxFuncCall[T, E]{&And[T, E]{}}, nil
	case 2:
		return &BlackBoxFuncCall[T, E]{&Xor[T, E]{}}, nil
	case 3:
		return &BlackBoxFuncCall[T, E]{&Range[T, E]{}}, nil
	case 10:
		return &BlackBoxFuncCall[T, E]{&Keccakf1600[T, E]{}}, nil
	case 19:
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
