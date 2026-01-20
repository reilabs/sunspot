package opcodes

import (
	"encoding/binary"
	"fmt"
	"io"
	bbo "sunspot/go/acir/brillig/black_box_ops"
	shr "sunspot/go/acir/shared"
)

type BrilligOpcode[T shr.ACIRField] struct {
	OpCode         BrilligOpcodeType
	BinaryFieldOp  *BinaryFieldOp
	BinaryIntOp    *BinaryIntOp
	Not            *Not
	Cast           *Cast
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

func (b *BrilligOpcode[T]) Equals(other BrilligOpcode[T]) bool {
	if b.OpCode != other.OpCode {
		return false
	}

	switch b.OpCode {
	case ACIRBrilligOpcodeBinaryFieldOp:
		return b.BinaryFieldOp.Equals(*other.BinaryFieldOp)
	case ACIRBrilligOpcodeBinaryIntOp:
		return b.BinaryIntOp.Equals(*other.BinaryIntOp)
	case ACIRBrilligOpcodeNot:
		return b.Not.Equals(*other.Not)
	case ACIRBrilligOpcodeCast:
		return b.Cast.Equals(*other.Cast)
	case ACIRBrilligOpcodeJumpIf:
		return b.JumpIf.Equals(*other.JumpIf)
	case ACIRBrilligOpcodeJump:
		return b.Jump.Equals(*other.Jump)
	case ACIRBrilligOpcodeCalldataCopy:
		return b.CalldataCopy.Equals(*other.CalldataCopy)
	case ACIRBrilligOpcodeCall:
		return b.Call.Equals(*other.Call)
	case ACIRBrilligOpcodeConst:
		return b.Const.Equals(*other.Const)
	case ACIRBrilligOpcodeIndirectConst:
		return b.IndirectConst.Equals(*other.IndirectConst)
	case ACIRBrilligOpcodeForeignCall:
		return b.ForeignCall.Equals(*other.ForeignCall)
	case ACIRBrilligOpcodeMov:
		return b.Mov.Equals(*other.Mov)
	case ACIRBrilligOpcodeConditionalMov:
		return b.ConditionalMov.Equals(*other.ConditionalMov)
	case ACIRBrilligOpcodeLoad:
		return b.Load.Equals(*other.Load)
	case ACIRBrilligOpcodeStore:
		return b.Store.Equals(*other.Store)
	case ACIRBrilligOpcodeBlackBoxOp:
		panic("BlackBoxOp equality not implemented")
		//return b.BlackBox.Equals(*other.BlackBox)
	case ACIRBrilligOpcodeTrap:
		return b.Trap.Equals(*other.Trap)
	case ACIRBrilligOpcodeStop:
		return b.Stop.Equals(*other.Stop)
	default:
		return false
	}
}
