package blackboxfunc

import (
	"encoding/binary"
	"io"
)

type BigIntDiv struct {
	Lhs    uint32
	Rhs    uint32
	Output uint32
}

func (a *BigIntDiv) UnmarshalReader(r io.Reader) error {
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
