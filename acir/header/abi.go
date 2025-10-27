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
