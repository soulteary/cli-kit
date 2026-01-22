package validator

import (
	"fmt"
	"strconv"
)

// ErrInvalidPort is returned when a port value is outside the valid range (1-65535)
var ErrInvalidPort = fmt.Errorf("port must be between 1 and 65535")

// ValidatePort validates that a port number is within the valid range (1-65535)
//
// Parameters:
//   - port: The port number to validate
//
// Returns:
//   - error: Returns ErrInvalidPort if port is outside valid range, nil otherwise
func ValidatePort(port int) error {
	if port < 1 || port > 65535 {
		return fmt.Errorf("%w: got %d", ErrInvalidPort, port)
	}
	return nil
}

// ValidatePortString validates a port string and converts it to an integer
//
// Parameters:
//   - portStr: The port string to validate and parse
//
// Returns:
//   - int: The parsed port number
//   - error: Returns error if port string is invalid or out of range
func ValidatePortString(portStr string) (int, error) {
	if portStr == "" {
		return 0, fmt.Errorf("port string cannot be empty")
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return 0, fmt.Errorf("invalid port format: %w", err)
	}

	if err := ValidatePort(port); err != nil {
		return 0, err
	}

	return port, nil
}
