package memory_op

import (
	"encoding/json"
	"sunspot/go/acir/msgpackutil"
	ops "sunspot/go/acir/opcodes"
	shr "sunspot/go/acir/shared"

	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/lookup/logderivlookup"
)

type MemoryOp[T shr.ACIRField, E constraint.Element] struct {
	BlockID uint32
	Memory  map[uint32]*logderivlookup.Table
	IsWrite bool
	Index   shr.Witness
	Value   shr.Witness
}

// MemoryOp's payload (after the outer Opcode dispatch): block_id and op
// (MemOp). The inner MemOp has fields operation (bool, serde-renamed "read"),
// index (Witness), and value (Witness).
func (m *MemoryOp[T, E]) UnmarshalReader(r *msgpackutil.Reader) error {
	return msgpackutil.ReadStruct(r, "MemoryOp", []msgpackutil.Field{
		{Name: "block_id", Decode: func(r *msgpackutil.Reader) error {
			v, err := r.ReadUint()
			if err != nil {
				return err
			}
			m.BlockID = uint32(v)
			return nil
		}},
		{Name: "op", Decode: func(r *msgpackutil.Reader) error {
			return msgpackutil.ReadStruct(r, "MemOp", []msgpackutil.Field{
				{Name: "read", Decode: func(r *msgpackutil.Reader) error {
					v, err := r.ReadBool()
					if err != nil {
						return err
					}
					m.IsWrite = v
					return nil
				}},
				{Name: "index", Decode: m.Index.UnmarshalReader},
				{Name: "value", Decode: m.Value.UnmarshalReader},
			})
		}},
	})
}

func (m *MemoryOp[T, E]) Equals(other ops.Opcode[E]) bool {
	o, ok := other.(*MemoryOp[T, E])
	if !ok {
		return false
	}
	return m.BlockID == o.BlockID &&
		m.IsWrite == o.IsWrite &&
		m.Index == o.Index &&
		m.Value == o.Value
}

func (o *MemoryOp[T, E]) Define(api frontend.Builder[E], witnesses map[shr.Witness]frontend.Variable) error {
	table := o.Memory[o.BlockID]
	indexVar := witnesses[o.Index]
	valueVar := witnesses[o.Value]

	if o.IsWrite {
		newTable := logderivlookup.New(api)
		// dummy insertion to find the length of the table
		tableLen := (*table).Insert(0)
		for i := 0; i < tableLen; i++ {
			isWritable := api.IsZero(api.Sub(indexVar, frontend.Variable(i)))
			updated := api.Select(isWritable, valueVar, (*table).Lookup(i)[0])
			newTable.Insert(updated)
		}
		o.Memory[o.BlockID] = &newTable
		return nil
	} else {
		api.AssertIsEqual((*table).Lookup(indexVar)[0], valueVar)
		return nil
	}

}

func (o MemoryOp[T, E]) MarshalJSON() ([]byte, error) {
	stringMap := make(map[string]interface{})
	stringMap["memory_op"] = o
	return json.Marshal(stringMap)
}

func (*MemoryOp[T, E]) SerdeName() string { return "MemoryOp" }
