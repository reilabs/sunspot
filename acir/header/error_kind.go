package header

import (
	"encoding/json"
	"fmt"
)

type ACIRErrorKind int

const (
	ACIRErrorKindString ACIRErrorKind = iota
	ACIRErrorKindFmtString
	ACIRErrorKindCustom
)

func (k *ACIRErrorKind) UnmarshalJSON(data []byte) error {
	var kindStr string
	if err := json.Unmarshal(data, &kindStr); err != nil {
		return err
	}

	switch kindStr {
	case "string":
		*k = ACIRErrorKindString
	case "fmtstring":
		*k = ACIRErrorKindFmtString
	case "custom":
		*k = ACIRErrorKindCustom
	default:
		return fmt.Errorf("unknown ACIR error kind: %s", kindStr)
	}
	return nil
}
