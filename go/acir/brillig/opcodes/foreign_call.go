package opcodes

import (
	"encoding/binary"
	"io"
	mem "sunspot/go/acir/brillig/memory"
)

type ForeignCall struct {
	Function              string
	Destinations          []mem.ValueOrArray
	DestinationValueTypes []mem.HeapValueType
	Inputs                []mem.ValueOrArray
	InputValueTypes       []mem.HeapValueType
}

func (f *ForeignCall) UnmarshalReader(r io.Reader) error {
	var functionLength uint64
	if err := binary.Read(r, binary.LittleEndian, &functionLength); err != nil {
		return err
	}
	functionBytes := make([]byte, functionLength)
	if _, err := r.Read(functionBytes); err != nil {
		return err
	}
	f.Function = string(functionBytes)

	var numDestinations uint64
	if err := binary.Read(r, binary.LittleEndian, &numDestinations); err != nil {
		return err
	}
	f.Destinations = make([]mem.ValueOrArray, numDestinations)
	for i := uint64(0); i < numDestinations; i++ {
		if err := f.Destinations[i].UnmarshalReader(r); err != nil {
			return err
		}
	}

	var numValueTypes uint64
	if err := binary.Read(r, binary.LittleEndian, &numValueTypes); err != nil {
		return err
	}

	f.DestinationValueTypes = make([]mem.HeapValueType, numValueTypes)
	for i := uint64(0); i < numValueTypes; i++ {
		if err := f.DestinationValueTypes[i].UnmarshalReader(r); err != nil {
			return err
		}
	}

	var numInputs uint64
	if err := binary.Read(r, binary.LittleEndian, &numInputs); err != nil {
		return err
	}
	f.Inputs = make([]mem.ValueOrArray, numInputs)
	for i := uint64(0); i < numInputs; i++ {
		if err := f.Inputs[i].UnmarshalReader(r); err != nil {
			return err
		}
	}

	var numInputValueTypes uint64
	if err := binary.Read(r, binary.LittleEndian, &numInputValueTypes); err != nil {
		return err
	}
	f.InputValueTypes = make([]mem.HeapValueType, numInputValueTypes)
	for i := uint64(0); i < numInputValueTypes; i++ {
		if err := f.InputValueTypes[i].UnmarshalReader(r); err != nil {
			return err
		}
	}

	return nil
}

func (f *ForeignCall) Equals(other ForeignCall) bool {
	if f.Function != other.Function {
		return false
	}

	if len(f.Destinations) != len(other.Destinations) ||
		len(f.DestinationValueTypes) != len(other.DestinationValueTypes) ||
		len(f.Inputs) != len(other.Inputs) ||
		len(f.InputValueTypes) != len(other.InputValueTypes) {
		return false
	}

	for i := range f.Destinations {
		if !f.Destinations[i].Equals(other.Destinations[i]) {
			return false
		}
	}

	for i := range f.DestinationValueTypes {
		if !f.DestinationValueTypes[i].Equals(other.DestinationValueTypes[i]) {
			return false
		}
	}

	for i := range f.Inputs {
		if !f.Inputs[i].Equals(other.Inputs[i]) {
			return false
		}
	}

	for i := range f.InputValueTypes {
		if !f.InputValueTypes[i].Equals(other.InputValueTypes[i]) {
			return false
		}
	}

	return true
}
