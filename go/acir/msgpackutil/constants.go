package msgpackutil

// MessagePack format markers — the type-tag byte that prefixes every value on
// the wire. See https://github.com/msgpack/msgpack/blob/master/spec.md.
//
// The variable-length-in-marker families (positive fixint, negative fixint,
// fixmap, fixarray, fixstr) embed length or value in the low bits of the
// marker byte itself; those come as Low/High range bounds plus the mask used
// to extract the embedded length. Everything else is a single fixed byte.
const (
	// Singleton markers.
	markerNil   = 0xc0
	markerFalse = 0xc2
	markerTrue  = 0xc3

	// Width-tagged scalars.
	markerBin8    = 0xc4
	markerBin16   = 0xc5
	markerBin32   = 0xc6
	markerFloat32 = 0xca
	markerFloat64 = 0xcb
	markerUint8   = 0xcc
	markerUint16  = 0xcd
	markerUint32  = 0xce
	markerUint64  = 0xcf
	markerInt8    = 0xd0
	markerInt16   = 0xd1
	markerInt32   = 0xd2
	markerInt64   = 0xd3

	// Width-tagged strings (length follows as 1/2/4 BE bytes).
	markerStr8  = 0xd9
	markerStr16 = 0xda
	markerStr32 = 0xdb

	// Width-tagged container headers (length follows as 2/4 BE bytes).
	markerArray16 = 0xdc
	markerArray32 = 0xdd
	markerMap16   = 0xde
	markerMap32   = 0xdf

	// Positive fixint: 0x00..0x7f, marker IS the value (uint7).
	markerPosFixintMax = 0x7f

	// Negative fixint: 0xe0..0xff, marker reinterpreted as int8.
	markerNegFixintLow = 0xe0

	// fixmap: 0x80..0x8f. Low 4 bits = entry count.
	markerFixmapLow  = 0x80
	markerFixmapHigh = 0x8f

	// fixarray: 0x90..0x9f. Low 4 bits = element count.
	markerFixarrayLow  = 0x90
	markerFixarrayHigh = 0x9f

	// fixstr: 0xa0..0xbf. Low 5 bits = byte length.
	markerFixstrLow  = 0xa0
	markerFixstrHigh = 0xbf

	// Length masks for the fix* families.
	markerFixContainerLenMask = 0x0f // fixmap, fixarray
	markerFixstrLenMask       = 0x1f
)
