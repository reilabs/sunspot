package blackboxfunc

import (
	"encoding/binary"
	"fmt"
	"io"
	shr "nr-groth16/acir/shared"

	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/lookup/logderivlookup"
	"github.com/consensys/gnark/std/math/uints"
	"github.com/google/btree"
)

const AES_BLOCK_SIZE = 16 // in bytes
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

	var IV [16]uints.U8

	for i := range IV {
		iv_val, err := a.Iv[i].ToVariable(witnesses)
		if err != nil {
			return err
		}
		IV[i] = uapi.ByteValueOf(iv_val)
	}

	plaintext := make([]uints.U8, len(a.Inputs))
	for i := range plaintext {
		pt_val, err := a.Inputs[i].ToVariable(witnesses)
		if err != nil {
			return err
		}
		plaintext[i] = uapi.ByteValueOf(pt_val)
	}

	key, err := a.expandKey(api, witnesses, t0, t1, t2, t3)
	if err != nil {
		return err
	}
	ciphertext, err := CBCEncrypt(api, key, IV, plaintext, t0, t1, t2, t3)
	if err != nil {
		return err
	}

	for i := range a.Outputs {
		api.AssertIsEqual(uapi.Value(ciphertext[i]), witnesses[a.Outputs[i]])
	}

	return nil
}

// CBCEncrypt: encrypts data in CBC mode in-place and returns the padded ciphertext.
func CBCEncrypt(
	api frontend.API,
	key [60]uints.U32,
	iv [16]uints.U8,
	plaintext []uints.U8,
	te0, te1, te2, te3 logderivlookup.Table,
) ([]uints.U8, error) {

	bytes, err := uints.NewBytes(api)
	if err != nil {
		return nil, err
	}
	plaintext = pad(plaintext)
	ciphertext := make([]uints.U8, len(plaintext))

	prevBlock := iv

	for blockStart := 0; blockStart < len(plaintext); blockStart += 16 {
		blockEnd := blockStart + 16
		var block [16]uints.U8

		for i := 0; i < 16; i++ {
			block[i] = bytes.Xor(plaintext[blockStart+i], prevBlock[i])
		}
		blockCipher, err := encrypt(api, key, block, te0, te1, te2, te3)
		if err != nil {
			return nil, fmt.Errorf("CBCEncrypt: block %d encrypt failed: %w", blockStart/16, err)
		}

		copy(ciphertext[blockStart:blockEnd], blockCipher[:])

		prevBlock = blockCipher
	}

	return ciphertext, nil
}

func encrypt(
	api frontend.API,
	key [60]uints.U32,
	inputs [16]uints.U8,
	te0, te1, te2, te3 logderivlookup.Table,
) ([16]uints.U8, error) {
	var outputs [16]uints.U8
	rk := key[:]
	uapi, err := uints.NewBinaryField[uints.U32](api)
	if err != nil {
		return outputs, err
	}

	const keyRounds = 10

	s0 := uapi.Xor(uapi.PackMSB(inputs[0:4]...), rk[0])
	s1 := uapi.Xor(uapi.PackMSB(inputs[4:8]...), rk[1])
	s2 := uapi.Xor(uapi.PackMSB(inputs[8:12]...), rk[2])
	s3 := uapi.Xor(uapi.PackMSB(inputs[12:16]...), rk[3])

	output := inputs

	r := keyRounds >> 1

	var t0 uints.U32
	var t1 uints.U32
	var t2 uints.U32
	var t3 uints.U32

	for {
		t0 = uapi.Xor(
			uapi.ValueOf(te0.Lookup(uapi.ToValue(uapi.Rshift(s0, 24)))[0]),
			uapi.ValueOf(te1.Lookup(uapi.ToValue(uapi.And(uapi.Rshift(s1, 16), uints.NewU32(0xff))))[0]),
			uapi.ValueOf(te2.Lookup(uapi.ToValue(uapi.And(uapi.Rshift(s2, 8), uints.NewU32(0xff))))[0]),
			uapi.ValueOf(te3.Lookup(uapi.ToValue(uapi.And(s3, uints.NewU32(0xff))))[0]),
			rk[4],
		)

		t1 = uapi.Xor(
			uapi.ValueOf(te0.Lookup(uapi.ToValue(uapi.Rshift(s1, 24)))[0]),
			uapi.ValueOf(te1.Lookup(uapi.ToValue(uapi.And(uapi.Rshift(s2, 16), uints.NewU32(0xff))))[0]),
			uapi.ValueOf(te2.Lookup(uapi.ToValue(uapi.And(uapi.Rshift(s3, 8), uints.NewU32(0xff))))[0]),
			uapi.ValueOf(te3.Lookup(uapi.ToValue(uapi.And(s0, uints.NewU32(0xff))))[0]),
			rk[5],
		)

		t2 = uapi.Xor(
			uapi.ValueOf(te0.Lookup(uapi.ToValue(uapi.Rshift(s2, 24)))[0]),
			uapi.ValueOf(te1.Lookup(uapi.ToValue(uapi.And(uapi.Rshift(s3, 16), uints.NewU32(0xff))))[0]),
			uapi.ValueOf(te2.Lookup(uapi.ToValue(uapi.And(uapi.Rshift(s0, 8), uints.NewU32(0xff))))[0]),
			uapi.ValueOf(te3.Lookup(uapi.ToValue(uapi.And(s1, uints.NewU32(0xff))))[0]),
			rk[6],
		)

		t3 = uapi.Xor(
			uapi.ValueOf(te0.Lookup(uapi.ToValue(uapi.Rshift(s3, 24)))[0]),
			uapi.ValueOf(te1.Lookup(uapi.ToValue(uapi.And(uapi.Rshift(s0, 16), uints.NewU32(0xff))))[0]),
			uapi.ValueOf(te2.Lookup(uapi.ToValue(uapi.And(uapi.Rshift(s1, 8), uints.NewU32(0xff))))[0]),
			uapi.ValueOf(te3.Lookup(uapi.ToValue(uapi.And(s2, uints.NewU32(0xff))))[0]),
			rk[7],
		)

		rk = rk[8:]

		r--
		if r == 0 {
			break
		}

		s0 = uapi.Xor(
			uapi.ValueOf(te0.Lookup(uapi.ToValue(uapi.Rshift(t0, 24)))[0]),
			uapi.ValueOf(te1.Lookup(uapi.ToValue(uapi.And(uapi.Rshift(t1, 16), uints.NewU32(0xff))))[0]),
			uapi.ValueOf(te2.Lookup(uapi.ToValue(uapi.And(uapi.Rshift(t2, 8), uints.NewU32(0xff))))[0]),
			uapi.ValueOf(te3.Lookup(uapi.ToValue(uapi.And(t3, uints.NewU32(0xff))))[0]),
			rk[0],
		)

		s1 = uapi.Xor(
			uapi.ValueOf(te0.Lookup(uapi.ToValue(uapi.Rshift(t1, 24)))[0]),
			uapi.ValueOf(te1.Lookup(uapi.ToValue(uapi.And(uapi.Rshift(t2, 16), uints.NewU32(0xff))))[0]),
			uapi.ValueOf(te2.Lookup(uapi.ToValue(uapi.And(uapi.Rshift(t3, 8), uints.NewU32(0xff))))[0]),
			uapi.ValueOf(te3.Lookup(uapi.ToValue(uapi.And(t0, uints.NewU32(0xff))))[0]),
			rk[1],
		)

		s2 = uapi.Xor(
			uapi.ValueOf(te0.Lookup(uapi.ToValue(uapi.Rshift(t2, 24)))[0]),
			uapi.ValueOf(te1.Lookup(uapi.ToValue(uapi.And(uapi.Rshift(t3, 16), uints.NewU32(0xff))))[0]),
			uapi.ValueOf(te2.Lookup(uapi.ToValue(uapi.And(uapi.Rshift(t0, 8), uints.NewU32(0xff))))[0]),
			uapi.ValueOf(te3.Lookup(uapi.ToValue(uapi.And(t1, uints.NewU32(0xff))))[0]),
			rk[2],
		)

		s3 = uapi.Xor(
			uapi.ValueOf(te0.Lookup(uapi.ToValue(uapi.Rshift(t3, 24)))[0]),
			uapi.ValueOf(te1.Lookup(uapi.ToValue(uapi.And(uapi.Rshift(t0, 16), uints.NewU32(0xff))))[0]),
			uapi.ValueOf(te2.Lookup(uapi.ToValue(uapi.And(uapi.Rshift(t1, 8), uints.NewU32(0xff))))[0]),
			uapi.ValueOf(te3.Lookup(uapi.ToValue(uapi.And(t2, uints.NewU32(0xff))))[0]),
			rk[3],
		)
	}

	s0 = uapi.Xor(
		uapi.And(uapi.ValueOf(te2.Lookup(uapi.ToValue(uapi.Rshift(t0, 24)))[0]), uints.NewU32(0xff000000)),
		uapi.And(uapi.ValueOf(te3.Lookup(uapi.ToValue(uapi.And(uapi.Rshift(t1, 16), uints.NewU32(0xff))))[0]), uints.NewU32(0x00ff0000)),
		uapi.And(uapi.ValueOf(te0.Lookup(uapi.ToValue(uapi.And(uapi.Rshift(t2, 8), uints.NewU32(0xff))))[0]), uints.NewU32(0x0000ff00)),
		uapi.And(uapi.ValueOf(te1.Lookup(uapi.ToValue(uapi.And(t3, uints.NewU32(0xff))))[0]), uints.NewU32(0x000000ff)),
		rk[0],
	)
	copy(output[0:4], uapi.UnpackMSB(s0))

	s1 = uapi.Xor(
		uapi.And(uapi.ValueOf(te2.Lookup(uapi.ToValue(uapi.Rshift(t1, 24)))[0]), uints.NewU32(0xff000000)),
		uapi.And(uapi.ValueOf(te3.Lookup(uapi.ToValue(uapi.And(uapi.Rshift(t2, 16), uints.NewU32(0xff))))[0]), uints.NewU32(0x00ff0000)),
		uapi.And(uapi.ValueOf(te0.Lookup(uapi.ToValue(uapi.And(uapi.Rshift(t3, 8), uints.NewU32(0xff))))[0]), uints.NewU32(0x0000ff00)),
		uapi.And(uapi.ValueOf(te1.Lookup(uapi.ToValue(uapi.And(t0, uints.NewU32(0xff))))[0]), uints.NewU32(0x000000ff)),
		rk[1],
	)
	copy(output[4:8], uapi.UnpackMSB(s1))

	s2 = uapi.Xor(
		uapi.And(uapi.ValueOf(te2.Lookup(uapi.ToValue(uapi.Rshift(t2, 24)))[0]), uints.NewU32(0xff000000)),
		uapi.And(uapi.ValueOf(te3.Lookup(uapi.ToValue(uapi.And(uapi.Rshift(t3, 16), uints.NewU32(0xff))))[0]), uints.NewU32(0x00ff0000)),
		uapi.And(uapi.ValueOf(te0.Lookup(uapi.ToValue(uapi.And(uapi.Rshift(t0, 8), uints.NewU32(0xff))))[0]), uints.NewU32(0x0000ff00)),
		uapi.And(uapi.ValueOf(te1.Lookup(uapi.ToValue(uapi.And(t1, uints.NewU32(0xff))))[0]), uints.NewU32(0x000000ff)),
		rk[2],
	)
	copy(output[8:12], uapi.UnpackMSB(s2))

	s3 = uapi.Xor(
		uapi.And(uapi.ValueOf(te2.Lookup(uapi.ToValue(uapi.Rshift(t3, 24)))[0]), uints.NewU32(0xff000000)),
		uapi.And(uapi.ValueOf(te3.Lookup(uapi.ToValue(uapi.And(uapi.Rshift(t0, 16), uints.NewU32(0xff))))[0]), uints.NewU32(0x00ff0000)),
		uapi.And(uapi.ValueOf(te0.Lookup(uapi.ToValue(uapi.And(uapi.Rshift(t1, 8), uints.NewU32(0xff))))[0]), uints.NewU32(0x0000ff00)),
		uapi.And(uapi.ValueOf(te1.Lookup(uapi.ToValue(uapi.And(t2, uints.NewU32(0xff))))[0]), uints.NewU32(0x000000ff)),
		rk[3],
	)
	copy(output[12:16], uapi.UnpackMSB(s3))

	return output, nil
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

func pad(input []uints.U8) []uints.U8 {
	sz := len(input)
	add := 16 - (sz % 16)
	v := make([]uints.U8, 0, sz+add)
	v = append(v, input...)
	for i := 0; i < add; i++ {
		v = append(v, uints.NewU8(uint8(add)))
	}
	return v
}
