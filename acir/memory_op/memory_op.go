package memory_op

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

type MemoryOp[T shr.ACIRField] struct {
	BlockID   uint32
	Operation exp.Expression[T]
	Index     exp.Expression[T]
	Value     exp.Expression[T]
	Predicate *exp.Expression[T]
}

func (m *MemoryOp[T]) UnmarshalReader(r io.Reader) error {
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
		m.Predicate = new(exp.Expression[T])
		if err := m.Predicate.UnmarshalReader(r); err != nil {
			return err
		}
	}
	return nil
}

func (m *MemoryOp[T]) Equals(other ops.Opcode) bool {
	mem_op, ok := other.(*MemoryOp[T])
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

func (*MemoryOp[T]) CollectConstantsAsWitnesses(start uint32, tree *btree.BTree) bool {
	return !(tree == nil)
}

func (o *MemoryOp[T]) Define(api frontend.API, witnesses map[shr.Witness]frontend.Variable) error {
	panic("MemoryInit opcode is not implemented yet") // TODO: Implement MemoryInit calculation
	//return o.MemoryInit.Calculate(api, witnesses)
}

func (o *MemoryOp[T]) FeedConstantsAsWitnesses() []*big.Int {
	values := make([]*big.Int, 0)
	return values
}

func (o *MemoryOp[T]) FillWitnessTree(tree *btree.BTree) bool {
	return !(tree == nil)
}

func (o MemoryOp[T]) MarshalJSON() ([]byte, error) {
	stringMap := make(map[string]interface{})
	stringMap["memory_op"] = o
	return json.Marshal(stringMap)
}
