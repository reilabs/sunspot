package expression

import (
	"nr-groth16/bn254"
	"os"
	"testing"
)

func TestExpressionUnmarshalReaderEmpty(t *testing.T) {
	file, err := os.Open("../../binaries/expression/expression/expression_empty.bin")
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}

	var expr Expression[*bn254.BN254Field]
	if err := expr.UnmarshalReader(file); err != nil {
		t.Fatalf("Failed to unmarshal expression: %v", err)
	}

	expectedExpression := Expression[*bn254.BN254Field]{
		MulTerms:           []MulTerm[*bn254.BN254Field]{},
		LinearCombinations: []LinearCombination[*bn254.BN254Field]{},
		Constant:           &bn254.BN254Field{},
	}

	if !expr.Equals(&expectedExpression) {
		t.Errorf("Expected expression to be %v, got %v", expectedExpression, expr)
	}

	defer file.Close()
}

func TestExpressionUnmarshalReaderWithLinearCombinations(t *testing.T) {
	file, err := os.Open("../../binaries/expression/expression/expression_linear_combinations.bin")
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}

	var expr Expression[*bn254.BN254Field]
	if err := expr.UnmarshalReader(file); err != nil {
		t.Fatalf("Failed to unmarshal expression: %v", err)
	}

	expectedExpression := Expression[*bn254.BN254Field]{
		MulTerms: []MulTerm[*bn254.BN254Field]{},
		LinearCombinations: []LinearCombination[*bn254.BN254Field]{
			{Term: &bn254.BN254Field{}, Witness: 0},
			{Term: &bn254.BN254Field{}, Witness: 1234},
			{Term: &bn254.BN254Field{}, Witness: 5678},
		},
		Constant: &bn254.BN254Field{},
	}

	if !expr.Equals(&expectedExpression) {
		t.Errorf("Expected expression to be %v, got %v", expectedExpression, expr)
	}

	defer file.Close()
}

func TestExpressionUnmarshalReaderWithMulTerms(t *testing.T) {
	file, err := os.Open("../../binaries/expression/expression/expression_mul_terms.bin")
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}

	var expr Expression[*bn254.BN254Field]
	if err := expr.UnmarshalReader(file); err != nil {
		t.Fatalf("Failed to unmarshal expression: %v", err)
	}

	expectedExpression := Expression[*bn254.BN254Field]{
		MulTerms: []MulTerm[*bn254.BN254Field]{
			{Term: &bn254.BN254Field{}, WitnessLeft: 0, WitnessRight: 1},
			{Term: &bn254.BN254Field{}, WitnessLeft: 1234, WitnessRight: 5678},
			{Term: &bn254.BN254Field{}, WitnessLeft: 5678, WitnessRight: 1234},
		},
		LinearCombinations: []LinearCombination[*bn254.BN254Field]{},
		Constant:           &bn254.BN254Field{},
	}

	if !expr.Equals(&expectedExpression) {
		t.Errorf("Expected expression to be %v, got %v", expectedExpression, expr)
	}

	defer file.Close()
}

func TestExpressionUnmarshalReaderMulTermsWithLinearCombinations(t *testing.T) {
	file, err := os.Open("../../binaries/expression/expression/expression_mul_terms_with_linear_combinations.bin")
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}

	var expr Expression[*bn254.BN254Field]
	if err := expr.UnmarshalReader(file); err != nil {
		t.Fatalf("Failed to unmarshal expression: %v", err)
	}

	expectedExpression := Expression[*bn254.BN254Field]{
		MulTerms: []MulTerm[*bn254.BN254Field]{
			{Term: &bn254.BN254Field{}, WitnessLeft: 0, WitnessRight: 1},
			{Term: &bn254.BN254Field{}, WitnessLeft: 1234, WitnessRight: 5678},
			{Term: &bn254.BN254Field{}, WitnessLeft: 5678, WitnessRight: 1234},
		},
		LinearCombinations: []LinearCombination[*bn254.BN254Field]{
			{Term: &bn254.BN254Field{}, Witness: 0},
			{Term: &bn254.BN254Field{}, Witness: 1234},
			{Term: &bn254.BN254Field{}, Witness: 5678},
		},
		Constant: &bn254.BN254Field{},
	}

	if !expr.Equals(&expectedExpression) {
		t.Errorf("Expected expression to be %v, got %v", expectedExpression, expr)
	}

	defer file.Close()
}
