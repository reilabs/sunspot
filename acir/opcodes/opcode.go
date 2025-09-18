package opcodes

import (
	"io"
	"math/big"
	shr "nr-groth16/acir/shared"

	"github.com/consensys/gnark/frontend"
	"github.com/google/btree"
)

type Opcode interface {
	UnmarshalReader(r io.Reader) error
	Equals(other Opcode) bool
	Define(api frontend.Builder, witnesses map[shr.Witness]frontend.Variable) error
	MarshalJSON() ([]byte, error)
	FillWitnessTree(tree *btree.BTree) bool
	CollectConstantsAsWitnesses(start uint32, tree *btree.BTree) bool
	FeedConstantsAsWitnesses() []*big.Int
}
