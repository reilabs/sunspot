package opcodes

import (
	"encoding/binary"
	"fmt"
	"io"
	bbo "nr-groth16/acir_decoder/brillig/black_box_ops"
	shr "nr-groth16/acir_decoder/shared"
)

type BrilligOpcode[T shr.ACIRField] struct {
	OpCode         BrilligOpcodeType
	BinaryFieldOp  *BinaryFieldOp
	BinaryIntOp    *BinaryIntOp
	Not            *Not
	Cast           *Cast
	JumpIfNot      *JumpIfNot
	JumpIf         *JumpIf
	Jump           *Jump
	CalldataCopy   *CallDataCopy
	Call           *Call
	Const          *Const[T]
	IndirectConst  *IndirectConst[T]
	ForeignCall    *ForeignCall
	Mov            *Mov
	ConditionalMov *ConditionalMov
	Load           *Load
	Store          *Store
	BlackBox       *bbo.BlackBoxOp
	Trap           *Trap
	Stop           *Stop
}

type BrilligOpcodeType uint32

const (
	ACIRBrilligOpcodeBinaryFieldOp BrilligOpcodeType = iota
	ACIRBrilligOpcodeBinaryIntOp
	ACIRBrilligOpcodeNot
	ACIRBrilligOpcodeCast
	ACIRBrilligOpcodeJumpIfNot
	ACIRBrilligOpcodeJumpIf
	ACIRBrilligOpcodeJump
	ACIRBrilligOpcodeCalldataCopy
	ACIRBrilligOpcodeCall
	ACIRBrilligOpcodeConst
	ACIRBrilligOpcodeIndirectConst
	ACIRBrilligOpcodeReturn
	ACIRBrilligOpcodeForeignCall
	ACIRBrilligOpcodeMov
	ACIRBrilligOpcodeConditionalMov
	ACIRBrilligOpcodeLoad
	ACIRBrilligOpcodeStore
	ACIRBrilligOpcodeBlackBoxOp
	ACIRBrilligOpcodeTrap
	ACIRBrilligOpcodeStop
)

func (b *BrilligOpcodeType) UnmarshalReader(r io.Reader) error {
	if err := binary.Read(r, binary.LittleEndian, b); err != nil {
		return err
	}

	if *b > ACIRBrilligOpcodeStop {
		return fmt.Errorf("invalid BrilligOpcodeType: %d", *b)
	}

	return nil
}

func (b *BrilligOpcode[T]) UnmarshalReader(r io.Reader) error {
	if err := b.OpCode.UnmarshalReader(r); err != nil {
		return err
	}

	switch b.OpCode {
	case ACIRBrilligOpcodeBinaryFieldOp:
		b.BinaryFieldOp = &BinaryFieldOp{}
		if err := b.BinaryFieldOp.UnmarshalReader(r); err != nil {
			return err
		}
	case ACIRBrilligOpcodeBinaryIntOp:
		b.BinaryIntOp = &BinaryIntOp{}
		if err := b.BinaryIntOp.UnmarshalReader(r); err != nil {
			return err
		}
	case ACIRBrilligOpcodeNot:
		b.Not = &Not{}
		if err := b.Not.UnmarshalReader(r); err != nil {
			return err
		}
	case ACIRBrilligOpcodeCast:
		b.Cast = &Cast{}
		if err := b.Cast.UnmarshalReader(r); err != nil {
			return err
		}
	case ACIRBrilligOpcodeJumpIfNot:
		b.JumpIfNot = &JumpIfNot{}
		if err := b.JumpIfNot.UnmarshalReader(r); err != nil {
			return err
		}
	case ACIRBrilligOpcodeJumpIf:
		b.JumpIf = &JumpIf{}
		if err := b.JumpIf.UnmarshalReader(r); err != nil {
			return err
		}
	case ACIRBrilligOpcodeJump:
		b.Jump = &Jump{}
		if err := b.Jump.UnmarshalReader(r); err != nil {
			return err
		}
	case ACIRBrilligOpcodeCalldataCopy:
		b.CalldataCopy = &CallDataCopy{}
		if err := b.CalldataCopy.UnmarshalReader(r); err != nil {
			return err
		}
	case ACIRBrilligOpcodeCall:
		b.Call = &Call{}
		if err := b.Call.UnmarshalReader(r); err != nil {
			return err
		}
	case ACIRBrilligOpcodeConst:
		b.Const = &Const[T]{}
		if err := b.Const.UnmarshalReader(r); err != nil {
			return err
		}
	case ACIRBrilligOpcodeIndirectConst:
		b.IndirectConst = &IndirectConst[T]{}
		if err := b.IndirectConst.UnmarshalReader(r); err != nil {
			return err
		}
	case ACIRBrilligOpcodeForeignCall:
		b.ForeignCall = &ForeignCall{}
		if err := b.ForeignCall.UnmarshalReader(r); err != nil {
			return err
		}
	case ACIRBrilligOpcodeMov:
		b.Mov = &Mov{}
		if err := b.Mov.UnmarshalReader(r); err != nil {
			return err
		}
	case ACIRBrilligOpcodeConditionalMov:
		b.ConditionalMov = &ConditionalMov{}
		if err := b.ConditionalMov.UnmarshalReader(r); err != nil {
			return err
		}
	case ACIRBrilligOpcodeLoad:
		b.Load = &Load{}
		if err := b.Load.UnmarshalReader(r); err != nil {
			return err
		}
	case ACIRBrilligOpcodeStore:
		b.Store = &Store{}
		if err := b.Store.UnmarshalReader(r); err != nil {
			return err
		}
	case ACIRBrilligOpcodeBlackBoxOp:
		b.BlackBox = &bbo.BlackBoxOp{}
		if err := b.BlackBox.UnmarshalReader(r); err != nil {
			return err
		}
	case ACIRBrilligOpcodeTrap:
		b.Trap = &Trap{}
		if err := b.Trap.UnmarshalReader(r); err != nil {
			return err
		}
	case ACIRBrilligOpcodeStop:
		b.Stop = &Stop{}
		if err := b.Stop.UnmarshalReader(r); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown opcode: %d", b.OpCode)
	}

	return nil
}
