package blackboxfunc

import (
	"encoding/binary"
	"fmt"
	"io"
	shr "nr-groth16/acir/shared"
)

type FunctionInput[T shr.ACIRField] struct {
	FunctionInputKind FunctionInputKind
	ConstantInput     *T
	Witness           *shr.Witness
	NumberOfBits      uint32
}

func (f *FunctionInput[T]) UnmarshalReader(r io.Reader) error {
	var kind FunctionInputKind
	if err := kind.UnmarshalReader(r); err != nil {
		return err
	}
	f.FunctionInputKind = kind

	switch f.FunctionInputKind {
	case ACIRFunctionInputKindConstant:
		var constant T
		if err := constant.UnmarshalReader(r); err != nil {
			return err
		}
		f.ConstantInput = &constant
		f.Witness = nil
	case ACIRFunctionInputKindWitness:
		var witness shr.Witness
		if err := witness.UnmarshalReader(r); err != nil {
			return err
		}
		f.Witness = &witness
		f.ConstantInput = nil
	default:
		return ACIRFunctionInputKindError{
			DecodedKind: f.FunctionInputKind,
		}
	}

	if err := binary.Read(r, binary.LittleEndian, &f.NumberOfBits); err != nil {
		return err
	}
	return nil
}

func (f *FunctionInput[T]) Equals(other *FunctionInput[T]) bool {
	if f.FunctionInputKind != other.FunctionInputKind {
		return false
	}

	switch f.FunctionInputKind {
	case ACIRFunctionInputKindConstant:
		if f.ConstantInput == nil || other.ConstantInput == nil {
			return false
		}
		return (*f.ConstantInput).Equals(*other.ConstantInput)
	case ACIRFunctionInputKindWitness:
		if f.Witness == nil || other.Witness == nil {
			return false
		}
		return *f.Witness == *other.Witness
	default:
		return false
	}
}

type FunctionInputKind uint32

const (
	ACIRFunctionInputKindConstant FunctionInputKind = iota
	ACIRFunctionInputKindWitness
)

func (f *FunctionInputKind) UnmarshalReader(r io.Reader) error {
	var kind uint32
	if err := binary.Read(r, binary.LittleEndian, &kind); err != nil {
		return err
	}

	if kind > uint32(ACIRFunctionInputKindWitness) {
		return ACIRFunctionInputKindError{
			DecodedKind: FunctionInputKind(kind),
		}
	}

	*f = FunctionInputKind(kind)
	return nil
}

type ACIRFunctionInputKindError struct {
	DecodedKind FunctionInputKind
}

func (e ACIRFunctionInputKindError) Error() string {
	return fmt.Sprintf("Invalid ACIR function input kind (can be either Constant or Witness) - received %d", e.DecodedKind)
}
