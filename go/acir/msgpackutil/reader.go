// Package msgpackutil implements just enough of MessagePack to decode
// values emitted by any of noir's three ACIR serialization formats:
// `Msgpack` (string-keyed fixmap), `MsgpackCompact` (positional fixarray),
// and `MsgpackTagged` (int-keyed fixmap or positional fixarray per-type).
//
// Decoders consume whichever shape is on the wire by peeking the first
// byte — see [ReadStruct].
package msgpackutil

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"math"
)

// Reader is a peekable byte stream backed by a bufio.Reader. Every
// decode routine takes a *Reader.
// Also tracks maximum read witness values via embedded [witnessTracker]
type Reader struct {
	r *bufio.Reader
	witnessTracker
}

// NewReader wraps r as a Reader. Callers consuming a noir wire stream are
// responsible for reading and validating the `acir::serialization::Format`
// envelope byte themselves before constructing the Reader; this constructor
// only buffers the underlying stream.
func NewReader(r io.Reader) *Reader {
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
	if b != markerNil {
		return fmt.Errorf("msgpack: expected nil (0x%02x), got 0x%02x", markerNil, b)
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
	case markerFalse:
		return false, nil
	case markerTrue:
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
	case b <= markerPosFixintMax:
		return uint64(b), nil
	case b == markerUint8:
		v, err := r.readByte()
		return uint64(v), err
	case b == markerUint16:
		buf, err := r.readN(2)
		if err != nil {
			return 0, err
		}
		return uint64(binary.BigEndian.Uint16(buf)), nil
	case b == markerUint32:
		buf, err := r.readN(4)
		if err != nil {
			return 0, err
		}
		return uint64(binary.BigEndian.Uint32(buf)), nil
	case b == markerUint64:
		buf, err := r.readN(8)
		if err != nil {
			return 0, err
		}
		return binary.BigEndian.Uint64(buf), nil
	case b == markerInt8:
		v, err := r.readByte()
		if err != nil {
			return 0, err
		}
		if int8(v) < 0 {
			return 0, fmt.Errorf("msgpack: negative int where uint expected")
		}
		return uint64(int8(v)), nil
	case b == markerInt16:
		buf, err := r.readN(2)
		if err != nil {
			return 0, err
		}
		v := int16(binary.BigEndian.Uint16(buf))
		if v < 0 {
			return 0, fmt.Errorf("msgpack: negative int where uint expected")
		}
		return uint64(v), nil
	case b == markerInt32:
		buf, err := r.readN(4)
		if err != nil {
			return 0, err
		}
		v := int32(binary.BigEndian.Uint32(buf))
		if v < 0 {
			return 0, fmt.Errorf("msgpack: negative int where uint expected")
		}
		return uint64(v), nil
	case b == markerInt64:
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
	case b <= markerPosFixintMax:
		return int64(b), nil
	case b >= markerNegFixintLow:
		return int64(int8(b)), nil
	case b == markerUint8:
		v, err := r.readByte()
		return int64(v), err
	case b == markerUint16:
		buf, err := r.readN(2)
		if err != nil {
			return 0, err
		}
		return int64(binary.BigEndian.Uint16(buf)), nil
	case b == markerUint32:
		buf, err := r.readN(4)
		if err != nil {
			return 0, err
		}
		return int64(binary.BigEndian.Uint32(buf)), nil
	case b == markerUint64:
		buf, err := r.readN(8)
		if err != nil {
			return 0, err
		}
		return int64(binary.BigEndian.Uint64(buf)), nil
	case b == markerInt8:
		v, err := r.readByte()
		return int64(int8(v)), err
	case b == markerInt16:
		buf, err := r.readN(2)
		if err != nil {
			return 0, err
		}
		return int64(int16(binary.BigEndian.Uint16(buf))), nil
	case b == markerInt32:
		buf, err := r.readN(4)
		if err != nil {
			return 0, err
		}
		return int64(int32(binary.BigEndian.Uint32(buf))), nil
	case b == markerInt64:
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
	case markerFloat32:
		buf, err := r.readN(4)
		if err != nil {
			return 0, err
		}
		return float64(math.Float32frombits(binary.BigEndian.Uint32(buf))), nil
	case markerFloat64:
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
	case b >= markerFixstrLow && b <= markerFixstrHigh:
		return int(b & markerFixstrLenMask), nil
	case b == markerStr8:
		v, err := r.readByte()
		return int(v), err
	case b == markerStr16:
		buf, err := r.readN(2)
		if err != nil {
			return 0, err
		}
		return int(binary.BigEndian.Uint16(buf)), nil
	case b == markerStr32:
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
	case markerBin8:
		v, err := r.readByte()
		if err != nil {
			return nil, err
		}
		n = int(v)
	case markerBin16:
		buf, err := r.readN(2)
		if err != nil {
			return nil, err
		}
		n = int(binary.BigEndian.Uint16(buf))
	case markerBin32:
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
	case b >= markerFixarrayLow && b <= markerFixarrayHigh:
		return int(b & markerFixContainerLenMask), nil
	case b == markerArray16:
		buf, err := r.readN(2)
		if err != nil {
			return 0, err
		}
		return int(binary.BigEndian.Uint16(buf)), nil
	case b == markerArray32:
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
	case b >= markerFixmapLow && b <= markerFixmapHigh:
		return int(b & markerFixContainerLenMask), nil
	case b == markerMap16:
		buf, err := r.readN(2)
		if err != nil {
			return 0, err
		}
		return int(binary.BigEndian.Uint16(buf)), nil
	case b == markerMap32:
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
	case b == markerNil, b == markerFalse, b == markerTrue:
		_, err = r.readByte()
		return err
	case b <= markerPosFixintMax, b >= markerNegFixintLow:
		_, err = r.readByte()
		return err
	case b >= markerFixstrLow && b <= markerFixstrHigh:
		_, err = r.ReadString()
		return err
	case b >= markerFixarrayLow && b <= markerFixarrayHigh, b == markerArray16, b == markerArray32:
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
	case b >= markerFixmapLow && b <= markerFixmapHigh, b == markerMap16, b == markerMap32:
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
	case b == markerBin8, b == markerBin16, b == markerBin32:
		_, err = r.ReadBytes()
		return err
	case b == markerFloat32, b == markerFloat64:
		_, err = r.ReadFloat64()
		return err
	case b == markerUint8, b == markerUint16, b == markerUint32, b == markerUint64:
		_, err = r.ReadUint()
		return err
	case b == markerInt8, b == markerInt16, b == markerInt32, b == markerInt64:
		_, err = r.ReadInt()
		return err
	case b == markerStr8, b == markerStr16, b == markerStr32:
		_, err = r.ReadString()
		return err
	default:
		return fmt.Errorf("msgpack: SkipValue: unsupported marker 0x%02x", b)
	}
}
