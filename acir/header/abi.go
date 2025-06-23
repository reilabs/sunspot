package header

type ACIRABI struct {
	Parameters []ACIRParameter          `json:"parameters"`
	ReturnType *ACIRReturnType          `json:"return_type"`
	ErrorTypes map[string]ACIRErrorType `json:"-"`
}
