// Package msgpackutil implements just enough of MessagePack to decode
// values emitted by Noir's `msgpack_tagged` serializer.
//
// The serializer can write a struct as either an int-keyed fixmap (the
// "Tagged" strategy) or a positional fixarray (the "Array" strategy),
// chosen per-type. Decoders consume whichever shape is on the wire by
// peeking the first byte — see [ReadStruct].
package msgpackutil

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"math"
)

// Reader is a peekable byte stream backed by a bufio.Reader. Every
// MsgpackTagged decode routine takes a *Reader.
// Also tracks maximum read witness values via embedded [witnessTracker]
type Reader struct {
	r *bufio.Reader
	witnessTracker
}

func NewReader(r io.Reader) *Reader {
	if br, ok := r.(*bufio.Reader); ok {
		return &Reader{r: br}
	}
	return &Reader{r: bufio.NewReader(r)}
}

// Peek returns the next type marker byte without consuming it.
func (r *Reader) Peek() (byte, error) {
	b, err := r.r.Peek(1)
	if err != nil {
		return 0, err
	}
	return b[0], nil
}

func (r *Reader) readByte() (byte, error) { return r.r.ReadByte() }

func (r *Reader) readN(n int) ([]byte, error) {
	buf := make([]byte, n)
	if _, err := io.ReadFull(r.r, buf); err != nil {
		return nil, err
	}
	return buf, nil
}

// ReadNil consumes a nil marker.
func (r *Reader) ReadNil() error {
	b, err := r.readByte()
	if err != nil {
		return err
	}
	if b != 0xc0 {
		return fmt.Errorf("msgpack: expected nil (0xc0), got 0x%02x", b)
	}
	return nil
}

// ReadBool consumes a bool.
func (r *Reader) ReadBool() (bool, error) {
	b, err := r.readByte()
	if err != nil {
		return false, err
	}
	switch b {
	case 0xc2:
		return false, nil
	case 0xc3:
		return true, nil
	default:
		return false, fmt.Errorf("msgpack: expected bool, got 0x%02x", b)
	}
}

// ReadUint decodes any unsigned integer encoding (positive fixint, uint8,
// uint16, uint32, uint64). Accepts non-negative fixed-width signed encodings
// too — useful for small enum tags that the encoder may emit as int.
func (r *Reader) ReadUint() (uint64, error) {
	b, err := r.readByte()
	if err != nil {
		return 0, err
	}
	switch {
	case b <= 0x7f:
		return uint64(b), nil
	case b == 0xcc:
		v, err := r.readByte()
		return uint64(v), err
	case b == 0xcd:
		buf, err := r.readN(2)
		if err != nil {
			return 0, err
		}
		return uint64(binary.BigEndian.Uint16(buf)), nil
	case b == 0xce:
		buf, err := r.readN(4)
		if err != nil {
			return 0, err
		}
		return uint64(binary.BigEndian.Uint32(buf)), nil
	case b == 0xcf:
		buf, err := r.readN(8)
		if err != nil {
			return 0, err
		}
		return binary.BigEndian.Uint64(buf), nil
	case b == 0xd0:
		v, err := r.readByte()
		if err != nil {
			return 0, err
		}
		if int8(v) < 0 {
			return 0, fmt.Errorf("msgpack: negative int where uint expected")
		}
		return uint64(int8(v)), nil
	case b == 0xd1:
		buf, err := r.readN(2)
		if err != nil {
			return 0, err
		}
		v := int16(binary.BigEndian.Uint16(buf))
		if v < 0 {
			return 0, fmt.Errorf("msgpack: negative int where uint expected")
		}
		return uint64(v), nil
	case b == 0xd2:
		buf, err := r.readN(4)
		if err != nil {
			return 0, err
		}
		v := int32(binary.BigEndian.Uint32(buf))
		if v < 0 {
			return 0, fmt.Errorf("msgpack: negative int where uint expected")
		}
		return uint64(v), nil
	case b == 0xd3:
		buf, err := r.readN(8)
		if err != nil {
			return 0, err
		}
		v := int64(binary.BigEndian.Uint64(buf))
		if v < 0 {
			return 0, fmt.Errorf("msgpack: negative int where uint expected")
		}
		return uint64(v), nil
	default:
		return 0, fmt.Errorf("msgpack: expected uint, got 0x%02x", b)
	}
}

// ReadU32 reads any uint encoding and narrows to uint32, erroring on overflow.
func (r *Reader) ReadU32() (uint32, error) {
	v, err := r.ReadUint()
	if err != nil {
		return 0, err
	}
	if v > math.MaxUint32 {
		return 0, fmt.Errorf("u32 field overflow: %d", v)
	}
	return uint32(v), nil
}

// ReadInt decodes a signed integer in any encoding.
func (r *Reader) ReadInt() (int64, error) {
	b, err := r.readByte()
	if err != nil {
		return 0, err
	}
	switch {
	case b <= 0x7f:
		return int64(b), nil
	case b >= 0xe0:
		return int64(int8(b)), nil
	case b == 0xcc:
		v, err := r.readByte()
		return int64(v), err
	case b == 0xcd:
		buf, err := r.readN(2)
		if err != nil {
			return 0, err
		}
		return int64(binary.BigEndian.Uint16(buf)), nil
	case b == 0xce:
		buf, err := r.readN(4)
		if err != nil {
			return 0, err
		}
		return int64(binary.BigEndian.Uint32(buf)), nil
	case b == 0xcf:
		buf, err := r.readN(8)
		if err != nil {
			return 0, err
		}
		return int64(binary.BigEndian.Uint64(buf)), nil
	case b == 0xd0:
		v, err := r.readByte()
		return int64(int8(v)), err
	case b == 0xd1:
		buf, err := r.readN(2)
		if err != nil {
			return 0, err
		}
		return int64(int16(binary.BigEndian.Uint16(buf))), nil
	case b == 0xd2:
		buf, err := r.readN(4)
		if err != nil {
			return 0, err
		}
		return int64(int32(binary.BigEndian.Uint32(buf))), nil
	case b == 0xd3:
		buf, err := r.readN(8)
		if err != nil {
			return 0, err
		}
		return int64(binary.BigEndian.Uint64(buf)), nil
	default:
		return 0, fmt.Errorf("msgpack: expected int, got 0x%02x", b)
	}
}

// ReadFloat64 decodes a float (accepts float32 or float64).
func (r *Reader) ReadFloat64() (float64, error) {
	b, err := r.readByte()
	if err != nil {
		return 0, err
	}
	switch b {
	case 0xca:
		buf, err := r.readN(4)
		if err != nil {
			return 0, err
		}
		return float64(math.Float32frombits(binary.BigEndian.Uint32(buf))), nil
	case 0xcb:
		buf, err := r.readN(8)
		if err != nil {
			return 0, err
		}
		return math.Float64frombits(binary.BigEndian.Uint64(buf)), nil
	default:
		return 0, fmt.Errorf("msgpack: expected float, got 0x%02x", b)
	}
}

// ReadString decodes a string (any width).
func (r *Reader) ReadString() (string, error) {
	n, err := r.readStringLen()
	if err != nil {
		return "", err
	}
	buf, err := r.readN(n)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

func (r *Reader) readStringLen() (int, error) {
	b, err := r.readByte()
	if err != nil {
		return 0, err
	}
	switch {
	case b >= 0xa0 && b <= 0xbf:
		return int(b & 0x1f), nil
	case b == 0xd9:
		v, err := r.readByte()
		return int(v), err
	case b == 0xda:
		buf, err := r.readN(2)
		if err != nil {
			return 0, err
		}
		return int(binary.BigEndian.Uint16(buf)), nil
	case b == 0xdb:
		buf, err := r.readN(4)
		if err != nil {
			return 0, err
		}
		return int(binary.BigEndian.Uint32(buf)), nil
	default:
		return 0, fmt.Errorf("msgpack: expected string, got 0x%02x", b)
	}
}

// ReadBytes decodes a bin* value.
func (r *Reader) ReadBytes() ([]byte, error) {
	b, err := r.readByte()
	if err != nil {
		return nil, err
	}
	var n int
	switch b {
	case 0xc4:
		v, err := r.readByte()
		if err != nil {
			return nil, err
		}
		n = int(v)
	case 0xc5:
		buf, err := r.readN(2)
		if err != nil {
			return nil, err
		}
		n = int(binary.BigEndian.Uint16(buf))
	case 0xc6:
		buf, err := r.readN(4)
		if err != nil {
			return nil, err
		}
		n = int(binary.BigEndian.Uint32(buf))
	default:
		return nil, fmt.Errorf("msgpack: expected bin, got 0x%02x", b)
	}
	return r.readN(n)
}

// ReadArrayLen reads any array marker and returns its length.
func (r *Reader) ReadArrayLen() (int, error) {
	b, err := r.readByte()
	if err != nil {
		return 0, err
	}
	switch {
	case b >= 0x90 && b <= 0x9f:
		return int(b & 0x0f), nil
	case b == 0xdc:
		buf, err := r.readN(2)
		if err != nil {
			return 0, err
		}
		return int(binary.BigEndian.Uint16(buf)), nil
	case b == 0xdd:
		buf, err := r.readN(4)
		if err != nil {
			return 0, err
		}
		return int(binary.BigEndian.Uint32(buf)), nil
	default:
		return 0, fmt.Errorf("msgpack: expected array, got 0x%02x", b)
	}
}

// ReadMapLen reads any map marker and returns its entry count.
func (r *Reader) ReadMapLen() (int, error) {
	b, err := r.readByte()
	if err != nil {
		return 0, err
	}
	switch {
	case b >= 0x80 && b <= 0x8f:
		return int(b & 0x0f), nil
	case b == 0xde:
		buf, err := r.readN(2)
		if err != nil {
			return 0, err
		}
		return int(binary.BigEndian.Uint16(buf)), nil
	case b == 0xdf:
		buf, err := r.readN(4)
		if err != nil {
			return 0, err
		}
		return int(binary.BigEndian.Uint32(buf)), nil
	default:
		return 0, fmt.Errorf("msgpack: expected map, got 0x%02x", b)
	}
}

// SkipValue consumes and discards one MessagePack value. Used to step over
// unknown tags in `allow_unknown_tags` products and reserved positions.
func (r *Reader) SkipValue() error {
	b, err := r.Peek()
	if err != nil {
		return err
	}
	switch {
	case b == 0xc0, b == 0xc2, b == 0xc3:
		_, err = r.readByte()
		return err
	case b <= 0x7f, b >= 0xe0:
		_, err = r.readByte()
		return err
	case b >= 0xa0 && b <= 0xbf:
		_, err = r.ReadString()
		return err
	case b >= 0x90 && b <= 0x9f, b == 0xdc, b == 0xdd:
		n, err := r.ReadArrayLen()
		if err != nil {
			return err
		}
		for i := 0; i < n; i++ {
			if err := r.SkipValue(); err != nil {
				return err
			}
		}
		return nil
	case b >= 0x80 && b <= 0x8f, b == 0xde, b == 0xdf:
		n, err := r.ReadMapLen()
		if err != nil {
			return err
		}
		for i := 0; i < 2*n; i++ {
			if err := r.SkipValue(); err != nil {
				return err
			}
		}
		return nil
	case b == 0xc4, b == 0xc5, b == 0xc6:
		_, err = r.ReadBytes()
		return err
	case b == 0xca, b == 0xcb:
		_, err = r.ReadFloat64()
		return err
	case b >= 0xcc && b <= 0xcf:
		_, err = r.ReadUint()
		return err
	case b >= 0xd0 && b <= 0xd3:
		_, err = r.ReadInt()
		return err
	case b == 0xd9, b == 0xda, b == 0xdb:
		_, err = r.ReadString()
		return err
	default:
		return fmt.Errorf("msgpack: SkipValue: unsupported marker 0x%02x", b)
	}
}
