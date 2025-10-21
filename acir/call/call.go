package call

import (
	"encoding/binary"
	"encoding/json"
	"io"
	"math/big"
	exp "nr-groth16/acir/expression"
	ops "nr-groth16/acir/opcodes"
	shr "nr-groth16/acir/shared"

	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/google/btree"
)

type Call[T shr.ACIRField, E constraint.Element] struct {
	ID        uint32
	Inputs    []shr.Witness
	Outputs   []shr.Witness
	Predicate *exp.Expression[T, E]
}

func (c *Call[T, E]) UnmarshalReader(r io.Reader) error {
	if err := binary.Read(r, binary.LittleEndian, &c.ID); err != nil {
		return err
	}

	var numInputs uint64
	if err := binary.Read(r, binary.LittleEndian, &numInputs); err != nil {
		return err
	}
	c.Inputs = make([]shr.Witness, numInputs)
	if err := binary.Read(r, binary.LittleEndian, &c.Inputs); err != nil {
		return err
	}

	var numOutputs uint64
	if err := binary.Read(r, binary.LittleEndian, &numOutputs); err != nil {
		return err
	}
	c.Outputs = make([]shr.Witness, numOutputs)
	if err := binary.Read(r, binary.LittleEndian, &c.Outputs); err != nil {
		return err
	}

	var predicateExists uint8
	if err := binary.Read(r, binary.LittleEndian, &predicateExists); err != nil {
		return err
	}
	if predicateExists == 1 {
		c.Predicate = new(exp.Expression[T, E])
		if err := c.Predicate.UnmarshalReader(r); err != nil {
			return err
		}
	}

	return nil
}

func (c *Call[T, E]) Equals(other ops.Opcode[E]) bool {
	value, ok := other.(*Call[T, E])
	if !ok || c.ID != value.ID {
		return false
	}

	if len(c.Inputs) != len(value.Inputs) || len(c.Outputs) != len(value.Outputs) {
		return false
	}

	for i := range c.Inputs {
		if c.Inputs[i] != value.Inputs[i] {
			return false
		}
	}

	for i := range c.Outputs {
		if c.Outputs[i] != value.Outputs[i] {
			return false
		}
	}

	if (c.Predicate == nil) != (value.Predicate == nil) {
		return false
	}

	if c.Predicate != nil && !c.Predicate.Equals(value.Predicate) {
		return false
	}

	return true
}

func (*Call[T, E]) CollectConstantsAsWitnesses(start uint32, tree *btree.BTree) bool {
	return tree != nil
}

func (o *Call[T, E]) Define(api frontend.Builder[E], witnesses map[shr.Witness]frontend.Variable) error {
	return nil
}

func (c *Call[T, E]) FeedConstantsAsWitnesses() []*big.Int {
	return make([]*big.Int, 0)
}

func (c *Call[T, E]) FillWitnessTree(tree *btree.BTree, index uint32) bool {
	for i := range c.Inputs {
		tree.ReplaceOrInsert(c.Inputs[i] + shr.Witness(index))
	}
	for i := range c.Outputs {
		tree.ReplaceOrInsert(c.Outputs[i] + shr.Witness(index))
	}
	return tree != nil
}

func (c *Call[T, E]) MarshalJSON() ([]byte, error) {
	stringMap := make(map[string]interface{})
	stringMap["circuit_call"] = c
	return json.Marshal(stringMap)
}
