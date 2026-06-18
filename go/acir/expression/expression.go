package expression

import (
	"encoding/json"
	"fmt"
	"sunspot/go/acir/msgpackutil"
	"sunspot/go/acir/opcodes"
	shr "sunspot/go/acir/shared"

	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
)

type Expression[T shr.ACIRField, E constraint.Element] struct {
	MulTerms           []MulTerm[T]           `json:"mul_terms"`           // Terms that are multiplied together
	LinearCombinations []LinearCombination[T] `json:"linear_combinations"` // Linear combinations of variables
	Constant           T                      `json:"constant"`
}

func (e *Expression[T, E]) Define(
	api frontend.Builder[E],
	witnesses map[shr.Witness]frontend.Variable,
) error {
	api.AssertIsEqual(e.Calculate(api, witnesses), 0)
	return nil
}

func (e *Expression[T, E]) UnmarshalReader(r *msgpackutil.Reader) error {
	e.Constant = shr.MakeNonNil(e.Constant)
	return msgpackutil.ReadStruct(r, expressionSchema, e.decode)
}

func (e *Expression[T, E]) decode(f msgpackutil.Field, r *msgpackutil.Reader) error {
	switch f.Tag {
	case 0:
		return msgpackutil.ReadVec(r, &e.MulTerms)
	case 1:
		return msgpackutil.ReadVec(r, &e.LinearCombinations)
	case 2:
		return e.Constant.UnmarshalReader(r)
	default:
		return fmt.Errorf("Expression: unknown field %s", f)
	}
}

var expressionSchema = msgpackutil.NewSchema(map[string]int{
	"mul_terms":           0,
	"linear_combinations": 1,
	"q_c":                 2,
})

func (e *Expression[T, E]) Equals(other opcodes.Opcode[E]) bool {
	value, ok := other.(*Expression[T, E])
	if !ok {
		return false
	}

	if len(e.MulTerms) != len(value.MulTerms) {
		return false
	}
	for i := range e.MulTerms {
		if !e.MulTerms[i].Equals(&value.MulTerms[i]) {
			return false
		}
	}

	if len(e.LinearCombinations) != len(value.LinearCombinations) {
		return false
	}
	for i := range e.LinearCombinations {
		if !e.LinearCombinations[i].Equals(&value.LinearCombinations[i]) {
			return false
		}
	}

	return e.Constant.Equals(value.Constant)
}

func (e *Expression[T, E]) Calculate(api frontend.API, witnesses map[shr.Witness]frontend.Variable) frontend.Variable {
	sum := e.Constant.ToFrontendVariable()
	for _, term := range e.MulTerms {
		sum = api.Add(sum, term.Calculate(api, witnesses))
	}
	for _, lc := range e.LinearCombinations {
		sum = api.Add(sum, lc.Calculate(api, witnesses))
	}

	return sum
}

func (e *Expression[T, E]) MarshalJSON() ([]byte, error) {
	stringMap := make(map[string]interface{})
	stringMap["assert_zero"] = e
	return json.Marshal(stringMap)
}

func (*Expression[T, E]) SerdeName() string { return "AssertZero" }
