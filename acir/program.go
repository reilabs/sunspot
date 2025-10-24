package acir

import (
	"encoding/binary"
	"fmt"
	"io"
	brl "sunpot/acir/brillig"
	shr "sunpot/acir/shared"

	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/google/btree"
)

// The Circuit resolver is a function type that takes a circuit id and returns a reference to the circuit
// at that index and an error if no suh circuit exists
type CircuitResolver[T shr.ACIRField, E constraint.Element] func(id uint32) (*Circuit[T, E], error)

// Program struct represents the circuits in an ACIR programme
type Program[T shr.ACIRField, E constraint.Element] struct {
	Functions              []Circuit[T, E]          `json:"functions"`
	UnconstrainedFunctions []brl.BrilligBytecode[T] `json:"unconstrained_functions"`
}

func (p *Program[T, E]) UnmarshalReader(r io.Reader) error {
	var funcCount uint64
	if err := binary.Read(r, binary.LittleEndian, &funcCount); err != nil {
		return err
	}
	p.Functions = make([]Circuit[T, E], funcCount)
	for i := uint64(0); i < funcCount; i++ {
		if err := p.Functions[i].UnmarshalReader(r); err != nil {
			return err
		}
	}

	var unconstrainedFuncCount uint64
	if err := binary.Read(r, binary.BigEndian, &unconstrainedFuncCount); err != nil {
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			return nil
		}
		return err
	}

	return nil
}

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

// GetWitnesses returns all the used witness indices in the programme
// and the starting index for the witness referenced by the outermost (main) circuit
func (p *Program[T, E]) GetWitnesses() (*btree.BTree, uint32, error) {
	witnessTree := btree.New(2)
	outerCircuitWitnessIndex, err := p.Functions[0].FillWitnessTree(witnessTree, makeResolver(*p), uint32(0))
	if err != nil {
		return nil, outerCircuitWitnessIndex, err
	}
	return witnessTree, outerCircuitWitnessIndex, nil
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
