package blackboxfunc

import (
	"encoding/binary"
	"fmt"
	"io"
	shr "nr-groth16/acir/shared"

	"github.com/consensys/gnark/frontend"
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

func (a *BigIntDiv) Equals(other *BigIntDiv) bool {
	return a.Lhs == other.Lhs && a.Rhs == other.Rhs && a.Output == other.Output
}

func (a *BigIntDiv) Define(api frontend.API, witnesses map[shr.Witness]frontend.Variable) error {
	lhs, ok := witnesses[shr.Witness(a.Lhs)]
	if !ok {
		return fmt.Errorf("witness for LHS not found")
	}
	rhs, ok := witnesses[shr.Witness(a.Rhs)]
	if !ok {
		return fmt.Errorf("witness for RHS not found")
	}
	output, ok := witnesses[shr.Witness(a.Output)]
	if !ok {
		return fmt.Errorf("witness for Output not found")
	}

	api.AssertIsEqual(output, api.Div(lhs, rhs))
	return nil
}
