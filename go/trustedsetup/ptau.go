package trustedsetup

// Parser for snarkjs Powers-of-Tau (.ptau) files, producing a gnark
// mpcsetup.SrsCommons

import (
	"crypto/subtle"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	curve "github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fp"
	mpcsetup "github.com/consensys/gnark/backend/groth16/bn254/mpcsetup"
	"golang.org/x/crypto/blake2b"
)

const ptauFieldElementSize = 32

// Full ceremony PTAU files indexed in
// https://github.com/iden3/snarkjs#7-prepare-phase-2.
// Available powers are 08..28. Power N supports 2^N constraints.
const (
	ptauBaseURL  = "https://storage.googleapis.com/zkevm/ptau/"
	ptauMinPower = 8
	ptauMaxPower = 28
)

// ptauBlake2b pins the blake2b-512 digests of the Hermez Phase-1 PTAU files, as
// published in the snarkjs README. Downloads and cache hits are rejected unless
// the file hashes to the pinned value. snarkjs uses blake2b, not sha256.
var ptauBlake2b = map[uint32]string{
	8:  "d6a8fb3a04feb600096c3b791f936a578c4e664d262e4aa24beed1b7a9a96aa5eb72864d628db247e9293384b74b36ffb52ca8d148d6e1b8b51e279fdf57b583",
	9:  "94f108a80e81b5d932d8e8c9e8fd7f46cf32457e31462deeeef37af1b71c2c1b3c71fb0d9b59c654ec266b042735f50311f9fd1d4cadce47ab234ad163157cb5",
	10: "6cfeb8cda92453099d20120bdd0e8a5c4e7706c2da9a8f09ccc157ed2464d921fd0437fb70db42104769efd7d6f3c1f964bcf448c455eab6f6c7d863e88a5849",
	11: "47c282116b892e5ac92ca238578006e31a47e7c7e70f0baa8b687f0a5203e28ea07bbbec765a98dcd654bad618475d4661bfaec3bd9ad2ed12e7abc251d94d33",
	12: "ded2694169b7b08e898f736d5de95af87c3f1a64594013351b1a796dbee393bd825f88f9468c84505ddd11eb0b1465ac9b43b9064aa8ec97f2b73e04758b8a4a",
	13: "58efc8bf2834d04768a3d7ffcd8e1e23d461561729beaac4e3e7a47829a1c9066d5320241e124a1a8e8aa6c75be0ba66f65bc8239a0542ed38e11276f6fdb4d9",
	14: "eeefbcf7c3803b523c94112023c7ff89558f9b8e0cf5d6cdcba3ade60f168af4a181c9c21774b94fbae6c90411995f7d854d02ebd93fb66043dbb06f17a831c1",
	15: "982372c867d229c236091f767e703253249a9b432c1710b4f326306bfa2428a17b06240359606cfe4d580b10a5a1f63fbed499527069c18ae17060472969ae6e",
	16: "6a6277a2f74e1073601b4f9fed6e1e55226917efb0f0db8a07d98ab01df1ccf43eb0e8c3159432acd4960e2f29fe84a4198501fa54c8dad9e43297453efec125",
	17: "6247a3433948b35fbfae414fa5a9355bfb45f56efa7ab4929e669264a0258976741dfbe3288bfb49828e5df02c2e633df38d2245e30162ae7e3bcca5b8b49345",
	18: "7e6a9c2e5f05179ddfc923f38f917c9e6831d16922a902b0b4758b8e79c2ab8a81bb5f29952e16ee6c5067ed044d7857b5de120a90704c1d3b637fd94b95b13e",
	19: "bca9d8b04242f175189872c42ceaa21e2951e0f0f272a0cc54fc37193ff6648600eaf1c555c70cdedfaf9fb74927de7aa1d33dc1e2a7f1a50619484989da0887",
	20: "89a66eb5590a1c94e3f1ee0e72acf49b1669e050bb5f93c73b066b564dca4e0c7556a52b323178269d64af325d8fdddb33da3a27c34409b821de82aa2bf1a27b",
	21: "9aef0573cef4ded9c4a75f148709056bf989f80dad96876aadeb6f1c6d062391f07a394a9e756d16f7eb233198d5b69407cca44594c763ab4a5b67ae73254678",
	22: "0d64f63dba1a6f11139df765cb690da69d9b2f469a1ddd0de5e4aa628abb28f787f04c6a5fb84a235ec5ea7f41d0548746653ecab0559add658a83502d1cb21b",
	23: "3063a0bd81d68711197c8820a92466d51aeac93e915f5136d74f63c394ee6d88c5e8016231ea6580bec02e25d491f319d92e77f5c7f46a9caa8f3b53c0ea544f",
	24: "fa404d140d5819d39984833ca5ec3632cd4995f81e82db402371a4de7c2eae8687c62bc632a95b0c6aadba3fb02680a94e09174b7233ccd26d78baca2647c733",
	25: "0377d860cdb09a8a31ea1b0b8c04335614c8206357181573bf294c25d5ca7dff72387224fbd868897e6769f7805b3dab02854aec6d69d7492883b5e4e5f35eeb",
	26: "418dee4a74b9592198bd8fd02ad1aea76f9cf3085f206dfd7d594c9e264ae919611b1459a1cc920c2f143417744ba9edd7b8d51e44be9452344a225ff7eead19",
	27: "10ffd99837c512ef99752436a54b9810d1ac8878d368fb4b806267bdd664b4abf276c9cd3c4b9039a1fa4315a0c326c0e8e9e8fe0eb588ffd4f9021bf7eae1a1",
	28: "55c77ce8562366c91e7cda394cf7b7c15a06c12d8c905e8b36ba9cf5e13eb37d1a429c589e8eaba4c591bc4b88a0e2828745a53e170eac300236f5c1a326f41a",
}

// ptauFileName returns the canonical filename used in the public ceremony for
// a given power. Power 28 is published without the _28 suffix.
func ptauFileName(power uint32) string {
	if power == 28 {
		return "powersOfTau28_hez_final.ptau"
	}
	return fmt.Sprintf("powersOfTau28_hez_final_%02d.ptau", power)
}

// minPtauPowerForConstraints returns the smallest power p such that 2^p >= n,
// clamped to the publicly available PTAU range.
func minPtauPowerForConstraints(n int) (uint32, error) {
	if n > 1<<ptauMaxPower {
		return 0, fmt.Errorf("circuit has %d constraints; no PTAU available above 2^%d", n, ptauMaxPower)
	}
	p := uint32(ptauMinPower)
	for (1 << p) < n {
		p++
	}
	return p, nil
}

// ensurePtau returns a local path to the public PTAU file for the given power,
// downloading and caching it under the user cache dir on first use. Every path
// it returns has had its blake2b-512 digest verified against ptauBlake2b.
func ensurePtau(power uint32) (string, error) {
	if power < ptauMinPower || power > ptauMaxPower {
		return "", fmt.Errorf("PTAU power %d outside available range [%d,%d]", power, ptauMinPower, ptauMaxPower)
	}
	wantHex, ok := ptauBlake2b[power]
	if !ok {
		return "", fmt.Errorf("no pinned blake2b for PTAU power %d", power)
	}
	want, err := hex.DecodeString(wantHex)
	if err != nil {
		return "", fmt.Errorf("invalid pinned hash for power %d: %w", power, err)
	}

	cacheRoot, err := os.UserCacheDir()
	if err != nil {
		return "", fmt.Errorf("locate user cache dir: %w", err)
	}
	dir := filepath.Join(cacheRoot, "sunspot", "ptau")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("create cache dir: %w", err)
	}

	name := ptauFileName(power)
	path := filepath.Join(dir, name)
	if _, err := os.Stat(path); err == nil {
		fmt.Printf("📦 Using cached PTAU: %s\n", path)
		if err := verifyFileBlake2b(path, want); err != nil {
			return "", fmt.Errorf("cached PTAU failed integrity check (delete %s and retry): %w", path, err)
		}
		fmt.Printf("🔐 Pinned blake2b matches.\n")
		return path, nil
	}

	url := ptauBaseURL + name
	fmt.Printf("⬇️  Downloading PTAU file from %s\n", url)

	tmp, err := os.CreateTemp(dir, name+".tmp-*")
	if err != nil {
		return "", fmt.Errorf("create temp file: %w", err)
	}
	tmpPath := tmp.Name()
	cleanup := func() {
		tmp.Close()
		os.Remove(tmpPath)
	}

	resp, err := http.Get(url)
	if err != nil {
		cleanup()
		return "", fmt.Errorf("GET %s: %w", url, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		cleanup()
		return "", fmt.Errorf("GET %s: HTTP %d", url, resp.StatusCode)
	}

	hasher, err := blake2b.New512(nil)
	if err != nil {
		cleanup()
		return "", fmt.Errorf("init blake2b: %w", err)
	}
	if _, err := io.Copy(io.MultiWriter(tmp, hasher), resp.Body); err != nil {
		cleanup()
		return "", fmt.Errorf("write PTAU: %w", err)
	}
	if err := tmp.Close(); err != nil {
		os.Remove(tmpPath)
		return "", fmt.Errorf("close PTAU temp file: %w", err)
	}
	if subtle.ConstantTimeCompare(hasher.Sum(nil), want) != 1 {
		os.Remove(tmpPath)
		return "", fmt.Errorf("downloaded PTAU blake2b mismatch (expected %s); refusing to install", wantHex)
	}
	if err := os.Rename(tmpPath, path); err != nil {
		os.Remove(tmpPath)
		return "", fmt.Errorf("install PTAU into cache: %w", err)
	}
	fmt.Printf("🔐 Pinned blake2b matches; cached PTAU at %s\n", path)
	return path, nil
}

// verifyFileBlake2b streams the file at path through blake2b-512 and compares
// the digest to want in constant time.
func verifyFileBlake2b(path string, want []byte) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	h, err := blake2b.New512(nil)
	if err != nil {
		return fmt.Errorf("init blake2b: %w", err)
	}
	if _, err := io.Copy(h, f); err != nil {
		return fmt.Errorf("read for hashing: %w", err)
	}
	if subtle.ConstantTimeCompare(h.Sum(nil), want) != 1 {
		return fmt.Errorf("blake2b mismatch")
	}
	return nil
}

type ptauSection struct {
	pos  int64
	size uint64
}

// readPtauSRS parses path as a snarkjs .ptau file and returns a populated
// gnark SrsCommons together with the PTAU header power (log2 of domain size).
// All G1/G2 points are required to lie in the prime-order subgroup.
func readPtauSRS(path string) (mpcsetup.SrsCommons, uint32, error) {
	f, err := os.Open(path)
	if err != nil {
		return mpcsetup.SrsCommons{}, 0, err
	}
	defer f.Close()

	magic := make([]byte, 4)
	if _, err := io.ReadFull(f, magic); err != nil {
		return mpcsetup.SrsCommons{}, 0, fmt.Errorf("read magic: %w", err)
	}
	if string(magic) != "ptau" {
		return mpcsetup.SrsCommons{}, 0, fmt.Errorf("not a .ptau file: magic = %q", string(magic))
	}
	if _, err := readULE32(f); err != nil {
		return mpcsetup.SrsCommons{}, 0, fmt.Errorf("read version: %w", err)
	}
	if _, err := readULE32(f); err != nil {
		return mpcsetup.SrsCommons{}, 0, fmt.Errorf("read section count: %w", err)
	}

	// PTAU has 7 defined sections; index 0 unused so we can use 1-based ids.
	sections := make(map[uint32]ptauSection)
	for i := 0; i < 7; i++ {
		id, err := readULE32(f)
		if err != nil {
			return mpcsetup.SrsCommons{}, 0, fmt.Errorf("section %d header: %w", i, err)
		}
		size, err := readULE64(f)
		if err != nil {
			return mpcsetup.SrsCommons{}, 0, fmt.Errorf("section %d size: %w", i, err)
		}
		pos, _ := f.Seek(0, io.SeekCurrent)
		sections[id] = ptauSection{pos: pos, size: size}
		if _, err := f.Seek(int64(size), io.SeekCurrent); err != nil {
			return mpcsetup.SrsCommons{}, 0, fmt.Errorf("section %d skip: %w", i, err)
		}
	}

	seek := func(id uint32) error {
		s, ok := sections[id]
		if !ok {
			return fmt.Errorf("missing PTAU section %d", id)
		}
		_, err := f.Seek(s.pos, io.SeekStart)
		return err
	}

	if err := seek(1); err != nil {
		return mpcsetup.SrsCommons{}, 0, err
	}
	n8, err := readULE32(f)
	if err != nil {
		return mpcsetup.SrsCommons{}, 0, fmt.Errorf("header n8: %w", err)
	}
	if n8 != ptauFieldElementSize {
		return mpcsetup.SrsCommons{}, 0, fmt.Errorf("unexpected field element size %d (expected %d, only bn254 is supported)", n8, ptauFieldElementSize)
	}
	if _, err := f.Seek(int64(n8), io.SeekCurrent); err != nil {
		return mpcsetup.SrsCommons{}, 0, fmt.Errorf("skip header prime: %w", err)
	}
	power, err := readULE32(f)
	if err != nil {
		return mpcsetup.SrsCommons{}, 0, fmt.Errorf("header power: %w", err)
	}
	domain := uint32(1) << power

	var srs mpcsetup.SrsCommons
	srs.G1.Tau = make([]curve.G1Affine, 2*domain-1)
	srs.G1.AlphaTau = make([]curve.G1Affine, domain)
	srs.G1.BetaTau = make([]curve.G1Affine, domain)
	srs.G2.Tau = make([]curve.G2Affine, domain)

	if err := seek(2); err != nil {
		return mpcsetup.SrsCommons{}, 0, err
	}
	if err := readG1Slice(f, srs.G1.Tau); err != nil {
		return mpcsetup.SrsCommons{}, 0, fmt.Errorf("tauG1: %w", err)
	}

	if err := seek(3); err != nil {
		return mpcsetup.SrsCommons{}, 0, err
	}
	if err := readG2Slice(f, srs.G2.Tau); err != nil {
		return mpcsetup.SrsCommons{}, 0, fmt.Errorf("tauG2: %w", err)
	}

	if err := seek(4); err != nil {
		return mpcsetup.SrsCommons{}, 0, err
	}
	if err := readG1Slice(f, srs.G1.AlphaTau); err != nil {
		return mpcsetup.SrsCommons{}, 0, fmt.Errorf("alphaTauG1: %w", err)
	}

	if err := seek(5); err != nil {
		return mpcsetup.SrsCommons{}, 0, err
	}
	if err := readG1Slice(f, srs.G1.BetaTau); err != nil {
		return mpcsetup.SrsCommons{}, 0, fmt.Errorf("betaTauG1: %w", err)
	}

	if err := seek(6); err != nil {
		return mpcsetup.SrsCommons{}, 0, err
	}
	if err := readG2Affine(f, &srs.G2.Beta); err != nil {
		return mpcsetup.SrsCommons{}, 0, fmt.Errorf("betaG2: %w", err)
	}

	return srs, power, nil
}

func readULE32(r io.Reader) (uint32, error) {
	var b [4]byte
	if _, err := io.ReadFull(r, b[:]); err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint32(b[:]), nil
}

func readULE64(r io.Reader) (uint64, error) {
	var b [8]byte
	if _, err := io.ReadFull(r, b[:]); err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint64(b[:]), nil
}

// readFpElement reads 32 LE-encoded Montgomery-form bytes into z. PTAU stores
// each fp.Element this way, matching gnark-crypto's internal limb layout.
func readFpElement(r io.Reader, z *fp.Element) error {
	var b [ptauFieldElementSize]byte
	if _, err := io.ReadFull(r, b[:]); err != nil {
		return err
	}
	z[0] = binary.LittleEndian.Uint64(b[0:8])
	z[1] = binary.LittleEndian.Uint64(b[8:16])
	z[2] = binary.LittleEndian.Uint64(b[16:24])
	z[3] = binary.LittleEndian.Uint64(b[24:32])
	return nil
}

func readG1Affine(r io.Reader, p *curve.G1Affine) error {
	if err := readFpElement(r, &p.X); err != nil {
		return err
	}
	if err := readFpElement(r, &p.Y); err != nil {
		return err
	}
	if !p.IsOnCurve() {
		return fmt.Errorf("G1 point not on curve")
	}
	if !p.IsInSubGroup() {
		return fmt.Errorf("G1 point not in prime-order subgroup")
	}
	return nil
}

func readG2Affine(r io.Reader, p *curve.G2Affine) error {
	if err := readFpElement(r, &p.X.A0); err != nil {
		return err
	}
	if err := readFpElement(r, &p.X.A1); err != nil {
		return err
	}
	if err := readFpElement(r, &p.Y.A0); err != nil {
		return err
	}
	if err := readFpElement(r, &p.Y.A1); err != nil {
		return err
	}
	if !p.IsOnCurve() {
		return fmt.Errorf("G2 point not on curve")
	}
	if !p.IsInSubGroup() {
		return fmt.Errorf("G2 point not in prime-order subgroup")
	}
	return nil
}

func readG1Slice(r io.Reader, out []curve.G1Affine) error {
	for i := range out {
		if err := readG1Affine(r, &out[i]); err != nil {
			return fmt.Errorf("index %d: %w", i, err)
		}
	}
	return nil
}

func readG2Slice(r io.Reader, out []curve.G2Affine) error {
	for i := range out {
		if err := readG2Affine(r, &out[i]); err != nil {
			return fmt.Errorf("index %d: %w", i, err)
		}
	}
	return nil
}

