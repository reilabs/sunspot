package blackboxfunc

import (
	"fmt"
	"io"
	"math/big"
	shr "nr-groth16/acir/shared"

	"github.com/consensys/gnark/frontend"
	"github.com/rs/zerolog/log"
)

type Range[T shr.ACIRField] struct {
	Input FunctionInput[T]
}

func (a *Range[T]) UnmarshalReader(r io.Reader) error {
	if err := a.Input.UnmarshalReader(r); err != nil {
		return err
	}
	return nil
}

func (a Range[T]) Equals(other BlackBoxFunction) bool {
	value, ok := other.(*Range[T])
	return ok && a.Input.Equals(&value.Input)
}

func (a Range[T]) Define(api frontend.API, witnesses map[shr.Witness]frontend.Variable) error {
	if a.Input.FunctionInputKind == ACIRFunctionInputKindConstant {
		return nil
	}

	witness := a.Input.Witness
	if witness == nil {
		return fmt.Errorf("witness is nil for Range function input")
	}

	w, ok := witnesses[*witness]
	if !ok {
		return fmt.Errorf("witness %v not found in witnesses map", *witness)
	}

	max_value := big.NewInt(1)
	max_value = max_value.Lsh(max_value, uint(a.Input.NumberOfBits)) // 2^n
	max_value = max_value.Sub(max_value, big.NewInt(1))              // 2^n - 1
	log.Trace().Msgf("IMPOSING RANGE CONSTRAINT: %s FOR %d", max_value.String(), a.Input.NumberOfBits)

	_ = max_value
	api.AssertIsLessOrEqual(w, max_value)
	api.AssertIsLessOrEqual(big.NewInt(0), w) // Ensure w is non-negative

	return nil
}
