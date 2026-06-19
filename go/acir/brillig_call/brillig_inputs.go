package brillig_call

import (
	exp "github.com/reilabs/sunspot/go/acir/expression"
	"github.com/reilabs/sunspot/go/acir/msgpackutil"
	shr "github.com/reilabs/sunspot/go/acir/shared"

	"github.com/consensys/gnark/constraint"
)

type BrilligInputs[T shr.ACIRField, E constraint.Element] struct {
	Single  *exp.Expression[T, E]
	Array   *[]exp.Expression[T, E]
	BlockID *uint32
}

func (inputs *BrilligInputs[T, E]) UnmarshalReader(r *msgpackutil.Reader) error {
	return msgpackutil.ReadEnum(r, "BrilligInputs", []msgpackutil.Field{
		{Name: "Single", Decode: func(r *msgpackutil.Reader) error {
			inputs.Single = new(exp.Expression[T, E])
			return inputs.Single.UnmarshalReader(r)
		}},
		{Name: "Array", Decode: func(r *msgpackutil.Reader) error {
			inputs.Array = new([]exp.Expression[T, E])
			return msgpackutil.ReadVec(r, inputs.Array)
		}},
		{Name: "MemoryArray", Decode: func(r *msgpackutil.Reader) error {
			v, err := r.ReadUint()
			if err != nil {
				return err
			}
			inputs.BlockID = new(uint32)
			*inputs.BlockID = uint32(v)
			return nil
		}},
	})
}
