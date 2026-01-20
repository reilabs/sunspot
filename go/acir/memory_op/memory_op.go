package memory_op

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	exp "sunspot/go/acir/expression"
	ops "sunspot/go/acir/opcodes"
	shr "sunspot/go/acir/shared"

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

	return m.Operation.Equals(&mem_op.Operation) && m.Index.Equals(&mem_op.Index) && m.Value.Equals(&mem_op.Value)

}

func (*MemoryOp[T, E]) CollectConstantsAsWitnesses(start uint32, tree *btree.BTree) bool {
	return tree != nil
}

func (o *MemoryOp[T, E]) Define(api frontend.Builder[E], witnesses map[shr.Witness]frontend.Variable) error {

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
			is_writable := api.IsZero(api.Sub(insertion_index, frontend.Variable(i)))
			updated := api.Select(is_writable, o.Value.Calculate(api, witnesses), (*table).Lookup(i)[0])
			newTable.Insert(updated)
		}

		o.Memory[o.BlockID] = &newTable
		return nil

	default:
		return fmt.Errorf("unknown memory operation: %d", o.Operation.Constant.ToBigInt().Uint64())
	}

	return nil
}

func (o *MemoryOp[T, E]) FillWitnessTree(tree *btree.BTree, index uint32) bool {
	return (o.Index.FillWitnessTree(tree, index) &&
		o.Operation.FillWitnessTree(tree, index) &&
		o.Value.FillWitnessTree(tree, index))
}

func (o MemoryOp[T, E]) MarshalJSON() ([]byte, error) {
	stringMap := make(map[string]interface{})
	stringMap["memory_op"] = o
	return json.Marshal(stringMap)
}
