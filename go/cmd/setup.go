package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup [ccs_file]",
	Short: "Generate a proving key (pk) and verifying key (vk) from a CCS file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ccsPath := args[0]
		// ✅ Assert correct file extension
		if filepath.Ext(ccsPath) != ".ccs" {
			return fmt.Errorf("invalid input file: %s (must end with .ccs)", ccsPath)
		}

		fmt.Printf("🔧 Loading CCS file: %s\n", ccsPath)

		// Open CCS file
		f, err := os.Open(ccsPath)
		if err != nil {
			return fmt.Errorf("failed to load CCS: %v", err)
		}
		defer f.Close()

		// Read CCS from file
		ccs := groth16.NewCS(ecc.BN254)
		if _, err := ccs.ReadFrom(f); err != nil {
			return fmt.Errorf("failed to read CCS: %w", err)
		}

		// Generate proving and verifying keys
		pk, vk, err := groth16.Setup(ccs)
		if err != nil {
			return fmt.Errorf("failed to setup Groth16: %v", err)
		}

		// Derive output file names
		base := ccsPath[:len(ccsPath)-len(filepath.Ext(ccsPath))]
		pkPath := base + ".pk"
		vkPath := base + ".vk"

		//  Write proving key
		pkFile, err := os.Create(pkPath)
		if err != nil {
			return fmt.Errorf("failed to create proving key file: %w", err)
		}
		defer pkFile.Close()

		if _, err := pk.WriteRawTo(pkFile); err != nil {
			return fmt.Errorf("failed to write proving key: %w", err)
		}
		fmt.Printf("💾 Proving key written to %s\n", pkPath)

		// Write verifying key
		vkFile, err := os.Create(vkPath)
		if err != nil {
			return fmt.Errorf("failed to create verifying key file: %w", err)
		}
		defer vkFile.Close()

		if _, err := vk.WriteRawTo(vkFile); err != nil {
			return fmt.Errorf("failed to write verifying key: %w", err)
		}
		fmt.Printf("💾 Verifying key written to %s\n", vkPath)

		fmt.Println("✅ Setup complete!")
		return nil
	},
}
