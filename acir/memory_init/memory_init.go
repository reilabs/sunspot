package memory_init

import (
	"encoding/binary"
	"encoding/json"
	"io"
	ops "sunspot/acir/opcodes"
	shr "sunspot/acir/shared"

	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/lookup/logderivlookup"
	"github.com/google/btree"
)

type MemoryInit[T shr.ACIRField, E constraint.Element] struct {
	BlockID   uint32
	Table     *logderivlookup.Table
	Init      []shr.Witness
	BlockType BlockKind
	CallData  *uint32
}

type BlockKind uint32

const (
	ACIRMemoryBlockMemory BlockKind = iota
	ACIRMemoryBlockCallData
	ACIRMemoryBlockReturnData
)

func (m *MemoryInit[T, E]) UnmarshalReader(r io.Reader) error {
	if err := binary.Read(r, binary.LittleEndian, &m.BlockID); err != nil {
		return err
	}

	var NumWitnesses uint64
	if err := binary.Read(r, binary.LittleEndian, &NumWitnesses); err != nil {
		return err
	}
	m.Init = make([]shr.Witness, NumWitnesses)
	for i := uint64(0); i < NumWitnesses; i++ {
		if err := m.Init[i].UnmarshalReader(r); err != nil {
			return err
		}
	}

	if err := binary.Read(r, binary.LittleEndian, &m.BlockType); err != nil {
		return err
	}
	if m.BlockType == ACIRMemoryBlockCallData {
		m.CallData = new(uint32)
		if err := binary.Read(r, binary.LittleEndian, m.CallData); err != nil {
			return err
		}
	}

	return nil
}

func (m *MemoryInit[T, E]) Equals(other ops.Opcode[E]) bool {
	value, ok := other.(*MemoryInit[T, E])
	if !ok || m.BlockID != value.BlockID || m.BlockType != value.BlockType {
		return false
	}

	if len(m.Init) != len(value.Init) {
		return false
	}
	for i := range m.Init {
		if !m.Init[i].Equals(&value.Init[i]) {
			return false
		}
	}

	if (m.CallData == nil && value.CallData != nil) || (m.CallData != nil && value.CallData == nil) {
		return false
	}
	if m.CallData != nil && value.CallData != nil && *m.CallData != *value.CallData {
		return false
	}

	return true
}

func (m *MemoryInit[T, E]) Define(api frontend.Builder[E], witnesses map[shr.Witness]frontend.Variable) error {
	for i := range m.Init {
		(*m.Table).Insert(witnesses[m.Init[i]])
	}
	return nil
}

func (m *MemoryInit[T, E]) FillWitnessTree(tree *btree.BTree, index uint32) bool {
	if tree == nil {
		return false
	}
	for _, entry := range m.Init {
		tree.ReplaceOrInsert(entry + shr.Witness(index))
	}
	return true
}

func (m *MemoryInit[T, E]) MarshalJSON() ([]byte, error) {
	stringMap := make(map[string]interface{})
	stringMap["assert_zero"] = m
	return json.Marshal(stringMap)
}
