package blackboxfunc

import (
	"encoding/binary"
	"encoding/json"
	"io"
	"math/big"
	"nr-groth16/acir/opcodes"
	shr "nr-groth16/acir/shared"

	"github.com/consensys/gnark/frontend"
	"github.com/google/btree"
)

type BlackBoxFunction interface {
	UnmarshalReader(r io.Reader) error
	Define(api frontend.API, witnesses map[shr.Witness]frontend.Variable) error
	Equals(other BlackBoxFunction) bool
}

// Struct that implements the Opcode interface
// Allows us to create generic behaviour for all black box functions
// TODO: revisit and think about this design
type BlackBoxFuncCall[T shr.ACIRField] struct {
	function BlackBoxFunction
}

func (b BlackBoxFuncCall[T]) CollectConstantsAsWitnesses(start uint32, tree *btree.BTree) bool {
	return true
}

func (b BlackBoxFuncCall[T]) Define(api frontend.API, witnesses map[shr.Witness]frontend.Variable) error {
	return b.function.Define(api, witnesses)
}

func (b BlackBoxFuncCall[T]) Equals(other opcodes.Opcode) bool {
	bbf, ok := other.(BlackBoxFuncCall[T])
	if !ok {
		return false
	}
	return b.function.Equals(bbf.function)
}

func (b BlackBoxFuncCall[T]) FeedConstantsAsWitnesses() []*big.Int {
	values := make([]*big.Int, 0)
	return values
}

func (b BlackBoxFuncCall[T]) FillWitnessTree(tree *btree.BTree) bool {
	return !(tree == nil)

}

func (b BlackBoxFuncCall[T]) MarshalJSON() ([]byte, error) {
	stringMap := make(map[string]interface{})
	stringMap["black_box_func_call"] = b
	return json.Marshal(stringMap)
}

func (b BlackBoxFuncCall[T]) UnmarshalReader(r io.Reader) error {
	return b.function.UnmarshalReader(r)
}

func NewBlackBoxFunction[T shr.ACIRField](r io.Reader) (*BlackBoxFuncCall[T], error) {
	var kind uint32
	if err := binary.Read(r, binary.LittleEndian, &kind); err != nil {
		return nil, err
	}
	switch kind {
	case 3:
		function := &Range[T]{}
		return &BlackBoxFuncCall[T]{function}, nil

	case 14:

		return &BlackBoxFuncCall[T]{&Keccakf1600[T]{}}, nil

	default:
		panic("unimplemented")
	}
}

type BlackBoxFuncKindError struct {
	Code uint32
}

func (e BlackBoxFuncKindError) Error() string {
	return "unknown black box function kind"
}
