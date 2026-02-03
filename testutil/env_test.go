package testutil

import (
	"errors"
	"os"
	"testing"
)

func isErrInvalidEnvKey(err error) bool {
	return errors.Is(err, ErrInvalidEnvKey)
}

func TestEnvManager(t *testing.T) {
	t.Run("Set and Restore", func(t *testing.T) {
		originalValue := os.Getenv("TEST_ENV_MANAGER")
		defer func() {
			if originalValue == "" {
				if err := os.Unsetenv("TEST_ENV_MANAGER"); err != nil {
					t.Logf("Failed to unset env var: %v", err)
				}
			} else {
				if err := os.Setenv("TEST_ENV_MANAGER", originalValue); err != nil {
					t.Logf("Failed to set env var: %v", err)
				}
			}
		}()

		// Clear the variable first
		if err := os.Unsetenv("TEST_ENV_MANAGER"); err != nil {
			t.Logf("Failed to unset env var: %v", err)
		}

		manager := NewEnvManager()
		defer manager.Cleanup()

		// Set a value
		if err := manager.Set("TEST_ENV_MANAGER", "test_value"); err != nil {
			t.Fatalf("Set() error = %v", err)
		}

		if got := os.Getenv("TEST_ENV_MANAGER"); got != "test_value" {
			t.Errorf("Set() failed, got = %q, want %q", got, "test_value")
		}

		// Restore
		if err := manager.Restore(); err != nil {
			t.Fatalf("Restore() error = %v", err)
		}

		if got := os.Getenv("TEST_ENV_MANAGER"); got != originalValue {
			t.Errorf("Restore() failed, got = %q, want %q", got, originalValue)
		}
	})

	t.Run("Unset and Restore", func(t *testing.T) {
		originalValue := os.Getenv("TEST_ENV_MANAGER_UNSET")
		defer func() {
			if originalValue == "" {
				if err := os.Unsetenv("TEST_ENV_MANAGER_UNSET"); err != nil {
					t.Logf("Failed to unset env var: %v", err)
				}
			} else {
				if err := os.Setenv("TEST_ENV_MANAGER_UNSET", originalValue); err != nil {
					t.Logf("Failed to set env var: %v", err)
				}
			}
		}()

		// Set a value first
		if err := os.Setenv("TEST_ENV_MANAGER_UNSET", "original"); err != nil {
			t.Fatalf("Failed to set env var: %v", err)
		}

		manager := NewEnvManager()
		defer manager.Cleanup()

		// Unset
		if err := manager.Unset("TEST_ENV_MANAGER_UNSET"); err != nil {
			t.Fatalf("Unset() error = %v", err)
		}

		if got := os.Getenv("TEST_ENV_MANAGER_UNSET"); got != "" {
			t.Errorf("Unset() failed, got = %q, want empty", got)
		}

		// Restore
		if err := manager.Restore(); err != nil {
			t.Fatalf("Restore() error = %v", err)
		}

		if got := os.Getenv("TEST_ENV_MANAGER_UNSET"); got != "original" {
			t.Errorf("Restore() failed, got = %q, want %q", got, "original")
		}
	})

	t.Run("SetMultiple", func(t *testing.T) {
		manager := NewEnvManager()
		defer manager.Cleanup()

		vars := map[string]string{
			"TEST_VAR1": "value1",
			"TEST_VAR2": "value2",
		}

		if err := manager.SetMultiple(vars); err != nil {
			t.Fatalf("SetMultiple() error = %v", err)
		}

		if got := os.Getenv("TEST_VAR1"); got != "value1" {
			t.Errorf("SetMultiple() failed for TEST_VAR1, got = %q, want %q", got, "value1")
		}
		if got := os.Getenv("TEST_VAR2"); got != "value2" {
			t.Errorf("SetMultiple() failed for TEST_VAR2, got = %q, want %q", got, "value2")
		}
	})

	t.Run("Cleanup", func(t *testing.T) {
		originalValue := os.Getenv("TEST_ENV_MANAGER_CLEANUP")
		defer func() {
			if originalValue == "" {
				if err := os.Unsetenv("TEST_ENV_MANAGER_CLEANUP"); err != nil {
					t.Logf("Failed to unset env var: %v", err)
				}
			} else {
				if err := os.Setenv("TEST_ENV_MANAGER_CLEANUP", originalValue); err != nil {
					t.Logf("Failed to set env var: %v", err)
				}
			}
		}()

		manager := NewEnvManager()
		if err := manager.Set("TEST_ENV_MANAGER_CLEANUP", "test"); err != nil {
			t.Fatalf("Set() error = %v", err)
		}
		manager.Cleanup()

		if got := os.Getenv("TEST_ENV_MANAGER_CLEANUP"); got != originalValue {
			t.Errorf("Cleanup() failed, got = %q, want %q", got, originalValue)
		}
	})

	t.Run("Clear", func(t *testing.T) {
		manager := NewEnvManager()
		defer manager.Cleanup()

		// Set some variables
		if err := manager.Set("TEST_CLEAR_VAR1", "value1"); err != nil {
			t.Fatalf("Set() error = %v", err)
		}
		if err := manager.Set("TEST_CLEAR_VAR2", "value2"); err != nil {
			t.Fatalf("Set() error = %v", err)
		}

		// Clear all managed variables
		if err := manager.Clear(); err != nil {
			t.Fatalf("Clear() error = %v", err)
		}

		// Verify variables are unset
		if got := os.Getenv("TEST_CLEAR_VAR1"); got != "" {
			t.Errorf("Clear() failed for TEST_CLEAR_VAR1, got = %q, want empty", got)
		}
		if got := os.Getenv("TEST_CLEAR_VAR2"); got != "" {
			t.Errorf("Clear() failed for TEST_CLEAR_VAR2, got = %q, want empty", got)
		}
	})

	t.Run("SetMultiple with error", func(t *testing.T) {
		manager := NewEnvManager()
		defer manager.Cleanup()

		// SetMultiple should handle empty map
		if err := manager.SetMultiple(map[string]string{}); err != nil {
			t.Errorf("SetMultiple() with empty map error = %v, want nil", err)
		}
	})

	t.Run("Set with empty key returns error", func(t *testing.T) {
		manager := NewEnvManager()
		defer manager.Cleanup()
		if err := manager.Set("", "value"); err == nil {
			t.Error("Set() with empty key want error, got nil")
		} else if err != ErrInvalidEnvKey && !isErrInvalidEnvKey(err) {
			t.Errorf("Set() with empty key want ErrInvalidEnvKey, got %v", err)
		}
	})

	t.Run("Unset with empty key returns error", func(t *testing.T) {
		manager := NewEnvManager()
		defer manager.Cleanup()
		if err := manager.Unset(""); err == nil {
			t.Error("Unset() with empty key want error, got nil")
		} else if err != ErrInvalidEnvKey && !isErrInvalidEnvKey(err) {
			t.Errorf("Unset() with empty key want ErrInvalidEnvKey, got %v", err)
		}
	})

	t.Run("SetMultiple with empty key returns error", func(t *testing.T) {
		manager := NewEnvManager()
		defer manager.Cleanup()
		if err := manager.SetMultiple(map[string]string{"": "v"}); err == nil {
			t.Error("SetMultiple() with empty key want error, got nil")
		} else if err != ErrInvalidEnvKey && !isErrInvalidEnvKey(err) {
			t.Errorf("SetMultiple() with empty key want ErrInvalidEnvKey, got %v", err)
		}
	})

	t.Run("Set with NUL key returns error", func(t *testing.T) {
		manager := NewEnvManager()
		defer manager.Cleanup()
		if err := manager.Set("\x00", "value"); err == nil {
			t.Error("Set() with NUL key want error, got nil")
		} else if !isErrInvalidEnvKey(err) {
			t.Errorf("Set() with NUL key want ErrInvalidEnvKey, got %v", err)
		}
	})

	t.Run("Clear with no variables", func(t *testing.T) {
		manager := NewEnvManager()
		defer manager.Cleanup()
		if err := manager.Clear(); err != nil {
			t.Errorf("Clear() with no variables = %v, want nil", err)
		}
	})

	t.Run("Restore with multiple variables", func(t *testing.T) {
		originalValue1 := os.Getenv("TEST_RESTORE_VAR1")
		originalValue2 := os.Getenv("TEST_RESTORE_VAR2")
		defer func() {
			if originalValue1 == "" {
				if err := os.Unsetenv("TEST_RESTORE_VAR1"); err != nil {
					t.Logf("Failed to unset env var: %v", err)
				}
			} else {
				if err := os.Setenv("TEST_RESTORE_VAR1", originalValue1); err != nil {
					t.Logf("Failed to set env var: %v", err)
				}
			}
			if originalValue2 == "" {
				if err := os.Unsetenv("TEST_RESTORE_VAR2"); err != nil {
					t.Logf("Failed to unset env var: %v", err)
				}
			} else {
				if err := os.Setenv("TEST_RESTORE_VAR2", originalValue2); err != nil {
					t.Logf("Failed to set env var: %v", err)
				}
			}
		}()

		manager := NewEnvManager()
		defer manager.Cleanup()

		// Set multiple variables
		if err := manager.Set("TEST_RESTORE_VAR1", "new_value1"); err != nil {
			t.Fatalf("Set() error = %v", err)
		}
		if err := manager.Set("TEST_RESTORE_VAR2", "new_value2"); err != nil {
			t.Fatalf("Set() error = %v", err)
		}

		// Restore
		if err := manager.Restore(); err != nil {
			t.Fatalf("Restore() error = %v", err)
		}

		// Verify restored values
		if got := os.Getenv("TEST_RESTORE_VAR1"); got != originalValue1 {
			t.Errorf("Restore() failed for TEST_RESTORE_VAR1, got = %q, want %q", got, originalValue1)
		}
		if got := os.Getenv("TEST_RESTORE_VAR2"); got != originalValue2 {
			t.Errorf("Restore() failed for TEST_RESTORE_VAR2, got = %q, want %q", got, originalValue2)
		}
	})

	t.Run("SetMultiple with error handling", func(t *testing.T) {
		manager := NewEnvManager()
		defer manager.Cleanup()

		// Test SetMultiple with multiple variables, one that might cause issues
		vars := map[string]string{
			"TEST_SETMULTIPLE_VAR1": "value1",
			"TEST_SETMULTIPLE_VAR2": "value2",
			"TEST_SETMULTIPLE_VAR3": "value3",
		}

		if err := manager.SetMultiple(vars); err != nil {
			t.Fatalf("SetMultiple() error = %v", err)
		}

		// Verify all variables are set
		for key, expectedValue := range vars {
			if got := os.Getenv(key); got != expectedValue {
				t.Errorf("SetMultiple() failed for %s, got = %q, want %q", key, got, expectedValue)
			}
		}
	})

	t.Run("Clear with multiple variables", func(t *testing.T) {
		manager := NewEnvManager()
		defer manager.Cleanup()

		// Set multiple variables
		vars := map[string]string{
			"TEST_CLEAR_VAR1": "value1",
			"TEST_CLEAR_VAR2": "value2",
			"TEST_CLEAR_VAR3": "value3",
		}

		if err := manager.SetMultiple(vars); err != nil {
			t.Fatalf("SetMultiple() error = %v", err)
		}

		// Clear all
		if err := manager.Clear(); err != nil {
			t.Fatalf("Clear() error = %v", err)
		}

		// Verify all are unset
		for key := range vars {
			if got := os.Getenv(key); got != "" {
				t.Errorf("Clear() failed for %s, got = %q, want empty", key, got)
			}
		}
	})
}
