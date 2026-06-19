package blackboxfunc

import (
	"fmt"
	"github.com/reilabs/sunspot/go/acir/msgpackutil"
	shr "github.com/reilabs/sunspot/go/acir/shared"

	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/math/uints"
	"github.com/consensys/gnark/std/permutation/keccakf"
)

type Keccakf1600[T shr.ACIRField, E constraint.Element] struct {
	Inputs  [25]FunctionInput[T]
	Outputs [25]shr.Witness
}

func (a *Keccakf1600[T, E]) UnmarshalReader(r *msgpackutil.Reader) error {
	return msgpackutil.ReadStruct(r, "Keccakf1600", []msgpackutil.Field{
		{Name: "inputs", Decode: func(r *msgpackutil.Reader) error { return msgpackutil.ReadArrayInto(r, a.Inputs[:]) }},
		{Name: "outputs", Decode: func(r *msgpackutil.Reader) error { return msgpackutil.ReadArrayInto(r, a.Outputs[:]) }},
	})
}

func (a *Keccakf1600[T, E]) Equals(other BlackBoxFunction[E]) bool {
	value, ok := other.(*Keccakf1600[T, E])
	if !ok || len(a.Inputs) != len(value.Inputs) {
		return false
	}

	for i := 0; i < 25; i++ {
		if !a.Inputs[i].Equals(&value.Inputs[i]) {
			return false
		}
	}

	for i := 0; i < 25; i++ {
		if a.Outputs[i] != value.Outputs[i] {
			return false
		}
	}

	return true
}

func (a *Keccakf1600[T, E]) Define(api frontend.Builder[E], witnesses map[shr.Witness]frontend.Variable) error {
	uapi, err := uints.New[uints.U64](api)
	if err != nil {
		return err
	}
	var keccak_inputs [25]uints.U64
	for i, input := range a.Inputs {
		v, err := input.ToVariable(witnesses)
		if err != nil {
			return fmt.Errorf("unable to get input as variable, index %d", i)
		}
		keccak_inputs[i] = uapi.ValueOf(v)
	}

	var keccak_outputs [25]uints.U64
	for i, output := range a.Outputs {
		v := witnesses[output]
		keccak_outputs[i] = uapi.ValueOf(v)
	}

	constrained_outputs := keccakf.Permute(uapi, keccak_inputs)

	for i := 0; i < 25; i++ {
		uapi.AssertEq(constrained_outputs[i], keccak_outputs[i])
	}
	return nil
}

func (*Keccakf1600[T, E]) SerdeName() string { return "Keccakf1600" }
