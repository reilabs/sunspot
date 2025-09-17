package memory_init

import (
	"encoding/binary"
	"encoding/json"
	"io"
	"math/big"
	ops "nr-groth16/acir/opcodes"
	shr "nr-groth16/acir/shared"

	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/lookup/logderivlookup"
	"github.com/google/btree"
)

type MemoryInit[T shr.ACIRField] struct {
	BlockID   uint32
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

func (m *MemoryInit[T]) UnmarshalReader(r io.Reader) error {
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

func (m *MemoryInit[T]) Equals(other ops.Opcode) bool {
	value, ok := other.(*MemoryInit[T])
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

func (m *MemoryInit[T]) CollectConstantsAsWitnesses(start uint32, tree *btree.BTree) bool {
	return tree != nil
}

func (m *MemoryInit[T]) Define(api frontend.API, witnesses map[shr.Witness]frontend.Variable) error {
	t := logderivlookup.New(api)
	for i := range m.Init {
		t.Insert(m.Init[i])
	}
	return nil
}

func (m *MemoryInit[T]) FeedConstantsAsWitnesses() []*big.Int {
	return make([]*big.Int, 0)
}

func (m *MemoryInit[T]) FillWitnessTree(tree *btree.BTree) bool {
	if tree == nil {
		return false
	}
	for _, entry := range m.Init {
		tree.ReplaceOrInsert(entry)
	}
	return true
}

func (m *MemoryInit[T]) MarshalJSON() ([]byte, error) {
	stringMap := make(map[string]interface{})
	stringMap["assert_zero"] = m
	return json.Marshal(stringMap)
}
