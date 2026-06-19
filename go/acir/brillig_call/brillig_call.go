package brillig_call

import (
	"encoding/json"
	exp "github.com/reilabs/sunspot/go/acir/expression"
	"github.com/reilabs/sunspot/go/acir/msgpackutil"
	ops "github.com/reilabs/sunspot/go/acir/opcodes"
	shr "github.com/reilabs/sunspot/go/acir/shared"

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
	return msgpackutil.ReadStruct(r, "BrilligCall", []msgpackutil.Field{
		{Name: "id", Decode: func(r *msgpackutil.Reader) error {
			v, err := r.ReadUint()
			if err != nil {
				return err
			}
			b.ID = uint32(v)
			return nil
		}},
		{Name: "inputs", Decode: func(r *msgpackutil.Reader) error { return msgpackutil.ReadVec(r, &b.Inputs) }},
		{Name: "outputs", Decode: func(r *msgpackutil.Reader) error { return msgpackutil.ReadVec(r, &b.Outputs) }},
		{Name: "predicate", Decode: b.Predicate.UnmarshalReader},
	})
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

func (*BrilligCall[T, E]) SerdeName() string { return "BrilligCall" }
