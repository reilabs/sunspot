package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sunspot/acir"

	"github.com/spf13/cobra"
)

var compileCmd = &cobra.Command{
	Use:   "compile [acir_file]",
	Short: "Compile an ACIR file into a CCS file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		acirPath := args[0]

		if filepath.Ext(acirPath) != ".json" {
			return fmt.Errorf("invalid input file: %s (must end with .json)", acirPath)
		}
		fmt.Printf("Loading ACIR file: %s\n", acirPath)

		acir, err := acir.LoadACIR[T, E](acirPath)

		if err != nil {
			return fmt.Errorf("failed to load ACIR: %v", err)
		}

		ccs, err := acir.Compile()
		if err != nil {
			return fmt.Errorf("failed to compile ACIR: %v", err)
		}

		fmt.Println("Compilation successful.")

		base := strings.TrimSuffix(acirPath, ".json")
		outPath := base + ".ccs"

		// Open output file for writing
		outFile, err := os.Create(outPath)
		if err != nil {
			return fmt.Errorf("failed to create CCS file: %w", err)
		}
		defer outFile.Close()

		// Write CCS to file using its WriteTo() method
		if _, err := ccs.WriteTo(outFile); err != nil {
			return fmt.Errorf("failed to write CCS file: %w", err)
		}

		fmt.Printf("💾 CCS written to %s\n", outPath)

		return nil
	},
}
