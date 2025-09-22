package blackboxfunc

import (
	"encoding/binary"
	"fmt"
	"io"
	"math/big"
	shr "nr-groth16/acir/shared"

	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
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

func AESShiftRows(state [16]frontend.Variable) [16]frontend.Variable {
	var shiftedState [16]frontend.Variable
	shiftedState[0], shiftedState[1], shiftedState[2], shiftedState[3] = state[0], state[1], state[2], state[3]
	shiftedState[4], shiftedState[5], shiftedState[6], shiftedState[7] = state[5], state[6], state[7], state[4]
	shiftedState[8], shiftedState[9], shiftedState[10], shiftedState[11] = state[10], state[11], state[8], state[9]
	shiftedState[12], shiftedState[13], shiftedState[14], shiftedState[15] = state[15], state[12], state[13], state[14]
	return shiftedState
}

func AESSubWord(api frontend.API, word [4]frontend.Variable) [4]frontend.Variable {
	var result [4]frontend.Variable
	for i := 0; i < 4; i++ {
		result[i] = word[i]
	}
	return result
}

func AESRotWord(api frontend.API, word [4]frontend.Variable) [4]frontend.Variable {
	var result [4]frontend.Variable
	result[0] = word[1]
	result[1] = word[2]
	result[2] = word[3]
	result[3] = word[0]
	return result
}

func AESRcon(api frontend.API, round int) frontend.Variable {
	rcon := []frontend.Variable{
		big.NewInt(0x01000000),
		big.NewInt(0x02000000),
		big.NewInt(0x04000000),
		big.NewInt(0x08000000),
		big.NewInt(0x10000000),
		big.NewInt(0x20000000),
		big.NewInt(0x40000000),
		big.NewInt(0x80000000),
		big.NewInt(0x1B000000),
		big.NewInt(0x36000000),
	}
	return rcon[round]
}

func AESKeyExpansion(api frontend.API, key [16]frontend.Variable) [11][16]frontend.Variable {
	var roundKeys [11][16]frontend.Variable
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			roundKeys[0][i*4+j] = key[i*4+j]
		}
	}

	for round := 1; round <= 10; round++ {
		var temp [4]frontend.Variable
		for i := 0; i < 4; i++ {
			temp[i] = roundKeys[round-1][i*4+3]
		}

		temp = AESSubWord(api, temp)
		temp = AESRotWord(api, temp)

		temp[0] = api.Xor(temp[0], AESRcon(api, round))

		for i := 0; i < 4; i++ {
			roundKeys[round][i*4] = api.Xor(roundKeys[round-1][i*4], temp[i])
		}

		// Generate the other 3 words of the round key
		for col := 1; col < 4; col++ {
			for row := 0; row < 4; row++ {
				roundKeys[round][row*4+col] = api.Xor(roundKeys[round][row*4+col-1], roundKeys[round-1][row*4+col])
			}
		}
	}

	return roundKeys
}

func AESAddRoundKey(api frontend.API, state [16]frontend.Variable, roundKey [16]frontend.Variable) [16]frontend.Variable {
	var result [16]frontend.Variable
	for i := 0; i < 16; i++ {
		result[i] = api.Xor(state[i], roundKey[i])
	}
	return result
}

func (a *AES128Encrypt[T, E]) Define(api frontend.Builder[E], witnesses map[shr.Witness]frontend.Variable) error {
	numBlocks := len(a.Inputs) / 16
	var state [16]frontend.Variable
	var results [16]frontend.Variable
	var key [16]frontend.Variable
	for i := 0; i < 16; i++ {
		key_val, err := a.Key[i].ToVariable(witnesses)
		if err != nil {
			return err
		}
		key[i] = key_val
	}
	roundKeys := AESKeyExpansion(api, key)
	for i := 0; i < numBlocks; i++ {
		block := a.Inputs[i*16 : (i+1)*16]

		//
		// Initialize the state with the input data
		//
		for j := 0; j < 16; j++ {
			state[j] = block[j]
		}

		for j := 0; j < 16; j++ {
			if i == 0 {
				state[j] = api.Xor(state[j], a.Iv[j])
			} else {
				state[j] = api.Xor(state[j], results[j])
			}
		}

		state = AESAddRoundKey(api, state, roundKeys[0])

		for i := 0; i < 10; i++ {
			// Sub Bytes

			state = AESShiftRows(state)

			if i < 9 {
				// Mix Columns
			}

			state = AESAddRoundKey(api, state, roundKeys[i+1])
		}

		for j := 0; j < 16; j++ {
			results[j] = state[j]
		}

		for j := 0; j < 16; j++ {
			expected := witnesses[a.Outputs[i*16+j]]
			if expected == nil {
				return fmt.Errorf("expected output for AES128Encrypt at index %d is nil", i*16+j)
			}
			api.AssertIsEqual(results[j], expected)
		}
	}

	return nil
}
