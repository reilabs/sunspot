package header

import (
	"fmt"
)

type ACIRABI struct {
	Parameters []ACIRParameter          `json:"parameters"`
	ReturnType *ACIRReturnType          `json:"return_type"`
	ErrorTypes map[string]ACIRErrorType `json:"-"`
}

type ParamInfo struct {
	visibility ACIRParameterVisibility
	Name       string
}

// Params flattens all ABI parameters into a list of inputs.
//
// The ACIR circuit representation can express complex parameter types
// (such as arrays, structs, and tuples). However, Groth16/Gnark expects
// a flat set of scalar inputs. This function recursively expands each
// parameter—breaking down complex composite types into individual elements.
func (a *ACIRABI) Params() []ParamInfo {
	var ret []ParamInfo
	for _, param := range a.Parameters {
		ret = append(ret, flattenParam(param.Visibility, param.Name, param.Type)...)
	}
	return ret
}

// flattenParam recursively flattens any ACIR parameter (scalar, array, or struct)
func flattenParam(vis ACIRParameterVisibility, name string, typ ACIRParameterType) []ParamInfo {
	var result []ParamInfo

	switch typ.Kind {
	case ACIRParameterKindArray:
		if typ.ArrayType == nil || typ.Length == nil {
			return []ParamInfo{{visibility: vis, Name: name}}
		}
		for i := 0; i < *typ.Length; i++ {
			elementName := fmt.Sprintf("%s[%d]", name, i)
			result = append(result, flattenParam(vis, elementName, *typ.ArrayType)...)
		}

	case ACIRParameterKindTuple:
		if typ.TupleFields == nil {
			return []ParamInfo{{visibility: vis, Name: name}}
		}
		for index, tupleField := range *typ.TupleFields {
			fieldName := fmt.Sprintf("%s_%d", name, index)
			result = append(result, flattenParam(vis, fieldName, tupleField)...)
		}

	case ACIRParameterKindStruct:
		if typ.Fields == nil {
			return []ParamInfo{{visibility: vis, Name: name}}
		}
		for _, field := range *typ.Fields {
			fieldName := fmt.Sprintf("%s.%s", name, field.Name)
			result = append(result, flattenParam(vis, fieldName, field.Type)...)
		}

	default:
		result = append(result, ParamInfo{
			visibility: vis,
			Name:       name,
		})
	}

	return result
}

func (a *ParamInfo) IsPublic() bool {
	return a.visibility == ACIRParameterVisibilityPublic
}
