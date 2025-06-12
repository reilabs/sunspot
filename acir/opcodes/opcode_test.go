package opcodes

import (
	"math/big"
	expr "nr-groth16/acir/expression"
	shr "nr-groth16/acir/shared"
	"nr-groth16/bn254"
	"testing"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/backend/witness"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/consensys/gnark/frontend/schema"
)

func TestExpressionExecution(t *testing.T) {
	builder, err := r1cs.NewBuilder(ecc.BN254.ScalarField(), frontend.CompileConfig{
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

	zero := bn254.Zero()
	one := bn254.One()
	exprOpCode := Opcode[*bn254.BN254Field]{
		Kind: ACIROpcodeAssertZero,
		Expression: &expr.Expression[*bn254.BN254Field]{
			MulTerms: []expr.MulTerm[*bn254.BN254Field]{
				{
					Term:         &one,
					WitnessLeft:  1,
					WitnessRight: 2,
				},
			},
			LinearCombinations: []expr.LinearCombination[*bn254.BN254Field]{
				{
					Term:    &one,
					Witness: 3,
				},
			},
			Constant: &zero,
		},
	}

	builder.AssertIsEqual(exprOpCode.Expression.Calculate(builder, witnessMap), 0)

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
