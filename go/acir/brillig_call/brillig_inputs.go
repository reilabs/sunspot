package brillig_call

import (
	"fmt"
	exp "sunspot/go/acir/expression"
	"sunspot/go/acir/msgpackutil"
	shr "sunspot/go/acir/shared"

	"github.com/consensys/gnark/constraint"
)

type BrilligInputs[T shr.ACIRField, E constraint.Element] struct {
	Single  *exp.Expression[T, E]
	Array   *[]exp.Expression[T, E]
	BlockID *uint32
}

func (inputs *BrilligInputs[T, E]) UnmarshalReader(r *msgpackutil.Reader) error {
	return msgpackutil.ReadEnum(r, inputs.decode)
}

func (inputs *BrilligInputs[T, E]) decode(tag int, r *msgpackutil.Reader) error {
	switch tag {
	case 0:
		inputs.Single = new(exp.Expression[T, E])
		return inputs.Single.UnmarshalReader(r)
	case 1:
		n, err := r.ReadArrayLen()
		if err != nil {
			return err
		}
		inputs.Array = new([]exp.Expression[T, E])
		*inputs.Array = make([]exp.Expression[T, E], n)
		for i := 0; i < n; i++ {
			if err := (*inputs.Array)[i].UnmarshalReader(r); err != nil {
				return err
			}
		}
		return nil
	case 2:
		v, err := r.ReadUint()
		if err != nil {
			return err
		}
		inputs.BlockID = new(uint32)
		*inputs.BlockID = uint32(v)
		return nil
	default:
		return fmt.Errorf("unknown BrilligInputsKind: %d", tag)
	}
}
