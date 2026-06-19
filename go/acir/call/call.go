package call

import (
	"encoding/json"
	exp "github.com/reilabs/sunspot/go/acir/expression"
	"github.com/reilabs/sunspot/go/acir/msgpackutil"
	ops "github.com/reilabs/sunspot/go/acir/opcodes"
	shr "github.com/reilabs/sunspot/go/acir/shared"

	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
)

type Call[T shr.ACIRField, E constraint.Element] struct {
	ID        uint32
	Inputs    []shr.Witness
	Outputs   []shr.Witness
	Predicate exp.Expression[T, E]
}

func (c *Call[T, E]) UnmarshalReader(r *msgpackutil.Reader) error {
	return msgpackutil.ReadStruct(r, "Call", []msgpackutil.Field{
		{Name: "id", Decode: func(r *msgpackutil.Reader) error {
			v, err := r.ReadUint()
			if err != nil {
				return err
			}
			c.ID = uint32(v)
			return nil
		}},
		{Name: "inputs", Decode: func(r *msgpackutil.Reader) error { return msgpackutil.ReadVec(r, &c.Inputs) }},
		{Name: "outputs", Decode: func(r *msgpackutil.Reader) error { return msgpackutil.ReadVec(r, &c.Outputs) }},
		{Name: "predicate", Decode: c.Predicate.UnmarshalReader},
	})
}
func (c *Call[T, E]) Equals(other ops.Opcode[E]) bool {
	value, ok := other.(*Call[T, E])
	if !ok || c.ID != value.ID {
		return false
	}

	if len(c.Inputs) != len(value.Inputs) || len(c.Outputs) != len(value.Outputs) {
		return false
	}

	for i := range c.Inputs {
		if c.Inputs[i] != value.Inputs[i] {
			return false
		}
	}

	for i := range c.Outputs {
		if c.Outputs[i] != value.Outputs[i] {
			return false
		}
	}

	return c.Predicate.Equals(&value.Predicate)
}

func (o *Call[T, E]) Define(api frontend.Builder[E], witnesses map[shr.Witness]frontend.Variable) error {
	return nil
}

func (c *Call[T, E]) MarshalJSON() ([]byte, error) {
	stringMap := make(map[string]interface{})
	stringMap["circuit_call"] = c
	return json.Marshal(stringMap)
}

func (*Call[T, E]) SerdeName() string { return "Call" }
