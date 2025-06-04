package blackboxfunc

import (
	"encoding/binary"
	"io"
)

type BigIntAdd struct {
	Lhs    uint32
	Rhs    uint32
	Output uint32
}

func (a *BigIntAdd) UnmarshalReader(r io.Reader) error {
	if err := binary.Read(r, binary.LittleEndian, &a.Lhs); err != nil {
		return err
	}
	if err := binary.Read(r, binary.LittleEndian, &a.Rhs); err != nil {
		return err
	}
	if err := binary.Read(r, binary.LittleEndian, &a.Output); err != nil {
		return err
	}
	return nil
}

func (a *BigIntAdd) Equals(other *BigIntAdd) bool {
	return a.Lhs == other.Lhs && a.Rhs == other.Rhs && a.Output == other.Output
}
