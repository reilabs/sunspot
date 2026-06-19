package acir

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"strconv"
	hdr "github.com/reilabs/sunspot/go/acir/header"
	"github.com/reilabs/sunspot/go/acir/msgpackutil"
	shr "github.com/reilabs/sunspot/go/acir/shared"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/consensys/gnark/frontend/schema"
)

// Struct representation of an ACIR programme
type ACIR[T shr.ACIRField, E constraint.Element] struct {
	NoirVersion  string                      `json:"noir_version"`
	Hash         uint64                      `json:"hash"`
	ABI          hdr.ACIRABI                 `json:"abi"`
	Program      Program[T, E]               `json:"program"`
	DebugSymbols string                      `json:"debug_symbols"`
	FileMap      map[string]hdr.ACIRFileData `json:"file_map"`
}

// Loads ACIR from disk and creates representation in memory
func LoadACIR[T shr.ACIRField, E constraint.Element](fileName string) (ACIR[T, E], error) {
	file, err := os.Open(fileName)
	if err != nil {
		return ACIR[T, E]{}, fmt.Errorf("failed to open ACIR file: %w", err)
	}
	defer file.Close()

	var acir ACIR[T, E]
	if err := json.NewDecoder(file).Decode(&acir); err != nil {
		return ACIR[T, E]{}, fmt.Errorf("failed to decode ACIR JSON: %w", err)
	}

	return acir, nil
}

// Construct an ACIR instance from json data
func (a *ACIR[T, E]) UnmarshalJSON(data []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	if version, ok := raw["noir_version"].(string); ok {
		a.NoirVersion = version
	} else {
		return fmt.Errorf("missing or invalid noir_version field in ACIR")
	}

	if hashStr, ok := raw["hash"].(string); ok {
		hash, err := strconv.ParseUint(hashStr, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid hash value in ACIR: %v", err)
		}
		a.Hash = hash
	} else {
		return fmt.Errorf("missing or invalid hash field in ACIR")
	}

	if abiData, ok := raw["abi"].(map[string]interface{}); ok {
		var abi hdr.ACIRABI
		abiBytes, err := json.Marshal(abiData)
		if err != nil {
			return fmt.Errorf("error marshalling ACIR ABI: %v", err)
		}

		if err := json.Unmarshal(abiBytes, &abi); err != nil {
			return fmt.Errorf("error unmarshalling ACIR ABI: %v", err)
		}
		a.ABI = abi
	} else {
		return fmt.Errorf("missing or invalid abi field in ACIR")
	}

	if bytecode, ok := raw["bytecode"].(string); ok {
		reader, err := decodeProgramBytecode(bytecode)
		if err != nil {
			return fmt.Errorf("error decoding bytecode: %v", err)
		}
		if err := a.Program.UnmarshalReader(reader); err != nil {
			return fmt.Errorf("error unmarshalling program bytecode: %v", err)
		}
	} else {
		return fmt.Errorf("missing or invalid bytecode field in ACIR")
	}

	if debugSymbols, ok := raw["debug_symbols"].(string); ok {
		a.DebugSymbols = debugSymbols
	} else {
		return fmt.Errorf("missing or invalid debug_symbols field in ACIR")
	}

	if fileMap, ok := raw["file_map"].(map[string]interface{}); ok {
		a.FileMap = make(map[string]hdr.ACIRFileData)
		for fileName, fileData := range fileMap {
			var file hdr.ACIRFileData
			fileBytes, err := json.Marshal(fileData)
			if err != nil {
				return fmt.Errorf("error marshalling file data for %s: %v", fileName, err)
			}
			if err := json.Unmarshal(fileBytes, &file); err != nil {
				return fmt.Errorf("error unmarshalling ACIR file data for %s: %v", fileName, err)
			}
			a.FileMap[fileName] = file
		}
	} else {
		return fmt.Errorf("missing or invalid file_map field in ACIR")
	}

	return nil
}

func decodeProgramBytecode(bytecode string) (*msgpackutil.Reader, error) {
	data, err := base64.StdEncoding.DecodeString(bytecode)
	if err != nil {
		return nil, fmt.Errorf("failed to decode bytecode: %w", err)
	}
	gz, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gz.Close()
	if err := msgpackutil.ConsumeFormatByte(gz); err != nil {
		return nil, err
	}
	return msgpackutil.NewReader(gz), nil
}

func (a *ACIR[T, E]) Compile() (constraint.ConstraintSystemGeneric[E], error) {
	// Implement the NewBuilder[E] function from gnark
	// This allows us to feed the builder into a circuit and call Compile
	// on the builder
	builder_generator := func(*big.Int, frontend.CompileConfig) (frontend.Builder[E], error) {
		builder, err := r1cs.NewBuilder[E](ecc.BN254.ScalarField(), frontend.CompileConfig{
			CompressThreshold: 300,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create R1CS builder: %w", err)
		}

		totalSlots, mainStart, err := a.Program.WitnessLayout()
		if err != nil {
			return nil, fmt.Errorf("failed to compute witness layout: %w", err)
		}

		witnessMap := make(map[shr.Witness]frontend.Variable, totalSlots)

		// Gnark expects public witnesses to be added before private ones. Noir's
		// public params occupy the first slots of the main circuit (i.e. starting
		// at mainStart, which sits after every transitively-called subcircuit).
		for index, param := range a.ABI.Params() {
			if param.Visibility == hdr.ACIRParameterVisibilityPublic {
				witnessMap[shr.Witness(uint32(index)+mainStart)] = builder.PublicVariable(
					schema.LeafInfo{
						FullName:   func() string { return param.Name },
						Visibility: schema.Public,
					},
				)
			}
		}

		// Allocate a private variable for every remaining slot. We don't inspect
		// opcodes to determine which witnesses are referenced — every index from
		// 0 to totalSlots-1 gets a variable.
		for i := uint32(0); i < totalSlots; i++ {
			if _, ok := witnessMap[shr.Witness(i)]; !ok {
				witnessMap[shr.Witness(i)] = builder.SecretVariable(
					schema.LeafInfo{
						FullName:   func() string { return fmt.Sprintf("__witness_%d", i) },
						Visibility: schema.Secret,
					},
				)
			}
		}

		err = a.Program.Define(builder, witnessMap)
		if err != nil {
			return nil, err
		}
		return builder, nil
	}

	return frontend.CompileGeneric(ecc.BN254.ScalarField(), builder_generator, &DummyCircuit{})

}

func (a *ACIR[T, E]) String() string {
	jsonData, err := json.MarshalIndent(a, "", "  ")
	if err != nil {
		return fmt.Sprintf("Error marshalling ACIR: %v", err)
	}
	return string(jsonData)
}

// We need the dummy circuit to feed in our custom builder
// This makes sure that `callDeferred` is actually called on our custom builder
// See desired behaviour [here](https://github.com/Consensys/gnark/blob/master/frontend/compile.go#L159)
// and notice how it is not called by custom constraint system builders [here](https://github.com/Consensys/gnark/blob/55b0e54d2ae15e886ad37300a8d2b00ad00a8023/frontend/cs/r1cs/builder.go#L278)
type DummyCircuit struct{}

func (a *DummyCircuit) Define(frontend.API) error {
	return nil
}
