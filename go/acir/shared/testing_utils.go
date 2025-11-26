package shared

import (
	"encoding/binary"
	"io"
	"testing"
)

func ParseThrough32bits(t *testing.T, r io.Reader) uint32 {
	var kind uint32
	if err := binary.Read(r, binary.LittleEndian, &kind); err != nil {
		t.Fatal("was not able to read type")
	}
	return kind
}
