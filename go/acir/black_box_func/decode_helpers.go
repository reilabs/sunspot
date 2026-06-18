package blackboxfunc

import (
	"fmt"
	"sunspot/go/acir/msgpackutil"
	shr "sunspot/go/acir/shared"
)

// readFunctionInputArray reads a fixarray of FunctionInput values into a
// caller-provided fixed-size backing array. Returns an error if the wire
// length doesn't match.
func readFunctionInputArray[T shr.ACIRField](r *msgpackutil.Reader, out []FunctionInput[T]) error {
	n, err := r.ReadArrayLen()
	if err != nil {
		return err
	}
	if n != len(out) {
		return fmt.Errorf("expected fixed-size FunctionInput array of %d, got %d", len(out), n)
	}
	for i := 0; i < n; i++ {
		if err := out[i].UnmarshalReader(r); err != nil {
			return err
		}
	}
	return nil
}

// readFunctionInputVec reads a fixarray of FunctionInput values into a
// new slice and stores it through *out.
func readFunctionInputVec[T shr.ACIRField](r *msgpackutil.Reader, out *[]FunctionInput[T]) error {
	n, err := r.ReadArrayLen()
	if err != nil {
		return err
	}
	*out = make([]FunctionInput[T], n)
	for i := 0; i < n; i++ {
		if err := (*out)[i].UnmarshalReader(r); err != nil {
			return err
		}
	}
	return nil
}
