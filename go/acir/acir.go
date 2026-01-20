package acir

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"os"
	"strconv"
	expression "sunspot/go/acir/expression"
	hdr "sunspot/go/acir/header"
	shr "sunspot/go/acir/shared"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/consensys/gnark/frontend/schema"
	"github.com/google/btree"
	"github.com/rs/zerolog/log"
)

// Struct representation of an ACIR programme
type ACIR[T shr.ACIRField, E constraint.Element] struct {
	NoirVersion         string                      `json:"noir_version"`
	Hash                uint64                      `json:"hash"`
	ABI                 hdr.ACIRABI                 `json:"abi"`
	Program             Program[T, E]               `json:"program"`
	DebugSymbols        string                      `json:"debug_symbols"`
	FileMap             map[string]hdr.ACIRFileData `json:"file_map"`
	ExpressionWidth     expression.ExpressionWidth  `json:"expression_width"`
	WitnessTree         *btree.BTree                `json:"-"`
	ConstantWitnessTree *btree.BTree                `json:"-"`
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
		// Decoding bytecode from hex string
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

	// 2. Now you can treat the value as bytes
	if ewVal, ok := raw["expression_width"]; ok {
		data, err := json.Marshal(ewVal)
		if err != nil {
			return fmt.Errorf("error marshalling expression_width: %w", err)
		}

		if err := json.Unmarshal(data, &a.ExpressionWidth); err != nil {
			return fmt.Errorf("error unmarshalling ACIR ABI (expression_width): %w", err)
		}
	}

	return nil
}

func decodeProgramBytecode(bytecode string) (reader io.Reader, err error) {
	data, err := base64.StdEncoding.DecodeString(bytecode)
	if err != nil {
		return nil, fmt.Errorf("failed to decode bytecode: %w", err)
	}
	// Decompress the bytecode using gzip
	reader, err = gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to create gzip reader: %w", err)
	}
	return reader, err
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

		witnessMap := make(map[shr.Witness]frontend.Variable)
		var outerCircuitIndex uint32
		a.WitnessTree, outerCircuitIndex, err = a.Program.GetWitnesses()
		if err != nil {
			return nil, fmt.Errorf("failed to get witness tree: %w", err)
		}

		if a.WitnessTree == nil {
			return nil, fmt.Errorf("witness tree is nil, cannot compile ACIR")
		}

		// Gnark expects for the public witnesses to be added first
		// But the Noir witnesses are visibility-agnostic
		// We add the public ones first but make sure they are indexed by their
		// Noir index in the witness map
		for index, param := range a.ABI.Params() {
			if param.Visibility == hdr.ACIRParameterVisibilityPublic {
				witnessMap[shr.Witness(index+int(outerCircuitIndex))] = builder.PublicVariable(
					schema.LeafInfo{
						FullName:   func() string { return param.Name },
						Visibility: schema.Public,
					},
				)
			}
		}

		// Now we traverse the witness tree and add the private witnesses
		a.WitnessTree.Ascend(func(it btree.Item) bool {
			witness, ok := it.(shr.Witness)
			if !ok {
				log.Warn().Msgf("Item in witness tree is not of type shr.Witness: %T", it)
				return true // Continue processing other items
			}
			if _, ok := witnessMap[witness]; !ok {
				witnessMap[witness] = builder.SecretVariable(
					schema.LeafInfo{
						FullName:   func() string { return fmt.Sprintf("__witness_%d", witness) },
						Visibility: schema.Secret,
					},
				)
			}
			return true
		})

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
