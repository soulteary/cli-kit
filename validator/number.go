package validator

import "fmt"

// ErrNotPositive is returned when a value is not positive (> 0)
var ErrNotPositive = fmt.Errorf("value must be positive (> 0)")

// ErrNegative is returned when a value is negative (< 0)
var ErrNegative = fmt.Errorf("value must be non-negative (>= 0)")

// ValidatePositive validates that an integer is positive (> 0)
//
// Parameters:
//   - value: The integer value to validate
//
// Returns:
//   - error: Returns ErrNotPositive if value <= 0, nil otherwise
func ValidatePositive(value int) error {
	if value <= 0 {
		return fmt.Errorf("%w: got %d", ErrNotPositive, value)
	}
	return nil
}

// ValidatePositiveInt64 validates that an int64 is positive (> 0)
//
// Parameters:
//   - value: The int64 value to validate
//
// Returns:
//   - error: Returns ErrNotPositive if value <= 0, nil otherwise
func ValidatePositiveInt64(value int64) error {
	if value <= 0 {
		return fmt.Errorf("%w: got %d", ErrNotPositive, value)
	}
	return nil
}

// ValidateNonNegative validates that an integer is non-negative (>= 0)
//
// Parameters:
//   - value: The integer value to validate
//
// Returns:
//   - error: Returns ErrNegative if value < 0, nil otherwise
func ValidateNonNegative(value int) error {
	if value < 0 {
		return fmt.Errorf("%w: got %d", ErrNegative, value)
	}
	return nil
}

// ValidateNonNegativeInt64 validates that an int64 is non-negative (>= 0)
//
// Parameters:
//   - value: The int64 value to validate
//
// Returns:
//   - error: Returns ErrNegative if value < 0, nil otherwise
func ValidateNonNegativeInt64(value int64) error {
	if value < 0 {
		return fmt.Errorf("%w: got %d", ErrNegative, value)
	}
	return nil
}

// ValidateInRange validates that an integer is within a specified range [min, max]
//
// Parameters:
//   - value: The integer value to validate
//   - min: The minimum allowed value (inclusive)
//   - max: The maximum allowed value (inclusive)
//
// Returns:
//   - error: Returns error if value is outside the range, nil otherwise
func ValidateInRange(value, min, max int) error {
	if value < min || value > max {
		return fmt.Errorf("value must be between %d and %d: got %d", min, max, value)
	}
	return nil
}

// ValidateInRangeInt64 validates that an int64 is within a specified range [min, max]
//
// Parameters:
//   - value: The int64 value to validate
//   - min: The minimum allowed value (inclusive)
//   - max: The maximum allowed value (inclusive)
//
// Returns:
//   - error: Returns error if value is outside the range, nil otherwise
func ValidateInRangeInt64(value, min, max int64) error {
	if value < min || value > max {
		return fmt.Errorf("value must be between %d and %d: got %d", min, max, value)
	}
	return nil
}
