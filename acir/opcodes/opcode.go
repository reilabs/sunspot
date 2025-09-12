package opcodes

import (
	"io"
	"math/big"
	shr "nr-groth16/acir/shared"

	"github.com/consensys/gnark/frontend"
	"github.com/google/btree"
)

// type Opcode[T shr.ACIRField] struct {
// 	Kind             OpcodeKind
// 	Expression       *exp.Expression[T]
// 	BlackBoxFuncCall *bbf.BlackBoxFuncCall[T]
// 	MemoryOp         *MemoryOp[T]
// 	MemoryInit       *MemoryInit[T]
// 	BrilligCall      *brl.BrilligCall[T]
// 	Call             *Call[T]
// }

type Opcode interface {
	UnmarshalReader(r io.Reader) error
	Equals(other Opcode) bool
	Define(api frontend.API, witnesses map[shr.Witness]frontend.Variable) error
	MarshalJSON() ([]byte, error)
	FillWitnessTree(tree *btree.BTree) bool
	CollectConstantsAsWitnesses(start uint32, tree *btree.BTree) bool
	FeedConstantsAsWitnesses() []*big.Int
}
type OpcodeKind uint32

const (
	ACIROpcodeAssertZero OpcodeKind = iota
	ACIROpcodeBlackBoxFuncCall
	ACIROpcodeMemoryOp
	ACIROpcodeMemoryInit
	ACIROpcodeBrilligCall
	ACIROpcodeCall
)
