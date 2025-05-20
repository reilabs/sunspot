package black_box_ops

import (
	"io"
	mem "nr-groth16/acir_decoder/brillig/memory"
)

type EmbeddedCurveAdd struct {
	Input1X        mem.MemoryAddress
	Input1Y        mem.MemoryAddress
	Input1Infinite mem.MemoryAddress
	Input2X        mem.MemoryAddress
	Input2Y        mem.MemoryAddress
	Input2Infinite mem.MemoryAddress
	Result         mem.HeapArray
}

func (e *EmbeddedCurveAdd) UnmarshalReader(r io.Reader) error {
	if err := e.Input1X.UnmarshalReader(r); err != nil {
		return err
	}

	if err := e.Input1Y.UnmarshalReader(r); err != nil {
		return err
	}

	if err := e.Input1Infinite.UnmarshalReader(r); err != nil {
		return err
	}

	if err := e.Input2X.UnmarshalReader(r); err != nil {
		return err
	}

	if err := e.Input2Y.UnmarshalReader(r); err != nil {
		return err
	}

	if err := e.Input2Infinite.UnmarshalReader(r); err != nil {
		return err
	}

	if err := e.Result.UnmarshalReader(r); err != nil {
		return err
	}

	return nil
}
