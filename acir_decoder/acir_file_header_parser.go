package acir_decoder

import (
	"encoding/json"
	"fmt"
)

type ACIRFile struct {
	NoirVersion  string                  `json:"noir_version"`
	Hash         uint64                  `json:"hash"`
	ABI          ACIRABI                 `json:"abi"`
	Bytecode     string                  `json:"bytecode"`
	DebugSymbols string                  `json:"debug_symbols"`
	FileMap      map[string]ACIRFileData `json:"file_map"`
	Names        []string                `json:"names"`
	BrilligNames []string                `json:"brillig_names"`
}

type ACIRABI struct {
	Parameters []ACIRParameter          `json:"parameters"`
	ReturnType *ACIRReturnType          `json:"return_type"`
	ErrorTypes map[string]ACIRErrorType `json:"error_types"`
}

type ACIRParameter struct {
	Name       string                  `json:"name"`
	Type       ACIRParameterType       `json:"type"`
	Visibility ACIRParameterVisibility `json:"visibility"`
}

type ACIRParameterVisibility int

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

type ACIRParameterKind int

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

type ACIRParameterTypeStructField struct {
	Name string            `json:"name"`
	Type ACIRParameterType `json:"type"`
}

type ACIRReturnType struct {
	Type       ACIRParameterType       `json:"type"`
	Visibility ACIRParameterVisibility `json:"visibility"`
}

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

type ACIRFileData struct {
	Source string `json:"source"`
	Path   string `json:"path"`
}
