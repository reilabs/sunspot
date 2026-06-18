package blackboxfunc

import (
	"fmt"
	"sunspot/go/acir/msgpackutil"
	shr "sunspot/go/acir/shared"

	"github.com/consensys/gnark/frontend"
)

// Function input represents a type that can be either a constant or a witness index
// An internal representation for the following docs
// https://noir-lang.github.io/noir/docs/acir/circuit/opcodes/struct.FunctionInput.html
type FunctionInput[T shr.ACIRField] struct {
	ConstantInput *T
	Witness       *shr.Witness
}

func (f *FunctionInput[T]) UnmarshalReader(r *msgpackutil.Reader) error {
	return msgpackutil.ReadEnum(r, functionInputSchema, f.decode)
}

func (f *FunctionInput[T]) decode(field msgpackutil.Field, r *msgpackutil.Reader) error {
	switch field.Tag {
	case 0:
		var constant T
		constant = shr.MakeNonNil(constant)
		if err := constant.UnmarshalReader(r); err != nil {
			return err
		}
		f.ConstantInput = &constant
		f.Witness = nil
		return nil
	case 1:
		var witness shr.Witness
		if err := witness.UnmarshalReader(r); err != nil {
			return err
		}
		f.Witness = &witness
		f.ConstantInput = nil
		return nil
	default:
		return fmt.Errorf("invalid ACIR function input kind (can be either Constant or Witness) - received %v", field)
	}
}

var functionInputSchema = msgpackutil.NewSchema(map[string]int{
	"Constant": 0, "Witness": 1,
})

func (f *FunctionInput[T]) Equals(other *FunctionInput[T]) bool {
	if f.IsWitness() != other.IsWitness() {
		return false
	}
	if f.IsWitness() {
		return *f.Witness == *other.Witness
	}
	if f.ConstantInput == nil || other.ConstantInput == nil {
		return f.ConstantInput == other.ConstantInput
	}
	return (*f.ConstantInput).Equals(*other.ConstantInput)
}

func (f *FunctionInput[T]) ToVariable(witnesses map[shr.Witness]frontend.Variable) (frontend.Variable, error) {
	if f.IsWitness() {
		if _, ok := witnesses[*f.Witness]; !ok {
			return nil, fmt.Errorf("witness %d not found in witnesses map", *f.Witness)
		}
		return witnesses[*f.Witness], nil
	}
	if f.ConstantInput != nil {
		return (*f.ConstantInput).ToFrontendVariable(), nil
	}
	return nil, fmt.Errorf("function input has neither constant nor witness set")
}

func (f *FunctionInput[T]) IsWitness() bool {
	return f.Witness != nil
}
