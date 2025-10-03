package acir

import (
	"nr-groth16/bn254"
	"testing"

	"github.com/consensys/gnark/constraint"

	ecc_bn254 "github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark/backend/groth16"
)

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

func TestACIRRange(t *testing.T) {
	testProveAndVerify(
		t,
		"../noir-samples/black_box_functions/range/target/range.json",
		"../noir-samples/black_box_functions/range/target/range.gz",
	)
}

func TestACIRMemory(t *testing.T) {
	testProveAndVerify(
		t,
		"../noir-samples/memory/target/memory.json",
		"../noir-samples/memory/target/memory.gz",
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

func TestACIRKeccakF1600(t *testing.T) {
	testProveAndVerify(
		t,
		"../noir-samples/black_box_functions/keccak_f1600/target/keccak_f1600.json",
		"../noir-samples/black_box_functions/keccak_f1600/target/keccak_f1600.gz",
	)
}

func TestACIRSHA256(t *testing.T) {
	testProveAndVerify(
		t,
		"../noir-samples/black_box_functions/sha256_compression/target/sha256.json",
		"../noir-samples/black_box_functions/sha256_compression/target/sha256.gz",
	)
}

func TestACIRSHA256Hash(t *testing.T) {
	testProveAndVerify(
		t,
		"../noir-samples/black_box_functions/sha256_hash/target/sha256_hash.json",
		"../noir-samples/black_box_functions/sha256_hash/target/sha256_hash.gz",
	)
}

func TestACIRBlake2s(t *testing.T) {
	testProveAndVerify(
		t,
		"../noir-samples/black_box_functions/blake2s/target/blake2s.json",
		"../noir-samples/black_box_functions/blake2s/target/blake2s.gz",
	)
}

func TestACIRBlake3(t *testing.T) {
	testProveAndVerify(
		t,
		"../noir-samples/black_box_functions/blake3/target/blake3.json",
		"../noir-samples/black_box_functions/blake3/target/blake3.gz",
	)
}

func TestACIREmbeddedCurveAdd(t *testing.T) {
	testProveAndVerify(
		t,
		"../noir-samples/black_box_functions/embedded_curve_add/target/embedded_curve_add.json",
		"../noir-samples/black_box_functions/embedded_curve_add/target/embedded_curve_add.gz",
	)
}
func TestACIRMultiscalarMultiplication(t *testing.T) {
	testProveAndVerify(
		t,
		"../noir-samples/black_box_functions/multiscalar_multiplication/target/multiscalar_multiplication.json",
		"../noir-samples/black_box_functions/multiscalar_multiplication/target/multiscalar_multiplication.gz",
	)
}

func TestACIRECDSASecp256k1(t *testing.T) {
	testProveAndVerify(
		t,
		"../noir-samples/black_box_functions/ecdsa_secp256k1/target/ecdsa_secp256k1.json",
		"../noir-samples/black_box_functions/ecdsa_secp256k1/target/ecdsa_secp256k1.gz",
	)
}

func TestACIRECDSASecp256k1Failing(t *testing.T) {
	testProveAndVerify(
		t,
		"../noir-samples/black_box_functions/ecdsa_secp256k1_failing/target/ecdsa_secp256k1_failing.json",
		"../noir-samples/black_box_functions/ecdsa_secp256k1_failing/target/ecdsa_secp256k1_failing.gz",
	)
}

func TestACIRECDSASecp256r1(t *testing.T) {
	testProveAndVerify(
		t,
		"../noir-samples/black_box_functions/ecdsa_secp256r1/target/ecdsa_secp256r1.json",
		"../noir-samples/black_box_functions/ecdsa_secp256r1/target/ecdsa_secp256r1.gz",
	)
}

func TestACIRPoseidon2(t *testing.T) {
	testProveAndVerify(
		t,
		"../noir-samples/black_box_functions/poseidon2/target/poseidon2.json",
		"../noir-samples/black_box_functions/poseidon2/target/poseidon2.gz",
	)
}

func TestACIRAES128(t *testing.T) {
	testProveAndVerify(
		t,
		"../noir-samples/black_box_functions/aes128encrypt/target/aes128encrypt.json",
		"../noir-samples/black_box_functions/aes128encrypt/target/aes128encrypt.gz",
	)
}

func TestACIRLCChecker(t *testing.T) {
	testProveAndVerify(
		t,
		"../noir-samples/expressions/lcchecker/target/lcchecker.json",
		"../noir-samples/expressions/lcchecker/target/lcchecker.gz",
	)
}

func TestACIRProveKitBasic(t *testing.T) {
	testProveAndVerify(
		t,
		"../noir-samples/real_world/provekit_basic/target/provekit_basic.json",
		"../noir-samples/real_world/provekit_basic/target/provekit_basic.gz",
	)
}

// Helper function for testing files,
// Provide circuit and witness path and compile to r1cs, proves and verifies in groth16
func testProveAndVerify(t *testing.T, acirPath string, witnessPath string) {
	type E = constraint.U64
	acir, err := LoadACIR[*bn254.BN254Field, E](acirPath)

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
