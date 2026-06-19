package trustedsetup

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	mpcsetup "github.com/consensys/gnark/backend/groth16/bn254/mpcsetup"
	cs_bn254 "github.com/consensys/gnark/constraint/bn254"
	"github.com/spf13/cobra"
)

var (
	trustedSetupBeacon     string
	trustedSetupDrandRound uint64
	trustedSetupPhase2Out  string
	trustedSetupPhase2In   []string
)

// Cmd is the `trusted-setup` subcommand. cmd/root.go wires this into the root
// cobra command; everything else in this package is implementation detail.
var Cmd = &cobra.Command{
	Use:   "trusted-setup [ccs_file]",
	Short: "Participate in a multi-party Groth16 Phase 2 trusted setup",
	Long: "Multi-party Phase 2 trusted setup against a CCS file. The Hermez Phase-1 PTAU file " +
		"matching the circuit size is downloaded automatically and integrity-checked against a " +
		"pinned blake2b digest. Use --phase2-out alone to start a ceremony; --phase2-in " +
		"(repeatable, in order) together with --phase2-out to verify the prior chain and add " +
		"your contribution; or --phase2-in alone (full ordered chain) to verify and seal the " +
		"final pk/vk pair.",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ccsPath := args[0]

		if filepath.Ext(ccsPath) != ".ccs" {
			return fmt.Errorf("invalid input file: %s (must end with .ccs)", ccsPath)
		}
		if len(trustedSetupPhase2In) == 0 && trustedSetupPhase2Out == "" {
			return fmt.Errorf("must specify --phase2-out (to start or extend a ceremony) or --phase2-in (to seal); single-machine setup is not supported, use `setup` for testing")
		}
		sealing := trustedSetupPhase2Out == ""
		if sealing {
			if (trustedSetupBeacon == "") == (trustedSetupDrandRound == 0) {
				return fmt.Errorf("for the final seal, exactly one of --beacon or --beacon-drand-round must be set")
			}
			if len(trustedSetupPhase2In) == 0 {
				return fmt.Errorf("sealing requires the full ordered chain of prior contributions via --phase2-in")
			}
		} else {
			if trustedSetupBeacon != "" || trustedSetupDrandRound != 0 {
				return fmt.Errorf("--beacon / --beacon-drand-round only apply to the final seal step (omit --phase2-out)")
			}
		}

		fmt.Printf("🔧 Loading CCS file: %s\n", ccsPath)
		ccsFile, err := os.Open(ccsPath)
		if err != nil {
			return fmt.Errorf("failed to load CCS: %v", err)
		}
		defer ccsFile.Close()

		ccs := groth16.NewCS(ecc.BN254)
		if _, err := ccs.ReadFrom(ccsFile); err != nil {
			return fmt.Errorf("failed to read CCS: %w", err)
		}
		r1cs, ok := ccs.(*cs_bn254.R1CS)
		if !ok {
			return fmt.Errorf("expected bn254 R1CS in CCS file, got %T", ccs)
		}

		power, err := minPtauPowerForConstraints(r1cs.GetNbConstraints())
		if err != nil {
			return err
		}
		fmt.Printf("🔎 Circuit has %d constraints; selecting PTAU power %d\n", r1cs.GetNbConstraints(), power)
		ptauPath, err := ensurePtau(power)
		if err != nil {
			return fmt.Errorf("failed to obtain PTAU: %w", err)
		}

		fmt.Printf("🔧 Loading PTAU: %s\n", ptauPath)
		srs, power, err := readPtauSRS(ptauPath)
		if err != nil {
			return fmt.Errorf("failed to read PTAU: %w", err)
		}
		fmt.Printf("    PTAU power=%d (supports up to %d constraints)\n", power, 1<<power)
		if 1<<power < r1cs.GetNbConstraints() {
			return fmt.Errorf("PTAU power %d supports %d constraints, but CCS has %d", power, 1<<power, r1cs.GetNbConstraints())
		}

		chain, err := loadPhase2Chain(trustedSetupPhase2In)
		if err != nil {
			return err
		}

		if sealing {
			beaconBytes, err := resolveBeacon()
			if err != nil {
				return err
			}
			fmt.Printf("🔍 Verifying chain of %d contribution(s) and sealing...\n", len(chain))
			pk, vk, err := mpcsetup.VerifyPhase2(r1cs, &srs, beaconBytes, chain...)
			if err != nil {
				return fmt.Errorf("chain verification failed: %w", err)
			}
			return writePkVk(ccsPath, pk, vk)
		}

		fmt.Println("⚙️  Running Phase 2 Initialize...")
		var p2 mpcsetup.Phase2
		p2.Initialize(r1cs, &srs)

		last := &p2
		if len(chain) > 0 {
			fmt.Printf("🔍 Verifying chain of %d prior contribution(s)...\n", len(chain))
			prev := &p2
			for i, c := range chain {
				if err := prev.Verify(c); err != nil {
					return fmt.Errorf("contribution %d (%s) failed verification: %w", i, trustedSetupPhase2In[i], err)
				}
				prev = c
			}
			last = prev
		}

		fmt.Println("⚙️  Contributing local entropy to Phase 2...")
		last.Contribute()

		fmt.Printf("💾 Writing Phase 2 contribution to %s\n", trustedSetupPhase2Out)
		f, err := os.Create(trustedSetupPhase2Out)
		if err != nil {
			return fmt.Errorf("failed to create Phase 2 file: %w", err)
		}
		defer f.Close()
		if _, err := last.WriteTo(f); err != nil {
			return fmt.Errorf("failed to write Phase 2 file: %w", err)
		}
		fmt.Println("✅ Phase 2 contribution written. Pass this file (alongside the prior chain) to the next contributor or the sealer.")
		return nil
	},
}

func loadPhase2Chain(paths []string) ([]*mpcsetup.Phase2, error) {
	chain := make([]*mpcsetup.Phase2, 0, len(paths))
	for _, path := range paths {
		fmt.Printf("🔧 Loading Phase 2 contribution: %s\n", path)
		f, err := os.Open(path)
		if err != nil {
			return nil, fmt.Errorf("failed to open Phase 2 file %s: %w", path, err)
		}
		c := new(mpcsetup.Phase2)
		_, readErr := c.ReadFrom(f)
		f.Close()
		if readErr != nil {
			return nil, fmt.Errorf("failed to read Phase 2 file %s: %w", path, readErr)
		}
		chain = append(chain, c)
	}
	return chain, nil
}

func resolveBeacon() ([]byte, error) {
	if trustedSetupDrandRound != 0 {
		fmt.Printf("🌐 Fetching drand round %d (api.drand.sh default chain)...\n", trustedSetupDrandRound)
		bz, pulse, err := fetchDrandBeacon(trustedSetupDrandRound)
		if err != nil {
			return nil, fmt.Errorf("fetch drand beacon: %w", err)
		}
		fmt.Printf("    drand round=%d\n", pulse.Round)
		fmt.Printf("    signature=%s\n", pulse.Signature)
		fmt.Printf("    randomness=%x  (sha256 of signature)\n", bz)
		return bz, nil
	}
	fmt.Printf("🔒 Sealing with raw beacon %q...\n", trustedSetupBeacon)
	return []byte(trustedSetupBeacon), nil
}

func writePkVk(ccsPath string, pk groth16.ProvingKey, vk groth16.VerifyingKey) error {
	base := ccsPath[:len(ccsPath)-len(filepath.Ext(ccsPath))]
	pkPath := base + ".pk"
	vkPath := base + ".vk"

	pkFile, err := os.Create(pkPath)
	if err != nil {
		return fmt.Errorf("failed to create proving key file: %w", err)
	}
	defer pkFile.Close()
	if _, err := pk.WriteRawTo(pkFile); err != nil {
		return fmt.Errorf("failed to write proving key: %w", err)
	}
	fmt.Printf("💾 Proving key written to %s\n", pkPath)

	vkFile, err := os.Create(vkPath)
	if err != nil {
		return fmt.Errorf("failed to create verifying key file: %w", err)
	}
	defer vkFile.Close()
	if _, err := vk.WriteRawTo(vkFile); err != nil {
		return fmt.Errorf("failed to write verifying key: %w", err)
	}
	fmt.Printf("💾 Verifying key written to %s\n", vkPath)
	fmt.Println("✅ Trusted setup complete!")
	return nil
}

func init() {
	Cmd.Flags().StringVar(&trustedSetupBeacon, "beacon", "",
		"raw beacon bytes for the final Seal step (escape hatch for offline ceremonies; prefer --beacon-drand-round)")
	Cmd.Flags().Uint64Var(&trustedSetupDrandRound, "beacon-drand-round", 0,
		"drand round (api.drand.sh default chain) to fetch and use as the Seal beacon; commit to a future round publicly before contributions begin")
	Cmd.Flags().StringVar(&trustedSetupPhase2Out, "phase2-out", "",
		"write the Phase 2 contribution to this path and exit; skips Seal so the file can be passed to further contributors")
	Cmd.Flags().StringSliceVar(&trustedSetupPhase2In, "phase2-in", nil,
		"ordered list of prior Phase 2 contribution files (comma-separated, or repeat the flag). When extending, the chain is verified before your contribution is added; when sealing, the chain is verified end-to-end before pk/vk are emitted")
}
