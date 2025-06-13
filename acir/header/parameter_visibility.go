package header

import (
	"encoding/json"
	"fmt"
)

type ACIRParameterVisibility uint32

const (
	ACIRParameterVisibilityPublic ACIRParameterVisibility = iota
	ACIRParameterVisibilityPrivate
	ACIRParameterVisibilityDatabus
)

func (v *ACIRParameterVisibility) UnmarshalJSON(data []byte) error {
	var visibilityStr string
	if err := json.Unmarshal(data, &visibilityStr); err != nil {
		return err
	}

	switch visibilityStr {
	case "public":
		*v = ACIRParameterVisibilityPublic
	case "private":
		*v = ACIRParameterVisibilityPrivate
	case "databus":
		*v = ACIRParameterVisibilityDatabus
	default:
		return fmt.Errorf("unknown ACIR parameter visibility: %s", visibilityStr)
	}
	return nil
}
