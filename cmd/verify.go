package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/backend/witness"
	"github.com/spf13/cobra"
)

var verifyCmd = &cobra.Command{
	Use:   "verify [verification_key_file]  [proof_file] [public_witness_file]",
	Short: "Verify a proof and public witness with a verification key",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		vkPath := args[0]
		proofPath := args[1]
		pubWitnessPath := args[2]

		if filepath.Ext(vkPath) != ".vk" {
			return fmt.Errorf("invalid verification key file: %s (must end with .vk)", vkPath)
		}

		fmt.Printf("🔑 Loading Proving Key: %s\n", vkPath)
		vkFile, err := os.Open(vkPath)
		if err != nil {
			return fmt.Errorf("failed to open verifying key: %v", err)
		}
		defer vkFile.Close()

		vk := groth16.NewVerifyingKey(ecc.BN254)
		if _, err := vk.ReadFrom(vkFile); err != nil {
			return fmt.Errorf("failed to read proving key: %w", err)
		}

		fmt.Printf("Loading Proof: %s\n", proofPath)

		proofFile, err := os.Open(proofPath)
		if err != nil {
			return fmt.Errorf("failed to open proof: %v", err)
		}
		defer vkFile.Close()

		proof := groth16.NewProof(ecc.BN254)

		if _, err := proof.ReadFrom(proofFile); err != nil {
			return fmt.Errorf("failed to read proof: %v", err)
		}

		fmt.Printf("Loading public witness: %s\n", pubWitnessPath)
		pubWitnessFile, err := os.Open(pubWitnessPath)
		if err != nil {
			return fmt.Errorf("failed to open public witness: %v", err)
		}
		defer pubWitnessFile.Close()

		publicWitness, err := witness.New(ecc.BN254.ScalarField())
		if err != nil {
			return fmt.Errorf("unable to initialise public witness: %w", err)
		}

		if _, err := publicWitness.ReadFrom(pubWitnessFile); err != nil {
			return fmt.Errorf("failed to read public witness: %w", err)
		}

		if err := groth16.Verify(proof, vk, publicWitness); err != nil {
			return fmt.Errorf("verification failed: %v", err)
		}
		fmt.Println("✅ Verification successful!")
		return nil
	},
}
