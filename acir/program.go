package acir

import (
	"encoding/binary"
	"fmt"
	"io"
	brl "nr-groth16/acir/brillig"
	shr "nr-groth16/acir/shared"

	"github.com/consensys/gnark/frontend"
	"github.com/google/btree"
	"github.com/rs/zerolog/log"
)

type Program[T shr.ACIRField] struct {
	Functions              []Circuit[T]             `json:"functions"`
	UnconstrainedFunctions []brl.BrilligBytecode[T] `json:"unconstrained_functions"`
}

func (p *Program[T]) UnmarshalReader(r io.Reader) error {
	var funcCount uint64
	if err := binary.Read(r, binary.LittleEndian, &funcCount); err != nil {
		return err
	}
	log.Trace().Msg("Unmarshalling program with " + fmt.Sprint(funcCount) + " circuits")
	p.Functions = make([]Circuit[T], funcCount)
	for i := uint64(0); i < funcCount; i++ {
		log.Trace().Msg("Unmarshalling circuit " + fmt.Sprint(i))
		if err := p.Functions[i].UnmarshalReader(r); err != nil {
			return err
		}
	}

	var unconstrainedFuncCount uint64
	if err := binary.Read(r, binary.BigEndian, &unconstrainedFuncCount); err != nil {
		if err == io.EOF {
			return nil
		}
		return err
	}

	log.Trace().Msg("Unmarshalling program with " + fmt.Sprintf("%x", unconstrainedFuncCount) + " unconstrained brillig bytecode functions")

	/*p.UnconstrainedFunctions = make([]brl.BrilligBytecode[T], unconstrainedFuncCount)
	for i := uint64(0); i < unconstrainedFuncCount; i++ {
		log.Trace().Msg("Unmarshalling unconstrained brillig bytecode function " + fmt.Sprint(i))
		if err := p.UnconstrainedFunctions[i].UnmarshalReader(r); err != nil {
			return err
		}
	}*/

	return nil
}

func (p *Program[T]) Define(api frontend.API, witnesses map[shr.Witness]frontend.Variable) error {
	for _, circuit := range p.Functions {
		if err := circuit.Define(api, witnesses); err != nil {
			return err
		}
	}
	return nil
}

func (p *Program[T]) GetWitnessTree() *btree.BTree {
	witnessTree := btree.New(2)
	for _, circuit := range p.Functions {
		circuit.FillWitnessTree(witnessTree)
	}
	return witnessTree
}
