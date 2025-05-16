package acir_decoder

import (
	shr "nr-groth16/acir_decoder/shared"
)

type Program[T shr.ACIRField] struct {
	Functions              []Circuit[T]
	UnconstrainedFunctions []BrilligBytecode[T]
}

type BrilligBytecode[T shr.ACIRField] struct {
}
