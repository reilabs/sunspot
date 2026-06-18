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
	return msgpackutil.ReadEnum(r, b.decode)
}

func (b *BrilligOutputs) decode(tag int, r *msgpackutil.Reader) error {
	switch tag {
	case 0:
		b.Single = new(shr.Witness)
		return b.Single.UnmarshalReader(r)
	case 1:

		n, err := r.ReadArrayLen()
		if err != nil {
			return err
		}
		b.Array = new([]shr.Witness)
		*b.Array = make([]shr.Witness, n)
		for i := 0; i < n; i++ {
			if err := (*b.Array)[i].UnmarshalReader(r); err != nil {
				return err
			}
		}
		return nil
	default:
		return fmt.Errorf("unknown BrilligOutputsKind: %d", tag)
	}
}
