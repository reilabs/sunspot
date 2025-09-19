package acir

import (
	"encoding/binary"
	"io"
	"math/big"
	brl "nr-groth16/acir/brillig"
	shr "nr-groth16/acir/shared"

	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/google/btree"
	"github.com/rs/zerolog/log"
)

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
	api frontend.API,
	witnesses map[shr.Witness]frontend.Variable,
) error {
	for _, circuit := range p.Functions {
		if err := circuit.Define(api, witnesses); err != nil {
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
