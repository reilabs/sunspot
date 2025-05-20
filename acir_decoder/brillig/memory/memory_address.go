package memory

import (
	"encoding/binary"
	"fmt"
	"io"
)

type MemoryAddress struct {
	Kind  MemoryAddressKind
	Value uint64
}

type MemoryAddressKind uint32

const (
	MemoryAddressKindDirect MemoryAddressKind = iota
	MemoryAddressKindRelative
)

func (m *MemoryAddress) UnmarshalReader(r io.Reader) error {
	if err := binary.Read(r, binary.LittleEndian, &m.Kind); err != nil {
		return err
	}

	if err := binary.Read(r, binary.LittleEndian, &m.Value); err != nil {
		return err
	}

	if m.Kind != MemoryAddressKindDirect && m.Kind != MemoryAddressKindRelative {
		return fmt.Errorf("invalid MemoryAddressKind: %d", m.Kind)
	}

	return nil
}
