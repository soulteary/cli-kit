package validator

import (
	"fmt"
	"net"
	"strings"
)

// ErrInvalidHostPort is returned when a host:port format is invalid
var ErrInvalidHostPort = fmt.Errorf("invalid host:port format")

// ValidateHostPort validates and parses a host:port address string
//
// Parameters:
//   - addr: The host:port address string (e.g., "localhost:6379", "192.168.1.1:8080")
//
// Returns:
//   - host: The host part of the address
//   - port: The port number
//   - error: Returns error if format is invalid or port is out of range
func ValidateHostPort(addr string) (host string, port int, err error) {
	if addr == "" {
		return "", 0, fmt.Errorf("%w: address cannot be empty", ErrInvalidHostPort)
	}

	host, portStr, err := net.SplitHostPort(addr)
	if err != nil {
		return "", 0, fmt.Errorf("%w: %w", ErrInvalidHostPort, err)
	}

	if host == "" {
		return "", 0, fmt.Errorf("%w: host cannot be empty", ErrInvalidHostPort)
	}

	port, err = ValidatePortString(portStr)
	if err != nil {
		return "", 0, fmt.Errorf("%w: %w", ErrInvalidHostPort, err)
	}

	return host, port, nil
}

// ParseHostPort is an alias for ValidateHostPort for consistency with standard library naming
func ParseHostPort(addr string) (host string, port int, err error) {
	return ValidateHostPort(addr)
}

// ValidateHostPortWithDefaults validates host:port and allows default host/port
//
// Parameters:
//   - addr: The host:port address string
//   - defaultHost: Default host if not specified (e.g., "localhost")
//   - defaultPort: Default port if not specified (must be valid port)
//
// Returns:
//   - host: The host part (or defaultHost)
//   - port: The port number (or defaultPort)
//   - error: Returns error if format is invalid
func ValidateHostPortWithDefaults(addr string, defaultHost string, defaultPort int) (host string, port int, err error) {
	if addr == "" {
		if err := ValidatePort(defaultPort); err != nil {
			return "", 0, fmt.Errorf("invalid default port: %w", err)
		}
		return defaultHost, defaultPort, nil
	}

	// Check if addr contains a colon
	if !strings.Contains(addr, ":") {
		// Only host specified, use default port
		if err := ValidatePort(defaultPort); err != nil {
			return "", 0, fmt.Errorf("invalid default port: %w", err)
		}
		return addr, defaultPort, nil
	}

	// Parse as host:port
	return ValidateHostPort(addr)
}
