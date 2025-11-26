package header

type ACIRParameter struct {
	Name       string                  `json:"name"`
	Type       ACIRParameterType       `json:"type"`
	Visibility ACIRParameterVisibility `json:"visibility"`
}
