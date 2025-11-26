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
	Use:   "deploy <path/to/file.vk> [rust-project-dir]",
	Short: "Builds a Solana gnark verification program with the given verification key",
	Args:  cobra.RangeArgs(1, 2),
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

		vkDir := filepath.Dir(absVkPath)                         // directory of VK file
		vkBase := filepath.Base(absVkPath)                       // e.g., example.vk
		vkName := vkBase[:len(vkBase)-len(filepath.Ext(vkBase))] // e.g., "example"

		// Paths for output files
		soPath := filepath.Join(vkDir, vkName+".so")
		keypairPath := filepath.Join(vkDir, vkName+"-keypair.json")

		// --- Step 2: Rust project directory ---
		rustDir := "."
		if len(args) == 2 {
			rustDir = args[1]
		}

		absRustDir, err := filepath.Abs(rustDir)
		if err != nil {
			log.Fatalf("Failed to resolve Rust project directory '%s': %v", rustDir, err)
		}
		fmt.Println("Using Rust project directory:", absRustDir)

		// Environment
		env := append(os.Environ(), "VK_PATH="+absVkPath)

		// --- Step 3: cargo-sbf build ---
		cargoSbf := exec.Command("cargo", "build-sbf", "--sbf-out-dir", vkDir)
		cargoSbf.Env = env
		cargoSbf.Dir = absRustDir
		cargoSbf.Stdout = os.Stdout
		cargoSbf.Stderr = os.Stderr

		fmt.Println("Running cargo-sbf build...")
		if err := cargoSbf.Run(); err != nil {
			log.Fatalf("cargo-sbf failed: %v", err)
		}

		// Rename verifier-bin.so
		originalSo := filepath.Join(vkDir, "verifier_bin.so")
		if err := os.Rename(originalSo, soPath); err != nil {
			log.Fatalf("Failed to rename .so: %v", err)
		}

		// Rename verifier-bin-keypair.json
		originalKeypair := filepath.Join(vkDir, "verifier_bin-keypair.json")
		if err := os.Rename(originalKeypair, keypairPath); err != nil {
			log.Fatalf("Failed to rename keypair.json: %v", err)
		}

		fmt.Println("Build completed successfully:")
		fmt.Println("  Program:", soPath)
		fmt.Println("  Keypair:", keypairPath)
	},
}
