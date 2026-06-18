package shared

import (
	"fmt"
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

// ReadWitnessVec reads a fixarray of Witness values into a new slice.
func ReadWitnessVec(r *msgpackutil.Reader, out *[]Witness) error {
	n, err := r.ReadArrayLen()
	if err != nil {
		return err
	}
	*out = make([]Witness, n)
	for i := 0; i < n; i++ {
		if err := (*out)[i].UnmarshalReader(r); err != nil {
			return err
		}
	}
	return nil
}

// ReadWitnessArray reads a fixarray of Witness values into a caller-provided
// fixed-size backing array.
func ReadWitnessArray(r *msgpackutil.Reader, out []Witness) error {
	n, err := r.ReadArrayLen()
	if err != nil {
		return err
	}
	if n != len(out) {
		return fmt.Errorf("expected fixed-size Witness array of %d, got %d", len(out), n)
	}
	for i := 0; i < n; i++ {
		if err := out[i].UnmarshalReader(r); err != nil {
			return err
		}
	}
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
