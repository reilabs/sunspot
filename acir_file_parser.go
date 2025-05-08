package main

import (
	"encoding/json"
	"fmt"
)

type ACIRFile struct {
	NoirVersion  string
	Hash         uint64
	ABI          ACIRABI
	Bytecode     string
	DebugSymbols string
	FileMap      map[string]ACIRFileData
	Names        []string
	BrilligNames []string
}

func (f *ACIRFile) UnmarshalJSON(data []byte) error {
	// Implement the JSON unmarshalling logic here
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	for k, v := range raw {
		fmt.Printf("Key: %s\n", k)
		switch k {
		case "noir_version":
			if err := json.Unmarshal(v, &f.NoirVersion); err != nil {
				return fmt.Errorf("failed to unmarshal noir_version: %w", err)
			}
		case "hash":
			if err := json.Unmarshal(v, &f.Hash); err != nil {
				return fmt.Errorf("failed to unmarshal hash: %w", err)
			}
		case "abi":

		case "bytecode":
			if err := json.Unmarshal(v, &f.Bytecode); err != nil {
				return fmt.Errorf("failed to unmarshal bytecode: %w", err)
			}
		case "debug_symbols":
			if err := json.Unmarshal(v, &f.DebugSymbols); err != nil {
				return fmt.Errorf("failed to unmarshal debug_symbols: %w", err)
			}
		case "file_map":
			if err := json.Unmarshal(v, &f.FileMap); err != nil {
				return fmt.Errorf("failed to unmarshal file_map: %w", err)
			}
		case "names":
			if err := json.Unmarshal(v, &f.Names); err != nil {
				return fmt.Errorf("failed to unmarshal names: %w", err)
			}
		case "brillig_names":
			if err := json.Unmarshal(v, &f.BrilligNames); err != nil {
				return fmt.Errorf("failed to unmarshal brillig_names: %w", err)
			}
		default:
			return fmt.Errorf("unknown key in JSON: %s", k)
		}
	}

	return nil
}

type ACIRABI struct {
	Parameters []ACIRParameter
	ReturnType *ACIRReturnType
	ErrorTypes map[string]ACIRErrorType
}

type ACIRParameter struct {
	Name       string
	Type       ACIRParameterType
	Visibility ACIRParameterVisibility
}

type ACIRParameterVisibility int

const (
	ACIRParameterVisibilityPublic ACIRParameterVisibility = iota
	ACIRParameterVisibilityPrivate
	ACIRParameterVisibilityDatabus
)

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

type ACIRParameterSign int

const (
	ACIRParameterSignUnsigned ACIRParameterSign = iota
	ACIRParameterSignSigned
)

type ACIRParameterType struct {
	Kind      ACIRParameterKind
	Length    int
	Sign      ACIRParameterSign
	Width     int
	ArrayType *ACIRParameterType
	Path      *string
	Fields    *[]ACIRParameterTypeStructField
}

type ACIRParameterTypeStructField struct {
	Name string
	Type ACIRParameterType
}

type ACIRReturnType struct {
	Type       ACIRParameterType
	Visibility ACIRParameterVisibility
}

type ACIRErrorKind int

const (
	ACIRErrorKindString ACIRErrorKind = iota
	ACIRErrorKindFmtString
	ACIRErrorKindCustom
)

type ACIRErrorType struct {
	Kind       ACIRErrorKind
	String     *string
	Length     *int
	ItemTypes  *[]ACIRParameterType
	CustomType *ACIRParameterType
}

type ACIRFileData struct {
	Source string `json:"source"`
	Path   string `json:"path"`
}
