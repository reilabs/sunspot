package memory_op

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	exp "nr-groth16/acir/expression"
	ops "nr-groth16/acir/opcodes"
	shr "nr-groth16/acir/shared"

	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/lookup/logderivlookup"
	"github.com/google/btree"
)

type MemoryOp[T shr.ACIRField, E constraint.Element] struct {
	BlockID   uint32
	Memory    map[uint32]*logderivlookup.Table
	Operation exp.Expression[T, E] // operation can be read (0) or write (1)
	Index     exp.Expression[T, E] // witness value of expression is the operation index
	Value     exp.Expression[T, E] // witness value of expression is the operation value
	// predicate is an arithmetic expression that disables the execution of the opcode when the expression evaluates to zero
	Predicate *exp.Expression[T, E]
}

func (m *MemoryOp[T, E]) UnmarshalReader(r io.Reader) error {
	if err := binary.Read(r, binary.LittleEndian, &m.BlockID); err != nil {
		return err
	}

	if err := m.Operation.UnmarshalReader(r); err != nil {
		return err
	}

	if err := m.Index.UnmarshalReader(r); err != nil {
		return err
	}

	if err := m.Value.UnmarshalReader(r); err != nil {
		return err
	}

	var predicateExists uint8
	if err := binary.Read(r, binary.LittleEndian, &predicateExists); err != nil {
		return err
	}
	if predicateExists == 1 {
		m.Predicate = new(exp.Expression[T, E])
		if err := m.Predicate.UnmarshalReader(r); err != nil {
			return err
		}
	}
	return nil
}

func (m *MemoryOp[T, E]) Equals(other ops.Opcode[E]) bool {
	mem_op, ok := other.(*MemoryOp[T, E])
	if !ok {
		return false
	}
	if m.BlockID != mem_op.BlockID {
		return false
	}

	if !m.Operation.Equals(&mem_op.Operation) || !m.Index.Equals(&mem_op.Index) || !m.Value.Equals(&mem_op.Value) {
		return false
	}

	if m.Predicate == nil && mem_op.Predicate == nil {
		return true
	}

	if m.Predicate == nil || mem_op.Predicate == nil {
		return false
	}

	return m.Predicate.Equals(mem_op.Predicate)
}

func (*MemoryOp[T, E]) CollectConstantsAsWitnesses(start uint32, tree *btree.BTree) bool {
	return tree != nil
}

func (o *MemoryOp[T, E]) Define(api frontend.Builder[E], witnesses map[shr.Witness]frontend.Variable) error {
	if o.Predicate != nil && o.Predicate.Constant.ToBigInt() == big.NewInt(0) {
		return nil
	}
	table := o.Memory[o.BlockID]
	switch o.Operation.Constant.ToBigInt().Uint64() { // a bit convoluted but we need a primitve type for switch to work
	case 0:
		api.AssertIsEqual((*table).Lookup(o.Index.Calculate(api, witnesses))[0], o.Value.Calculate(api, witnesses))

	case 1:
		insertion_index := o.Index.Calculate(api, witnesses)
		newTable := logderivlookup.New(api)

		// dummy insertion to find the length of the table
		table_length := (*table).Insert(0)

		for i := 0; i < table_length; i++ {
			if insertion_index == i {
				newTable.Insert(o.Value.Calculate(api, witnesses))
			} else {
				newTable.Insert((*table).Lookup(i)[0])
			}
		}

		o.Memory[o.BlockID] = &newTable
		return nil

	default:
		return fmt.Errorf("unknown memory operation: %d", o.Operation.Constant.ToBigInt().Uint64())
	}

	return nil
}

func (o *MemoryOp[T, E]) FeedConstantsAsWitnesses() []*big.Int {
	return make([]*big.Int, 0)
}

func (o *MemoryOp[T, E]) FillWitnessTree(tree *btree.BTree, index uint32) bool {
	return (o.Predicate == nil || o.Predicate.FillWitnessTree(tree, index)) &&
		o.Index.FillWitnessTree(tree, index) &&
		o.Operation.FillWitnessTree(tree, index) &&
		o.Value.FillWitnessTree(tree, index)
}

func (o MemoryOp[T, E]) MarshalJSON() ([]byte, error) {
	stringMap := make(map[string]interface{})
	stringMap["memory_op"] = o
	return json.Marshal(stringMap)
}
