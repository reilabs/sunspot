package shared

import (
	"sunspot/go/acir/msgpackutil"

	"github.com/google/btree"
)

type Witness uint32

// On the wire a Witness is a single-field tuple struct that serde
// flattens to its inner u32. Witnesses encode positionally as plain
// MessagePack uint values. We also notify the Reader's witness tracker
// so Circuit can size its witness vector without an encoded count.
func (w *Witness) UnmarshalReader(r *msgpackutil.Reader) error {
	v, err := r.ReadU32()
	if err != nil {
		return err
	}
	*w = Witness(v)
	r.ObserveWitness(v)
	return nil
}

func (w Witness) Less(other btree.Item) bool {
	otherWitness, ok := other.(Witness)
	if !ok {
		return false
	}
	return w < otherWitness
}

func (w *Witness) Equals(other *Witness) bool {
	if other == nil {
		return false
	}
	return *w == *other
}
