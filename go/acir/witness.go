package acir

import (
	"compress/gzip"
	"encoding/binary"
	"math/big"
	"os"
	hdr "sunspot/go/acir/header"
	shr "sunspot/go/acir/shared"

	"fmt"

	"github.com/consensys/gnark/backend/witness"
	"github.com/rs/zerolog/log"
	"github.com/tidwall/btree"
)

// WitnessStacks stores witnesses from `nargo execute` in postorder based on circuit calls.
// For each execution, witnesses of called subcircuits are stored before their caller,
// so the main circuit’s witnesses appear last.
type WitnessStacks[T shr.ACIRField] map[uint64]btree.Map[shr.Witness, T]

// Loads the witness stacks from a compressed file
func LoadWitnessStacksFromFile[T shr.ACIRField](filePath string, modulus *big.Int) (WitnessStacks[T], error) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Error().Err(err).Str("file", filePath).Msg("Failed to open witness file")
		return WitnessStacks[T]{}, err
	}
	defer file.Close()

	reader, err := gzip.NewReader(file)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create gzip reader")
		return WitnessStacks[T]{}, err
	}
	defer reader.Close()

	var witnesses WitnessStacks[T]
	var stacksNum uint64
	if err := binary.Read(reader, binary.LittleEndian, &stacksNum); err != nil {
		log.Error().Err(err).Msg("Failed to read number of stacks")
		return WitnessStacks[T]{}, err
	}

	witnesses = make(WitnessStacks[T], stacksNum)
	for i := uint64(0); i < stacksNum; i++ {
		var stackIndex uint32
		if err := binary.Read(reader, binary.LittleEndian, &stackIndex); err != nil {
			return WitnessStacks[T]{}, err
		}

		var witnessMap btree.Map[shr.Witness, T]
		var mapSize uint64
		if err := binary.Read(reader, binary.LittleEndian, &mapSize); err != nil {
			return WitnessStacks[T]{}, err
		}
		for j := uint64(0); j < mapSize; j++ {
			var witness shr.Witness
			if err := binary.Read(reader, binary.LittleEndian, &witness); err != nil {
				return WitnessStacks[T]{}, err
			}

			var value T
			value = shr.MakeNonNil(value)
			if err := value.UnmarshalReader(reader); err != nil {
				return WitnessStacks[T]{}, err
			}

			witnessMap.Set(witness, value)
		}
		witnesses[i] = witnessMap
	}
	return witnesses, nil
}

// Constructs a gnark witness for the constraint system we generate when we call
// acir.compile()
// The trick here is that Gnark wants the public variables to be at the beginning of the witness vector,
// whereas noir witness stacks don't care about which variables are public
func (acir *ACIR[T, E]) GetWitness(fileName string, field *big.Int) (witness.Witness, error) {
	witnesses, err := LoadWitnessStacksFromFile[T](fileName, field)
	if err != nil {
		return nil, fmt.Errorf("failed to load witness stack from file %s: %w", fileName, err)
	}

	witness, err := witness.New(field)
	if err != nil {
		return nil, fmt.Errorf("failed to create new witness: %w", err)
	}

	params := acir.ABI.Params()
	values := make(chan any)

	// Calculate the number of private and public variables
	countPublic := 0
	countPrivate := 0
	for _, param := range params {
		if param.Visibility == hdr.ACIRParameterVisibilityPublic {
			countPublic++
		}
	}

	for _, itemStack := range witnesses {
		itemStackCount := itemStack.Len()
		countPrivate += itemStackCount
	}

	countPrivate -= countPublic

	go func() {
		// Add the public variables to the beginning of the witness vector.
		for index, param := range params {
			if param.Visibility == hdr.ACIRParameterVisibilityPublic {
				outerCircuitStack := witnesses[uint64(len(witnesses)-1)]
				if value, ok := outerCircuitStack.Get(shr.Witness(index)); ok {
					values <- value.ToFrontendVariable()
				} else {
					log.Warn().Msgf("Public parameter %s not found in outermost circuit stack", param.Name)
				}

			}
		}
		for i := 0; i < len(witnesses); i++ {
			partialWitness := witnesses[uint64(i)]
			for it := partialWitness.Iter(); it.Next(); {
				witnessKey := it.Key()
				skipKey := false
				// For the outermost circuit, we skip the witness values
				// that have already been added as part of the public variables
				if i == len(witnesses)-1 {
					for index, param := range params {
						if witnessKey == shr.Witness(index) && param.Visibility == hdr.ACIRParameterVisibilityPublic {
							skipKey = true
							break
						}
					}
					if acir.WitnessTree != nil && !acir.WitnessTree.Has(witnessKey) {
						skipKey = true
					}
				}
				if skipKey {
					continue
				}
				witnessValue := it.Value()
				values <- witnessValue.ToFrontendVariable()
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
