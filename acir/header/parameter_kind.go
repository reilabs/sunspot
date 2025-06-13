package header

import (
	"encoding/json"
	"fmt"
)

type ACIRParameterKind uint32

const (
	ACIRParameterKindField ACIRParameterKind = iota
	ACIRParameterKindBoolean
	ACIRParameterKindInteger
	ACIRParameterKindFloat
	ACIRParameterKindString
	ACIRParameterKindArray
	ACIRParameterKindTuple
	ACIRParameterKindStruct
)

func (k *ACIRParameterKind) UnmarshalJSON(data []byte) error {
	var kindStr string
	if err := json.Unmarshal(data, &kindStr); err != nil {
		return err
	}

	switch kindStr {
	case "field":
		*k = ACIRParameterKindField
	case "boolean":
		*k = ACIRParameterKindBoolean
	case "integer":
		*k = ACIRParameterKindInteger
	case "float":
		*k = ACIRParameterKindFloat
	case "string":
		*k = ACIRParameterKindString
	case "array":
		*k = ACIRParameterKindArray
	case "tuple":
		*k = ACIRParameterKindTuple
	case "struct":
		*k = ACIRParameterKindStruct
	default:
		return fmt.Errorf("unknown ACIR parameter kind: %s", kindStr)
	}
	return nil
}
