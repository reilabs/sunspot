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

type Blake2s[T shr.ACIRField, E constraint.Element] struct {
	Inputs  []FunctionInput[T]
	Outputs [32]shr.Witness
}

func (a *Blake2s[T, E]) UnmarshalReader(r io.Reader) error {
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

func (a *Blake2s[T, E]) Equals(other BlackBoxFunction[E]) bool {
	value, ok := other.(*Blake2s[T, E])
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

func (a *Blake2s[T, E]) FillWitnessTree(tree *btree.BTree) bool {
	return tree != nil
}

// Define builds the circuit constraints.
func (a *Blake2s[T, E]) Define(api frontend.Builder[E], witnesses map[shr.Witness]frontend.Variable) error {
	uapi, err := uints.New[uints.U32](api)
	if err != nil {
		return fmt.Errorf("unable to create 32 bit operation API for blake2s")
	}
	bytes_api, err := uints.NewBytes(api)
	if err != nil {
		return fmt.Errorf("unable to create byte operation API for blake2s")
	}

	data := make([]uints.U8, len(a.Inputs))

	for i := range len(a.Inputs) {
		input_var, err := a.Inputs[i].ToVariable(witnesses)
		if err != nil {
			return fmt.Errorf("blake input %d not found in witness map", i)
		}
		data[i] = bytes_api.ValueOf(input_var)
	}
	constrained_output, err := Blake2Permute(api, uapi, SplitIntoBlocks16(uapi, data), len(a.Inputs))
	if err != nil {
		return err
	}

	for i := range 4 {
		output_word := make([]uints.U8, 4)
		for j := range 4 {
			output_word[j] = bytes_api.ValueOf(witnesses[a.Outputs[4*i+j]])
		}

		uapi.AssertEq(constrained_output[i], uapi.PackLSB(output_word...))
	}

	return nil
}
func Blake2Permute(api frontend.API, uapi *uints.BinaryField[uints.U32], data [][]uints.U32, ll int) ([]uints.U32, error) {
	var err error
	h := GetIV()
	h[0] = uapi.Xor(h[0], uints.NewU32(0x01010000), uints.NewU32(32))

	if len(data) > 1 {
		for i := 0; i < len(data)-1; i++ {
			t := ((i + 1) * 16)
			h, err = F(api, uapi, h, data[i], uints.NewU64(uint64(t)), false)
			if err != nil {
				return nil, fmt.Errorf("error in F function in blake 2: %s", err)
			}
		}
	}
	h, err = F(api, uapi, h, data[len(data)-1], uints.NewU64(uint64(ll)), true)
	if err != nil {
		return nil, fmt.Errorf("error in the final F function in blake 2: %s", err)
	}
	return h[0:4], nil

}

func F(api frontend.API, uapi *uints.BinaryField[uints.U32], h []uints.U32, m []uints.U32, t uints.U64, f bool) ([]uints.U32, error) {
	v := make([]uints.U32, 16)
	uapi64, err := uints.NewBinaryField[uints.U64](api)

	if err != nil {
		return nil, fmt.Errorf("unable to create 64 bit operation api in blake2s")
	}
	copy(v[0:8], h[0:8])
	copy(v[8:16], GetIV())

	tBytes := uapi64.UnpackLSB(t)
	lowerBytes := uapi.PackLSB(tBytes[0:4]...)
	upperBytes := uapi.PackLSB(tBytes[4:8]...)
	v[12] = uapi.Xor(v[12], lowerBytes)
	v[13] = uapi.Xor(v[13], upperBytes)

	if f {
		v[14] = uapi.Xor(v[14], uints.NewU32(0xFFFFFFFF))
	}

	for i := 0; i < 10; i++ {
		s := make([]uint8, 16)
		copy(s[0:16], GetSigma(i)[0:16])

		v = G(uapi, v, 0, 4, 8, 12, m[s[0]], m[s[1]])
		v = G(uapi, v, 1, 5, 9, 13, m[s[2]], m[s[3]])
		v = G(uapi, v, 2, 6, 10, 14, m[s[4]], m[s[5]])
		v = G(uapi, v, 3, 7, 11, 15, m[s[6]], m[s[7]])

		v = G(uapi, v, 0, 5, 10, 15, m[s[8]], m[s[9]])
		v = G(uapi, v, 1, 6, 11, 12, m[s[10]], m[s[11]])
		v = G(uapi, v, 2, 7, 8, 13, m[s[12]], m[s[13]])
		v = G(uapi, v, 3, 4, 9, 14, m[s[14]], m[s[15]])

	}

	for i := 0; i < 8; i++ {
		h[i] = uapi.Xor(h[i], v[i], v[i+8])
	}

	return h[0:8], nil
}

// FUNCTION G( v[0..15], a, b, c, d, x, y )
// |
// |   v[a] := (v[a] + v[b] + x) mod 2**w
// |   v[d] := (v[d] ^ v[a]) >>> R1
// |   v[c] := (v[c] + v[d])     mod 2**w
// |   v[b] := (v[b] ^ v[c]) >>> R2
// |   v[a] := (v[a] + v[b] + y) mod 2**w
// |   v[d] := (v[d] ^ v[a]) >>> R3
// |   v[c] := (v[c] + v[d])     mod 2**w
// |   v[b] := (v[b] ^ v[c]) >>> R4
// |
// |   RETURN v[0..15]
// |
// END FUNCTION.
func G(uapi *uints.BinaryField[uints.U32], v []uints.U32, a, b, c, d uint32, x, y uints.U32) []uints.U32 {
	v[a] = uapi.Add(v[a], v[b], x)
	v[d] = RightRotation(uapi, uapi.Xor(v[d], v[a]), 16)

	v[c] = uapi.Add(v[c], v[d])
	v[b] = RightRotation(uapi, uapi.Xor(v[b], v[c]), 12)

	v[a] = uapi.Add(uapi.Add(v[a], v[b]), y)
	v[d] = RightRotation(uapi, uapi.Xor(v[d], v[a]), 8)

	v[c] = uapi.Add(v[c], v[d])
	v[b] = RightRotation(uapi, uapi.Xor(v[b], v[c]), 7)

	return v
}

func RightRotation(uapi *uints.BinaryField[uints.U32], x uints.U32, n int) uints.U32 {
	return uapi.Lrot(x, 32-n)
}

func LshiftU32(bf *uints.BinaryField[uints.U32], a uints.U32, c int) uints.U32 {
	if c <= 0 {
		return a
	}
	if c >= 32 {
		return uints.NewU32(0)
	}
	rot := bf.Lrot(a, c)
	mask := ^((uint32(1) << c) - 1) // upper (32-c) bits = 1
	return bf.And(rot, uints.NewU32(mask))
}

var SIGMA = [][]uint8{
	{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
	{14, 10, 4, 8, 9, 15, 13, 6, 1, 12, 0, 2, 11, 7, 5, 3},
	{11, 8, 12, 0, 5, 2, 15, 13, 10, 14, 3, 6, 7, 1, 9, 4},
	{7, 9, 3, 1, 13, 12, 11, 14, 2, 6, 5, 10, 4, 0, 15, 8},
	{9, 0, 5, 7, 2, 4, 10, 15, 14, 1, 11, 12, 6, 8, 3, 13},
	{2, 12, 6, 10, 0, 11, 8, 3, 4, 13, 7, 5, 15, 14, 1, 9},
	{12, 5, 1, 15, 14, 13, 4, 10, 0, 7, 6, 3, 9, 2, 8, 11},
	{13, 11, 7, 14, 12, 1, 3, 9, 5, 0, 15, 4, 8, 6, 2, 10},
	{6, 15, 14, 9, 11, 3, 0, 8, 12, 2, 13, 7, 1, 4, 10, 5},
	{10, 2, 8, 4, 7, 6, 1, 5, 15, 11, 9, 14, 3, 12, 13, 0},
}

func GetIV() []uints.U32 {
	return uints.NewU32Array([]uint32{
		0x6A09E667,
		0xBB67AE85,
		0x3C6EF372,
		0xA54FF53A,
		0x510E527F,
		0x9B05688C,
		0x1F83D9AB,
		0x5BE0CD19,
	})
}

// GetSigma returns the SIGMA slice for a given round
func GetSigma(round int) []uint8 {
	idx := round % len(SIGMA) // wrap around if round >= 10
	result := make([]uint8, len(SIGMA[idx]))
	copy(result, SIGMA[idx])
	return result
}

// SplitIntoBlocks16 splits data into 16-word blocks (64 bytes),
// each word is a big-endian uint32. Pads with zeros if needed.
func SplitIntoBlocks16(uapi_32 *uints.BinaryField[uints.U32], data []uints.U8) [][]uints.U32 {
	const wordsPerBlock = 16
	const bytesPerBlock = wordsPerBlock * 4

	// Round up to nearest multiple of 64
	paddedLen := ((len(data) + bytesPerBlock - 1) / bytesPerBlock) * bytesPerBlock
	padded := make([]uints.U8, paddedLen)
	copy(padded, data)

	blocks := make([][]uints.U32, paddedLen/bytesPerBlock)

	for i := 0; i < len(blocks); i++ {
		block := make([]uints.U32, wordsPerBlock)
		for j := 0; j < wordsPerBlock; j++ {
			base := i*bytesPerBlock + j*4

			word := uapi_32.PackMSB(
				getByte(padded, base+3), // least significant
				getByte(padded, base+2),
				getByte(padded, base+1),
				getByte(padded, base), // most significant
			)
			block[j] = word
		}
		blocks[i] = block
	}

	return blocks
}

// helper: safe byte fetch
func getByte(data []uints.U8, idx int) uints.U8 {
	if idx < 0 || idx >= len(data) || data[idx].Val == nil {
		return uints.NewU8(0)
	}
	return data[idx]
}
