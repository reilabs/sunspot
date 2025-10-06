package blackboxfunc

import (
	"encoding/binary"
	"io"
	shr "nr-groth16/acir/shared"

	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/lookup/logderivlookup"
	"github.com/consensys/gnark/std/math/uints"
	"github.com/google/btree"
)

type AES128Encrypt[T shr.ACIRField, E constraint.Element] struct {
	Inputs  []FunctionInput[T]
	Iv      [16]FunctionInput[T]
	Key     [16]FunctionInput[T]
	Outputs []shr.Witness
}

func (a *AES128Encrypt[T, E]) UnmarshalReader(r io.Reader) error {
	InputsNum := uint64(0)
	if err := binary.Read(r, binary.LittleEndian, &InputsNum); err != nil {
		return err
	}
	for i := uint64(0); i < InputsNum; i++ {
		var input FunctionInput[T]
		if err := input.UnmarshalReader(r); err != nil {
			return err
		}
		a.Inputs = append(a.Inputs, input)
	}

	for i := uint32(0); i < 16; i++ {
		if err := a.Iv[i].UnmarshalReader(r); err != nil {
			return err
		}
	}

	for i := uint32(0); i < 16; i++ {
		if err := a.Key[i].UnmarshalReader(r); err != nil {
			return err
		}
	}

	OutputsNum := uint64(0)
	if err := binary.Read(r, binary.LittleEndian, &OutputsNum); err != nil {
		return err
	}
	for i := uint64(0); i < OutputsNum; i++ {
		var witness shr.Witness
		if err := witness.UnmarshalReader(r); err != nil {
			return err
		}
		a.Outputs = append(a.Outputs, witness)
	}

	return nil
}

func (a *AES128Encrypt[T, E]) Equals(other BlackBoxFunction[E]) bool {
	value, ok := other.(*AES128Encrypt[T, E])
	if !ok || len(a.Inputs) != len(value.Inputs) {
		return false
	}
	for i := range a.Inputs {
		if !a.Inputs[i].Equals(&value.Inputs[i]) {
			return false
		}
	}

	for i := uint32(0); i < 16; i++ {
		if !a.Iv[i].Equals(&value.Iv[i]) {
			return false
		}
	}

	for i := uint32(0); i < 16; i++ {
		if !a.Key[i].Equals(&value.Key[i]) {
			return false
		}
	}

	if len(a.Outputs) != len(value.Outputs) {
		return false
	}
	for i := range a.Outputs {
		if !a.Outputs[i].Equals(&value.Outputs[i]) {
			return false
		}
	}

	return true
}

func (a *AES128Encrypt[T, E]) FillWitnessTree(tree *btree.BTree) bool {
	return tree != nil
}

func (a *AES128Encrypt[T, E]) Define(api frontend.Builder[E], witnesses map[shr.Witness]frontend.Variable) error {
	uapi, err := uints.NewBinaryField[uints.U32](api)
	if err != nil {
		return err
	}
	t0 := logderivlookup.New(api)
	t1 := logderivlookup.New(api)
	t2 := logderivlookup.New(api)
	t3 := logderivlookup.New(api)
	for i := 0; i < 256; i++ {
		t0.Insert(TE[0][i])
		t1.Insert(TE[1][i])
		t2.Insert(TE[2][i])
		t3.Insert(TE[3][i])
	}
	key, err := a.expandKey(api, witnesses, t0, t1, t2, t3)
	api.Println(uapi.ToValue(key[3]))
	if err != nil {
		return err
	}

	return nil
}

func (a *AES128Encrypt[T, E]) expandKey(api frontend.Builder[E], witnesses map[shr.Witness]frontend.Variable, TE0, TE1, TE2, TE3 logderivlookup.Table) ([60]uints.U32, error) {
	var rk [60]uints.U32

	var keyBytes [16]uints.U8
	uapi, err := uints.NewBinaryField[uints.U32](api)

	if err != nil {
		return rk, err
	}

	for i := range a.Key {
		val, err := a.Key[i].ToVariable(witnesses)
		if err != nil {
			return rk, err
		}
		keyBytes[i] = uapi.ByteValueOf(val)
	}

	rk[0] = uapi.PackMSB(keyBytes[0:4]...)
	rk[1] = uapi.PackMSB(keyBytes[4:8]...)
	rk[2] = uapi.PackMSB(keyBytes[8:12]...)
	rk[3] = uapi.PackMSB(keyBytes[12:16]...)

	start := 0
	for i := range 10 {
		temp := rk[3+start]

		te2index := uapi.ToValue(uapi.And(uapi.Rshift(temp, 16), uints.NewU32(0xff)))
		te2Val := uapi.ValueOf(TE2.Lookup(te2index)[0])

		te3index := uapi.ToValue(uapi.And(uapi.Rshift(temp, 8), uints.NewU32(0xff)))
		te3Val := uapi.ValueOf(TE3.Lookup(te3index)[0])

		te0index := uapi.ToValue(uapi.And(temp, uints.NewU32(0xff)))
		te0Val := uapi.ValueOf(TE0.Lookup(te0index)[0])

		te1index := uapi.ToValue(uapi.Rshift(temp, 24))
		te1Val := uapi.ValueOf(TE1.Lookup(te1index)[0])

		rk[4+start] = uapi.Xor(rk[0+start], uapi.And(te2Val, uints.NewU32(0xff000000)),
			uapi.And(te3Val, uints.NewU32(0x00ff0000)),
			uapi.And(te0Val, uints.NewU32(0x0000ff00)),
			uapi.And(te1Val, uints.NewU32(0x000000ff)),
			uints.NewU32(RCON[i]),
		)

		rk[5+start] = uapi.Xor(rk[1+start], rk[4+start])
		rk[6+start] = uapi.Xor(rk[2+start], rk[5+start])
		rk[7+start] = uapi.Xor(rk[3+start], rk[6+start])

		start += 4
	}
	return rk, nil
}

func ToVariables(uapi uints.BinaryField[uints.U32], ints []uints.U32) []frontend.Variable {
	ret := make([]frontend.Variable, len(ints))

	for i := range ret {
		ret[i] = uapi.ToValue(ints[i])
	}
	return ret
}
