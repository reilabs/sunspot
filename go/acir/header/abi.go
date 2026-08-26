package header

import "fmt"

type ACIRABI struct {
	Parameters []ACIRParameter          `json:"parameters"`
	ReturnType *ACIRReturnType          `json:"return_type"`
	ErrorTypes map[string]ACIRErrorType `json:"-"`
}

type ParamInfo struct {
	Visibility ACIRParameterVisibility
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

// Returns flattens the ABI's return type into a list of scalar outputs.
func (a *ACIRABI) Returns() []ParamInfo {
	if a.ReturnType == nil {
		return nil
	}
	return flattenParam(a.ReturnType.Visibility, "return", a.ReturnType.Type)
}

// flattenParam recursively flattens any ACIR parameter (scalar, array, or struct)
func flattenParam(vis ACIRParameterVisibility, name string, typ ACIRParameterType) []ParamInfo {
	var result []ParamInfo

	switch typ.Kind {
	case ACIRParameterKindString:
		if typ.Length == nil {
			return []ParamInfo{{Visibility: vis, Name: name}}
		}
		for i := 0; i < *typ.Length; i++ {
			result = append(result, ParamInfo{Visibility: vis, Name: fmt.Sprintf("%s[%d]", name, i)})
		}

	case ACIRParameterKindArray:
		if typ.ArrayType == nil || typ.Length == nil {
			return []ParamInfo{{Visibility: vis, Name: name}}
		}
		for i := 0; i < *typ.Length; i++ {
			elementName := fmt.Sprintf("%s[%d]", name, i)
			result = append(result, flattenParam(vis, elementName, *typ.ArrayType)...)
		}

	case ACIRParameterKindTuple:
		if typ.TupleFields == nil {
			return []ParamInfo{{Visibility: vis, Name: name}}
		}
		for index, tupleField := range *typ.TupleFields {
			fieldName := fmt.Sprintf("%s_%d", name, index)
			result = append(result, flattenParam(vis, fieldName, tupleField)...)
		}

	case ACIRParameterKindStruct:
		if typ.Fields == nil {
			return []ParamInfo{{Visibility: vis, Name: name}}
		}
		for _, field := range *typ.Fields {
			fieldName := fmt.Sprintf("%s.%s", name, field.Name)
			result = append(result, flattenParam(vis, fieldName, field.Type)...)
		}

	default:
		result = append(result, ParamInfo{
			Visibility: vis,
			Name:       name,
		})
	}

	return result
}
