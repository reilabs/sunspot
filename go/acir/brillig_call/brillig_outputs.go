package brillig_call

import (
	"fmt"
	"sunspot/go/acir/msgpackutil"
	shr "sunspot/go/acir/shared"
)

type BrilligOutputs struct {
	Single *shr.Witness
	Array  *[]shr.Witness
}

// BrilligOutputs: 0 = Simple(Witness), 1 = Array(Vec<Witness>).
func (b *BrilligOutputs) UnmarshalReader(r *msgpackutil.Reader) error {
	return msgpackutil.ReadEnum(r, brilligOutputsSchema, b.decode)
}

func (b *BrilligOutputs) decode(f msgpackutil.Field, r *msgpackutil.Reader) error {
	switch f.Tag {
	case 0:
		b.Single = new(shr.Witness)
		return b.Single.UnmarshalReader(r)
	case 1:
		b.Array = new([]shr.Witness)
		return msgpackutil.ReadVec(r, b.Array)
	default:
		return fmt.Errorf("unknown BrilligOutputsKind: %v", f)
	}
}

var brilligOutputsSchema = msgpackutil.NewSchema(map[string]int{
	"Simple": 0, "Array": 1,
})
