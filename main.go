package main

import (
	"fmt"
	"math/big"

	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/backend/witness"
	"github.com/consensys/gnark/constraint"
	cs_bn254 "github.com/consensys/gnark/constraint/bn254"
)

// type MyCircuit struct {
// 	X frontend.Variable
// 	Y frontend.Variable
// }

type AssertZeroBlueprint struct {
	PublicInputsIDBase uint32
	LIsPublic          []bool
}

func (bp AssertZeroBlueprint) CalldataSize() int {
	return 2 * len(bp.LIsPublic)
}

func (bp AssertZeroBlueprint) NbConstraints() int {
	return 1
}

func (bp AssertZeroBlueprint) NbOutputs(instr constraint.Instruction) int {
	return 0
}

func (bp AssertZeroBlueprint) UpdateInstructionTree(inst constraint.Instruction, tree constraint.InstructionTree) constraint.Level {
	return 0
}

func (bp AssertZeroBlueprint) CompressR1C(c *constraint.R1C, to *[]uint32) {
	// Compress each term into the instruction's serialized slice
	*to = make([]uint32, bp.CalldataSize())
	for i, term := range c.L {
		(*to)[2*i] = uint32(term.CID)
		if bp.LIsPublic[i] {
			(*to)[2*i+1] = term.VID + bp.PublicInputsIDBase
		} else {
			(*to)[2*i+1] = term.VID
		}
	}
}

func (bp AssertZeroBlueprint) DecompressR1C(c *constraint.R1C, inst constraint.Instruction) {
	payload := inst.Calldata
	if len(payload) < bp.CalldataSize() {
		panic("invalid calldata size for AssertZeroBlueprint")
	}
	fmt.Println("Decompressing R1C with payload: ", payload)
	c.L = []constraint.Term{}
	c.R = []constraint.Term{
		{CID: constraint.CoeffIdOne, VID: 0},
	}
	c.O = []constraint.Term{
		{CID: constraint.CoeffIdZero, VID: 0},
	}

	for i := 0; i < len(payload)/2; i++ {
		c.L = append(c.L, constraint.Term{
			CID: payload[2*i],
			VID: payload[2*i+1],
		})
	}
}

// This is what gnark uses to build the constraint system
// func (c *MyCircuit) Define(api frontend.API) error {
// 	sum := api.Add(c.X, c.Y)
// 	api.AssertIsEqual(sum, 0)
// 	return nil
// }

func main() {
	builder := cs_bn254.NewR1CS(100)
	_ = builder.AddPublicVariable("ONE")
	xid := builder.AddSecretVariable("x")
	yid := builder.AddSecretVariable("y")
	zid := builder.AddPublicVariable("z")

	fmt.Printf("xid: %d, yid: %d, zid: %d\n", xid, yid, zid)
	r1c := constraint.R1C{
		L: []constraint.Term{
			{
				CID: constraint.CoeffIdOne,
				VID: uint32(xid),
			},
			{
				CID: constraint.CoeffIdOne,
				VID: uint32(yid),
			},
			{
				CID: constraint.CoeffIdMinusOne,
				VID: uint32(zid),
			},
		},
		R: []constraint.Term{
			{
				CID: constraint.CoeffIdOne,
				VID: 0,
			},
		},
		O: []constraint.Term{
			{
				CID: constraint.CoeffIdZero,
				VID: 0,
			},
		},
	}

	bp := AssertZeroBlueprint{
		PublicInputsIDBase: uint32(builder.GetNbSecretVariables()),
		LIsPublic:          []bool{false, false, true},
	}
	blueprintID := builder.AddBlueprint(bp)

	builder.AddR1C(r1c, blueprintID)

	witness, err := witness.New(fr.Modulus())
	if err != nil {
		panic(err)
	}

	values := make(chan any)

	go func() {
		values <- big.NewInt(123)
		values <- big.NewInt(12)
		values <- big.NewInt(135)

		close(values)
	}()

	fmt.Println("Circuit: ", builder)
	fmt.Println("Public: ", builder.GetNbPublicVariables())
	fmt.Println("Secret: ", builder.GetNbSecretVariables())
	err = witness.Fill(1, 2, values)
	if err != nil {
		panic(err)
	}

	fmt.Println("Witnesses: ", witness)

	pk, vk, err := groth16.Setup(builder)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Primary Key: %v\n", pk)
	fmt.Printf("Verification Key: %v\n", vk)

	proof, err := groth16.Prove(builder, pk, witness)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Proof: %v\n", proof)

	publicWitness, err := witness.Public()
	if err != nil {
		panic(err)
	}

	err = groth16.Verify(proof, vk, publicWitness)
	if err != nil {
		panic("❌ Verification failed!")
	} else {
		fmt.Println("✅ Verification succeeded!")
	}
	// Define the circuit constraints (without assigning values)
	/*var circuit MyCircuit

	ccs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)
	if err != nil {
		panic(err)
	}

	// Setup proving & verifying keys
	pk, vk, err := groth16.Setup(ccs)
	if err != nil {
		panic(err)
	}

	// Create witness with actual values
	assignment := MyCircuit{
		X: -3,
		Y: 3,
	}
	witness, err := frontend.NewWitness(&assignment, ecc.BN254.ScalarField())
	if err != nil {
		panic(err)
	}

	publicWitness, _ := witness.Public()

	// Generate the proof
	proof, err := groth16.Prove(ccs, pk, witness)
	if err != nil {
		panic(err)
	}

	// Verify the proof
	err = groth16.Verify(proof, vk, publicWitness)
	if err != nil {
		panic("❌ Verification failed!")
	} else {
		fmt.Println("✅ Verification succeeded!")
	}*/
}
