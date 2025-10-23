package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"sunpot/acir"
	"sunpot/bn254"

	"github.com/consensys/gnark-crypto/ecc"
	ecc_bn254 "github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/constraint"
	"github.com/spf13/cobra"
)

var proveCmd = &cobra.Command{
	Use:   "prove [acir_file] [witness_file] [ccs_file] [pk_file]",
	Short: "Generate a Groth16 proof and public witness from a witness, CCS, and proving key",
	Args:  cobra.ExactArgs(4),
	RunE: func(cmd *cobra.Command, args []string) error {
		acirPath := args[0]
		witnessPath := args[1]
		ccsPath := args[2]
		pkPath := args[3]

		// ✅ Validate file extensions

		if filepath.Ext(acirPath) != ".json" {
			return fmt.Errorf("invalid input file: %s (must end with .json)", acirPath)
		}
		if filepath.Ext(witnessPath) != ".gz" {
			return fmt.Errorf("invalid witness file: %s (must end with .gz)", witnessPath)
		}
		if filepath.Ext(ccsPath) != ".ccs" {
			return fmt.Errorf("invalid CCS file: %s (must end with .ccs)", ccsPath)
		}
		if filepath.Ext(pkPath) != ".pk" {
			return fmt.Errorf("invalid proving key file: %s (must end with .pk)", pkPath)
		}

		fmt.Printf("Loading ACIR file: %s\n", acirPath)

		type E = constraint.U64
		type T = *bn254.BN254Field
		acir, err := acir.LoadACIR[T, E](acirPath)

		if err != nil {
			return fmt.Errorf("failed to load ACIR: %v", err)
		}

		fmt.Printf("🔧 Loading CCS: %s\n", ccsPath)
		ccsFile, err := os.Open(ccsPath)
		if err != nil {
			return fmt.Errorf("failed to open CCS: %v", err)
		}
		defer ccsFile.Close()

		ccs := groth16.NewCS(ecc.BN254)
		if _, err := ccs.ReadFrom(ccsFile); err != nil {
			return fmt.Errorf("failed to read CCS: %w", err)
		}

		fmt.Printf("🔑 Loading Proving Key: %s\n", pkPath)
		pkFile, err := os.Open(pkPath)
		if err != nil {
			return fmt.Errorf("failed to open proving key: %v", err)
		}
		defer pkFile.Close()

		pk := groth16.NewProvingKey(ecc.BN254)
		if _, err := pk.ReadFrom(pkFile); err != nil {
			return fmt.Errorf("failed to read proving key: %w", err)
		}

		fmt.Printf("📄 Loading Witness: %s\n", witnessPath)
		witness, err := acir.GetWitness(witnessPath, ecc_bn254.ID.ScalarField())
		if err != nil {
			return fmt.Errorf("failed to get witness: %v", err)
		}

		fmt.Println("⚙️  Generating Groth16 proof...")
		proof, err := groth16.Prove(ccs, pk, witness)
		if err != nil {
			return fmt.Errorf("failed to generate proof: %v", err)
		}

		// Derive output proof file path
		base := ccsPath[:len(ccsPath)-len(filepath.Ext(ccsPath))]

		proofPath := base + ".proof"

		fmt.Printf("💾 Writing proof to %s\n", proofPath)
		proofFile, err := os.Create(proofPath)
		if err != nil {
			return fmt.Errorf("failed to create proof file: %w", err)
		}
		defer proofFile.Close()

		if _, err := proof.WriteTo(proofFile); err != nil {
			return fmt.Errorf("failed to write proof: %w", err)
		}

		pwPath := base + ".pw"

		pw, err := witness.Public()
		if err != nil {
			return fmt.Errorf("failed to create public witness: %w", err)
		}

		fmt.Printf("💾 Writing public witness to %s\n", pwPath)
		pwFile, err := os.Create(pwPath)
		if err != nil {
			return fmt.Errorf("failed to create public witness file: %w", err)
		}
		defer pwFile.Close()

		if _, err := pw.WriteTo(pwFile); err != nil {
			return fmt.Errorf("failed to write public witness: %w", err)
		}

		fmt.Println("✅ Proof generation complete!")
		return nil
	},
}
