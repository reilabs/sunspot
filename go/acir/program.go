package acir

import (
	"fmt"
	"sunspot/go/acir/msgpackutil"
	shr "sunspot/go/acir/shared"

	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
)

// The Circuit resolver is a function type that takes a circuit id and returns a reference to the circuit
// at that index and an error if no suh circuit exists
type CircuitResolver[T shr.ACIRField, E constraint.Element] func(id uint32) (*Circuit[T, E], error)

// Program struct represents the circuits in an ACIR programme
type Program[T shr.ACIRField, E constraint.Element] struct {
	Functions []Circuit[T, E] `json:"functions"`
}

func (p *Program[T, E]) UnmarshalReader(r *msgpackutil.Reader) error {
	return msgpackutil.ReadStruct(r, programSchema, p.decode)
}

func (p *Program[T, E]) decode(f msgpackutil.Field, r *msgpackutil.Reader) error {
	switch f.Tag {
	case 0:
		n, err := r.ReadArrayLen()
		if err != nil {
			return err
		}
		p.Functions = make([]Circuit[T, E], n)
		for i := 0; i < n; i++ {
			if err := p.Functions[i].UnmarshalReader(r); err != nil {
				return fmt.Errorf("function %d: %w", i, err)
			}
		}
		return nil
	case 1:
		// skip unconstrained functions
		return r.SkipValue()
	default:
		return fmt.Errorf("Program: unknown field: %v", f)
	}
}

// Program serde field schema (noir acvm-repo/acir/src/circuit/mod.rs).
var programSchema = msgpackutil.NewSchema(map[string]int{
	"functions":               0,
	"unconstrained_functions": 1,
})

// Define adds constraints to the ACIR programme
func (p *Program[T, E]) Define(
	api frontend.Builder[E],
	witnesses map[shr.Witness]frontend.Variable,
) error {
	// We only call define on the first (main) circuit because it will recursively define
	// any circuits that it calls
	index := uint32(0)
	if _, _, err := p.Functions[0].Define(api, witnesses, makeResolver(*p), &index); err != nil {
		return err
	}
	return nil
}

// WitnessLayout walks the call tree to determine the dense slot count and the
// global slot at which the main circuit's local witnesses begin. The witness
// space is laid out in postorder: each circuit reserves CurrentWitnessIndex+1
// contiguous slots, with all transitively-called subcircuits placed before
// the caller.
func (p *Program[T, E]) WitnessLayout() (totalSlots, mainStart uint32, err error) {
	mainStart, err = p.Functions[0].countSubcircuitSlots(makeResolver(*p), 0)
	if err != nil {
		return 0, 0, err
	}
	totalSlots = mainStart + p.Functions[0].CurrentWitnessIndex + 1
	return totalSlots, mainStart, nil
}

// Resolver takes a progamme and an index and returns the circuit
// the programme has stored at that index
func resolver[T shr.ACIRField, E constraint.Element](p Program[T, E], id uint32) (*Circuit[T, E], error) {
	if id >= uint32(len(p.Functions)) {
		return nil, fmt.Errorf("unable to get circuit, index %d out of range", id)
	}
	c := p.Functions[id]
	return &c, nil

}

// We call this inside the main programme function to get a function
// by which we can get the circuit from its index
func makeResolver[T shr.ACIRField, E constraint.Element](p Program[T, E]) func(uint32) (*Circuit[T, E], error) {
	return func(id uint32) (*Circuit[T, E], error) {
		return resolver(p, id)
	}
}
