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

type ItemStack[T shr.ACIRField] map[uint32]btree.Map[shr.Witness, T]

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

	log.Trace().Msg("Number of item stacks in witness: " + fmt.Sprint(stacksNum))

	witnessStack.ItemStack = make(ItemStack[T], stacksNum)
	for i := uint64(0); i < stacksNum; i++ {
		var stackIndex uint32
		if err := binary.Read(reader, binary.LittleEndian, &stackIndex); err != nil {
			return WitnessStack[T]{}, err
		}
		log.Trace().Msgf("Reading stack %d", stackIndex)

		var witnessMap btree.Map[shr.Witness, T]
		var mapSize uint64
		if err := binary.Read(reader, binary.LittleEndian, &mapSize); err != nil {
			return WitnessStack[T]{}, err
		}
		log.Trace().Msgf("Reading %d witnesses for stack %d", mapSize, stackIndex)
		for j := uint64(0); j < mapSize; j++ {
			var witness shr.Witness
			if err := binary.Read(reader, binary.LittleEndian, &witness); err != nil {
				return WitnessStack[T]{}, err
			}
			log.Trace().Msgf("Reading witness %d for stack %d", witness, stackIndex)

			var value T
			value = shr.MakeNonNil(value)
			if err := value.UnmarshalReader(reader); err != nil {
				return WitnessStack[T]{}, err
			}
			log.Trace().Msgf("Reading value %s for witness %d in stack %d", value.String(), witness, stackIndex)

			witnessMap.Set(witness, value)
		}
		witnessStack.ItemStack[stackIndex] = witnessMap
	}
	return witnessStack, nil
}

func (acir *ACIR[T]) GetWitness(fileName string, field *big.Int) (witness.Witness, error) {
	witnessStack, err := LoadWitnessStackFromFile[T](fileName, field)
	if err != nil {
		return nil, fmt.Errorf("failed to load witness stack from file %s: %w", fileName, err)
	}

	witness, err := witness.New(field)
	if err != nil {
		return nil, fmt.Errorf("failed to create new witness: %w", err)
	}
	log.Trace().Msg("Starting to fill witness with public and private parameters")

	values := make(chan any)
	countPublic := 0
	countPrivate := 0
	for _, param := range acir.ABI.Parameters {
		if param.Visibility == hdr.ACIRParameterVisibilityPublic {
			countPublic++
		}
	}
	log.Trace().Msgf("Number of public parameters: %d", countPublic)

	for stackIndex, itemStack := range witnessStack.ItemStack {
		log.Trace().Msgf("Processing stack %d with %d items", stackIndex, itemStack.Len())
		itemStackCount := itemStack.Len()
		for it := itemStack.Iter(); it.Next(); {
			witnessKey := it.Key()
			if acir.WitnessTree != nil && !acir.WitnessTree.Has(witnessKey) {
				log.Warn().Msgf("Witness key %d not found in witness tree or is zero, skipping", witnessKey)
				itemStackCount--
				continue
			}
		}
		countPrivate += itemStackCount
	}

	countPrivate += acir.ConstantWitnessTree.Len()

	countPrivate -= countPublic
	log.Trace().Msgf("Number of private parameters: %d", countPrivate)

	go func() {
		for index, param := range acir.ABI.Parameters {
			if param.Visibility == hdr.ACIRParameterVisibilityPublic {
				log.Trace().Msgf("Sending public parameter %s", param.Name)
				for stackIndex, itemStack := range witnessStack.ItemStack {
					log.Trace().Msgf("Processing stack %d for public parameter %s", stackIndex, param.Name)
					if value, ok := itemStack.Get(shr.Witness(index)); ok {
						log.Trace().Msgf("Sending value %s for public parameter %s in stack %d", value.String(), param.Name, stackIndex)
						values <- value.ToFrontendVariable()
						break // Only send the first occurrence of the public parameter
					} else {
						log.Warn().Msgf("Public parameter %s not found in stack %d", param.Name, stackIndex)
					}
				}
			}
		}
		log.Trace().Msg("Finished sending public parameters")

		for _, itemStack := range witnessStack.ItemStack {
			log.Trace().Msgf("Processing private items in stack %d", itemStack.Len())
			for it := itemStack.Iter(); it.Next(); {
				witnessKey := it.Key()
				skipKey := false
				for index, param := range acir.ABI.Parameters {
					if witnessKey == shr.Witness(index) && param.Visibility == hdr.ACIRParameterVisibilityPublic {
						log.Trace().Msgf("Skipping public parameter %s in private witness processing", param.Name)
						skipKey = true
						break
					}
				}
				if acir.WitnessTree != nil && !acir.WitnessTree.Has(witnessKey) {
					log.Warn().Msgf("Witness key %d not found in witness tree, skipping", witnessKey)
					skipKey = true
				}
				if skipKey {
					continue
				}
				witnessValue := it.Value()
				log.Trace().Msgf("Sending private witness %d with value %s", witnessKey, witnessValue.String())
				values <- witnessValue.ToFrontendVariable()
			}
		}
		log.Trace().Msg("Finished sending private parameters")

		data := acir.Program.FeedConstantsAsWitnesses()
		for _, value := range data {
			log.Trace().Msgf("Sending constant value %s", value.String())
			values <- value
		}

		log.Trace().Msg("Finished sending constant values")

		close(values)
	}()

	err = witness.Fill(countPublic, countPrivate, values)
	if err != nil {
		return nil, fmt.Errorf("failed to fill witness: %w", err)
	}

	return witness, nil
}
