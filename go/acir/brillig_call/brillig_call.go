package brillig_call

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

type BrilligCall[T shr.ACIRField, E constraint.Element] struct {
	ID        uint32
	Inputs    []BrilligInputs[T, E]
	Outputs   []BrilligOutputs
	Predicate exp.Expression[T, E]
}

func (b *BrilligCall[T, E]) UnmarshalReader(r *msgpackutil.Reader) error {
	return msgpackutil.ReadStruct(r, b.decode)
}

func (b *BrilligCall[T, E]) decode(tag int, r *msgpackutil.Reader) error {
	switch tag {
	case 0:
		v, err := r.ReadUint()
		if err != nil {
			return err
		}
		b.ID = uint32(v)
		return nil
	case 1:
		n, err := r.ReadArrayLen()
		if err != nil {
			return err
		}
		b.Inputs = make([]BrilligInputs[T, E], n)
		for i := 0; i < n; i++ {
			if err := b.Inputs[i].UnmarshalReader(r); err != nil {
				return err
			}
		}
		return nil
	case 2:
		n, err := r.ReadArrayLen()
		if err != nil {
			return err
		}
		b.Outputs = make([]BrilligOutputs, n)
		for i := 0; i < n; i++ {
			if err := b.Outputs[i].UnmarshalReader(r); err != nil {
				return err
			}
		}
		return nil
	case 3:
		return b.Predicate.UnmarshalReader(r)
	default:
		return fmt.Errorf("brillig_call: unknown field tag %d", tag)
	}
}
func (o *BrilligCall[T, E]) Equals(other ops.Opcode[E]) bool {
	// Function exists for purposes of satisfying trait bound
	// Trait function only used in tests that are not exercised on this type
	panic("unimplemented")
}
func (o *BrilligCall[T, E]) Define(api frontend.Builder[E], witnesses map[shr.Witness]frontend.Variable) error {
	// do nothing: brillig calls are unconstrained
	return nil
}

func (o *BrilligCall[T, E]) MarshalJSON() ([]byte, error) {
	stringMap := make(map[string]interface{})
	stringMap["brillig_call"] = o
	return json.Marshal(stringMap)
}
