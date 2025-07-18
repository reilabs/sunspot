package blackboxfunc

import (
	"encoding/binary"
	"fmt"
	"io"
	shr "nr-groth16/acir/shared"

	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/math/uints"
	"github.com/consensys/gnark/std/permutation/sha2"
)

type SHA256Compression[T shr.ACIRField] struct {
	Inputs     [16]FunctionInput[T]
	HashValues [8]FunctionInput[T]
	Outputs    [8]shr.Witness
}

func (a *SHA256Compression[T]) UnmarshalReader(r io.Reader) error {
	for i := 0; i < 16; i++ {
		if err := a.Inputs[i].UnmarshalReader(r); err != nil {
			return err
		}
	}

	for i := 0; i < 8; i++ {
		if err := a.HashValues[i].UnmarshalReader(r); err != nil {
			return err
		}
	}

	if err := binary.Read(r, binary.LittleEndian, &a.Outputs); err != nil {
		return err
	}

	return nil
}

func (a *SHA256Compression[T]) Equals(other *SHA256Compression[T]) bool {
	if len(a.Inputs) != len(other.Inputs) || len(a.HashValues) != len(other.HashValues) {
		return false
	}

	for i := 0; i < 16; i++ {
		if !a.Inputs[i].Equals(&other.Inputs[i]) {
			return false
		}
	}

	for i := 0; i < 8; i++ {
		if !a.HashValues[i].Equals(&other.HashValues[i]) {
			return false
		}
	}

	for i := 0; i < 8; i++ {
		if a.Outputs[i] != other.Outputs[i] {
			return false
		}
	}

	return true
}

func (a *SHA256Compression[T]) Define(api frontend.API, witnesses map[shr.Witness]frontend.Variable) error {
	var old_state [8]uints.U32
	for i := 0; i < 8; i++ {
		variable, err := a.HashValues[i].ToVariable(witnesses)
		if err != nil {
			return err
		}
		var values []uint32
		api.Compiler().ToCanonicalVariable(variable).Compress(&values)
		old_state[i] = uints.NewU32(values[0])
	}

	var inputs [64]uints.U8
	for i := 0; i < 16; i++ {
		variable, err := a.Inputs[i].ToVariable(witnesses)
		if err != nil {
			return err
		}
		var values []uint32
		api.Compiler().ToCanonicalVariable(variable).Compress(&values)

		inputs[i*4] = uints.NewU8(uint8(values[0]))
		inputs[i*4+1] = uints.NewU8(uint8(values[0] >> 8))
		inputs[i*4+2] = uints.NewU8(uint8(values[0] >> 16))
		inputs[i*4+3] = uints.NewU8(uint8(values[0] >> 24))
	}

	binaryField, err := uints.New[uints.U32](api)
	if err != nil {
		return err
	}
	new_hash := sha2.Permute(binaryField, old_state, inputs)
	for i := 0; i < 8; i++ {
		val, ok := api.ConstantValue(new_hash[i])
		if !ok {
			return fmt.Errorf("failed to get constant value for new_hash[%d]", i)
		}
		api.AssertIsEqual(val, witnesses[a.Outputs[i]])
	}

	return nil
}
