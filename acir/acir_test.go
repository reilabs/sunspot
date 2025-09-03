package acir

import (
	"fmt"
	"nr-groth16/bn254"
	"testing"

	ecc_bn254 "github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark/backend/groth16"
)

func testProveAndVerify(t *testing.T, acirPath string, witnessPath string) {
	acir, err := LoadACIR[*bn254.BN254Field](acirPath)

	fmt.Println("acir", acir)
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

	witness, err := acir.GetWitness(witnessPath, ecc_bn254.ID.ScalarField())
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

func TestACIRSumABExecuted(t *testing.T) {
	testProveAndVerify(
		t,
		"../noir-samples/expressions/sum_a_b/target/sum_a_b.json",
		"../noir-samples/expressions/sum_a_b/target/sum_a_b.gz",
	)
}

func TestACIRLinearEquationExecuted(t *testing.T) {
	testProveAndVerify(
		t,
		"../noir-samples/expressions/linear_equation/target/linear_equation.json",
		"../noir-samples/expressions/linear_equation/target/linear_equation.gz",
	)
}

func TestACIRSquareEquationExecuted(t *testing.T) {
	testProveAndVerify(
		t,
		"../noir-samples/expressions/square_equation/target/square_equation.json",
		"../noir-samples/expressions/square_equation/target/square_equation.gz",
	)
}

func TestACIRRockPaperScissorsExecuted(t *testing.T) {
	testProveAndVerify(
		t,
		"../noir-samples/expressions/rock_paper_scissors/target/rock_paper_scissors.json",
		"../noir-samples/expressions/rock_paper_scissors/target/rock_paper_scissors.gz",
	)
}

func TestACIRPolynomial(t *testing.T) {
	testProveAndVerify(
		t,
		"../noir-samples/expressions/polynomial/target/polynomial.json",
		"../noir-samples/expressions/polynomial/target/polynomial.gz",
	)
}

// func TestACIRKeccakF1600(t *testing.T) {
// 	testProveAndVerify(
// 		t,
// 		"../noir-samples/black_box_functions/keccak_f1600/target/keccak_f1600.json",
// 		"../noir-samples/black_box_functions/keccak_f1600/target/keccak_f1600.gz",
// 	)
// }

func TestACIRRange(t *testing.T) {
	testProveAndVerify(
		t,
		"../noir-samples/black_box_functions/range/target/range.json",
		"../noir-samples/black_box_functions/range/target/range.gz",
	)
}

func TestACIRLCChecker(t *testing.T) {
	testProveAndVerify(
		t,
		"../noir-samples/expressions/lcchecker/target/lcchecker.json",
		"../noir-samples/expressions/lcchecker/target/lcchecker.gz",
	)
}

func TestACIRAnd(t *testing.T) {
	testProveAndVerify(
		t,
		"../noir-samples/black_box_functions/and/target/and.json",
		"../noir-samples/black_box_functions/and/target/and.gz",
	)
}

func TestACIRXor(t *testing.T) {
	testProveAndVerify(
		t,
		"../noir-samples/black_box_functions/xor/target/xor.json",
		"../noir-samples/black_box_functions/xor/target/xor.gz",
	)
}

// func TestACIRZKVoting(t *testing.T) {
// 	testProveAndVerify(
// 		t,
// 		"../noir-samples/real_world/zk-noir-voting/circuits/target/circuits.json",
// 		"../noir-samples/real_world/zk-noir-voting/circuits/target/circuits.gz",
// 	)
// }

// func TestACIRProveKitBasic(t *testing.T) {
// 	testProveAndVerify(
// 		t,
// 		"../noir-samples/real_world/ProveKit/noir-examples/basic/target/basic.json",
// 		"../noir-samples/real_world/ProveKit/noir-examples/basic/target/basic.gz",
// 	)
// }

// func TestACIRProveKitBasic2(t *testing.T) {
// 	testProveAndVerify(
// 		t,
// 		"../noir-samples/real_world/ProveKit/noir-examples/basic-2/target/basic.json",
// 		"../noir-samples/real_world/ProveKit/noir-examples/basic-2/target/basic.gz",
// 	)
// }

// func TestACIRProveKitBasic3(t *testing.T) {
// 	testProveAndVerify(
// 		t,
// 		"../noir-samples/real_world/ProveKit/noir-examples/basic-3/target/basic.json",
// 		"../noir-samples/real_world/ProveKit/noir-examples/basic-3/target/basic.gz",
// 	)
// }

// func TestACIRPoseidonVar(t *testing.T) {
// 	testProveAndVerify(
// 		t,
// 		"../noir-samples/real_world/ProveKit/noir-examples/poseidon-var/target/basic.json",
// 		"../noir-samples/real_world/ProveKit/noir-examples/poseidon-var/target/basic.gz",
// 	)
// }

/*func TestACIRAES128Encrypt(t *testing.T) {
	acir, err := LoadACIR[*bn254.BN254Field]("../noir-samples/aes128encrypt/target/aes128encrypt.json")
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

	witness, err := acir.GetWitness("../noir-samples/aes128encrypt/target/aes128encrypt.gz", ecc_bn254.ID.ScalarField())
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
}*/
