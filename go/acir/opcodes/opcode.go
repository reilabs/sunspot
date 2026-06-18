package opcodes

import (
	"sunspot/go/acir/msgpackutil"
	shr "sunspot/go/acir/shared"

	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
)

type Opcode[E constraint.Element] interface {
	msgpackutil.EnumVariant
	UnmarshalReader(r *msgpackutil.Reader) error
	Equals(other Opcode[E]) bool
	Define(api frontend.Builder[E], witnesses map[shr.Witness]frontend.Variable) error
	MarshalJSON() ([]byte, error)
}
