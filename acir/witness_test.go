package acir

import (
	"nr-groth16/bn254"
	"testing"

	ecc_bn254 "github.com/consensys/gnark-crypto/ecc/bn254"
)

func TestACIRWitnessSumAB(t *testing.T) {
	witnessStack, err := LoadWitnessFromFile[*bn254.BN254Field](
		"../noir-samples/sum_a_b/target/sum_a_b.gz",
		ecc_bn254.ID.ScalarField(),
	)
	if err != nil {
		t.Fatalf("Failed to load witness from file: %v", err)
	}

	if len(witnessStack.ItemStack) == 0 {
		t.Fatal("Witness stack is empty")
	}

	t.Logf("Loaded witness stack with %d items", len(witnessStack.ItemStack))
}

func TestACIRWitnessSquareEquation(t *testing.T) {
	witnessStack, err := LoadWitnessFromFile[*bn254.BN254Field](
		"../noir-samples/square_equation/target/square_equation.gz",
		ecc_bn254.ID.ScalarField(),
	)
	if err != nil {
		t.Fatalf("Failed to load witness from file: %v", err)
	}

	if len(witnessStack.ItemStack) == 0 {
		t.Fatal("Witness stack is empty")
	}

	t.Logf("Loaded witness stack with %d items", len(witnessStack.ItemStack))
}

func TestACIRWitnessRockPaperScissors(t *testing.T) {
	witnessStack, err := LoadWitnessFromFile[*bn254.BN254Field](
		"../noir-samples/rock_paper_scissors/target/rock_paper_scissors.gz",
		ecc_bn254.ID.ScalarField(),
	)
	if err != nil {
		t.Fatalf("Failed to load witness from file: %v", err)
	}

	t.Logf("Loaded witness stack with %d items", len(witnessStack.ItemStack))
}
