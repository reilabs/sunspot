package header

import (
	"encoding/json"
	"fmt"
)

type ACIRErrorType struct {
	Kind       ACIRErrorKind
	String     *string
	Length     *int
	ItemTypes  *[]ACIRParameterType
	CustomType *ACIRParameterType
}

func (t *ACIRErrorType) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if raw["kind"] == nil {
		return fmt.Errorf("missing kind field in ACIRErrorType")
	}

	var Kind ACIRErrorKind
	if err := json.Unmarshal(raw["kind"], &Kind); err != nil {
		return err
	}
	t.Kind = Kind

	switch Kind {
	case ACIRErrorKindString:
		if raw["string"] == nil {
			return fmt.Errorf("missing string field for string error type")
		}
		var str string
		if err := json.Unmarshal(raw["string"], &str); err != nil {
			return err
		}
		t.String = &str
	case ACIRErrorKindFmtString:
		if raw["length"] == nil {
			return fmt.Errorf("missing length field for fmtstring error type")
		}
		var length int
		if err := json.Unmarshal(raw["length"], &length); err != nil {
			return err
		}
		t.Length = &length

		if raw["item_types"] == nil {
			return fmt.Errorf("missing item_types field for fmtstring error type")
		}

		var itemTypes []ACIRParameterType
		if err := json.Unmarshal(raw["item_types"], &itemTypes); err != nil {
			return err
		}
		t.ItemTypes = &itemTypes
	case ACIRErrorKindCustom:
		if raw["type"] == nil {
			return fmt.Errorf("missing type field for custom error type")
		}
		var customType ACIRParameterType
		if err := json.Unmarshal(raw["type"], &customType); err != nil {
			return err
		}
		t.CustomType = &customType

	default:
		return fmt.Errorf("unknown ACIR error kind: %d", Kind)
	}

	return nil
}
