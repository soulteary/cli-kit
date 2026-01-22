package testutil

import (
	"fmt"
	"os"
)

// EnvManager manages environment variables for testing
// It saves original values and can restore them after tests
type EnvManager struct {
	original map[string]string
}

// NewEnvManager creates a new environment variable manager
func NewEnvManager() *EnvManager {
	return &EnvManager{
		original: make(map[string]string),
	}
}

// Set sets an environment variable and saves the original value
func (m *EnvManager) Set(key, value string) error {
	// Save original value if not already saved
	if _, exists := m.original[key]; !exists {
		m.original[key] = os.Getenv(key)
	}
	return os.Setenv(key, value)
}

// Unset unsets an environment variable and saves the original value
func (m *EnvManager) Unset(key string) error {
	// Save original value if not already saved
	if _, exists := m.original[key]; !exists {
		m.original[key] = os.Getenv(key)
	}
	return os.Unsetenv(key)
}

// Restore restores all environment variables to their original values
func (m *EnvManager) Restore() error {
	for key, value := range m.original {
		if value == "" {
			if err := os.Unsetenv(key); err != nil {
				return fmt.Errorf("failed to unset %q: %w", key, err)
			}
		} else {
			if err := os.Setenv(key, value); err != nil {
				return fmt.Errorf("failed to set %q: %w", key, err)
			}
		}
	}
	return nil
}

// Cleanup is a convenience method for use with defer
// It restores all environment variables to their original values
func (m *EnvManager) Cleanup() {
	_ = m.Restore()
}

// SetMultiple sets multiple environment variables at once
func (m *EnvManager) SetMultiple(vars map[string]string) error {
	for key, value := range vars {
		if err := m.Set(key, value); err != nil {
			return err
		}
	}
	return nil
}

// Clear clears all managed environment variables (unsets them)
func (m *EnvManager) Clear() error {
	for key := range m.original {
		if err := m.Unset(key); err != nil {
			return err
		}
	}
	return nil
}
