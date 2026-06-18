package msgpackutil

import "fmt"

// ReadStruct decodes a struct that may be on the wire in any of three shapes:
//   - int-keyed fixmap (Tagged): callback gets f.Tag = the wire key
//   - string-keyed fixmap (Msgpack): callback gets f.Name = the wire key, and
//     f.Tag = schema.TagFor(name) if a schema entry exists (else -1)
//   - positional fixarray (MsgpackCompact / Tagged-as-Array): callback gets
//     f.Tag = the index
//
// Use this for every struct-shaped decode.
func ReadStruct(r *Reader, schema Schema, decode func(f Field, r *Reader) error) error {
	b, err := r.Peek()
	if err != nil {
		return err
	}
	switch {
	case b >= markerFixmapLow && b <= markerFixmapHigh, b == markerMap16, b == markerMap32:
		n, err := r.ReadMapLen()
		if err != nil {
			return err
		}
		for i := 0; i < n; i++ {
			f, err := readFieldKey(r, schema)
			if err != nil {
				return fmt.Errorf("product: read field key at entry %d: %w", i, err)
			}
			if err := decode(f, r); err != nil {
				return fmt.Errorf("product: field %v: %w", f, err)
			}
		}
		return nil
	case b >= markerFixarrayLow && b <= markerFixarrayHigh, b == markerArray16, b == markerArray32:
		n, err := r.ReadArrayLen()
		if err != nil {
			return err
		}
		for i := 0; i < n; i++ {
			if err := decode(Field{Tag: i}, r); err != nil {
				return fmt.Errorf("product: index %d: %w", i, err)
			}
		}
		return nil
	default:
		return fmt.Errorf("product: expected map or array, got 0x%02x", b)
	}
}

// ReadEnum decodes a enum-typed value. All three formats serialize
// enum-with-payload as a single-entry map `{variant: payload}`. The key is a
// string variant name in modes Msgpack/MsgpackCompact and an int tag in mode
// MsgpackTagged; this function peeks the key to determine which.
//
// For unit variants in mode MsgpackTagged the payload slot is nil — callbacks
// that don't read anything must still consume that nil via r.ReadNil().
func ReadEnum(r *Reader, schema Schema, decode func(f Field, r *Reader) error) error {
	n, err := r.ReadMapLen()
	if err != nil {
		return fmt.Errorf("enum: %w", err)
	}
	if n != 1 {
		return fmt.Errorf("enum: expected single-entry map, got %d entries", n)
	}
	f, err := readFieldKey(r, schema)
	if err != nil {
		return fmt.Errorf("enum: read variant key: %w", err)
	}
	if err := decode(f, r); err != nil {
		return fmt.Errorf("enum %v: %w", f, err)
	}
	return nil
}

// readFieldKey consumes either an int (any uint encoding) or a string (any
// width) and packages it as a Field. String keys are resolved against schema
// when possible so per-type decoders can switch on Tag uniformly.
func readFieldKey(r *Reader, schema Schema) (Field, error) {
	b, err := r.Peek()
	if err != nil {
		return Field{}, err
	}
	if isStringMarker(b) {
		name, err := r.ReadString()
		if err != nil {
			return Field{}, err
		}
		f := Field{Tag: -1, Name: name}
		if t, ok := schema.TagFor(name); ok {
			f.Tag = t
		}
		return f, nil
	}
	tag, err := r.ReadUint()
	if err != nil {
		return Field{}, err
	}
	return Field{Tag: int(tag)}, nil
}

func isStringMarker(b byte) bool {
	return (b >= markerFixstrLow && b <= markerFixstrHigh) || b == markerStr8 || b == markerStr16 || b == markerStr32
}
