package main

import (
	"fmt"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
)

type MyCircuit struct {
	X frontend.Variable
	Y frontend.Variable
}

// This is what gnark uses to build the constraint system
func (c *MyCircuit) Define(api frontend.API) error {
	sum := api.Add(c.X, c.Y)
	api.AssertIsEqual(sum, 0)
	return nil
}

func main() {
	// Define the circuit constraints (without assigning values)
	var circuit MyCircuit

	ccs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)
	if err != nil {
		panic(err)
	}

	// Setup proving & verifying keys
	pk, vk, err := groth16.Setup(ccs)
	if err != nil {
		panic(err)
	}

	// Create witness with actual values
	assignment := MyCircuit{
		X: -3,
		Y: 3,
	}
	witness, err := frontend.NewWitness(&assignment, ecc.BN254.ScalarField())
	if err != nil {
		panic(err)
	}

	publicWitness, _ := witness.Public()

	// Generate the proof
	proof, err := groth16.Prove(ccs, pk, witness)
	if err != nil {
		panic(err)
	}

	// Verify the proof
	err = groth16.Verify(proof, vk, publicWitness)
	if err != nil {
		panic("❌ Verification failed!")
	} else {
		fmt.Println("✅ Verification succeeded!")
	}
}
