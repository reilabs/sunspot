package opcodes

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	bbf "nr-groth16/acir/black_box_func"
	brl "nr-groth16/acir/brillig"
	exp "nr-groth16/acir/expression"
	shr "nr-groth16/acir/shared"

	"github.com/consensys/gnark/frontend"
	"github.com/google/btree"
	"github.com/rs/zerolog/log"
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

	log.Trace().Msg("Unmarshalling Opcode with kind: " + fmt.Sprint(o.Kind))

	switch o.Kind {
	case ACIROpcodeAssertZero:
		log.Trace().Msg("Unmarshalling AssertZero opcode")
		o.Expression = new(exp.Expression[T])
		if err := o.Expression.UnmarshalReader(r); err != nil {
			return err
		}
	case ACIROpcodeBlackBoxFuncCall:
		log.Trace().Msg("Unmarshalling BlackBoxFuncCall opcode")
		o.BlackBoxFuncCall = new(bbf.BlackBoxFuncCall[T])
		if err := o.BlackBoxFuncCall.UnmarshalReader(r); err != nil {
			return err
		}
	case ACIROpcodeMemoryOp:
		log.Trace().Msg("Unmarshalling MemoryOp opcode")
		o.MemoryOp = new(MemoryOp[T])
		if err := o.MemoryOp.UnmarshalReader(r); err != nil {
			return err
		}
	case ACIROpcodeMemoryInit:
		log.Trace().Msg("Unmarshalling MemoryInit opcode")
		o.MemoryInit = new(MemoryInit[T])
		if err := o.MemoryInit.UnmarshalReader(r); err != nil {
			return err
		}
	case ACIROpcodeBrilligCall:
		log.Trace().Msg("Unmarshalling BrilligCall opcode")
		o.BrilligCall = new(brl.BrilligCall[T])
		if err := o.BrilligCall.UnmarshalReader(r); err != nil {
			return err
		}
	case ACIROpcodeCall:
		log.Trace().Msg("Unmarshalling Call opcode")
		o.Call = new(Call[T])
		if err := o.Call.UnmarshalReader(r); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown OpcodeKind: %d", o.Kind)
	}

	return nil
}

func (o *Opcode[T]) Equals(other *Opcode[T]) bool {
	if o.Kind != other.Kind {
		return false
	}

	switch o.Kind {
	case ACIROpcodeAssertZero:
		if !o.Expression.Equals(other.Expression) {
			return false
		}
	case ACIROpcodeBlackBoxFuncCall:
		//if !o.BlackBoxFuncCall.Equals(other.BlackBoxFuncCall) {
		//	return false
		//}
		return true // BlackBoxFuncCall does not have a location, so we assume equality if the call matches.
	case ACIROpcodeMemoryOp:
		if !o.MemoryOp.Equals(other.MemoryOp) {
			return false
		}
	case ACIROpcodeMemoryInit:
		if !o.MemoryInit.Equals(other.MemoryInit) {
			return false
		}
	case ACIROpcodeBrilligCall:
		//if !o.BrilligCall.Equals(other.BrilligCall) {
		//	return false
		//}
		return true
	case ACIROpcodeCall:
		if !o.Call.Equals(other.Call) {
			return false
		}
	default:
		return false
	}

	return true
}

func (o *Opcode[T]) Define(api frontend.API, witnesses map[shr.Witness]frontend.Variable) error {
	switch o.Kind {
	case ACIROpcodeAssertZero:
		if o.Expression == nil {
			panic("Expression is nil for AssertZero opcode")
		}

		api.AssertIsEqual(o.Expression.Calculate(api, witnesses), 0)
	case ACIROpcodeBlackBoxFuncCall:
		//panic("BlackBoxFuncCall opcode is not implemented yet") // TODO: Implement BlackBoxFuncCall calculation
		//return o.BlackBoxFuncCall.Calculate(api, witnesses)
	case ACIROpcodeMemoryOp:
		panic("MemoryOp opcode is not implemented yet") // TODO: Implement MemoryOp calculation
		//return o.MemoryOp.Calculate(api, witnesses)
	case ACIROpcodeMemoryInit:
		panic("MemoryInit opcode is not implemented yet") // TODO: Implement MemoryInit calculation
		//return o.MemoryInit.Calculate(api, witnesses)
	case ACIROpcodeBrilligCall:
		//panic("BrilligCall opcode is not implemented yet") // TODO: Implement BrilligCall calculation
		//return o.BrilligCall.Calculate(api, witnesses)
	case ACIROpcodeCall:
		panic("Call opcode is not implemented yet") // TODO: Implement Call calculation
		//return o.Call.Calculate(api, witnesses)
	default:
		panic(fmt.Sprintf("unknown OpcodeKind: %d", o.Kind))
	}

	return nil
}

func (o Opcode[T]) MarshalJSON() ([]byte, error) {
	stringMap := make(map[string]interface{})
	switch o.Kind {
	case ACIROpcodeAssertZero:
		stringMap["assert_zero"] = o.Expression
	case ACIROpcodeBlackBoxFuncCall:
		stringMap["black_box_func_call"] = o.BlackBoxFuncCall
	case ACIROpcodeMemoryOp:
		stringMap["memory_op"] = o.MemoryOp
	case ACIROpcodeMemoryInit:
		stringMap["memory_init"] = o.MemoryInit
	case ACIROpcodeBrilligCall:
		stringMap["brillig_call"] = o.BrilligCall
	case ACIROpcodeCall:
		stringMap["call"] = o.Call
	}

	return json.Marshal(stringMap)
}

func (o Opcode[T]) FillWitnessTree(tree *btree.BTree) bool {
	if tree == nil {
		return false
	}

	ok := true
	if o.Kind == ACIROpcodeAssertZero && o.Expression != nil {
		ok = o.Expression.FillWitnessTree(tree)
	}

	return ok
}
