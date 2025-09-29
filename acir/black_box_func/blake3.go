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

const (
	OUT_LEN   = 32
	BLOCK_LEN = 64
	CHUNK_LEN = 1024
)

const (
	CHUNK_START = 1 << 0
	CHUNK_END   = 1 << 1
	PARENT      = 1 << 2
	ROOT        = 1 << 3
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
	bytesApi, err := uints.NewBytes(api)
	if err != nil {
		return fmt.Errorf("unable to create 64 bit operation API")
	}

	uapi, err := uints.NewBinaryField[uints.U32](api)
	if err != nil {
		return fmt.Errorf("unable to create 32 bit operation API")
	}

	uapi64, err := uints.NewBinaryField[uints.U64](api)
	if err != nil {
		return fmt.Errorf("unable to create 64 bit operation API")
	}

	data := make([]uints.U8, len(a.Inputs))

	for i := range len(a.Inputs) {
		input_var, err := a.Inputs[i].ToVariable(witnesses)
		if err != nil {
			return fmt.Errorf("blake input %d not found in witness map", i)
		}
		data[i] = bytesApi.ValueOf(input_var)
	}
	hasher := NewHasher()

	hasher.Update(api, *uapi, *uapi64, data)
	zeroByte := uints.NewU8(0)
	output := make([]uints.U8, 32)

	for i := range output {
		output[i] = zeroByte
	}

	hasher.Finalize(api, *uapi, output)

	for i := range 32 {
		bytesApi.AssertIsEqual(bytesApi.ValueOf(witnesses[a.Outputs[i]]), output[i])
	}
	return nil
}

func (a *Blake3[T, E]) FillWitnessTree(tree *btree.BTree) bool {
	return tree != nil
}

func Blake3Compress(api frontend.API, uapi uints.BinaryField[uints.U32], h [8]uints.U32, m [16]uints.U32, t uints.U64, len, flags uints.U32) ([8]uints.U32, error) {
	var v [16]uints.U32
	var ret [8]uints.U32
	uapi64, err := uints.NewBinaryField[uints.U64](api)
	if err != nil {
		return ret, fmt.Errorf("unable to create 64 bit operation api in blake3")
	}
	copy(v[0:8], h[0:8])
	IV := GetIV()
	for i := range 4 {
		v[8+i] = IV[i]
	}

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
		v = G(&uapi, v, 3, 4, 9, 14, m[14], m[15])
		m = permuteMessage(m)
	}

	for i := range 8 {
		v[i] = uapi.Xor(v[i], v[i+8])
		v[i+8] = uapi.Xor(v[i+8], h[i])
	}

	for i := range ret {
		ret[i] = v[i]
	}
	return ret, nil

}

func permuteMessage(input [16]uints.U32) [16]uints.U32 {
	perm := []int{2, 6, 3, 10, 7, 0, 4, 13, 1, 11, 12, 5, 9, 14, 15, 8}
	if len(input) != len(perm) {
		panic("input length must be 16")
	}
	var output [16]uints.U32
	for i, p := range perm {
		output[i] = input[p]
	}
	return output
}

type Output struct {
	inputChainingValue [8]uints.U32
	blockWords         [16]uints.U32
	counter            uint64
	blockLen           uints.U32
	flags              uints.U32
}

func (o *Output) ChainingValue(api frontend.API, uapi uints.BinaryField[uints.U32]) ([8]uints.U32, error) {
	output, err := Blake3Compress(api, uapi, o.inputChainingValue, o.blockWords, uints.NewU64(o.counter), o.blockLen, o.flags)
	if err != nil {
		return output, err
	}
	return output, nil
}
func (o *Output) RootOutputBytes(api frontend.API, uapi uints.BinaryField[uints.U32], out []uints.U8) error {
	outputBlockCounter := uint64(0)
	for i := 0; i < len(out); i += 2 * OUT_LEN {
		end := i + 2*OUT_LEN
		if end > len(out) {
			end = len(out)
		}

		o.blockWords = PadTo16Words(o.blockWords[:])
		words, err := Blake3Compress(api, uapi, o.inputChainingValue, o.blockWords, uints.NewU64(outputBlockCounter), o.blockLen, uapi.Or(o.flags, uints.NewU32(ROOT)))
		if err != nil {
			return err
		}

		blockBytes := make([]uints.U8, 16*4)
		for j, w := range words {
			copy(blockBytes[4*j:4*j+4], uapi.UnpackLSB(w))
		}
		copy(out[i:end], blockBytes[:end-i])
		outputBlockCounter++
	}
	return nil
}

// ---------------- ChunkState ----------------
type ChunkState struct {
	chainingValue    [8]uints.U32
	chunkCounter     uint64
	block            [BLOCK_LEN]uints.U8
	blockLen         int
	blocksCompressed int
	flags            uints.U32
}

func NewChunkState(keyWords [8]uints.U32, chunkCounter uint64, flags uints.U32) ChunkState {
	return ChunkState{chainingValue: keyWords, chunkCounter: chunkCounter, flags: flags}
}
func (c *ChunkState) Len() int {
	return BLOCK_LEN*c.blocksCompressed + c.blockLen
}
func (c *ChunkState) startFlag() uints.U32 {
	if c.blocksCompressed == 0 {
		return uints.NewU32(CHUNK_START)
	}
	return uints.NewU32(0)
}

func (c *ChunkState) Update(api frontend.API, uapi uints.BinaryField[uints.U32], input []uints.U8) error {
	var err error
	for len(input) > 0 {
		if c.blockLen == BLOCK_LEN {
			var blockWords [16]uints.U32
			for i := range blockWords {
				blockWords[i] = uapi.PackLSB(c.block[4*i : 4*i+4]...)
			}
			c.chainingValue, err = Blake3Compress(api, uapi, c.chainingValue, blockWords, uints.NewU64(c.chunkCounter), uints.NewU32(BLOCK_LEN), uapi.Or(c.flags, c.startFlag()))
			if err != nil {
				return err
			}
			c.blocksCompressed++
			c.blockLen = 0
			for i := range c.block {
				c.block[i] = uints.NewU8(0)
			}
		}
		want := BLOCK_LEN - c.blockLen
		take := want
		if take > len(input) {
			take = len(input)
		}
		copy(c.block[c.blockLen:], input[:take])
		c.blockLen += take
		input = input[take:]
	}
	return nil
}

func (c *ChunkState) Output(uapi uints.BinaryField[uints.U32]) *Output {
	var blockWords [16]uints.U32
	for i := range blockWords {
		blockWords[i] = uapi.PackLSB(c.block[4*i : 4*i+4]...)
	}
	return &Output{inputChainingValue: c.chainingValue, blockWords: blockWords, counter: c.chunkCounter, blockLen: uints.NewU32(uint32(c.blockLen)), flags: uapi.Or(c.flags, c.startFlag(), uints.NewU32(CHUNK_END))}
}

func parentOutput(uapi uints.BinaryField[uints.U32], left, right [8]uints.U32, keyWords [8]uints.U32, flags uints.U32) *Output {
	var blockWords [16]uints.U32
	copy(blockWords[0:8], left[:])
	copy(blockWords[8:16], right[:])
	return &Output{inputChainingValue: keyWords, blockWords: blockWords, counter: 0, blockLen: uints.NewU32(BLOCK_LEN), flags: uapi.Or(uints.NewU32(PARENT), flags)}
}

func parentCV(api frontend.API, uapi uints.BinaryField[uints.U32], left, right [8]uints.U32, keyWords [8]uints.U32, flags uints.U32) ([8]uints.U32, error) {
	chainingValue, err := parentOutput(uapi, left, right, keyWords, flags).ChainingValue(api, uapi)
	if err != nil {
		return chainingValue, err
	}
	return chainingValue, nil
}

// ---------------- Hasher ----------------
type Hasher struct {
	chunkState ChunkState
	keyWords   [8]uints.U32
	cvStack    [54][8]uints.U32
	cvStackLen int
	flags      uints.U32
}

func NewHasher() *Hasher {
	IV := GetIV()
	return &Hasher{
		chunkState: NewChunkState(IV, 0, uints.NewU32(0)), keyWords: [8]uints.U32(IV), flags: uints.NewU32(0)}
}
func (h *Hasher) pushStack(cv [8]uints.U32) {
	h.cvStack[h.cvStackLen] = cv
	h.cvStackLen++
}
func (h *Hasher) popStack() [8]uints.U32 {
	h.cvStackLen--
	return h.cvStack[h.cvStackLen]
}
func (h *Hasher) addChunkChainingValue(api frontend.API, uapi uints.BinaryField[uints.U32], newCV [8]uints.U32, totalChunks uint64) error {
	for totalChunks&1 == 0 {
		cvSlice, err := parentCV(api, uapi, h.popStack(), newCV, h.keyWords, h.flags)
		if err != nil {
			return err
		}
		newCV = [8]uints.U32(cvSlice)
		totalChunks >>= 1
	}
	h.pushStack(newCV)
	return nil
}

func (h *Hasher) Update(api frontend.API, uapi uints.BinaryField[uints.U32], uapi_64 uints.BinaryField[uints.U64], input []uints.U8) error {
	for len(input) > 0 {
		if h.chunkState.Len() == CHUNK_LEN {
			chunkCV, err := h.chunkState.Output(uapi).ChainingValue(api, uapi)
			if err != nil {
				return err
			}
			totalChunks := h.chunkState.chunkCounter + 1
			h.addChunkChainingValue(api, uapi, [8]uints.U32(chunkCV), totalChunks)
			h.chunkState = NewChunkState(h.keyWords, totalChunks, h.flags)
		}
		want := CHUNK_LEN - h.chunkState.Len()
		take := want

		if take > len(input) {
			take = len(input)
		}
		h.chunkState.Update(api, uapi, input[:take])
		input = input[take:]
	}
	return nil
}

func (h *Hasher) Finalize(api frontend.API, uapi uints.BinaryField[uints.U32], out []uints.U8) error {
	output := h.chunkState.Output(uapi)
	for i := h.cvStackLen - 1; i >= 0; i-- {
		chainingValue, err := output.ChainingValue(api, uapi)
		if err != nil {
			return err
		}
		output = parentOutput(uapi, h.cvStack[i], [8]uints.U32(chainingValue), h.keyWords, h.flags)
	}
	output.RootOutputBytes(api, uapi, out)
	return nil
}

func PadTo16Words(data []uints.U32) [16]uints.U32 {
	var padded [16]uints.U32

	// Fill from data
	for i := 0; i < 16; i++ {
		if data[i][0].Val == nil {
			padded[i] = uints.NewU32(0)
		} else {
			padded[i] = data[i]
		}
	}
	return padded
}
