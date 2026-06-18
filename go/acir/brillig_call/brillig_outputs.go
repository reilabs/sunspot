package brillig_call

import (
	"sunspot/go/acir/msgpackutil"
	shr "sunspot/go/acir/shared"
)

type BrilligOutputs struct {
	Single *shr.Witness
	Array  *[]shr.Witness
}

// BrilligOutputs: Simple(Witness) or Array(Vec<Witness>).
func (b *BrilligOutputs) UnmarshalReader(r *msgpackutil.Reader) error {
	return msgpackutil.ReadEnum(r, "BrilligOutputs", []msgpackutil.Field{
		{Name: "Simple", Decode: func(r *msgpackutil.Reader) error {
			b.Single = new(shr.Witness)
			return b.Single.UnmarshalReader(r)
		}},
		{Name: "Array", Decode: func(r *msgpackutil.Reader) error {
			b.Array = new([]shr.Witness)
			return msgpackutil.ReadVec(r, b.Array)
		}},
	})
}
