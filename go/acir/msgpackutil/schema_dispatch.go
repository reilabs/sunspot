package msgpackutil

import (
	"fmt"
	"reflect"
)

// FieldDecoder reads one struct field's payload (or one enum variant's
// payload).
type FieldDecoder = func(*Reader) error

// Field pairs a wire name with the decoder for its payload. The order of
// Field entries in a call to ReadStruct/ReadEnum mirrors the noir struct or
// enum definition; the position is used to dispatch positional
// (MsgpackCompact) and tag-keyed (MsgpackTagged) wire shapes.
type Field struct {
	Name   string
	Decode FieldDecoder
}

// SkipField is the decoder to register for wire fields that should be
// consumed and discarded — e.g. MemoryInit's block_type,
// Program's unconstrained_functions.
func SkipField(r *Reader) error { return r.SkipValue() }

// Unmarshaler is anything that decodes itself from a Reader.
type Unmarshaler interface {
	UnmarshalReader(r *Reader) error
}

// EnumVariant is the constraint for ReadDispatchedEnum exemplars: each
// variant identifies its own wire name and decodes its own payload.
type EnumVariant interface {
	Unmarshaler
	SerdeName() string
}

// ReadStruct decodes a struct that may be on the wire in any of three shapes:
// string-keyed map (Msgpack), int-keyed map (MsgpackTagged), or positional
// fixarray (MsgpackCompact / Tagged-as-Array). name is used only for error
// messages.
func ReadStruct(r *Reader, name string, fields []Field) error {
	b, err := r.Peek()
	if err != nil {
		return err
	}
	switch {
	case isMapMarker(b):
		n, err := r.ReadMapLen()
		if err != nil {
			return err
		}
		for i := 0; i < n; i++ {
			if err := dispatchKeyed(r, name, "field", fields); err != nil {
				return fmt.Errorf("%s: entry %d: %w", name, i, err)
			}
		}
		return nil
	case isArrayMarker(b):
		n, err := r.ReadArrayLen()
		if err != nil {
			return err
		}
		for i := 0; i < n; i++ {
			if i >= len(fields) {
				return fmt.Errorf("%s: positional index %d exceeds schema length %d", name, i, len(fields))
			}
			if err := invoke(r, name, fields[i]); err != nil {
				return err
			}
		}
		return nil
	default:
		return fmt.Errorf("%s: expected map or array, got 0x%02x", name, b)
	}
}

// ReadEnum decodes a sum type (single-entry map of {variant: payload}). Use
// for enums whose branches are inlined in the surrounding type
// (FunctionInput, BrilligInputs, BrilligOutputs). For enums where each variant
// is a distinct Go type, ReadDispatchedEnum is shorter.
func ReadEnum(r *Reader, name string, variants []Field) error {
	n, err := r.ReadMapLen()
	if err != nil {
		return fmt.Errorf("%s: %w", name, err)
	}
	if n != 1 {
		return fmt.Errorf("%s: expected single-entry map, got %d entries", name, n)
	}
	return dispatchKeyed(r, name, "variant", variants)
}

// ReadDispatchedEnum decodes an enum where each variant is its own Go type
// implementing EnumVariant. The framework allocates a fresh instance of the
// matching exemplar's type, lets it decode its own payload, and stores it
// via setter. exemplars provide both the wire names (SerdeName) and the Go
// types (via reflection).
func ReadDispatchedEnum[V EnumVariant](r *Reader, name string, exemplars []V, setter func(V)) error {
	fields := make([]Field, len(exemplars))
	for i, ex := range exemplars {
		t := reflect.TypeOf(ex).Elem()
		fields[i] = Field{
			Name: ex.SerdeName(),
			Decode: func(r *Reader) error {
				v := reflect.New(t).Interface().(V)
				if err := v.UnmarshalReader(r); err != nil {
					return err
				}
				setter(v)
				return nil
			},
		}
	}
	return ReadEnum(r, name, fields)
}

// dispatchKeyed consumes a map key (string or int) and invokes the matching
// field's decoder.
func dispatchKeyed(r *Reader, parent, kind string, fields []Field) error {
	b, err := r.Peek()
	if err != nil {
		return err
	}
	if isStringMarker(b) {
		s, err := r.ReadString()
		if err != nil {
			return err
		}
		for _, f := range fields {
			if f.Name == s {
				return invoke(r, parent, f)
			}
		}
		return fmt.Errorf("%s: unknown %s %q", parent, kind, s)
	}
	tag, err := r.ReadUint()
	if err != nil {
		return err
	}
	if tag >= uint64(len(fields)) {
		return fmt.Errorf("%s: unknown %s tag %d", parent, kind, tag)
	}
	return invoke(r, parent, fields[tag])
}

func invoke(r *Reader, parent string, f Field) error {
	if err := f.Decode(r); err != nil {
		return fmt.Errorf("%s.%s: %w", parent, f.Name, err)
	}
	return nil
}

func isStringMarker(b byte) bool {
	return (b >= markerFixstrLow && b <= markerFixstrHigh) || b == markerStr8 || b == markerStr16 || b == markerStr32
}

func isMapMarker(b byte) bool {
	return (b >= markerFixmapLow && b <= markerFixmapHigh) || b == markerMap16 || b == markerMap32
}

func isArrayMarker(b byte) bool {
	return (b >= markerFixarrayLow && b <= markerFixarrayHigh) || b == markerArray16 || b == markerArray32
}
