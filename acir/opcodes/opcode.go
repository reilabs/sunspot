package opcodes

import (
	"io"
	shr "sunpot/acir/shared"

	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/google/btree"
)

type Opcode[E constraint.Element] interface {
	UnmarshalReader(r io.Reader) error
	Equals(other Opcode[E]) bool
	Define(api frontend.Builder[E], witnesses map[shr.Witness]frontend.Variable) error
	MarshalJSON() ([]byte, error)
	FillWitnessTree(tree *btree.BTree, index uint32) bool
}
