package blackboxfunc

import (
	"fmt"
	"sunspot/go/acir/msgpackutil"
	shr "sunspot/go/acir/shared"

	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/math/uints"
	"github.com/consensys/gnark/std/permutation/sha2"
)

type SHA256Compression[T shr.ACIRField, E constraint.Element] struct {
	Inputs     [16]FunctionInput[T]
	HashValues [8]FunctionInput[T]
	Outputs    [8]shr.Witness
}

func (a *SHA256Compression[T, E]) decode(tag int, r *msgpackutil.Reader) error {
	switch tag {
	case 0:
		return readFunctionInputArray(r, a.Inputs[:])
	case 1:
		return readFunctionInputArray(r, a.HashValues[:])
	case 2:
		return shr.ReadWitnessArray(r, a.Outputs[:])
	default:
		return fmt.Errorf("Sha256Compression: unknown field tag %d", tag)
	}
}

func (a *SHA256Compression[T, E]) Equals(other BlackBoxFunction[E]) bool {
	value, ok := other.(*SHA256Compression[T, E])
	if !ok || len(a.Inputs) != len(value.Inputs) || len(a.HashValues) != len(value.HashValues) {
		return false
	}

	for i := 0; i < 16; i++ {
		if !a.Inputs[i].Equals(&value.Inputs[i]) {
			return false
		}
	}

	for i := 0; i < 8; i++ {
		if !a.HashValues[i].Equals(&value.HashValues[i]) {
			return false
		}
	}

	for i := 0; i < 8; i++ {
		if a.Outputs[i] != value.Outputs[i] {
			return false
		}
	}

	return true
}

func (a *SHA256Compression[T, E]) Define(api frontend.Builder[E], witnesses map[shr.Witness]frontend.Variable) error {
	uapi, err := uints.New[uints.U32](api)
	var old_state [8]uints.U32
	for i := 0; i < 8; i++ {
		variable, err := a.HashValues[i].ToVariable(witnesses)
		if err != nil {
			return err
		}
		old_state[i] = uapi.ValueOf(variable)
	}

	var inputs [64]uints.U8
	for i := 0; i < 16; i++ {
		variable, err := a.Inputs[i].ToVariable(witnesses)
		if err != nil {
			return err
		}
		copy(inputs[i*4:i*4+4], uapi.UnpackMSB(uapi.ValueOf(variable)))
	}

	if err != nil {
		return err
	}
	new_hash := sha2.Permute(uapi, old_state, inputs)
	for i := 0; i < 8; i++ {
		uapi.AssertEq(new_hash[i], uapi.ValueOf(witnesses[a.Outputs[i]]))
	}

	return nil
}
