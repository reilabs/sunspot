package msgpackutil

import "fmt"

type Unmarshaler interface {
	UnmarshalReader(r *Reader) error
}

// ReadVec decodes a fixarray of T into a freshly allocated slice stored
// through *out. T must have a pointer-receiver UnmarshalReader method.
func ReadVec[T any, PT interface {
	*T
	Unmarshaler
}](r *Reader, out *[]T) error {
	n, err := r.ReadArrayLen()
	if err != nil {
		return err
	}
	*out = make([]T, n)
	for i := range *out {
		if err := PT(&(*out)[i]).UnmarshalReader(r); err != nil {
			return err
		}
	}
	return nil
}

// ReadArrayInto decodes a fixarray of T into a caller-provided fixed-size
// slice. Errors if the wire length doesn't match len(out).
func ReadArrayInto[T any, PT interface {
	*T
	Unmarshaler
}](r *Reader, out []T) error {
	n, err := r.ReadArrayLen()
	if err != nil {
		return err
	}
	if n != len(out) {
		return fmt.Errorf("expected fixed-size array of %d, got %d", len(out), n)
	}
	for i := range out {
		if err := PT(&out[i]).UnmarshalReader(r); err != nil {
			return err
		}
	}
	return nil
}
