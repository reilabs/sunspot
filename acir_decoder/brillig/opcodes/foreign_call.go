package opcodes

import (
	"encoding/binary"
	"io"
	mem "nr-groth16/acir_decoder/brillig/memory"
)

type ForeignCall struct {
	Function              string
	Destinations          []mem.ValueOrArray
	DestinationValueTypes []mem.HeapValueType
	Inputs                []mem.ValueOrArray
	InputValueTypes       []mem.HeapValueType
}

func (f *ForeignCall) UnmarshalReader(r io.Reader) error {
	var functionLength uint32
	if err := binary.Read(r, binary.LittleEndian, &functionLength); err != nil {
		return err
	}
	functionBytes := make([]byte, functionLength)
	if _, err := r.Read(functionBytes); err != nil {
		return err
	}
	f.Function = string(functionBytes)

	var numDestinations uint32
	if err := binary.Read(r, binary.LittleEndian, &numDestinations); err != nil {
		return err
	}
	f.Destinations = make([]mem.ValueOrArray, numDestinations)
	for i := uint32(0); i < numDestinations; i++ {
		if err := f.Destinations[i].UnmarshalReader(r); err != nil {
			return err
		}
	}

	var numValueTypes uint32
	if err := binary.Read(r, binary.LittleEndian, &numValueTypes); err != nil {
		return err
	}

	f.DestinationValueTypes = make([]mem.HeapValueType, numValueTypes)
	for i := uint32(0); i < numValueTypes; i++ {
		if err := f.DestinationValueTypes[i].UnmarshalReader(r); err != nil {
			return err
		}
	}

	var numInputs uint32
	if err := binary.Read(r, binary.LittleEndian, &numInputs); err != nil {
		return err
	}
	f.Inputs = make([]mem.ValueOrArray, numInputs)
	for i := uint32(0); i < numInputs; i++ {
		if err := f.Inputs[i].UnmarshalReader(r); err != nil {
			return err
		}
	}

	var numInputValueTypes uint32
	if err := binary.Read(r, binary.LittleEndian, &numInputValueTypes); err != nil {
		return err
	}
	f.InputValueTypes = make([]mem.HeapValueType, numInputValueTypes)
	for i := uint32(0); i < numInputValueTypes; i++ {
		if err := f.InputValueTypes[i].UnmarshalReader(r); err != nil {
			return err
		}
	}

	return nil
}
