package brillig

import (
	"encoding/binary"
	"encoding/json"
	"io"
	"math/big"
	exp "nr-groth16/acir/expression"
	ops "nr-groth16/acir/opcodes"
	shr "nr-groth16/acir/shared"

	"github.com/consensys/gnark/frontend"
	"github.com/google/btree"
)

type BrilligCall[T shr.ACIRField] struct {
	ID        uint32
	Inputs    []BrilligInputs[T]
	Outputs   []BrilligOutputs
	Predicate *exp.Expression[T]
}

func (b *BrilligCall[T]) UnmarshalReader(r io.Reader) error {
	if err := binary.Read(r, binary.LittleEndian, &b.ID); err != nil {
		return err
	}

	var numInputs uint64
	if err := binary.Read(r, binary.LittleEndian, &numInputs); err != nil {
		return err
	}
	b.Inputs = make([]BrilligInputs[T], numInputs)
	for i := uint64(0); i < numInputs; i++ {
		if err := b.Inputs[i].UnmarshalReader(r); err != nil {
			return err
		}
	}

	var numOutputs uint64
	if err := binary.Read(r, binary.LittleEndian, &numOutputs); err != nil {
		return err
	}
	b.Outputs = make([]BrilligOutputs, numOutputs)
	for i := uint64(0); i < numOutputs; i++ {
		if err := b.Outputs[i].UnmarshalReader(r); err != nil {
			return err
		}
	}

	var predicateExists uint8
	if err := binary.Read(r, binary.LittleEndian, &predicateExists); err != nil {
		return err
	}
	if predicateExists == 1 {
		b.Predicate = new(exp.Expression[T])
		if err := b.Predicate.UnmarshalReader(r); err != nil {
			return err
		}
	} else {
		b.Predicate = nil
	}

	return nil
}

func (o *BrilligCall[T]) Equals(other ops.Opcode) bool {
	panic("unimplemented")
}
func (o *BrilligCall[T]) Define(api frontend.API, witnesses map[shr.Witness]frontend.Variable) error {
	// do nothing: brillig calls are unconstrained
	return nil
}

func (o *BrilligCall[T]) MarshalJSON() ([]byte, error) {
	stringMap := make(map[string]interface{})
	stringMap["brillig_call"] = o
	return json.Marshal(stringMap)
}

func (o *BrilligCall[T]) FillWitnessTree(tree *btree.BTree) bool {
	return !(tree == nil)
}

func (o *BrilligCall[T]) CollectConstantsAsWitnesses(start uint32, tree *btree.BTree) bool {
	return !(tree == nil)
}

func (o *BrilligCall[T]) FeedConstantsAsWitnesses() []*big.Int {
	values := make([]*big.Int, 0)
	return values
}
