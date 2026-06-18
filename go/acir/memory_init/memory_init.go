package memory_init

import (
	"encoding/json"
	"fmt"
	"sunspot/go/acir/msgpackutil"
	ops "sunspot/go/acir/opcodes"
	shr "sunspot/go/acir/shared"

	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/lookup/logderivlookup"
)

type MemoryInit[T shr.ACIRField, E constraint.Element] struct {
	BlockID uint32
	Table   *logderivlookup.Table
	Init    []shr.Witness
}

// On the wire MemoryInit's payload is a 3-tagged-field struct: 0=block_id (BlockId, transparent u32),
// 1=init (Vec<Witness>), 2=block_type (enum). We ignore block_type in this backend
func (m *MemoryInit[T, E]) UnmarshalReader(r *msgpackutil.Reader) error {
	return msgpackutil.ReadStruct(r, m.decode)
}

func (m *MemoryInit[T, E]) decode(tag int, r *msgpackutil.Reader) error {
	switch tag {
	case 0:
		v, err := r.ReadU32()
		if err != nil {
			return err
		}
		m.BlockID = v
		return nil
	case 1:
		return shr.ReadWitnessVec(r, &m.Init)
	case 2:
		// tag 2 encodes the blocktype, which for this backend is irrelevant
		return r.SkipValue()
	default:
		return fmt.Errorf("memory_init: unknown field tag %d", tag)
	}
}

func (m *MemoryInit[T, E]) Equals(other ops.Opcode[E]) bool {
	value, ok := other.(*MemoryInit[T, E])
	if !ok || m.BlockID != value.BlockID {
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

	return true
}

func (m *MemoryInit[T, E]) Define(api frontend.Builder[E], witnesses map[shr.Witness]frontend.Variable) error {
	for i := range m.Init {
		(*m.Table).Insert(witnesses[m.Init[i]])
	}
	return nil
}

func (m *MemoryInit[T, E]) MarshalJSON() ([]byte, error) {
	stringMap := make(map[string]interface{})
	stringMap["assert_zero"] = m
	return json.Marshal(stringMap)
}
