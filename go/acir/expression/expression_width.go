package expression

import (
	"encoding/json"
	"fmt"
)

type ExpressionWidth struct {
	Kind  ExpressionWidthKind
	Width *uint64
}

type ExpressionWidthKind uint32

const (
	ACIRExpressionWidthUnbounded ExpressionWidthKind = iota
	ACIRExpressionWidthBounded
)

func (e *ExpressionWidth) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	for key, value := range raw {
		switch key {
		case "Bounded":
			e.Kind = ACIRExpressionWidthBounded
			var bounded struct {
				Width uint64 `json:"width"`
			}
			if err := json.Unmarshal(value, &bounded); err != nil {
				return err
			}
			e.Width = &bounded.Width

		case "Unbounded":
			e.Kind = ACIRExpressionWidthUnbounded
			e.Width = nil

		default:
			return fmt.Errorf("unknown ExpressionWidth variant: %s", key)
		}
		return nil
	}
	return fmt.Errorf("invalid ExpressionWidth: %s", string(data))
}

func (e *ExpressionWidth) Equals(other *ExpressionWidth) bool {
	if e.Kind != other.Kind {
		return false
	}

	if e.Width == nil && other.Width == nil {
		return true
	}

	if e.Width == nil || other.Width == nil {
		return false
	}

	return *e.Width == *other.Width
}

func (e *ExpressionWidth) MarshalJSON() ([]byte, error) {
	fieldsMap := make(map[string]interface{})
	switch e.Kind {
	case ACIRExpressionWidthUnbounded:
		fieldsMap["Unbounded"] = nil
	case ACIRExpressionWidthBounded:
		fieldsMap["Bounded"] = *e.Width
	}
	return json.Marshal(fieldsMap)
}
