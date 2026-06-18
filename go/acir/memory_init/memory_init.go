package memory_init

import (
	"encoding/json"
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

// MemoryInit's payload: block_id (BlockId, transparent u32),
// init (Vec<Witness>), block_type (enum). We ignore block_type in this backend.
func (m *MemoryInit[T, E]) UnmarshalReader(r *msgpackutil.Reader) error {
	return msgpackutil.ReadStruct(r, "MemoryInit", []msgpackutil.Field{
		{Name: "block_id", Decode: func(r *msgpackutil.Reader) error {
			v, err := r.ReadU32()
			if err != nil {
				return err
			}
			m.BlockID = v
			return nil
		}},
		{Name: "init", Decode: func(r *msgpackutil.Reader) error { return msgpackutil.ReadVec(r, &m.Init) }},
		{Name: "block_type", Decode: msgpackutil.SkipField},
	})
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

func (*MemoryInit[T, E]) SerdeName() string { return "MemoryInit" }
