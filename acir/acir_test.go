package acir

import (
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"nr-groth16/bn254"
	"os"
	"testing"

	ecc_bn254 "github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/rs/zerolog"
)

func TestACIRSumAB(t *testing.T) {
	file, err := os.Open("../noir-samples/sum_a_b/target/sum_a_b.json")
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}

	defer file.Close()

	decoder := json.NewDecoder(file)

	var acir ACIR[*bn254.BN254Field]
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

	// inputs := map[string]*big.Int{
	// 	"x": big.NewInt(1),
	// 	"y": big.NewInt(2),
	// 	"z": big.NewInt(3),
	// }

	// witness, err := acir.GenerateWitness(inputs, fr.Modulus())
	// if err != nil {
	// 	t.Fatalf("Failed to generate witness: %v", err)
	// }

	fmt.Println("Reading witnesses from file")

	witness, err := acir.GetWitnessFromFile("../noir-samples/sum_a_b/target/sum_a_b", fr.Modulus())
	if err != nil {
		t.Fatalf("Failed to get witness from file: %v", err)
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
	file, err := os.Open("../noir-samples/sum_a_b/target/sum_a_b.json")
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}

	defer file.Close()

	decoder := json.NewDecoder(file)

	var acir ACIR[*bn254.BN254Field]
	if err := decoder.Decode(&acir); err != nil {
		t.Fatalf("Failed to decode ACIR: %v", err)
	}

	witnessStack, err := LoadWitnessFromFile[*bn254.BN254Field](
		"../noir-samples/sum_a_b/target/sum_a_b.gz",
		ecc_bn254.ID.ScalarField(),
	)
	if err != nil {
		t.Fatalf("Failed to load witness from file: %v", err)
	}

	ccs, err := acir.CompileExecuted(witnessStack)
	if err != nil {
		t.Fatalf("Failed to compile ACIR: %v", err)
	}

	pk, vk, err := groth16.Setup(ccs)
	if err != nil {
		t.Fatalf("Failed to setup Groth16: %v", err)
	}

	witness, err := acir.GetWitness(witnessStack, ecc_bn254.ID.ScalarField())
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

func TestACIRLinearEquation(t *testing.T) {
	zerolog.SetGlobalLevel(zerolog.TraceLevel)

	log.SetOutput(zerolog.ConsoleWriter{
		Out: os.Stdout})
	file, err := os.Open("../noir-samples/linear_equation/target/linear_equation.json")
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}

	defer file.Close()

	decoder := json.NewDecoder(file)

	var acir ACIR[*bn254.BN254Field]
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

	witnessMap := map[string]*big.Int{
		"x": big.NewInt(1),
		"a": big.NewInt(1),
		"b": big.NewInt(0),
		"y": big.NewInt(1),
	}

	witness, err := acir.GenerateWitness(witnessMap, fr.Modulus())
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

func TestACIRSquareEquation(t *testing.T) {
	file, err := os.Open("../noir-samples/square_equation/target/square_equation.json")
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}

	defer file.Close()

	decoder := json.NewDecoder(file)

	var acir ACIR[*bn254.BN254Field]
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
		"x": big.NewInt(-1),
		"a": big.NewInt(1),
		"b": big.NewInt(2),
		"c": big.NewInt(1),
	}

	witness, err := acir.GenerateWitness(inputs, fr.Modulus())
	if err != nil {
		t.Fatalf("Failed to generate witness: %v", err)
	}

	// witness, err := acir.GetWitnessFromFile("../noir-samples/sum_a_b/target/sum_a_b", fr.Modulus())
	// if err != nil {
	// 	t.Fatalf("Failed to get witness from file: %v", err)
	// }

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

/*func TestACIRRockPaperScissorsExecuted(t *testing.T) {
	witnessStack, err := LoadWitnessFromFile[*bn254.BN254Field](
		"../noir-samples/rock_paper_scissors/target/rock_paper_scissors.gz",
		ecc_bn254.ID.ScalarField(),
	)
	if err != nil {
		t.Fatalf("Failed to load witness from file: %v", err)
	}
	file, err := os.Open("../noir-samples/rock_paper_scissors/target/rock_paper_scissors.json")
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}

	defer file.Close()

	decoder := json.NewDecoder(file)

	var acir ACIR[*bn254.BN254Field]
	if err := decoder.Decode(&acir); err != nil {
		t.Fatalf("Failed to decode ACIR: %v", err)
	}

	ccs, err := acir.CompileExecuted(witnessStack)
	if err != nil {
		t.Fatalf("Failed to compile ACIR: %v", err)
	}

	pk, vk, err := groth16.Setup(ccs)
	if err != nil {
		t.Fatalf("Failed to setup Groth16: %v", err)
	}

	witness, err := acir.GetWitness(witnessStack, ecc_bn254.ID.ScalarField())
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
