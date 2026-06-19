package trustedsetup

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// drand "default" chain on api.drand.sh: chained mainnet, 30s pulses,
// pedersen-bls-chained scheme on BLS12-381. Chain hash:
//   8990e7a9aaed2ffed73dbd7092123d6f289930540d7651336225dc172e51b2ce
// Group public key (G1 compressed, hex):
//   868f005eb8e6e4ca0a47c8a77ceaa5309a47978a7c71bc5cce96366b5d7a569937c529eeda66c7293784a9402801af31
// The JSON response does not include a randomness field; by convention,
// randomness = sha256(signature).
const drandAPIBase = "https://api.drand.sh/v2/beacons/default/rounds"

type drandPulse struct {
	Round             uint64 `json:"round"`
	Signature         string `json:"signature"`
	PreviousSignature string `json:"previous_signature"`
}

// fetchDrandBeacon retrieves the published drand pulse for `round` from the
// public default chain and returns its derived randomness (sha256 of the BLS
// signature) for use as a trusted-setup Seal beacon. The returned pulse is
// intended to be logged so auditors can re-fetch the same round and re-derive
// the same Seal output.
//
// NOTE: this does not verify the BLS12-381 signature locally — trust in the
// fetched bytes rests on TLS to api.drand.sh. Anyone auditing the seal should
// independently fetch the same round and verify its signature against the
// drand group public key before accepting the resulting pk/vk.
func fetchDrandBeacon(round uint64) ([]byte, *drandPulse, error) {
	if round == 0 {
		return nil, nil, fmt.Errorf("drand round must be > 0")
	}
	url := fmt.Sprintf("%s/%d", drandAPIBase, round)
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, nil, fmt.Errorf("drand fetch %s: %w", url, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, nil, fmt.Errorf("drand fetch %s: HTTP %d: %s", url, resp.StatusCode, string(body))
	}
	var pulse drandPulse
	if err := json.NewDecoder(resp.Body).Decode(&pulse); err != nil {
		return nil, nil, fmt.Errorf("drand decode: %w", err)
	}
	if pulse.Round != round {
		return nil, nil, fmt.Errorf("drand returned round %d, asked for %d", pulse.Round, round)
	}
	sigBytes, err := hex.DecodeString(pulse.Signature)
	if err != nil {
		return nil, nil, fmt.Errorf("drand signature hex decode: %w", err)
	}
	if len(sigBytes) == 0 {
		return nil, nil, fmt.Errorf("drand returned empty signature for round %d", round)
	}
	randomness := sha256.Sum256(sigBytes)
	return randomness[:], &pulse, nil
}
