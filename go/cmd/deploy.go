package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

var deployCmd = &cobra.Command{
	Use:   "deploy <path/to/file.vk>",
	Short: "Builds a Solana gnark verification program with the given verification key",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		// --- Step 1: VK file ---
		vkPath := args[0]

		absVkPath, err := filepath.Abs(vkPath)
		if err != nil {
			log.Fatalf("Failed to resolve VK file path '%s': %v", vkPath, err)
		}
		if _, err := os.Stat(absVkPath); err != nil {
			log.Fatalf("VK file does not exist: %s", absVkPath)
		}
		fmt.Println("Using VK file:", absVkPath)

		vkDir := filepath.Dir(absVkPath)
		vkBase := filepath.Base(absVkPath)
		vkName := vkBase[:len(vkBase)-len(filepath.Ext(vkBase))]

		// Output paths
		soPath := filepath.Join(vkDir, vkName+".so")
		keypairPath := filepath.Join(vkDir, vkName+"-keypair.json")

		// --- Step 2: GNARK_VERIFIER_BIN ---
		verifierDir := os.Getenv("GNARK_VERIFIER_BIN")
		if verifierDir == "" {
			log.Fatalf("Environment variable GNARK_VERIFIER_BIN is not set.\n" +
				"Please set it to the Rust verifier-bin crate directory.")
		}

		absVerifierDir, err := filepath.Abs(verifierDir)
		if err != nil {
			log.Fatalf("Failed to resolve GNARK_VERIFIER_BIN '%s': %v", verifierDir, err)
		}
		if _, err := os.Stat(absVerifierDir); err != nil {
			log.Fatalf("GNARK_VERIFIER_BIN directory does not exist: %s", absVerifierDir)
		}

		fmt.Println("Using verifier-bin crate directory:", absVerifierDir)

		// --- Step 3: cargo build-sbf ---
		env := append(os.Environ(), "VK_PATH="+absVkPath)

		cargoSbf := exec.Command("cargo", "build-sbf", "--sbf-out-dir", vkDir)
		cargoSbf.Env = env
		cargoSbf.Dir = absVerifierDir
		cargoSbf.Stdout = os.Stdout
		cargoSbf.Stderr = os.Stderr

		fmt.Println("Running cargo build-sbf...")
		if err := cargoSbf.Run(); err != nil {
			log.Fatalf("cargo build-sbf failed: %v", err)
		}

		// --- Step 4: Rename outputs ---
		originalSo := filepath.Join(vkDir, "verifier_bin.so")
		if err := os.Rename(originalSo, soPath); err != nil {
			log.Fatalf("Failed to rename .so: %v", err)
		}

		originalKeypair := filepath.Join(vkDir, "verifier_bin-keypair.json")
		if err := os.Rename(originalKeypair, keypairPath); err != nil {
			log.Fatalf("Failed to rename keypair.json: %v", err)
		}

		fmt.Println("Build completed successfully:")
		fmt.Println("  Program:", soPath)
		fmt.Println("  Keypair:", keypairPath)
	},
}
