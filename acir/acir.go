package acir

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/binary"
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
	"github.com/rs/zerolog/log"
)

type ACIR[T shr.ACIRField] struct {
	NoirVersion  string                      `json:"noir_version"`
	Hash         uint64                      `json:"hash"`
	ABI          hdr.ACIRABI                 `json:"abi"`
	Program      Program[T]                  `json:"program"`
	DebugSymbols string                      `json:"debug_symbols"`
	FileMap      map[string]hdr.ACIRFileData `json:"file_map"`
	Names        []string                    `json:"names"`
	BrilligNames []string                    `json:"brillig_names"`
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

func (a *ACIR[T]) CompileExecuted(witness WitnessStack[T]) (constraint.ConstraintSystem, error) {
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

	for index, stack := range witness.ItemStack {
		log.Trace().Msgf("Processing item stack %d with %d witnesses", index, stack.Len())
		iter := stack.Iter()
		for iter.Next() {
			witnessItem := iter.Key()
			log.Trace().Msgf("Processing witness item %d in stack %d", witnessItem, index)
			if _, ok := witnessMap[witnessItem]; !ok {
				log.Trace().Msgf("Adding item stack %d to witness map, with witness key %d", index, witnessItem)
				witnessMap[witnessItem] = builder.SecretVariable(
					schema.LeafInfo{
						FullName:   func() string { return fmt.Sprintf("item_stack_%d", index) },
						Visibility: schema.Secret,
					},
				)
			}
		}
	}

	log.Trace().Msg("Adding internal variable to builder" + fmt.Sprint(witnessMap))

	a.Program.Define(builder, witnessMap)

	return builder.Compile()
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

	for index, param := range a.ABI.Parameters {
		if param.Visibility == hdr.ACIRParameterVisibilityPrivate {
			log.Trace().Msg("Adding private parameter to witness map: " + param.Name + " at index " + fmt.Sprint(index))
			witnessMap[shr.Witness(index)] = builder.SecretVariable(
				schema.LeafInfo{
					FullName:   func() string { return param.Name },
					Visibility: schema.Secret,
				},
			)
		}
	}

	for index, param := range a.ABI.Parameters {
		if param.Visibility == hdr.ACIRParameterVisibilityDatabus {
			log.Trace().Msg("Adding databus parameter to witness map: " + param.Name + " at index " + fmt.Sprint(index))
		}
	}

	builder.InternalVariable(1)

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

func (a *ACIR[T]) GetWitnessFromFile(filePath string, field *big.Int) (witness.Witness, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open witness file: %w", err)
	}

	defer file.Close()

	var padding [12]byte
	if _, err := file.Read(padding[:]); err != nil {
		return nil, fmt.Errorf("failed to read padding: %w", err)
	}

	var witnessCount uint64
	if err := binary.Read(file, binary.LittleEndian, &witnessCount); err != nil {
		return nil, fmt.Errorf("failed to read witness count: %w", err)
	}

	witnessMap := make(map[string]*big.Int)
	for i := uint64(0); i < witnessCount; i++ {
		var witnessIndex uint32
		err := binary.Read(file, binary.LittleEndian, &witnessIndex)
		if err == io.EOF {
			break // End of file reached
		} else if err != nil {
			return nil, fmt.Errorf("failed to read witness index: %w", err)
		}
		witnessIndex += 1 // Adjusting to 1-based index

		var length uint64
		if err := binary.Read(file, binary.LittleEndian, &length); err != nil {
			return nil, fmt.Errorf("failed to read witness value length: %w", err)
		}

		valueBytes := make([]byte, length)
		if _, err := file.Read(valueBytes); err != nil {
			return nil, fmt.Errorf("failed to read value bytes: %w", err)
		}

		value, _ := new(big.Int).SetString(string(valueBytes), 10)
		if witnessIndex > uint32(len(a.ABI.Parameters)) {
			return nil, fmt.Errorf("witness index %d out of bounds for parameters length %d", witnessIndex, len(a.ABI.Parameters))
		}

		witnessMap[a.ABI.Parameters[witnessIndex-1].Name] = value
	}

	return a.GenerateWitness(witnessMap, field)
}

func (a *ACIR[T]) String() string {
	jsonData, err := json.MarshalIndent(a, "", "  ")
	if err != nil {
		return fmt.Sprintf("Error marshalling ACIR: %v", err)
	}
	return string(jsonData)
}
