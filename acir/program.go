package acir

import (
	"encoding/binary"
	"io"
	brl "nr-groth16/acir/brillig"
	shr "nr-groth16/acir/shared"
)

type Program[T shr.ACIRField] struct {
	Functions              []Circuit[T]
	UnconstrainedFunctions []brl.BrilligBytecode[T]
}

func (p *Program[T]) UnmarshalReader(r io.Reader) error {
	var funcCount uint64
	if err := binary.Read(r, binary.LittleEndian, &funcCount); err != nil {
		return err
	}
	p.Functions = make([]Circuit[T], funcCount)
	for i := uint64(0); i < funcCount; i++ {
		if err := p.Functions[i].UnmarshalReader(r); err != nil {
			return err
		}
	}

	var unconstrainedFuncCount uint64
	if err := binary.Read(r, binary.LittleEndian, &unconstrainedFuncCount); err != nil {
		if err == io.EOF {
			return nil
		}
		return err
	}

	p.UnconstrainedFunctions = make([]brl.BrilligBytecode[T], unconstrainedFuncCount)
	for i := uint64(0); i < unconstrainedFuncCount; i++ {
		if err := p.UnconstrainedFunctions[i].UnmarshalReader(r); err != nil {
			return err
		}
	}

	return nil
}
