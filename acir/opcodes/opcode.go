package opcodes

import (
	"io"
	"math/big"
	shr "nr-groth16/acir/shared"

	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/google/btree"
)

type Opcode[E constraint.Element] interface {
	UnmarshalReader(r io.Reader) error
	Equals(other Opcode[E]) bool
	Define(api frontend.API, witnesses map[shr.Witness]frontend.Variable) error
	MarshalJSON() ([]byte, error)
	FillWitnessTree(tree *btree.BTree) bool
	CollectConstantsAsWitnesses(start uint32, tree *btree.BTree) bool
	FeedConstantsAsWitnesses() []*big.Int
}
