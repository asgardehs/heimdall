package heimdall

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// ValidateValue checks that value is compatible with the schema's declared type.
func ValidateValue(schema ConfigSchema, value string) error {
	switch schema.Type {
	case "string", "path", "secret":
		// Any string is valid. "secret" is a display hint, not a storage constraint.
		return nil
	case "number":
		if _, err := strconv.ParseFloat(value, 64); err != nil {
			return &ConfigError{
				Code:    ErrConfigValidation,
				Message: fmt.Sprintf("expected number, got %q", value),
			}
		}
	case "boolean":
		if value != "true" && value != "false" {
			return &ConfigError{
				Code:    ErrConfigValidation,
				Message: fmt.Sprintf("expected true or false, got %q", value),
			}
		}
	case "array":
		var arr []any
		if err := json.Unmarshal([]byte(value), &arr); err != nil {
			return &ConfigError{
				Code:    ErrConfigValidation,
				Message: fmt.Sprintf("expected JSON array, got %q", value),
			}
		}
	case "enum":
		if value == "" && !schema.Required {
			return nil
		}
		for _, c := range schema.Choices {
			if value == c {
				return nil
			}
		}
		return &ConfigError{
			Code:    ErrConfigValidation,
			Message: fmt.Sprintf("expected one of %v, got %q", schema.Choices, value),
		}
	default:
		return nil
	}
	return nil
}
