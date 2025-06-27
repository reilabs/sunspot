package acir

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	hdr "nr-groth16/acir/header"
	shr "nr-groth16/acir/shared"
	"os"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/witness"
	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/consensys/gnark/frontend/schema"
	"github.com/google/btree"
	"github.com/rs/zerolog/log"
)

type ACIR[T shr.ACIRField] struct {
	NoirVersion  string                      `json:"noir_version"`
	Hash         uint64                      `json:"hash"`
	ABI          hdr.ACIRABI                 `json:"abi"`
	Program      Program[T]                  `json:"program"`
	DebugSymbols string                      `json:"debug_symbols"`
	FileMap      map[string]hdr.ACIRFileData `json:"file_map"`
	WitnessTree  *btree.BTree                `json:"-"` // Optional, can be nil
	Names        []string                    `json:"names"`
	BrilligNames []string                    `json:"brillig_names"`
}

func LoadACIR[T shr.ACIRField](fileName string) (ACIR[T], error) {
	file, err := os.Open(fileName)
	if err != nil {
		return ACIR[T]{}, fmt.Errorf("failed to open ACIR file: %w", err)
	}
	defer file.Close()

	var acir ACIR[T]
	if err := json.NewDecoder(file).Decode(&acir); err != nil {
		return ACIR[T]{}, fmt.Errorf("failed to decode ACIR JSON: %w", err)
	}

	return acir, nil
}

func (a *ACIR[T]) UnmarshalJSON(data []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	if version, ok := raw["noir_version"].(string); ok {
		a.NoirVersion = version
	} else {
		return fmt.Errorf("missing or invalid noir_version field in ACIR")
	}

	if hash, ok := raw["hash"].(float64); ok {
		a.Hash = uint64(hash)
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

	if names, ok := raw["names"].([]interface{}); ok {
		for _, name := range names {
			if str, ok := name.(string); ok {
				a.Names = append(a.Names, str)
			} else {
				return fmt.Errorf("invalid name in names array: %v", name)
			}
		}
	} else {
		return fmt.Errorf("missing or invalid names field in ACIR")
	}

	if brilligNames, ok := raw["brillig_names"].([]interface{}); ok {
		for _, name := range brilligNames {
			if str, ok := name.(string); ok {
				a.BrilligNames = append(a.BrilligNames, str)
			} else {
				return fmt.Errorf("invalid name in brillig_names array: %v", name)
			}
		}
	} else {
		return fmt.Errorf("missing or invalid brillig_names field in ACIR")
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

func (a *ACIR[T]) Compile() (constraint.ConstraintSystem, error) {
	builder, err := r1cs.NewBuilder(ecc.BN254.ScalarField(), frontend.CompileConfig{
		CompressThreshold: 300,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create R1CS builder: %w", err)
	}

	witnessMap := make(map[shr.Witness]frontend.Variable)
	for index, param := range a.ABI.Parameters {
		if param.Visibility == hdr.ACIRParameterVisibilityPublic {
			log.Trace().Msg("Adding public parameter to witness map: " + param.Name + " at index " + fmt.Sprint(index))
			witnessMap[shr.Witness(index)] = builder.PublicVariable(
				schema.LeafInfo{
					FullName:   func() string { return param.Name },
					Visibility: schema.Public,
				},
			)
		}
	}

	a.WitnessTree = a.Program.GetWitnessTree()
	if a.WitnessTree == nil {
		return nil, fmt.Errorf("witness tree is nil, cannot compile ACIR")
	}
	log.Trace().Msg("Processing witness tree with " + fmt.Sprint(a.WitnessTree.Len()))

	a.WitnessTree.Ascend(func(it btree.Item) bool {
		witness, ok := it.(shr.Witness)
		if !ok {
			log.Warn().Msgf("Item in witness tree is not of type shr.Witness: %T", it)
			return true // Continue processing other items
		}
		log.Trace().Msgf("Processing witness item %d", it)
		if _, ok := witnessMap[witness]; !ok {
			log.Trace().Msgf("Adding witness to witness map, with key %d", witness)
			witnessMap[witness] = builder.SecretVariable(
				schema.LeafInfo{
					FullName:   func() string { return fmt.Sprintf("__witness_%d", witness) },
					Visibility: schema.Secret,
				},
			)
		}
		return true
	})

	log.Trace().Msg("Adding internal variable to builder" + fmt.Sprint(witnessMap))

	a.Program.Define(builder, witnessMap)

	return builder.Compile()
}

func (a *ACIR[T]) GenerateWitness(inputs map[string]*big.Int, field *big.Int) (witness.Witness, error) {
	witness, err := witness.New(field)
	if err != nil {
		return nil, err
	}

	values := make(chan any)
	countPublic := 0
	countPrivate := 0

	for _, param := range a.ABI.Parameters {
		if param.Visibility == hdr.ACIRParameterVisibilityPublic {
			countPublic++
		} else if param.Visibility == hdr.ACIRParameterVisibilityPrivate {
			countPrivate++
		}
	}

	go func() {
		for _, param := range a.ABI.Parameters {
			if param.Visibility == hdr.ACIRParameterVisibilityPublic {
				values <- inputs[param.Name]
			}
		}

		for _, param := range a.ABI.Parameters {
			if param.Visibility == hdr.ACIRParameterVisibilityPrivate {
				values <- inputs[param.Name]
			}
		}
		close(values)
	}()

	err = witness.Fill(countPublic, countPrivate, values)
	if err != nil {
		return nil, fmt.Errorf("failed to fill witness: %w", err)
	}

	return witness, nil
}

func (a *ACIR[T]) String() string {
	jsonData, err := json.MarshalIndent(a, "", "  ")
	if err != nil {
		return fmt.Sprintf("Error marshalling ACIR: %v", err)
	}
	return string(jsonData)
}
