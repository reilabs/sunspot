package opcodes

import (
	"encoding/binary"
	"fmt"
	"io"
)

type OpcodeLocation struct {
	Kind         OpcodeLocationKind
	ACIRAddress  *uint64
	ACIRIndex    *uint64
	BrilligIndex *uint64
}

type OpcodeLocationKind uint32

const (
	ACIROpcodeLocationKindACIR OpcodeLocationKind = iota
	ACIROpcodeLocationKindBrillig
)

func (o *OpcodeLocation) UnmarshalReader(r io.Reader) error {
	if err := binary.Read(r, binary.LittleEndian, &o.Kind); err != nil {
		return err
	}

	switch o.Kind {
	case ACIROpcodeLocationKindACIR:
		o.ACIRAddress = new(uint64)
		if err := binary.Read(r, binary.LittleEndian, o.ACIRAddress); err != nil {
			return err
		}
	case ACIROpcodeLocationKindBrillig:
		o.ACIRIndex = new(uint64)
		if err := binary.Read(r, binary.LittleEndian, o.ACIRIndex); err != nil {
			return err
		}

		o.BrilligIndex = new(uint64)
		if err := binary.Read(r, binary.LittleEndian, o.BrilligIndex); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown OpcodeLocation Kind: %d", o.Kind)
	}

	return nil
}
