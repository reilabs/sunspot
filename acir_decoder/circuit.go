package acir_decoder

import (
	ops "nr-groth16/acir_decoder/opcodes"
	shr "nr-groth16/acir_decoder/shared"

	"github.com/google/btree"
)

type Circuit[T shr.ACIRField] struct {
	CurrentWitnessIndex uint32
	Opcodes             []ops.Opcode[T]
	ExpressionWidth     ExpressionWidth
	PrivateParameters   []btree.BTree
	PublicParameters    []PublicInputs
	ReturnValues        []PublicInputs
	AssertMessages      map[OpcodeLocation]AssertionPayload[T]
}
