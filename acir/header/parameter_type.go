package header

import (
	"encoding/json"
	"fmt"
)

type ACIRParameterType struct {
	Kind        ACIRParameterKind
	Length      *int
	Sign        *ACIRParameterSign
	Width       *int
	ArrayType   *ACIRParameterType
	TupleFields *[]ACIRParameterType
	Path        *string
	Fields      *[]ACIRParameterTypeStructField
}

func (t *ACIRParameterType) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if raw["kind"] == nil {
		return fmt.Errorf("missing kind field in ACIRParameterType")
	}

	var Kind ACIRParameterKind
	if err := json.Unmarshal(raw["kind"], &Kind); err != nil {
		return err
	}
	t.Kind = Kind

	switch Kind {
	case ACIRParameterKindString:
		if raw["length"] == nil {
			return fmt.Errorf("missing length field for string type")
		}
		var length int
		if err := json.Unmarshal(raw["length"], &length); err != nil {
			return err
		}
		t.Length = &length
	case ACIRParameterKindInteger:
		if raw["width"] == nil {
			return fmt.Errorf("missing width field for integer type")
		}
		var width int
		if err := json.Unmarshal(raw["width"], &width); err != nil {
			return err
		}
		t.Width = &width

		if raw["sign"] == nil {
			return fmt.Errorf("missing sign field for integer type")
		}

		var sign ACIRParameterSign
		if err := json.Unmarshal(raw["sign"], &sign); err != nil {
			return err
		}

		t.Sign = &sign
	case ACIRParameterKindArray:
		if raw["type"] == nil {
			return fmt.Errorf("missing type field for array type")
		}

		var arrayType ACIRParameterType
		if err := json.Unmarshal(raw["type"], &arrayType); err != nil {
			return err
		}
		t.ArrayType = &arrayType

		if raw["length"] == nil {
			return fmt.Errorf("missing length field for array type")
		}

		var length int
		if err := json.Unmarshal(raw["length"], &length); err != nil {
			return err
		}
		t.Length = &length
	case ACIRParameterKindTuple:
		if raw["fields"] == nil {
			return fmt.Errorf("missing fields field for tuple type")
		}
		var fields []ACIRParameterType
		if err := json.Unmarshal(raw["fields"], &fields); err != nil {
			return err
		}
		t.TupleFields = &fields
	case ACIRParameterKindStruct:
		if raw["path"] == nil {
			return fmt.Errorf("missing path field for struct type")
		}
		var path string
		if err := json.Unmarshal(raw["path"], &path); err != nil {
			return err
		}
		t.Path = &path

		if raw["fields"] == nil {
			return fmt.Errorf("missing fields field for struct type")
		}
		var fields []ACIRParameterTypeStructField
		if err := json.Unmarshal(raw["fields"], &fields); err != nil {
			return err
		}
		t.Fields = &fields
	}

	return nil
}
