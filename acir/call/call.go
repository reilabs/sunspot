package call

import (
	"encoding/binary"
	"io"
	exp "nr-groth16/acir/expression"
	shr "nr-groth16/acir/shared"

	"github.com/consensys/gnark/constraint"
)

type Call[T shr.ACIRField, E constraint.Element] struct {
	ID        uint32
	Inputs    []shr.Witness
	Outputs   []shr.Witness
	Predicate *exp.Expression[T, E]
}

func (c *Call[T, E]) UnmarshalReader(r io.Reader) error {
	if err := binary.Read(r, binary.LittleEndian, &c.ID); err != nil {
		return err
	}

	var numInputs uint64
	if err := binary.Read(r, binary.LittleEndian, &numInputs); err != nil {
		return err
	}
	c.Inputs = make([]shr.Witness, numInputs)
	if err := binary.Read(r, binary.LittleEndian, &c.Inputs); err != nil {
		return err
	}

	var numOutputs uint64
	if err := binary.Read(r, binary.LittleEndian, &numOutputs); err != nil {
		return err
	}
	c.Outputs = make([]shr.Witness, numOutputs)
	if err := binary.Read(r, binary.LittleEndian, &c.Outputs); err != nil {
		return err
	}

	var predicateExists uint8
	if err := binary.Read(r, binary.LittleEndian, &predicateExists); err != nil {
		return err
	}
	if predicateExists == 1 {
		c.Predicate = new(exp.Expression[T, E])
		if err := c.Predicate.UnmarshalReader(r); err != nil {
			return err
		}
	}

	return nil
}

func (c *Call[T, E]) Equals(other *Call[T, E]) bool {
	if c.ID != other.ID {
		return false
	}

	if len(c.Inputs) != len(other.Inputs) || len(c.Outputs) != len(other.Outputs) {
		return false
	}

	for i := range c.Inputs {
		if c.Inputs[i] != other.Inputs[i] {
			return false
		}
	}

	for i := range c.Outputs {
		if c.Outputs[i] != other.Outputs[i] {
			return false
		}
	}

	if (c.Predicate == nil) != (other.Predicate == nil) {
		return false
	}

	if c.Predicate != nil && !c.Predicate.Equals(other.Predicate) {
		return false
	}

	return true
}
