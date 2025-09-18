package expression

import (
	"math/big"
	shr "nr-groth16/acir/shared"
	"nr-groth16/bn254"
	"os"
	"testing"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/backend/witness"
	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/consensys/gnark/frontend/schema"
)

func TestExpressionUnmarshalReaderEmpty(t *testing.T) {
	type T = *bn254.BN254Field
	type E = constraint.U64
	file, err := os.Open("../../binaries/expression/expression/expression_empty.bin")
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}

	var expr Expression[T, E]
	if err := expr.UnmarshalReader(file); err != nil {
		t.Fatalf("Failed to unmarshal expression: %v", err)
	}

	expectedExpression := Expression[T, E]{
		MulTerms:           []MulTerm[T]{},
		LinearCombinations: []LinearCombination[T]{},
		Constant:           bn254.Zero(),
	}

	if !expr.Equals(&expectedExpression) {
		t.Errorf("Expected expression to be %v, got %v", expectedExpression, expr)
	}

	defer file.Close()
}

func TestExpressionUnmarshalReaderWithLinearCombinations(t *testing.T) {
	type T = *bn254.BN254Field
	type E = constraint.U64
	file, err := os.Open("../../binaries/expression/expression/expression_linear_combinations.bin")
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}

	var expr Expression[T, E]
	if err := expr.UnmarshalReader(file); err != nil {
		t.Fatalf("Failed to unmarshal expression: %v", err)
	}

	expectedExpression := Expression[T, E]{
		MulTerms: []MulTerm[T]{},
		LinearCombinations: []LinearCombination[T]{
			{Term: bn254.One(), Witness: 0},
			{Term: bn254.One(), Witness: 1234},
			{Term: bn254.One(), Witness: 5678},
		},
		Constant: bn254.Zero(),
	}

	if !expr.Equals(&expectedExpression) {
		t.Errorf("Expected expression to be %v, got %v", expectedExpression, expr)
	}

	defer file.Close()
}

func TestExpressionUnmarshalReaderWithMulTerms(t *testing.T) {
	type T = *bn254.BN254Field
	type E = constraint.U64
	file, err := os.Open("../../binaries/expression/expression/expression_mul_terms.bin")
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}

	var expr Expression[T, E]
	if err := expr.UnmarshalReader(file); err != nil {
		t.Fatalf("Failed to unmarshal expression: %v", err)
	}

	expectedExpression := Expression[T, E]{
		MulTerms: []MulTerm[T]{
			{Term: bn254.One(), WitnessLeft: 0, WitnessRight: 1},
			{Term: bn254.One(), WitnessLeft: 1234, WitnessRight: 5678},
			{Term: bn254.One(), WitnessLeft: 5678, WitnessRight: 1234},
		},
		LinearCombinations: []LinearCombination[T]{},
		Constant:           bn254.Zero(),
	}

	if !expr.Equals(&expectedExpression) {
		t.Errorf("Expected expression to be %v, got %v", expectedExpression, expr)
	}

	defer file.Close()
}

func TestExpressionUnmarshalReaderMulTermsWithLinearCombinations(t *testing.T) {
	type T = *bn254.BN254Field
	type E = constraint.U64
	file, err := os.Open("../../binaries/expression/expression/expression_mul_terms_with_linear_combinations.bin")
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}

	var expr Expression[T, E]
	if err := expr.UnmarshalReader(file); err != nil {
		t.Fatalf("Failed to unmarshal expression: %v", err)
	}

	expectedExpression := Expression[T, E]{
		MulTerms: []MulTerm[T]{
			{Term: bn254.One(), WitnessLeft: 0, WitnessRight: 1},
			{Term: bn254.One(), WitnessLeft: 1234, WitnessRight: 5678},
			{Term: bn254.One(), WitnessLeft: 5678, WitnessRight: 1234},
		},
		LinearCombinations: []LinearCombination[T]{
			{Term: bn254.One(), Witness: 0},
			{Term: bn254.One(), Witness: 1234},
			{Term: bn254.One(), Witness: 5678},
		},
		Constant: bn254.Zero(),
	}

	if !expr.Equals(&expectedExpression) {
		t.Errorf("Expected expression to be %v, got %v", expectedExpression, expr)
	}

	defer file.Close()
}

func TestExpressionExecutionEmpty(t *testing.T) {
	type T = *bn254.BN254Field
	type E = constraint.U64
	builder, err := r1cs.NewBuilder[E](ecc.BN254.ScalarField(), frontend.CompileConfig{
		CompressThreshold: 300,
	})
	if err != nil {
		t.Fatalf("Failed to create R1CS builder: %v", err)
	}

	expression := Expression[T, E]{
		MulTerms:           []MulTerm[T]{},
		LinearCombinations: []LinearCombination[T]{},
		Constant:           bn254.Zero(),
	}

	builder.AssertIsEqual(expression.Calculate(builder, map[shr.Witness]frontend.Variable{}), 0)

	ccs, err := builder.Compile()
	if err != nil {
		t.Fatalf("Failed to compile R1CS: %v", err)
	}

	witness, err := witness.New(fr.Modulus())
	if err != nil {
		t.Fatalf("Failed to create witness: %v", err)
	}

	pk, vk, err := groth16.Setup(ccs)
	if err != nil {
		t.Fatalf("Failed to setup Groth16: %v", err)
	}

	proof, err := groth16.Prove(ccs, pk, witness)
	if err != nil {
		t.Fatalf("Failed to prove Groth16: %v", err)
	}

	publicWitness, err := witness.Public()
	if err != nil {
		t.Fatalf("Failed to get public witness: %v", err)
	}

	err = groth16.Verify(proof, vk, publicWitness)
	if err != nil {
		t.Fatalf("Verification failed: %v", err)
	} else {
		t.Logf("Verification succeeded!")
	}
}

func TestExpressionExecution(t *testing.T) {
	type T = *bn254.BN254Field
	type E = constraint.U64
	builder, err := r1cs.NewBuilder[E](ecc.BN254.ScalarField(), frontend.CompileConfig{
		CompressThreshold: 300,
	})
	if err != nil {
		t.Fatalf("Failed to create R1CS builder: %v", err)
	}

	pubVarX := builder.PublicVariable(schema.LeafInfo{
		FullName:   func() string { return "testPublicVariableX" },
		Visibility: schema.Public,
	})

	pubVarZ := builder.PublicVariable(schema.LeafInfo{
		FullName:   func() string { return "testPublicVariableZ" },
		Visibility: schema.Public,
	})

	privVarY := builder.SecretVariable(schema.LeafInfo{
		FullName:   func() string { return "testSecretVariableY" },
		Visibility: schema.Secret,
	})

	witnessMap := map[shr.Witness]frontend.Variable{
		1: pubVarX,
		2: privVarY,
		3: pubVarZ,
	}

	expression := &Expression[T, E]{
		MulTerms: []MulTerm[T]{
			{
				Term:         bn254.One(),
				WitnessLeft:  1,
				WitnessRight: 2,
			},
		},
		LinearCombinations: []LinearCombination[T]{
			{
				Term:    bn254.One(),
				Witness: 3,
			},
		},
		Constant: bn254.Zero(),
	}

	builder.AssertIsEqual(expression.Calculate(builder, witnessMap), 0)

	ccs, err := builder.Compile()
	if err != nil {
		t.Fatalf("Failed to compile R1CS: %v", err)
	}

	if ccs.GetNbPublicVariables() != 3 {
		t.Errorf("Expected 3 public variables, got %d", ccs.GetNbPublicVariables())
	}
	if ccs.GetNbSecretVariables() != 1 {
		t.Errorf("Expected 1 secret variables, got %d", ccs.GetNbSecretVariables())
	}

	witness, err := witness.New(fr.Modulus())
	if err != nil {
		t.Fatalf("Failed to create witness: %v", err)
	}

	values := make(chan any)

	go func() {
		values <- big.NewInt(2)
		values <- big.NewInt(-850)
		values <- big.NewInt(425)
		close(values)
	}()

	err = witness.Fill(2, 1, values)
	if err != nil {
		t.Fatalf("Failed to fill witness: %v", err)
	}

	pk, vk, err := groth16.Setup(ccs)
	if err != nil {
		t.Fatalf("Failed to setup Groth16: %v", err)
	}

	proof, err := groth16.Prove(ccs, pk, witness)
	if err != nil {
		t.Fatalf("Failed to prove Groth16: %v", err)
	}

	publicWitness, err := witness.Public()
	if err != nil {
		t.Fatalf("Failed to get public witness: %v", err)
	}

	err = groth16.Verify(proof, vk, publicWitness)
	if err != nil {
		t.Fatalf("Verification failed: %v", err)
	} else {
		t.Logf("Verification succeeded!")
	}
}
