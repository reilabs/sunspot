package brillig

import (
	"encoding/binary"
	"encoding/json"
	"io"
	exp "sunspot/go/acir/expression"
	ops "sunspot/go/acir/opcodes"
	shr "sunspot/go/acir/shared"

	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/google/btree"
)

type BrilligCall[T shr.ACIRField, E constraint.Element] struct {
	ID        uint32
	Inputs    []BrilligInputs[T, E]
	Outputs   []BrilligOutputs
	Predicate *exp.Expression[T, E]
}

func (b *BrilligCall[T, E]) UnmarshalReader(r io.Reader) error {
	if err := binary.Read(r, binary.LittleEndian, &b.ID); err != nil {
		return err
	}

	var numInputs uint64
	if err := binary.Read(r, binary.LittleEndian, &numInputs); err != nil {
		return err
	}
	b.Inputs = make([]BrilligInputs[T, E], numInputs)
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
		b.Predicate = new(exp.Expression[T, E])
		if err := b.Predicate.UnmarshalReader(r); err != nil {
			return err
		}
	} else {
		b.Predicate = nil
	}

	return nil
}

func (o *BrilligCall[T, E]) Equals(other ops.Opcode[E]) bool {
	panic("unimplemented")
}
func (o *BrilligCall[T, E]) Define(api frontend.Builder[E], witnesses map[shr.Witness]frontend.Variable) error {
	// do nothing: brillig calls are unconstrained
	return nil
}

func (o *BrilligCall[T, E]) MarshalJSON() ([]byte, error) {
	stringMap := make(map[string]interface{})
	stringMap["brillig_call"] = o
	return json.Marshal(stringMap)
}

func (o *BrilligCall[T, E]) FillWitnessTree(tree *btree.BTree, index uint32) bool {
	return tree != nil
}
