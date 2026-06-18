package memory_op

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

type MemoryOp[T shr.ACIRField, E constraint.Element] struct {
	BlockID uint32
	Memory  map[uint32]*logderivlookup.Table
	IsWrite bool
	Index   shr.Witness
	Value   shr.Witness
}

// On the wire MemoryOp's payload (after the outer Opcode dispatch) is a
// 2-tagged-field struct: 0=block_id, 1=op (MemOp). The inner MemOp has
// three tagged fields: 0=operation (bool), 1=index (Witness), 2=value
// (Witness). With EncodingStrategy::Array both are positional fixarrays.
func (m *MemoryOp[T, E]) UnmarshalReader(r *msgpackutil.Reader) error {
	return msgpackutil.ReadStruct(r, m.decode)
}

func (m *MemoryOp[T, E]) decode(tag int, r *msgpackutil.Reader) error {
	switch tag {
	case 0:
		v, err := r.ReadUint()
		if err != nil {
			return err
		}
		m.BlockID = uint32(v)
		return nil
	case 1:
		return msgpackutil.ReadStruct(r, func(fieldTag int, r *msgpackutil.Reader) error {
			switch fieldTag {
			case 0:
				var err error
				m.IsWrite, err = r.ReadBool()
				return err
			case 1:
				return m.Index.UnmarshalReader(r)
			case 2:
				return m.Value.UnmarshalReader(r)
			default:
				return fmt.Errorf("mem_op: unknown field tag %d", fieldTag)
			}
		})
	default:
		return fmt.Errorf("memory_op: unknown field tag %d", tag)
	}
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
