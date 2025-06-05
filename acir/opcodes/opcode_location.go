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

func (o *OpcodeLocation) Equals(other *OpcodeLocation) bool {
	if o.Kind != other.Kind {
		return false
	}

	switch o.Kind {
	case ACIROpcodeLocationKindACIR:
		if o.ACIRAddress == nil || other.ACIRAddress == nil {
			return o.ACIRAddress == other.ACIRAddress
		}
		return *o.ACIRAddress == *other.ACIRAddress

	case ACIROpcodeLocationKindBrillig:
		if o.ACIRIndex == nil || other.ACIRIndex == nil || o.BrilligIndex == nil || other.BrilligIndex == nil {
			return o.ACIRIndex == other.ACIRIndex && o.BrilligIndex == other.BrilligIndex
		}
		return *o.ACIRIndex == *other.ACIRIndex && *o.BrilligIndex == *other.BrilligIndex

	default:
		return false
	}
}
