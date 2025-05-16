package opcodes

import (
	"encoding/binary"
	"fmt"
	"io"
	bbf "nr-groth16/acir_decoder/black_box_func"
	brl "nr-groth16/acir_decoder/brillig"
	exp "nr-groth16/acir_decoder/expression"
	shr "nr-groth16/acir_decoder/shared"
)

type Opcode[T shr.ACIRField] struct {
	Kind             OpcodeKind
	Expression       *exp.Expression[T]
	BlackBoxFuncCall *bbf.BlackBoxFuncCall[T]
	MemoryOp         *MemoryOp[T]
	MemoryInit       *MemoryInit[T]
	BrilligCall      *brl.BrilligCall[T]
	Call             *Call[T]
}

type OpcodeKind uint32

const (
	ACIROpcodeAssertZero OpcodeKind = iota
	ACIROpcodeBlackBoxFuncCall
	ACIROpcodeMemoryOp
	ACIROpcodeMemoryInit
	ACIROpcodeBrilligCall
	ACIROpcodeCall
)

func (o *Opcode[T]) UnmarshalReader(r io.Reader) error {
	if err := binary.Read(r, binary.LittleEndian, &o.Kind); err != nil {
		return err
	}

	switch o.Kind {
	case ACIROpcodeAssertZero:
		o.Expression = new(exp.Expression[T])
		if err := o.Expression.UnmarshalReader(r); err != nil {
			return err
		}
	case ACIROpcodeBlackBoxFuncCall:
		o.BlackBoxFuncCall = new(bbf.BlackBoxFuncCall[T])
		if err := o.BlackBoxFuncCall.UnmarshalReader(r); err != nil {
			return err
		}
	case ACIROpcodeMemoryOp:
		o.MemoryOp = new(MemoryOp[T])
		if err := o.MemoryOp.UnmarshalReader(r); err != nil {
			return err
		}
	case ACIROpcodeMemoryInit:
		o.MemoryInit = new(MemoryInit[T])
		if err := o.MemoryInit.UnmarshalReader(r); err != nil {
			return err
		}
	case ACIROpcodeBrilligCall:
		o.BrilligCall = new(brl.BrilligCall[T])
		if err := o.BrilligCall.UnmarshalReader(r); err != nil {
			return err
		}
	case ACIROpcodeCall:
		o.Call = new(Call[T])
		if err := o.Call.UnmarshalReader(r); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown OpcodeKind: %d", o.Kind)
	}

	return nil
}
