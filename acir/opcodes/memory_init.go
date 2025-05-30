package opcodes

import (
	"encoding/binary"
	"io"
	shr "nr-groth16/acir/shared"
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

func (m *MemoryInit[T]) Equals(other *MemoryInit[T]) bool {
	if m.BlockID != other.BlockID || m.BlockType != other.BlockType {
		return false
	}

	if len(m.Init) != len(other.Init) {
		return false
	}
	for i := range m.Init {
		if !m.Init[i].Equals(&other.Init[i]) {
			return false
		}
	}

	if (m.CallData == nil && other.CallData != nil) || (m.CallData != nil && other.CallData == nil) {
		return false
	}
	if m.CallData != nil && other.CallData != nil && *m.CallData != *other.CallData {
		return false
	}

	return true
}
