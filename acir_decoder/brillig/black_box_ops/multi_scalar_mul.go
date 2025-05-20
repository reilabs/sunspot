package black_box_ops

import (
	"io"
	mem "nr-groth16/acir_decoder/brillig/memory"
)

type MultiScalarMul struct {
	Points  mem.HeapVector
	Scalars mem.HeapVector
	Outputs mem.HeapArray
}

func (m *MultiScalarMul) UnmarshalReader(r io.Reader) error {
	if err := m.Points.UnmarshalReader(r); err != nil {
		return err
	}

	if err := m.Scalars.UnmarshalReader(r); err != nil {
		return err
	}

	if err := m.Outputs.UnmarshalReader(r); err != nil {
		return err
	}

	return nil
}
