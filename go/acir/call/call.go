package call

import (
	"encoding/json"
	"fmt"
	exp "sunspot/go/acir/expression"
	"sunspot/go/acir/msgpackutil"
	ops "sunspot/go/acir/opcodes"
	shr "sunspot/go/acir/shared"

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
	return msgpackutil.ReadStruct(r, c.decode)
}

func (c *Call[T, E]) decode(tag int, r *msgpackutil.Reader) error {
	switch tag {
	case 0:
		v, err := r.ReadUint()
		if err != nil {
			return err
		}
		c.ID = uint32(v)
		return nil
	case 1:
		n, err := r.ReadArrayLen()
		if err != nil {
			return err
		}
		c.Inputs = make([]shr.Witness, n)
		for i := 0; i < n; i++ {
			if err := c.Inputs[i].UnmarshalReader(r); err != nil {
				return err
			}
		}
		return nil
	case 2:
		n, err := r.ReadArrayLen()
		if err != nil {
			return err
		}
		c.Outputs = make([]shr.Witness, n)
		for i := 0; i < n; i++ {
			if err := c.Outputs[i].UnmarshalReader(r); err != nil {
				return err
			}
		}
		return nil
	case 3:
		return c.Predicate.UnmarshalReader(r)
	default:
		return fmt.Errorf("call: unknown field tag %d", tag)
	}
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
