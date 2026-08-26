package acir

import (
	"compress/gzip"
	"math/big"
	"os"

	hdr "github.com/reilabs/sunspot/go/acir/header"
	"github.com/reilabs/sunspot/go/acir/msgpackutil"
	shr "github.com/reilabs/sunspot/go/acir/shared"

	"fmt"

	"github.com/consensys/gnark/backend/witness"
	"github.com/google/btree"
	"github.com/rs/zerolog/log"
)

// StackItem pairs a circuit's witness map with the function it corresponds to.
type StackItem[T shr.ACIRField] struct {
	CircuitIndex uint32
	WitnessMap   map[shr.Witness]T
}

// WitnessStack stores witnesses from `nargo execute` in postorder based on circuit calls.
// For each execution, witnesses of called subcircuits are stored before their caller,
// so the main circuit’s witnesses appear last.
type WitnessStack[T shr.ACIRField] []StackItem[T]

// LoadWitnessStackFromFile reads and decodes a witness stack file emitted
// by `nargo execute`. The wire envelope mirrors the bytecode side: gzip →
// format byte → msgpack-tagged payload. The payload is a `WitnessStack`
// struct containing a single tagged field 0 = Vec<StackItem>, where each
// StackItem has 0 = index (u32) and 1 = witness (WitnessMap, a fixmap of
// witness → FieldElement).
func LoadWitnessStackFromFile[T shr.ACIRField](filePath string, modulus *big.Int) (WitnessStack[T], error) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Error().Err(err).Str("file", filePath).Msg("Failed to open witness file")
		return WitnessStack[T]{}, err
	}
	defer file.Close()

	gz, err := gzip.NewReader(file)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create gzip reader")
		return WitnessStack[T]{}, err
	}
	defer gz.Close()

	if err := msgpackutil.ConsumeFormatByte(gz); err != nil {
		return WitnessStack[T]{}, fmt.Errorf("witness: %w", err)
	}
	r := msgpackutil.NewReader(gz)

	var witnesses WitnessStack[T]
	err = msgpackutil.ReadStruct(r, "WitnessStack", []msgpackutil.Field{
		{Name: "stack", Decode: func(r *msgpackutil.Reader) error {
			n, err := r.ReadArrayLen()
			if err != nil {
				return err
			}
			witnesses = make(WitnessStack[T], 0, n)
			for i := 0; i < n; i++ {
				stackItem, err := readStackItem[T](r)
				if err != nil {
					return fmt.Errorf("stack item %d: %w", i, err)
				}
				witnesses = append(witnesses, stackItem)
			}
			return nil
		}},
	})
	if err != nil {
		return WitnessStack[T]{}, err
	}
	return witnesses, nil
}

func readStackItem[T shr.ACIRField](r *msgpackutil.Reader) (StackItem[T], error) {
	var stackItem StackItem[T]
	err := msgpackutil.ReadStruct(r, "StackItem", []msgpackutil.Field{
		{Name: "index", Decode: func(r *msgpackutil.Reader) error {
			v, err := r.ReadUint()
			if err != nil {
				return err
			}
			stackItem.CircuitIndex = uint32(v)
			return nil
		}},
		{Name: "witness", Decode: func(r *msgpackutil.Reader) error { return readWitnessMap(r, &stackItem.WitnessMap) }},
	})
	return stackItem, err
}

// WitnessMap is a single-field tuple struct wrapping BTreeMap<Witness, F>,
// which serializes as a msgpack `fixmap` of int-keyed entries.
func readWitnessMap[T shr.ACIRField](r *msgpackutil.Reader, dst *map[shr.Witness]T) error {
	n, err := r.ReadMapLen()
	if err != nil {
		return err
	}
	*dst = make(map[shr.Witness]T, n)
	for i := 0; i < n; i++ {
		var w shr.Witness
		if err := w.UnmarshalReader(r); err != nil {
			return err
		}
		var value T
		value = shr.MakeNonNil(value)
		if err := value.UnmarshalReader(r); err != nil {
			return err
		}
		(*dst)[w] = value
	}
	return nil
}

// Constructs a gnark witness for the constraint system we generate when we call
// acir.compile()
// The trick here is that Gnark wants the public variables to be at the beginning of the witness vector,
// whereas the noir witness stack doesn't care about which variables are public
func (acir *ACIR[T, E]) GetWitness(fileName string, field *big.Int) (witness.Witness, error) {
	witnessStack, err := LoadWitnessStackFromFile[T](fileName, field)
	if err != nil {
		return nil, fmt.Errorf("failed to load witness stack from file %s: %w", fileName, err)
	}

	witness, err := witness.New(field)
	if err != nil {
		return nil, fmt.Errorf("failed to create new witness: %w", err)
	}

	if len(witnessStack) == 0 {
		return nil, fmt.Errorf("witness stack in %s is empty", fileName)
	}

	publicWitnesses := acir.PublicWitnesses()
	publicSlots := make(map[shr.Witness]struct{}, len(publicWitnesses))
	for _, pub := range publicWitnesses {
		publicSlots[pub.MainIndex] = struct{}{}
	}

	values := make(chan any)

	// Calculate the number of private and public variables
	countPublic := len(publicWitnesses)
	countPrivate := 0

	// Drive the count from the constraint system (one variable per slot in
	// 0..=CurrentWitnessIndex of every circuit) rather than the witness file,
	// which may omit slots that no opcode references.
	for _, stackItem := range witnessStack {
		c := &acir.Program.Functions[stackItem.CircuitIndex]
		countPrivate += int(c.CurrentWitnessIndex) + 1
	}

	countPrivate -= countPublic

	go func() {
		// Add the public variables to the beginning of the witness vector.
		outerStackItem := witnessStack[len(witnessStack)-1]
		for _, pub := range publicWitnesses {
			value, ok := outerStackItem.WitnessMap[pub.MainIndex]
			if !ok {
				log.Warn().Msgf("Public witness %s (slot %d) not found in outermost circuit witness map", pub.Name, pub.MainIndex)
				values <- 0
				continue
			}
			values <- value.ToFrontendVariable()
		}
		for i := 0; i < len(witnessStack); i++ {
			stackItem := witnessStack[i]
			c := &acir.Program.Functions[stackItem.CircuitIndex]
			for j := uint32(0); j <= c.CurrentWitnessIndex; j++ {
				witnessKey := shr.Witness(j)
				// For the outermost circuit, we skip the witness values
				// that have already been added as part of the public variables
				if i == len(witnessStack)-1 {
					if _, isPublic := publicSlots[witnessKey]; isPublic {
						continue
					}
				}
				// Slots not present in the witness file are filled with zero, matching
				// barretenberg's witness_map_to_witness_vector behavior.
				witnessValue, ok := stackItem.WitnessMap[witnessKey]
				if !ok {
					values <- 0
					continue
				}
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

// PublicWitness names a witness slot of the main circuit that must be allocated as a public gnark variable.
type PublicWitness struct {
	MainIndex shr.Witness
	Name      string
}

// PublicWitnesses returns the main circuit's public witness slots in the order gnark expects.
func (a *ACIR[T, E]) PublicWitnesses() []PublicWitness {
	if len(a.Program.Functions) == 0 {
		return nil
	}
	main := &a.Program.Functions[0]

	publicParameters := sortedWitnesses(&main.PublicParameters)

	var returnValues []shr.Witness
	if a.ABI.ReturnType != nil && a.ABI.ReturnType.Visibility == hdr.ACIRParameterVisibilityPublic {
		returnValues = sortedWitnesses(&main.ReturnValues)
	}

	params := a.ABI.Params()
	returns := a.ABI.Returns()

	publicWitnesses := make([]PublicWitness, 0, len(publicParameters)+len(returnValues))
	for _, w := range publicParameters {
		publicWitnesses = append(publicWitnesses, PublicWitness{MainIndex: w, Name: params[w].Name})
	}
	for i, w := range returnValues {
		publicWitnesses = append(publicWitnesses, PublicWitness{MainIndex: w, Name: returns[i].Name})
	}

	return publicWitnesses
}

// sortedWitnesses flattens one of the ACIR bytecode's witness sets into a sorted slice.
func sortedWitnesses(tree *btree.BTree) []shr.Witness {
	witnesses := make([]shr.Witness, 0, tree.Len())
	tree.Ascend(func(item btree.Item) bool {
		w, ok := item.(shr.Witness)
		if !ok {
			return false
		}
		witnesses = append(witnesses, w)
		return true
	})
	return witnesses
}
