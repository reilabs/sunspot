package acir

import (
	"encoding/json"
	"math/big"
	"nr-groth16/bn254"
	"os"
	"testing"

	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/consensys/gnark/backend/groth16"
)

func TestACIRSumAB(t *testing.T) {
	file, err := os.Open("../noir-samples/sum_a_b/target/sum_a_b.json")
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}

	defer file.Close()

	decoder := json.NewDecoder(file)

	var acir ACIR[bn254.BN254Field]
	if err := decoder.Decode(&acir); err != nil {
		t.Fatalf("Failed to decode ACIR: %v", err)
	}

	ccs, err := acir.Compile()
	if err != nil {
		t.Fatalf("Failed to compile ACIR: %v", err)
	}

	pk, vk, err := groth16.Setup(ccs)
	if err != nil {
		t.Fatalf("Failed to setup Groth16: %v", err)
	}

	inputs := map[string]*big.Int{
		"x": big.NewInt(1),
		"y": big.NewInt(2),
		"z": big.NewInt(3),
	}

	witness, err := acir.GenerateWitness(inputs, fr.Modulus())
	if err != nil {
		t.Fatalf("Failed to generate witness: %v", err)
	}

	proof, err := groth16.Prove(ccs, pk, witness)
	if err != nil {
		t.Fatalf("Failed to generate proof: %v", err)
	}

	publicWitness, err := witness.Public()
	if err != nil {
		t.Fatalf("Failed to get public witness: %v", err)
	}

	if err := groth16.Verify(proof, vk, publicWitness); err != nil {
		t.Fatalf("Verification failed: %v", err)
	} else {
		t.Logf("Verification succeeded!")
	}
}
