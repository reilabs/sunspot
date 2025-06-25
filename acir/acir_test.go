package acir

import (
	"nr-groth16/bn254"
	"testing"

	ecc_bn254 "github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark/backend/groth16"
)

func TestACIRSumABExecuted(t *testing.T) {
	acir, err := LoadACIR[*bn254.BN254Field]("../noir-samples/sum_a_b/target/sum_a_b.json")
	if err != nil {
		t.Fatalf("Failed to load ACIR: %v", err)
	}

	ccs, err := acir.Compile()
	if err != nil {
		t.Fatalf("Failed to compile ACIR: %v", err)
	}

	pk, vk, err := groth16.Setup(ccs)
	if err != nil {
		t.Fatalf("Failed to setup Groth16: %v", err)
	}

	witness, err := acir.GetWitness("../noir-samples/sum_a_b/target/sum_a_b.gz", ecc_bn254.ID.ScalarField())
	if err != nil {
		t.Fatalf("Failed to get witness: %v", err)
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

func TestACIRLinearEquationExecuted(t *testing.T) {
	acir, err := LoadACIR[*bn254.BN254Field]("../noir-samples/linear_equation/target/linear_equation.json")
	if err != nil {
		t.Fatalf("Failed to load ACIR: %v", err)
	}

	ccs, err := acir.Compile()
	if err != nil {
		t.Fatalf("Failed to compile ACIR: %v", err)
	}

	pk, vk, err := groth16.Setup(ccs)
	if err != nil {
		t.Fatalf("Failed to setup Groth16: %v", err)
	}

	witness, err := acir.GetWitness("../noir-samples/linear_equation/target/linear_equation.gz", ecc_bn254.ID.ScalarField())
	if err != nil {
		t.Fatalf("Failed to get witness: %v", err)
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

func TestACIRSquareEquationExecuted(t *testing.T) {
	acir, err := LoadACIR[*bn254.BN254Field]("../noir-samples/square_equation/target/square_equation.json")
	if err != nil {
		t.Fatalf("Failed to load ACIR: %v", err)
	}

	ccs, err := acir.Compile()
	if err != nil {
		t.Fatalf("Failed to compile ACIR: %v", err)
	}

	pk, vk, err := groth16.Setup(ccs)
	if err != nil {
		t.Fatalf("Failed to setup Groth16: %v", err)
	}

	witness, err := acir.GetWitness("../noir-samples/square_equation/target/square_equation.gz", ecc_bn254.ID.ScalarField())
	if err != nil {
		t.Fatalf("Failed to get witness: %v", err)
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

func TestACIRRockPaperScissorsExecuted(t *testing.T) {
	acir, err := LoadACIR[*bn254.BN254Field]("../noir-samples/rock_paper_scissors/target/rock_paper_scissors.json")
	if err != nil {
		t.Fatalf("Failed to load ACIR: %v", err)
	}

	ccs, err := acir.Compile()
	if err != nil {
		t.Fatalf("Failed to compile ACIR: %v", err)
	}

	pk, vk, err := groth16.Setup(ccs)
	if err != nil {
		t.Fatalf("Failed to setup Groth16: %v", err)
	}

	witness, err := acir.GetWitness("../noir-samples/rock_paper_scissors/target/rock_paper_scissors.gz", ecc_bn254.ID.ScalarField())
	if err != nil {
		t.Fatalf("Failed to get witness: %v", err)
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

func TestACIRPolynomial(t *testing.T) {
	acir, err := LoadACIR[*bn254.BN254Field]("../noir-samples/polynomial/target/polynomial.json")
	if err != nil {
		t.Fatalf("Failed to load ACIR: %v", err)
	}

	ccs, err := acir.Compile()
	if err != nil {
		t.Fatalf("Failed to compile ACIR: %v", err)
	}

	pk, vk, err := groth16.Setup(ccs)
	if err != nil {
		t.Fatalf("Failed to setup Groth16: %v", err)
	}

	witness, err := acir.GetWitness("../noir-samples/polynomial/target/polynomial.gz", ecc_bn254.ID.ScalarField())
	if err != nil {
		t.Fatalf("Failed to get witness: %v", err)
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
