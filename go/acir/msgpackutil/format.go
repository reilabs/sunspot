package msgpackutil

import (
	"fmt"
	"io"
)

// ConsumeFormatByte reads and validates the noir
// `acir::serialization::Format` envelope byte preceding the msgpack payload:
// 2=Msgpack, 3=MsgpackCompact, 4=MsgpackTagged.
func ConsumeFormatByte(r io.Reader) error {
	var b [1]byte
	if _, err := io.ReadFull(r, b[:]); err != nil {
		return fmt.Errorf("read format byte: %w", err)
	}
	switch b[0] {
	case 2, 3, 4:
		return nil
	default:
		return fmt.Errorf("unsupported ACIR serialization format byte 0x%02x", b[0])
	}
}
