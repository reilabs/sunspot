package blackboxfunc

import (
	"encoding/binary"
	"fmt"
	"io"
	shr "nr-groth16/acir/shared"

	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/math/uints"
	"github.com/google/btree"
)

type Blake3[T shr.ACIRField, E constraint.Element] struct {
	Inputs  []FunctionInput[T]
	Outputs [32]shr.Witness
}

func (a *Blake3[T, E]) UnmarshalReader(r io.Reader) error {
	NumInputs := uint64(0)
	if err := binary.Read(r, binary.LittleEndian, &NumInputs); err != nil {
		return err
	}

	a.Inputs = make([]FunctionInput[T], NumInputs)
	for i := uint64(0); i < NumInputs; i++ {
		if err := a.Inputs[i].UnmarshalReader(r); err != nil {
			return err
		}
	}

	if err := binary.Read(r, binary.LittleEndian, &a.Outputs); err != nil {
		return err
	}

	return nil
}

func (a *Blake3[T, E]) Equals(other BlackBoxFunction[E]) bool {
	value, ok := other.(*Blake3[T, E])
	if !ok || len(a.Inputs) != len(value.Inputs) {
		return false
	}
	for i := range a.Inputs {
		if !a.Inputs[i].Equals(&value.Inputs[i]) {
			return false
		}
	}

	for i := 0; i < 32; i++ {
		if a.Outputs[i] != value.Outputs[i] {
			return false
		}
	}

	return true
}

func (a *Blake3[T, E]) Define(api frontend.Builder[E], witnesses map[shr.Witness]frontend.Variable) error {
	panic("not yet implemented")
}

func (a *Blake3[T, E]) FillWitnessTree(tree *btree.BTree) bool {
	return tree != nil
}

func Blake3Compress(api frontend.API, uapi uints.BinaryField[uints.U32], h, m []uints.U32, t uints.U64, len, flags uints.U32) ([]uints.U32, error) {

	v := make([]uints.U32, 16)
	uapi64, err := uints.NewBinaryField[uints.U64](api)

	if err != nil {
		return nil, fmt.Errorf("unable to create 64 bit operation api in blake3")
	}
	copy(v[0:8], h[0:8])
	copy(v[8:11], GetIV()[0:3])

	tBytes := uapi64.UnpackLSB(t)
	lowerBytes := uapi.PackLSB(tBytes[0:4]...)
	upperBytes := uapi.PackLSB(tBytes[4:8]...)
	v[12] = lowerBytes
	v[13] = upperBytes
	v[14] = len
	v[15] = flags

	for range 7 {
		v = G(&uapi, v, 0, 4, 8, 12, m[0], m[1])
		v = G(&uapi, v, 1, 5, 9, 13, m[2], m[3])
		v = G(&uapi, v, 2, 6, 10, 14, m[4], m[5])
		v = G(&uapi, v, 3, 7, 11, 15, m[6], m[7])

		v = G(&uapi, v, 0, 5, 10, 15, m[8], m[9])
		v = G(&uapi, v, 1, 6, 11, 12, m[10], m[11])
		v = G(&uapi, v, 2, 7, 8, 13, m[12], m[13])
		v = G(&uapi, v, 3, 4, 9, 14, m[15], m[15])
		m = permuteMessage(m)
	}

	for i := range 8 {
		v[i] = uapi.Xor(v[i], v[i+8])
		v[i+8] = uapi.Xor(v[i+8], h[i])
	}

	return v, nil

}

func permuteMessage(input []uints.U32) []uints.U32 {
	perm := []int{2, 6, 3, 10, 7, 0, 4, 13, 1, 11, 12, 5, 9, 14, 15, 8}
	if len(input) != len(perm) {
		panic("input length must be 16")
	}
	output := make([]uints.U32, 16)
	for i, p := range perm {
		output[p] = input[i]
	}
	return output
}
