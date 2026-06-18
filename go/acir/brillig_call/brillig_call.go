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
	return msgpackutil.ReadStruct(r, brilligCallSchema, b.decode)
}

func (b *BrilligCall[T, E]) decode(f msgpackutil.Field, r *msgpackutil.Reader) error {
	switch f.Tag {
	case 0:
		v, err := r.ReadUint()
		if err != nil {
			return err
		}
		b.ID = uint32(v)
		return nil
	case 1:
		return msgpackutil.ReadVec(r, &b.Inputs)
	case 2:
		return msgpackutil.ReadVec(r, &b.Outputs)
	case 3:
		return b.Predicate.UnmarshalReader(r)
	default:
		return fmt.Errorf("BrilligCall: unknown field %s", f)
	}
}

var brilligCallSchema = msgpackutil.NewSchema(map[string]int{
	"id": 0, "inputs": 1, "outputs": 2, "predicate": 3,
})
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
