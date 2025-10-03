package poseidon2

import (
	"math/big"

	"github.com/consensys/gnark/frontend"
)

func Permute(api frontend.API, state []frontend.Variable) {
	RC := parseTwoDimensionArray(rc)
	matrix := parseOneDimensionArray(internal_matrix)

	matFull4(api, state)

	for i := range 4 {
		addRoundConstants(api, state, RC, i)
		Sbox(api, state)
		matFull4(api, state)

	}

	for i := 4; i < 60; i++ {
		state[0] = api.Add(state[0], RC[i][0])
		state[0] = SboxOne(api, state[0])
		internalMMultiplication(api, state, matrix)
	}

	for i := 60; i < 64; i++ {
		addRoundConstants(api, state, RC, i)
		Sbox(api, state)
		matFull4(api, state)
	}

}
func internalMMultiplication(api frontend.API, input []frontend.Variable, matrix []*big.Int) {
	sum := frontend.Variable(0)

	for i := range input {
		sum = api.Add(sum, input[i])
	}

	for i := 0; i < len(input); i++ {
		input[i] = api.Mul(input[i], matrix[i])
		input[i] = api.Add(input[i], sum)
	}

}

func addRoundConstants(api frontend.API, state []frontend.Variable, RC [][]*big.Int, round int) {
	for i := range 4 {
		state[i] = api.Add(state[i], RC[round][i])
	}
}

func matFull4(api frontend.API, state []frontend.Variable) {
	if len(state) != 4 {
		panic("matFull4 requires state of length 4")
	}
	t0 := api.Add(state[0], state[1])
	t1 := api.Add(state[2], state[3])

	t2 := api.Add(Double(api, state[1]), t1)
	t3 := api.Add(Double(api, state[3]), t0)

	t4 := api.Add(Double(api, Double(api, t1)), t3)
	t5 := api.Add(Double(api, Double(api, t0)), t2)

	t6 := api.Add(t3, t5)
	t7 := api.Add(t2, t4)

	state[0] = t6
	state[1] = t5
	state[2] = t7
	state[3] = t4
}

func Sbox(api frontend.API, state []frontend.Variable) {
	for i := range state {
		state[i] = SboxOne(api, state[i])
	}
}

func SboxOne(api frontend.API, x frontend.Variable) frontend.Variable {
	x2 := api.Mul(x, x)
	x4 := api.Mul(x2, x2)
	x5 := api.Mul(x, x4)
	return x5
}

func Double(api frontend.API, x frontend.Variable) frontend.Variable {
	return api.Add(x, x)
}
