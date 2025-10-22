package acir

import (
	"compress/gzip"
	"encoding/binary"
	"math/big"
	hdr "nr-groth16/acir/header"
	shr "nr-groth16/acir/shared"
	"os"

	"fmt"

	"github.com/consensys/gnark/backend/witness"
	"github.com/rs/zerolog/log"
	"github.com/tidwall/btree"
)

type WitnessStack[T shr.ACIRField] struct {
	ItemStack ItemStack[T]
}

type ItemStack[T shr.ACIRField] map[uint64]btree.Map[shr.Witness, T]

type WitnessMap[T shr.ACIRField] btree.Map[shr.Witness, T]

func LoadWitnessStackFromFile[T shr.ACIRField](filePath string, modulus *big.Int) (WitnessStack[T], error) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Error().Err(err).Str("file", filePath).Msg("Failed to open witness file")
		return WitnessStack[T]{}, err
	}
	defer file.Close()

	reader, err := gzip.NewReader(file)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create gzip reader")
		return WitnessStack[T]{}, err
	}
	defer reader.Close()

	var witnessStack WitnessStack[T]
	var stacksNum uint64
	if err := binary.Read(reader, binary.LittleEndian, &stacksNum); err != nil {
		log.Error().Err(err).Msg("Failed to read number of stacks")
		return WitnessStack[T]{}, err
	}

	witnessStack.ItemStack = make(ItemStack[T], stacksNum)
	for i := uint64(0); i < stacksNum; i++ {
		var stackIndex uint32
		if err := binary.Read(reader, binary.LittleEndian, &stackIndex); err != nil {
			return WitnessStack[T]{}, err
		}
		var witnessMap btree.Map[shr.Witness, T]
		var mapSize uint64
		if err := binary.Read(reader, binary.LittleEndian, &mapSize); err != nil {
			return WitnessStack[T]{}, err
		}
		for j := uint64(0); j < mapSize; j++ {
			var witness shr.Witness
			if err := binary.Read(reader, binary.LittleEndian, &witness); err != nil {
				return WitnessStack[T]{}, err
			}

			var value T
			value = shr.MakeNonNil(value)
			if err := value.UnmarshalReader(reader); err != nil {
				return WitnessStack[T]{}, err
			}

			witnessMap.Set(witness, value)
		}
		witnessStack.ItemStack[i] = witnessMap
	}
	return witnessStack, nil
}

func (acir *ACIR[T, E]) GetWitness(fileName string, field *big.Int) (witness.Witness, error) {
	witnessStack, err := LoadWitnessStackFromFile[T](fileName, field)
	if err != nil {
		return nil, fmt.Errorf("failed to load witness stack from file %s: %w", fileName, err)
	}

	witness, err := witness.New(field)
	if err != nil {
		return nil, fmt.Errorf("failed to create new witness: %w", err)
	}

	values := make(chan any)
	countPublic := 0
	countPrivate := 0
	for _, param := range acir.ABI.Parameters {
		if param.Visibility == hdr.ACIRParameterVisibilityPublic {
			countPublic++
		}
	}

	for _, itemStack := range witnessStack.ItemStack {
		itemStackCount := itemStack.Len()
		for it := itemStack.Iter(); it.Next(); {
			witnessKey := it.Key()
			if acir.WitnessTree != nil && !acir.WitnessTree.Has(witnessKey) {
				itemStackCount--
				continue
			}
		}
		countPrivate += itemStackCount
	}

	countPrivate -= countPublic

	go func() {

		for index, param := range acir.ABI.Parameters {
			if param.Visibility == hdr.ACIRParameterVisibilityPublic {
				outerCircuitStack := witnessStack.ItemStack[uint64(len(witnessStack.ItemStack)-1)]
				if value, ok := outerCircuitStack.Get(shr.Witness(index)); ok {
					values <- value.ToFrontendVariable()
				} else {
					log.Warn().Msgf("Public parameter %s not found in outermost circuit stack", param.Name)
				}

			}
		}
		for i := 0; i < len(witnessStack.ItemStack); i++ {
			itemStack := witnessStack.ItemStack[uint64(i)]
			for it := itemStack.Iter(); it.Next(); {
				witnessKey := it.Key()
				skipKey := false
				// For the outermost circuit, we skip the witness values
				// that have already been added as part of the public variables
				if i == len(witnessStack.ItemStack)-1 {
					for index, param := range acir.ABI.Parameters {
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
