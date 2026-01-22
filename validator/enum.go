package validator

import (
	"fmt"
	"strings"
)

// ErrInvalidEnumValue is returned when a value is not in the allowed enum list
var ErrInvalidEnumValue = fmt.Errorf("value is not in allowed enum values")

// ValidateEnum validates that a value is in the allowed enum list
//
// Parameters:
//   - value: The value to validate
//   - allowedValues: List of allowed values
//   - caseSensitive: Whether comparison should be case-sensitive
//
// Returns:
//   - error: Returns ErrInvalidEnumValue if value is not in allowed list, nil otherwise
func ValidateEnum(value string, allowedValues []string, caseSensitive bool) error {
	if len(allowedValues) == 0 {
		return fmt.Errorf("allowed values list cannot be empty")
	}

	if value == "" {
		return fmt.Errorf("value cannot be empty")
	}

	for _, allowed := range allowedValues {
		if caseSensitive {
			if value == allowed {
				return nil
			}
		} else {
			if strings.EqualFold(value, allowed) {
				return nil
			}
		}
	}

	return fmt.Errorf("%w: %q, allowed values: %v", ErrInvalidEnumValue, value, allowedValues)
}

// ValidateEnumCaseInsensitive is a convenience function for case-insensitive enum validation
func ValidateEnumCaseInsensitive(value string, allowedValues []string) error {
	return ValidateEnum(value, allowedValues, false)
}

// ValidateEnumCaseSensitive is a convenience function for case-sensitive enum validation
func ValidateEnumCaseSensitive(value string, allowedValues []string) error {
	return ValidateEnum(value, allowedValues, true)
}
