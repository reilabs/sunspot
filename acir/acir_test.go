package acir

import (
	"encoding/json"
	"nr-groth16/bn254"
	"os"
	"testing"
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
}
