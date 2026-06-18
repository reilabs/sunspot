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
	return msgpackutil.ReadEnum(r, brilligInputsSchema, inputs.decode)
}

func (inputs *BrilligInputs[T, E]) decode(f msgpackutil.Field, r *msgpackutil.Reader) error {
	switch f.Tag {
	case 0:
		inputs.Single = new(exp.Expression[T, E])
		return inputs.Single.UnmarshalReader(r)
	case 1:
		inputs.Array = new([]exp.Expression[T, E])
		return msgpackutil.ReadVec(r, inputs.Array)
	case 2:
		v, err := r.ReadUint()
		if err != nil {
			return err
		}
		inputs.BlockID = new(uint32)
		*inputs.BlockID = uint32(v)
		return nil
	default:
		return fmt.Errorf("unknown BrilligInputsKind: %v", f)
	}
}

var brilligInputsSchema = msgpackutil.NewSchema(map[string]int{
	"Single": 0, "Array": 1, "MemoryArray": 2,
})
