package header

import (
	"encoding/json"
	"fmt"
)

type ACIRParameterSign int

const (
	ACIRParameterSignUnsigned ACIRParameterSign = iota
	ACIRParameterSignSigned
)

func (s *ACIRParameterSign) UnmarshalJSON(data []byte) error {
	var signStr string
	if err := json.Unmarshal(data, &signStr); err != nil {
		return err
	}

	switch signStr {
	case "unsigned":
		*s = ACIRParameterSignUnsigned
	case "signed":
		*s = ACIRParameterSignSigned
	default:
		return fmt.Errorf("unknown ACIR parameter sign: %s", signStr)
	}
	return nil
}
