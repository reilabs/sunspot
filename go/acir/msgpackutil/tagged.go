package msgpackutil

import "fmt"

// ReadStruct decodes a struct that may be on the wire as either a fixmap
// (int-keyed, `EncodingStrategy::Tagged`) or a fixarray (positional,
// `EncodingStrategy::Array`). It peeks the type marker, dispatches to the
// matching shape, and calls `decode` once per entry with the field tag
// (key in tagged form, index in array form).
//
// Use this for every struct-shaped decode, regardless of which strategy
// the encoder chose — the decoder probes per-type at runtime, so a single
// buffer can mix shapes across nested types freely.
func ReadStruct(r *Reader, decode func(tag int, r *Reader) error) error {
	b, err := r.Peek()
	if err != nil {
		return err
	}
	switch {
	case b >= 0x80 && b <= 0x8f, b == 0xde, b == 0xdf:
		n, err := r.ReadMapLen()
		if err != nil {
			return err
		}
		for i := 0; i < n; i++ {
			tag, err := r.ReadUint()
			if err != nil {
				return fmt.Errorf("product: read field tag at entry %d: %w", i, err)
			}
			if err := decode(int(tag), r); err != nil {
				return fmt.Errorf("product: field tag %d: %w", tag, err)
			}
		}
		return nil
	case b >= 0x90 && b <= 0x9f, b == 0xdc, b == 0xdd:
		n, err := r.ReadArrayLen()
		if err != nil {
			return err
		}
		for i := 0; i < n; i++ {
			if err := decode(i, r); err != nil {
				return fmt.Errorf("product: index %d: %w", i, err)
			}
		}
		return nil
	default:
		return fmt.Errorf("product: expected map or array, got 0x%02x", b)
	}
}

// ReadEnum decodes a sum-typed value. `MsgpackTagged` always emits enum
// variants as a single-entry fixmap `{tag: payload}` regardless of the
// active EncodingStrategy. `decode` is called with the variant tag; the
// payload value sits next in the stream for the callback to consume.
//
// For unit variants the payload is nil — callbacks that don't read anything
// must still consume that nil. Use [ReadEnumUnit] when every variant is
// guaranteed to be unit.
func ReadEnum(r *Reader, decode func(tag int, r *Reader) error) error {
	n, err := r.ReadMapLen()
	if err != nil {
		return fmt.Errorf("enum: %w", err)
	}
	if n != 1 {
		return fmt.Errorf("enum: expected single-entry map, got %d entries", n)
	}
	tag, err := r.ReadUint()
	if err != nil {
		return fmt.Errorf("enum: read variant tag: %w", err)
	}
	if err := decode(int(tag), r); err != nil {
		return fmt.Errorf("enum tag %d: %w", tag, err)
	}
	return nil
}
