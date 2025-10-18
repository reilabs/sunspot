package acir

import (
	"encoding/binary"
	"fmt"
	"io"
	"math/big"
	brl "nr-groth16/acir/brillig"
	shr "nr-groth16/acir/shared"

	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/google/btree"
	"github.com/rs/zerolog/log"
)

type CircuitResolver[T shr.ACIRField, E constraint.Element] func(id uint32) (*Circuit[T, E], error)

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

func (p *Program[T, E]) Define(
	api frontend.Builder[E],
	witnesses map[shr.Witness]frontend.Variable,
) error {

	resolver := func(id uint32) (*Circuit[T, E], error) {
		if id >= uint32(len(p.Functions)) {
			return nil, fmt.Errorf("unable to get circuit, index %d out of range", id)
		}
		c := p.Functions[id]
		return &c, nil

	}

	for _, circuit := range p.Functions {
		if err := circuit.Define(api, witnesses, resolver); err != nil {
			return err
		}
	}
	return nil
}

func (p *Program[T, E]) GetWitnessTree() (*btree.BTree, *btree.BTree) {
	witnessTree := btree.New(2)
	for _, circuit := range p.Functions {
		circuit.FillWitnessTree(witnessTree)
	}

	constantsTree := btree.New(2)
	start, ok := witnessTree.Max().(shr.Witness)
	if !ok {
		log.Error().Msg("Failed to get max witness ID from witness tree")
		return nil, nil
	}

	for _, circuit := range p.Functions {
		circuit.CollectConstantsAsWitnesses(uint32(start)+1, constantsTree)
	}

	return witnessTree, constantsTree
}

func (p *Program[T, E]) FeedConstantsAsWitnesses() []*big.Int {
	values := make([]*big.Int, 0)

	for _, circuit := range p.Functions {
		values = append(values, circuit.FeedConstantsAsWitnesses()...)
	}

	return values
}
