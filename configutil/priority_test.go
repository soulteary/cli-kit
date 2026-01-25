package configutil

import (
	"flag"
	"fmt"
	"os"
	"testing"
	"time"
)

// setEnv sets an environment variable and panics on error
func setEnv(t *testing.T, key, value string) {
	t.Helper()
	if err := os.Setenv(key, value); err != nil {
		t.Fatalf("os.Setenv(%q, %q) failed: %v", key, value, err)
	}
}

// unsetEnv unsets an environment variable and panics on error
func unsetEnv(t *testing.T, key string) {
	t.Helper()
	if err := os.Unsetenv(key); err != nil {
		t.Fatalf("os.Unsetenv(%q) failed: %v", key, err)
	}
}

func TestResolveString(t *testing.T) {
	t.Run("CLI flag has highest priority", func(t *testing.T) {
		fs := flag.NewFlagSet("test", flag.ContinueOnError)
		fs.String("test-flag", "", "test flag")
		setEnv(t, "TEST_ENV", "env_value")
		defer unsetEnv(t, "TEST_ENV")

		if err := fs.Parse([]string{"--test-flag", "cli_value"}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}

		got := ResolveString(fs, "test-flag", "TEST_ENV", "default", false)
		if got != "cli_value" {
			t.Errorf("ResolveString() = %v, want %v", got, "cli_value")
		}
	})

	t.Run("Environment variable has priority over default", func(t *testing.T) {
		fs := flag.NewFlagSet("test", flag.ContinueOnError)
		fs.String("test-flag", "", "test flag")
		setEnv(t, "TEST_ENV", "env_value")
		defer unsetEnv(t, "TEST_ENV")

		if err := fs.Parse([]string{}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}

		got := ResolveString(fs, "test-flag", "TEST_ENV", "default", false)
		if got != "env_value" {
			t.Errorf("ResolveString() = %v, want %v", got, "env_value")
		}
	})

	t.Run("Default value used when neither CLI nor ENV set", func(t *testing.T) {
		fs := flag.NewFlagSet("test", flag.ContinueOnError)
		fs.String("test-flag", "", "test flag")
		if err := fs.Parse([]string{}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}

		got := ResolveString(fs, "test-flag", "TEST_ENV", "default", false)
		if got != "default" {
			t.Errorf("ResolveString() = %v, want %v", got, "default")
		}
	})

	t.Run("Trimmed option works", func(t *testing.T) {
		fs := flag.NewFlagSet("test", flag.ContinueOnError)
		fs.String("test-flag", "", "test flag")
		setEnv(t, "TEST_ENV", "  env_value  ")
		defer unsetEnv(t, "TEST_ENV")

		if err := fs.Parse([]string{}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}

		got := ResolveString(fs, "test-flag", "TEST_ENV", "default", true)
		if got != "env_value" {
			t.Errorf("ResolveString() with trimmed = %v, want %v", got, "env_value")
		}
	})

	t.Run("Empty ENV value falls back to default", func(t *testing.T) {
		fs := flag.NewFlagSet("test", flag.ContinueOnError)
		fs.String("test-flag", "", "test flag")
		setEnv(t, "TEST_ENV", "")
		defer unsetEnv(t, "TEST_ENV")

		if err := fs.Parse([]string{}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}

		got := ResolveString(fs, "test-flag", "TEST_ENV", "default", false)
		if got != "default" {
			t.Errorf("ResolveString() with empty ENV = %v, want %v", got, "default")
		}
	})
}

func TestResolveInt(t *testing.T) {
	t.Run("CLI flag has highest priority", func(t *testing.T) {
		fs := flag.NewFlagSet("test", flag.ContinueOnError)
		fs.String("test-flag", "", "test flag")
		setEnv(t, "TEST_ENV", "42")
		defer unsetEnv(t, "TEST_ENV")

		if err := fs.Parse([]string{"--test-flag", "8080"}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}

		got := ResolveInt(fs, "test-flag", "TEST_ENV", 0, true)
		if got != 8080 {
			t.Errorf("ResolveInt() = %v, want %v", got, 8080)
		}
	})

	t.Run("Environment variable has priority over default", func(t *testing.T) {
		fs := flag.NewFlagSet("test", flag.ContinueOnError)
		fs.String("test-flag", "", "test flag")
		setEnv(t, "TEST_ENV", "42")
		defer unsetEnv(t, "TEST_ENV")

		if err := fs.Parse([]string{}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}

		got := ResolveInt(fs, "test-flag", "TEST_ENV", 0, true)
		if got != 42 {
			t.Errorf("ResolveInt() = %v, want %v", got, 42)
		}
	})

	t.Run("Default value used when neither CLI nor ENV set", func(t *testing.T) {
		fs := flag.NewFlagSet("test", flag.ContinueOnError)
		fs.String("test-flag", "", "test flag")
		if err := fs.Parse([]string{}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}

		got := ResolveInt(fs, "test-flag", "TEST_ENV", 100, true)
		if got != 100 {
			t.Errorf("ResolveInt() = %v, want %v", got, 100)
		}
	})

	t.Run("AllowZero=false treats zero as not set", func(t *testing.T) {
		fs := flag.NewFlagSet("test", flag.ContinueOnError)
		fs.String("test-flag", "", "test flag")
		setEnv(t, "TEST_ENV", "0")
		defer unsetEnv(t, "TEST_ENV")

		if err := fs.Parse([]string{}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}

		got := ResolveInt(fs, "test-flag", "TEST_ENV", 100, false)
		if got != 100 {
			t.Errorf("ResolveInt() with allowZero=false = %v, want %v", got, 100)
		}
	})

	t.Run("AllowZero=true allows zero from ENV", func(t *testing.T) {
		fs := flag.NewFlagSet("test", flag.ContinueOnError)
		fs.String("test-flag", "", "test flag")
		setEnv(t, "TEST_ENV", "0")
		defer unsetEnv(t, "TEST_ENV")

		if err := fs.Parse([]string{}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}

		got := ResolveInt(fs, "test-flag", "TEST_ENV", 100, true)
		if got != 0 {
			t.Errorf("ResolveInt() with allowZero=true = %v, want %v", got, 0)
		}
	})

	t.Run("CLI flag zero value is always used", func(t *testing.T) {
		fs := flag.NewFlagSet("test", flag.ContinueOnError)
		fs.String("test-flag", "", "test flag")
		setEnv(t, "TEST_ENV", "42")
		defer unsetEnv(t, "TEST_ENV")

		if err := fs.Parse([]string{"--test-flag", "0"}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}

		got := ResolveInt(fs, "test-flag", "TEST_ENV", 100, false)
		if got != 0 {
			t.Errorf("ResolveInt() with CLI zero = %v, want %v", got, 0)
		}
	})

	t.Run("Invalid ENV value falls back to default", func(t *testing.T) {
		fs := flag.NewFlagSet("test", flag.ContinueOnError)
		fs.String("test-flag", "", "test flag")
		setEnv(t, "TEST_ENV", "not_a_number")
		defer unsetEnv(t, "TEST_ENV")

		if err := fs.Parse([]string{}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}

		got := ResolveInt(fs, "test-flag", "TEST_ENV", 100, true)
		if got != 100 {
			t.Errorf("ResolveInt() with invalid ENV = %v, want %v", got, 100)
		}
	})
}

func TestResolveInt64(t *testing.T) {
	t.Run("CLI flag has highest priority", func(t *testing.T) {
		fs := flag.NewFlagSet("test", flag.ContinueOnError)
		fs.String("test-flag", "", "test flag")
		setEnv(t, "TEST_ENV", "42")
		defer unsetEnv(t, "TEST_ENV")

		if err := fs.Parse([]string{"--test-flag", "9223372036854775807"}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}

		got := ResolveInt64(fs, "test-flag", "TEST_ENV", 0, true)
		if got != 9223372036854775807 {
			t.Errorf("ResolveInt64() = %v, want %v", got, int64(9223372036854775807))
		}
	})

	t.Run("Environment variable has priority over default", func(t *testing.T) {
		fs := flag.NewFlagSet("test", flag.ContinueOnError)
		fs.String("test-flag", "", "test flag")
		setEnv(t, "TEST_ENV", "10485760")
		defer unsetEnv(t, "TEST_ENV")

		if err := fs.Parse([]string{}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}

		got := ResolveInt64(fs, "test-flag", "TEST_ENV", 0, true)
		if got != 10485760 {
			t.Errorf("ResolveInt64() = %v, want %v", got, int64(10485760))
		}
	})

	t.Run("Default value used when neither CLI nor ENV set", func(t *testing.T) {
		fs := flag.NewFlagSet("test", flag.ContinueOnError)
		fs.String("test-flag", "", "test flag")
		if err := fs.Parse([]string{}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}

		got := ResolveInt64(fs, "test-flag", "TEST_ENV", 1048576, true)
		if got != 1048576 {
			t.Errorf("ResolveInt64() = %v, want %v", got, int64(1048576))
		}
	})

	t.Run("AllowZero=false treats zero as not set", func(t *testing.T) {
		fs := flag.NewFlagSet("test", flag.ContinueOnError)
		fs.String("test-flag", "", "test flag")
		setEnv(t, "TEST_ENV", "0")
		defer unsetEnv(t, "TEST_ENV")

		if err := fs.Parse([]string{}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}

		got := ResolveInt64(fs, "test-flag", "TEST_ENV", 1048576, false)
		if got != 1048576 {
			t.Errorf("ResolveInt64() with allowZero=false = %v, want %v", got, int64(1048576))
		}
	})

	t.Run("AllowZero=true allows zero from ENV", func(t *testing.T) {
		fs := flag.NewFlagSet("test", flag.ContinueOnError)
		fs.String("test-flag", "", "test flag")
		setEnv(t, "TEST_ENV", "0")
		defer unsetEnv(t, "TEST_ENV")

		if err := fs.Parse([]string{}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}

		got := ResolveInt64(fs, "test-flag", "TEST_ENV", 1048576, true)
		if got != 0 {
			t.Errorf("ResolveInt64() with allowZero=true = %v, want %v", got, int64(0))
		}
	})

	t.Run("CLI flag zero value is always used", func(t *testing.T) {
		fs := flag.NewFlagSet("test", flag.ContinueOnError)
		fs.String("test-flag", "", "test flag")
		setEnv(t, "TEST_ENV", "42")
		defer unsetEnv(t, "TEST_ENV")

		if err := fs.Parse([]string{"--test-flag", "0"}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}

		got := ResolveInt64(fs, "test-flag", "TEST_ENV", 1048576, false)
		if got != 0 {
			t.Errorf("ResolveInt64() with CLI zero = %v, want %v", got, int64(0))
		}
	})

	t.Run("Invalid ENV value falls back to default", func(t *testing.T) {
		fs := flag.NewFlagSet("test", flag.ContinueOnError)
		fs.String("test-flag", "", "test flag")
		setEnv(t, "TEST_ENV", "not_a_number")
		defer unsetEnv(t, "TEST_ENV")

		if err := fs.Parse([]string{}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}

		got := ResolveInt64(fs, "test-flag", "TEST_ENV", 1048576, true)
		if got != 1048576 {
			t.Errorf("ResolveInt64() with invalid ENV = %v, want %v", got, int64(1048576))
		}
	})
}

func TestResolveInt64WithValidation(t *testing.T) {
	validator := func(i int64) error {
		if i < 1024 || i > 1073741824 {
			return fmt.Errorf("value must be between 1KB and 1GB")
		}
		return nil
	}

	t.Run("Valid CLI value is used", func(t *testing.T) {
		fs := flag.NewFlagSet("test", flag.ContinueOnError)
		fs.String("test-flag", "", "test flag")
		setEnv(t, "TEST_ENV", "10485760")
		defer unsetEnv(t, "TEST_ENV")

		if err := fs.Parse([]string{"--test-flag", "1048576"}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}

		got, err := ResolveInt64WithValidation(fs, "test-flag", "TEST_ENV", 4096, true, validator)
		if err != nil {
			t.Errorf("ResolveInt64WithValidation() error = %v", err)
		}
		if got != 1048576 {
			t.Errorf("ResolveInt64WithValidation() = %v, want %v", got, int64(1048576))
		}
	})

	t.Run("Invalid CLI value falls back to ENV", func(t *testing.T) {
		fs := flag.NewFlagSet("test", flag.ContinueOnError)
		fs.String("test-flag", "", "test flag")
		setEnv(t, "TEST_ENV", "10485760")
		defer unsetEnv(t, "TEST_ENV")

		if err := fs.Parse([]string{"--test-flag", "500"}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}

		got, err := ResolveInt64WithValidation(fs, "test-flag", "TEST_ENV", 4096, true, validator)
		if err != nil {
			t.Errorf("ResolveInt64WithValidation() error = %v", err)
		}
		if got != 10485760 {
			t.Errorf("ResolveInt64WithValidation() = %v, want %v", got, int64(10485760))
		}
	})

	t.Run("All sources invalid returns error", func(t *testing.T) {
		fs := flag.NewFlagSet("test", flag.ContinueOnError)
		fs.String("test-flag", "", "test flag")
		setEnv(t, "TEST_ENV", "500")
		defer unsetEnv(t, "TEST_ENV")

		if err := fs.Parse([]string{"--test-flag", "100"}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}

		got, err := ResolveInt64WithValidation(fs, "test-flag", "TEST_ENV", 100, true, validator)
		if err == nil {
			t.Error("ResolveInt64WithValidation() error = nil, want error")
		}
		if got != 100 {
			t.Errorf("ResolveInt64WithValidation() = %v, want %v", got, int64(100))
		}
	})

	t.Run("AllowZero=false with zero ENV value", func(t *testing.T) {
		fs := flag.NewFlagSet("test", flag.ContinueOnError)
		fs.String("test-flag", "", "test flag")
		setEnv(t, "TEST_ENV", "0")
		defer unsetEnv(t, "TEST_ENV")

		if err := fs.Parse([]string{}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}

		got, err := ResolveInt64WithValidation(fs, "test-flag", "TEST_ENV", 4096, false, validator)
		if err != nil {
			t.Errorf("ResolveInt64WithValidation() error = %v", err)
		}
		if got != 4096 {
			t.Errorf("ResolveInt64WithValidation() = %v, want %v", got, int64(4096))
		}
	})
}

func TestResolveBool(t *testing.T) {
	t.Run("CLI flag has highest priority", func(t *testing.T) {
		fs := flag.NewFlagSet("test", flag.ContinueOnError)
		fs.Bool("test-flag", false, "test flag")
		setEnv(t, "TEST_ENV", "false")
		defer unsetEnv(t, "TEST_ENV")

		if err := fs.Parse([]string{"--test-flag"}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}

		got := ResolveBool(fs, "test-flag", "TEST_ENV", false)
		if got != true {
			t.Errorf("ResolveBool() = %v, want %v", got, true)
		}
	})

	t.Run("Environment variable has priority over default", func(t *testing.T) {
		fs := flag.NewFlagSet("test", flag.ContinueOnError)
		fs.Bool("test-flag", false, "test flag")
		setEnv(t, "TEST_ENV", "true")
		defer unsetEnv(t, "TEST_ENV")

		if err := fs.Parse([]string{}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}

		got := ResolveBool(fs, "test-flag", "TEST_ENV", false)
		if got != true {
			t.Errorf("ResolveBool() = %v, want %v", got, true)
		}
	})

	t.Run("Default value used when neither CLI nor ENV set", func(t *testing.T) {
		fs := flag.NewFlagSet("test", flag.ContinueOnError)
		fs.Bool("test-flag", false, "test flag")
		if err := fs.Parse([]string{}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}

		got := ResolveBool(fs, "test-flag", "TEST_ENV", true)
		if got != true {
			t.Errorf("ResolveBool() = %v, want %v", got, true)
		}
	})

	t.Run("ENV false value is respected", func(t *testing.T) {
		fs := flag.NewFlagSet("test", flag.ContinueOnError)
		fs.Bool("test-flag", false, "test flag")
		setEnv(t, "TEST_ENV", "false")
		defer unsetEnv(t, "TEST_ENV")

		if err := fs.Parse([]string{}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}

		got := ResolveBool(fs, "test-flag", "TEST_ENV", true)
		if got != false {
			t.Errorf("ResolveBool() = %v, want %v", got, false)
		}
	})

	t.Run("Invalid ENV value falls back to default", func(t *testing.T) {
		fs := flag.NewFlagSet("test", flag.ContinueOnError)
		fs.Bool("test-flag", false, "test flag")
		setEnv(t, "TEST_ENV", "maybe")
		defer unsetEnv(t, "TEST_ENV")

		if err := fs.Parse([]string{}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}

		got := ResolveBool(fs, "test-flag", "TEST_ENV", true)
		if got != true {
			t.Errorf("ResolveBool() with invalid ENV = %v, want %v", got, true)
		}
	})
}

func TestResolveDuration(t *testing.T) {
	t.Run("CLI flag has highest priority", func(t *testing.T) {
		fs := flag.NewFlagSet("test", flag.ContinueOnError)
		fs.String("test-flag", "", "test flag")
		setEnv(t, "TEST_ENV", "5m")
		defer unsetEnv(t, "TEST_ENV")

		if err := fs.Parse([]string{"--test-flag", "10s"}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}

		got := ResolveDuration(fs, "test-flag", "TEST_ENV", time.Second)
		if got != 10*time.Second {
			t.Errorf("ResolveDuration() = %v, want %v", got, 10*time.Second)
		}
	})

	t.Run("Environment variable has priority over default", func(t *testing.T) {
		fs := flag.NewFlagSet("test", flag.ContinueOnError)
		fs.String("test-flag", "", "test flag")
		setEnv(t, "TEST_ENV", "5m")
		defer unsetEnv(t, "TEST_ENV")

		if err := fs.Parse([]string{}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}

		got := ResolveDuration(fs, "test-flag", "TEST_ENV", time.Second)
		if got != 5*time.Minute {
			t.Errorf("ResolveDuration() = %v, want %v", got, 5*time.Minute)
		}
	})

	t.Run("Default value used when neither CLI nor ENV set", func(t *testing.T) {
		fs := flag.NewFlagSet("test", flag.ContinueOnError)
		fs.String("test-flag", "", "test flag")
		if err := fs.Parse([]string{}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}

		defaultVal := 1 * time.Hour
		got := ResolveDuration(fs, "test-flag", "TEST_ENV", defaultVal)
		if got != defaultVal {
			t.Errorf("ResolveDuration() = %v, want %v", got, defaultVal)
		}
	})

	t.Run("Invalid ENV value falls back to default", func(t *testing.T) {
		fs := flag.NewFlagSet("test", flag.ContinueOnError)
		fs.String("test-flag", "", "test flag")
		setEnv(t, "TEST_ENV", "invalid_duration")
		defer unsetEnv(t, "TEST_ENV")

		if err := fs.Parse([]string{}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}

		defaultVal := 1 * time.Hour
		got := ResolveDuration(fs, "test-flag", "TEST_ENV", defaultVal)
		if got != defaultVal {
			t.Errorf("ResolveDuration() with invalid ENV = %v, want %v", got, defaultVal)
		}
	})
}

func TestResolveIntAsString(t *testing.T) {
	t.Run("CLI flag has highest priority", func(t *testing.T) {
		fs := flag.NewFlagSet("test", flag.ContinueOnError)
		fs.String("test-flag", "", "test flag")
		setEnv(t, "TEST_ENV", "42")
		defer unsetEnv(t, "TEST_ENV")

		if err := fs.Parse([]string{"--test-flag", "8080"}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}

		got := ResolveIntAsString(fs, "test-flag", "TEST_ENV", 0, true)
		if got != "8080" {
			t.Errorf("ResolveIntAsString() = %v, want %v", got, "8080")
		}
	})

	t.Run("Environment variable has priority over default", func(t *testing.T) {
		fs := flag.NewFlagSet("test", flag.ContinueOnError)
		fs.String("test-flag", "", "test flag")
		setEnv(t, "TEST_ENV", "42")
		defer unsetEnv(t, "TEST_ENV")

		if err := fs.Parse([]string{}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}

		got := ResolveIntAsString(fs, "test-flag", "TEST_ENV", 0, true)
		if got != "42" {
			t.Errorf("ResolveIntAsString() = %v, want %v", got, "42")
		}
	})

	t.Run("Default value used when neither CLI nor ENV set", func(t *testing.T) {
		fs := flag.NewFlagSet("test", flag.ContinueOnError)
		fs.String("test-flag", "", "test flag")
		if err := fs.Parse([]string{}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}

		got := ResolveIntAsString(fs, "test-flag", "TEST_ENV", 100, true)
		if got != "100" {
			t.Errorf("ResolveIntAsString() = %v, want %v", got, "100")
		}
	})
}

func TestResolveStringWithValidator(t *testing.T) {
	validator := func(s string) bool {
		return len(s) > 3
	}

	t.Run("Valid CLI value is used", func(t *testing.T) {
		fs := flag.NewFlagSet("test", flag.ContinueOnError)
		fs.String("test-flag", "", "test flag")
		setEnv(t, "TEST_ENV", "valid_env")
		defer unsetEnv(t, "TEST_ENV")

		if err := fs.Parse([]string{"--test-flag", "valid_cli"}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}

		got := ResolveStringWithValidator(fs, "test-flag", "TEST_ENV", "default", false, validator)
		if got != "valid_cli" {
			t.Errorf("ResolveStringWithValidator() = %v, want %v", got, "valid_cli")
		}
	})

	t.Run("Invalid CLI value falls back to default", func(t *testing.T) {
		fs := flag.NewFlagSet("test", flag.ContinueOnError)
		fs.String("test-flag", "", "test flag")
		setEnv(t, "TEST_ENV", "valid_env")
		defer unsetEnv(t, "TEST_ENV")

		if err := fs.Parse([]string{"--test-flag", "bad"}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}

		got := ResolveStringWithValidator(fs, "test-flag", "TEST_ENV", "default", false, validator)
		if got != "default" {
			t.Errorf("ResolveStringWithValidator() = %v, want %v", got, "default")
		}
	})

	t.Run("Valid ENV value is used", func(t *testing.T) {
		fs := flag.NewFlagSet("test", flag.ContinueOnError)
		fs.String("test-flag", "", "test flag")
		setEnv(t, "TEST_ENV", "valid_env")
		defer unsetEnv(t, "TEST_ENV")

		if err := fs.Parse([]string{}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}

		got := ResolveStringWithValidator(fs, "test-flag", "TEST_ENV", "default", false, validator)
		if got != "valid_env" {
			t.Errorf("ResolveStringWithValidator() = %v, want %v", got, "valid_env")
		}
	})

	t.Run("Invalid ENV value falls back to default", func(t *testing.T) {
		fs := flag.NewFlagSet("test", flag.ContinueOnError)
		fs.String("test-flag", "", "test flag")
		setEnv(t, "TEST_ENV", "bad")
		defer unsetEnv(t, "TEST_ENV")

		if err := fs.Parse([]string{}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}

		got := ResolveStringWithValidator(fs, "test-flag", "TEST_ENV", "default", false, validator)
		if got != "default" {
			t.Errorf("ResolveStringWithValidator() = %v, want %v", got, "default")
		}
	})

	t.Run("Empty ENV value falls back to default", func(t *testing.T) {
		fs := flag.NewFlagSet("test", flag.ContinueOnError)
		fs.String("test-flag", "", "test flag")
		setEnv(t, "TEST_ENV", "")
		defer unsetEnv(t, "TEST_ENV")

		if err := fs.Parse([]string{}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}

		got := ResolveStringWithValidator(fs, "test-flag", "TEST_ENV", "default", false, validator)
		if got != "default" {
			t.Errorf("ResolveStringWithValidator() = %v, want %v", got, "default")
		}
	})

	t.Run("Trimmed ENV value", func(t *testing.T) {
		fs := flag.NewFlagSet("test", flag.ContinueOnError)
		fs.String("test-flag", "", "test flag")
		setEnv(t, "TEST_ENV", "  valid_env  ")
		defer unsetEnv(t, "TEST_ENV")

		if err := fs.Parse([]string{}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}

		got := ResolveStringWithValidator(fs, "test-flag", "TEST_ENV", "default", true, validator)
		if got != "valid_env" {
			t.Errorf("ResolveStringWithValidator() = %v, want %v", got, "valid_env")
		}
	})
}

func TestResolveStringNonEmpty(t *testing.T) {
	t.Run("Non-empty CLI value is used", func(t *testing.T) {
		fs := flag.NewFlagSet("test", flag.ContinueOnError)
		fs.String("test-flag", "", "test flag")
		setEnv(t, "TEST_ENV", "env_value")
		defer unsetEnv(t, "TEST_ENV")

		if err := fs.Parse([]string{"--test-flag", "cli_value"}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}

		got := ResolveStringNonEmpty(fs, "test-flag", "TEST_ENV", "default", false)
		if got != "cli_value" {
			t.Errorf("ResolveStringNonEmpty() = %v, want %v", got, "cli_value")
		}
	})

	t.Run("Empty CLI value falls back to ENV", func(t *testing.T) {
		fs := flag.NewFlagSet("test", flag.ContinueOnError)
		fs.String("test-flag", "", "test flag")
		setEnv(t, "TEST_ENV", "env_value")
		defer unsetEnv(t, "TEST_ENV")

		if err := fs.Parse([]string{"--test-flag", ""}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}

		got := ResolveStringNonEmpty(fs, "test-flag", "TEST_ENV", "default", false)
		if got != "env_value" {
			t.Errorf("ResolveStringNonEmpty() = %v, want %v", got, "env_value")
		}
	})

	t.Run("Empty ENV value falls back to default", func(t *testing.T) {
		fs := flag.NewFlagSet("test", flag.ContinueOnError)
		fs.String("test-flag", "", "test flag")
		setEnv(t, "TEST_ENV", "")
		defer unsetEnv(t, "TEST_ENV")

		if err := fs.Parse([]string{}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}

		got := ResolveStringNonEmpty(fs, "test-flag", "TEST_ENV", "default", false)
		if got != "default" {
			t.Errorf("ResolveStringNonEmpty() = %v, want %v", got, "default")
		}
	})

	t.Run("Trimmed CLI value", func(t *testing.T) {
		fs := flag.NewFlagSet("test", flag.ContinueOnError)
		fs.String("test-flag", "", "test flag")
		setEnv(t, "TEST_ENV", "env_value")
		defer unsetEnv(t, "TEST_ENV")

		if err := fs.Parse([]string{"--test-flag", "  cli_value  "}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}

		got := ResolveStringNonEmpty(fs, "test-flag", "TEST_ENV", "default", true)
		// Note: flag parsing doesn't trim, so we get the value as-is
		// But trimmed option only affects ENV, not CLI
		if got != "  cli_value  " {
			t.Errorf("ResolveStringNonEmpty() = %v, want %v", got, "  cli_value  ")
		}
	})

	t.Run("Trimmed ENV value", func(t *testing.T) {
		fs := flag.NewFlagSet("test", flag.ContinueOnError)
		fs.String("test-flag", "", "test flag")
		setEnv(t, "TEST_ENV", "  env_value  ")
		defer unsetEnv(t, "TEST_ENV")

		if err := fs.Parse([]string{}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}

		got := ResolveStringNonEmpty(fs, "test-flag", "TEST_ENV", "default", true)
		if got != "env_value" {
			t.Errorf("ResolveStringNonEmpty() = %v, want %v", got, "env_value")
		}
	})
}

func TestResolveStringWithValidation(t *testing.T) {
	validator := func(s string) error {
		if len(s) < 3 {
			return fmt.Errorf("value too short")
		}
		return nil
	}

	t.Run("Valid CLI value is used", func(t *testing.T) {
		fs := flag.NewFlagSet("test", flag.ContinueOnError)
		fs.String("test-flag", "", "test flag")
		setEnv(t, "TEST_ENV", "valid_env")
		defer unsetEnv(t, "TEST_ENV")

		if err := fs.Parse([]string{"--test-flag", "valid_cli"}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}

		got, err := ResolveStringWithValidation(fs, "test-flag", "TEST_ENV", "default", false, validator)
		if err != nil {
			t.Errorf("ResolveStringWithValidation() error = %v", err)
		}
		if got != "valid_cli" {
			t.Errorf("ResolveStringWithValidation() = %v, want %v", got, "valid_cli")
		}
	})

	t.Run("Invalid CLI value falls back to ENV", func(t *testing.T) {
		fs := flag.NewFlagSet("test", flag.ContinueOnError)
		fs.String("test-flag", "", "test flag")
		setEnv(t, "TEST_ENV", "valid_env")
		defer unsetEnv(t, "TEST_ENV")

		if err := fs.Parse([]string{"--test-flag", "ab"}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}

		got, err := ResolveStringWithValidation(fs, "test-flag", "TEST_ENV", "default", false, validator)
		if err != nil {
			t.Errorf("ResolveStringWithValidation() error = %v", err)
		}
		if got != "valid_env" {
			t.Errorf("ResolveStringWithValidation() = %v, want %v", got, "valid_env")
		}
	})

	t.Run("All sources invalid returns error", func(t *testing.T) {
		fs := flag.NewFlagSet("test", flag.ContinueOnError)
		fs.String("test-flag", "", "test flag")
		setEnv(t, "TEST_ENV", "ab")
		defer unsetEnv(t, "TEST_ENV")

		if err := fs.Parse([]string{"--test-flag", "xy"}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}

		got, err := ResolveStringWithValidation(fs, "test-flag", "TEST_ENV", "de", false, validator)
		if err == nil {
			t.Error("ResolveStringWithValidation() error = nil, want error")
		}
		if got != "de" {
			t.Errorf("ResolveStringWithValidation() = %v, want %v", got, "de")
		}
	})

	t.Run("Empty ENV value falls back to default", func(t *testing.T) {
		fs := flag.NewFlagSet("test", flag.ContinueOnError)
		fs.String("test-flag", "", "test flag")
		setEnv(t, "TEST_ENV", "")
		defer unsetEnv(t, "TEST_ENV")

		if err := fs.Parse([]string{}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}

		got, err := ResolveStringWithValidation(fs, "test-flag", "TEST_ENV", "default", false, validator)
		if err != nil {
			t.Errorf("ResolveStringWithValidation() error = %v", err)
		}
		if got != "default" {
			t.Errorf("ResolveStringWithValidation() = %v, want %v", got, "default")
		}
	})

	t.Run("Trimmed ENV value", func(t *testing.T) {
		fs := flag.NewFlagSet("test", flag.ContinueOnError)
		fs.String("test-flag", "", "test flag")
		setEnv(t, "TEST_ENV", "  valid_env  ")
		defer unsetEnv(t, "TEST_ENV")

		if err := fs.Parse([]string{}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}

		got, err := ResolveStringWithValidation(fs, "test-flag", "TEST_ENV", "default", true, validator)
		if err != nil {
			t.Errorf("ResolveStringWithValidation() error = %v", err)
		}
		if got != "valid_env" {
			t.Errorf("ResolveStringWithValidation() = %v, want %v", got, "valid_env")
		}
	})
}

func TestResolveIntWithValidation(t *testing.T) {
	validator := func(i int) error {
		if i < 1 || i > 100 {
			return fmt.Errorf("value must be between 1 and 100")
		}
		return nil
	}

	t.Run("Valid CLI value is used", func(t *testing.T) {
		fs := flag.NewFlagSet("test", flag.ContinueOnError)
		fs.String("test-flag", "", "test flag")
		setEnv(t, "TEST_ENV", "50")
		defer unsetEnv(t, "TEST_ENV")

		if err := fs.Parse([]string{"--test-flag", "42"}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}

		got, err := ResolveIntWithValidation(fs, "test-flag", "TEST_ENV", 10, true, validator)
		if err != nil {
			t.Errorf("ResolveIntWithValidation() error = %v", err)
		}
		if got != 42 {
			t.Errorf("ResolveIntWithValidation() = %v, want %v", got, 42)
		}
	})

	t.Run("Invalid CLI value falls back to ENV", func(t *testing.T) {
		fs := flag.NewFlagSet("test", flag.ContinueOnError)
		fs.String("test-flag", "", "test flag")
		setEnv(t, "TEST_ENV", "50")
		defer unsetEnv(t, "TEST_ENV")

		if err := fs.Parse([]string{"--test-flag", "200"}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}

		got, err := ResolveIntWithValidation(fs, "test-flag", "TEST_ENV", 10, true, validator)
		if err != nil {
			t.Errorf("ResolveIntWithValidation() error = %v", err)
		}
		if got != 50 {
			t.Errorf("ResolveIntWithValidation() = %v, want %v", got, 50)
		}
	})

	t.Run("All sources invalid returns error", func(t *testing.T) {
		fs := flag.NewFlagSet("test", flag.ContinueOnError)
		fs.String("test-flag", "", "test flag")
		setEnv(t, "TEST_ENV", "200")
		defer unsetEnv(t, "TEST_ENV")

		if err := fs.Parse([]string{"--test-flag", "300"}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}

		got, err := ResolveIntWithValidation(fs, "test-flag", "TEST_ENV", 300, true, validator)
		if err == nil {
			t.Error("ResolveIntWithValidation() error = nil, want error")
		}
		if got != 300 {
			t.Errorf("ResolveIntWithValidation() = %v, want %v", got, 300)
		}
	})

	t.Run("AllowZero=false with zero ENV value", func(t *testing.T) {
		fs := flag.NewFlagSet("test", flag.ContinueOnError)
		fs.String("test-flag", "", "test flag")
		setEnv(t, "TEST_ENV", "0")
		defer unsetEnv(t, "TEST_ENV")

		if err := fs.Parse([]string{}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}

		got, err := ResolveIntWithValidation(fs, "test-flag", "TEST_ENV", 50, false, validator)
		if err != nil {
			t.Errorf("ResolveIntWithValidation() error = %v", err)
		}
		if got != 50 {
			t.Errorf("ResolveIntWithValidation() = %v, want %v", got, 50)
		}
	})
}

func TestResolveStringSlice(t *testing.T) {
	t.Run("CLI flag has highest priority", func(t *testing.T) {
		fs := flag.NewFlagSet("test", flag.ContinueOnError)
		fs.String("test-flag", "", "test flag")
		setEnv(t, "TEST_ENV", "env1,env2,env3")
		defer unsetEnv(t, "TEST_ENV")

		if err := fs.Parse([]string{"--test-flag", "cli_value"}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}

		got := ResolveStringSlice(fs, "test-flag", "TEST_ENV", []string{"default"}, ",")
		if len(got) != 1 || got[0] != "cli_value" {
			t.Errorf("ResolveStringSlice() = %v, want %v", got, []string{"cli_value"})
		}
	})

	t.Run("Environment variable has priority over default", func(t *testing.T) {
		fs := flag.NewFlagSet("test", flag.ContinueOnError)
		fs.String("test-flag", "", "test flag")
		setEnv(t, "TEST_ENV", "env1,env2,env3")
		defer unsetEnv(t, "TEST_ENV")

		if err := fs.Parse([]string{}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}

		got := ResolveStringSlice(fs, "test-flag", "TEST_ENV", []string{"default"}, ",")
		expected := []string{"env1", "env2", "env3"}
		if len(got) != len(expected) {
			t.Errorf("ResolveStringSlice() length = %d, want %d", len(got), len(expected))
		}
		for i, v := range got {
			if v != expected[i] {
				t.Errorf("ResolveStringSlice()[%d] = %v, want %v", i, v, expected[i])
			}
		}
	})

	t.Run("Default value used when neither CLI nor ENV set", func(t *testing.T) {
		fs := flag.NewFlagSet("test", flag.ContinueOnError)
		fs.String("test-flag", "", "test flag")
		if err := fs.Parse([]string{}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}

		got := ResolveStringSlice(fs, "test-flag", "TEST_ENV", []string{"default1", "default2"}, ",")
		expected := []string{"default1", "default2"}
		if len(got) != len(expected) {
			t.Errorf("ResolveStringSlice() length = %d, want %d", len(got), len(expected))
		}
	})

	t.Run("Custom separator works", func(t *testing.T) {
		fs := flag.NewFlagSet("test", flag.ContinueOnError)
		fs.String("test-flag", "", "test flag")
		setEnv(t, "TEST_ENV", "env1;env2;env3")
		defer unsetEnv(t, "TEST_ENV")

		if err := fs.Parse([]string{}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}

		got := ResolveStringSlice(fs, "test-flag", "TEST_ENV", []string{"default"}, ";")
		expected := []string{"env1", "env2", "env3"}
		if len(got) != len(expected) {
			t.Errorf("ResolveStringSlice() length = %d, want %d", len(got), len(expected))
		}
	})

	t.Run("Empty separator defaults to comma", func(t *testing.T) {
		fs := flag.NewFlagSet("test", flag.ContinueOnError)
		fs.String("test-flag", "", "test flag")
		setEnv(t, "TEST_ENV", "env1,env2")
		defer unsetEnv(t, "TEST_ENV")

		if err := fs.Parse([]string{}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}

		got := ResolveStringSlice(fs, "test-flag", "TEST_ENV", []string{"default"}, "")
		if len(got) != 2 {
			t.Errorf("ResolveStringSlice() length = %d, want %d", len(got), 2)
		}
	})

	t.Run("Whitespace trimming in ENV values", func(t *testing.T) {
		fs := flag.NewFlagSet("test", flag.ContinueOnError)
		fs.String("test-flag", "", "test flag")
		setEnv(t, "TEST_ENV", " env1 , env2 , env3 ")
		defer unsetEnv(t, "TEST_ENV")

		if err := fs.Parse([]string{}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}

		got := ResolveStringSlice(fs, "test-flag", "TEST_ENV", []string{"default"}, ",")
		expected := []string{"env1", "env2", "env3"}
		if len(got) != len(expected) {
			t.Errorf("ResolveStringSlice() length = %d, want %d", len(got), len(expected))
		}
		for i, v := range got {
			if v != expected[i] {
				t.Errorf("ResolveStringSlice()[%d] = %v, want %v", i, v, expected[i])
			}
		}
	})
}

func TestResolveStringSliceMulti(t *testing.T) {
	t.Run("CLI flag values have highest priority", func(t *testing.T) {
		fs := flag.NewFlagSet("test", flag.ContinueOnError)
		fs.String("test-flag", "", "test flag")
		setEnv(t, "TEST_ENV", "env1,env2")
		defer unsetEnv(t, "TEST_ENV")

		if err := fs.Parse([]string{"--test-flag", "cli1"}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}

		// Simulating multi-value flag where values were collected
		currentFlagValue := []string{"cli1", "cli2"}
		got := ResolveStringSliceMulti(fs, "test-flag", "TEST_ENV", currentFlagValue, []string{"default"}, ",")
		if len(got) != 2 || got[0] != "cli1" || got[1] != "cli2" {
			t.Errorf("ResolveStringSliceMulti() = %v, want %v", got, []string{"cli1", "cli2"})
		}
	})

	t.Run("ENV used when CLI flag not set", func(t *testing.T) {
		fs := flag.NewFlagSet("test", flag.ContinueOnError)
		fs.String("test-flag", "", "test flag")
		setEnv(t, "TEST_ENV", "env1,env2,env3")
		defer unsetEnv(t, "TEST_ENV")

		if err := fs.Parse([]string{}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}

		currentFlagValue := []string{}
		got := ResolveStringSliceMulti(fs, "test-flag", "TEST_ENV", currentFlagValue, []string{"default"}, ",")
		expected := []string{"env1", "env2", "env3"}
		if len(got) != len(expected) {
			t.Errorf("ResolveStringSliceMulti() length = %d, want %d", len(got), len(expected))
		}
	})

	t.Run("Default value used when neither CLI nor ENV set", func(t *testing.T) {
		fs := flag.NewFlagSet("test", flag.ContinueOnError)
		fs.String("test-flag", "", "test flag")
		if err := fs.Parse([]string{}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}

		currentFlagValue := []string{}
		got := ResolveStringSliceMulti(fs, "test-flag", "TEST_ENV", currentFlagValue, []string{"default1", "default2"}, ",")
		expected := []string{"default1", "default2"}
		if len(got) != len(expected) {
			t.Errorf("ResolveStringSliceMulti() length = %d, want %d", len(got), len(expected))
		}
	})
}

func TestResolveEnum(t *testing.T) {
	allowedValues := []string{"DEFAULT", "REMOTE_FIRST", "ONLY_LOCAL"}

	t.Run("Valid CLI value", func(t *testing.T) {
		fs := flag.NewFlagSet("test", flag.ContinueOnError)
		fs.String("test-flag", "", "test flag")
		setEnv(t, "TEST_ENV", "REMOTE_FIRST")
		defer unsetEnv(t, "TEST_ENV")

		if err := fs.Parse([]string{"--test-flag", "DEFAULT"}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}

		got, err := ResolveEnum(fs, "test-flag", "TEST_ENV", "ONLY_LOCAL", allowedValues, true)
		if err != nil {
			t.Errorf("ResolveEnum() error = %v", err)
		}
		if got != "DEFAULT" {
			t.Errorf("ResolveEnum() = %v, want %v", got, "DEFAULT")
		}
	})

	t.Run("Case insensitive", func(t *testing.T) {
		fs := flag.NewFlagSet("test", flag.ContinueOnError)
		fs.String("test-flag", "", "test flag")
		setEnv(t, "TEST_ENV", "default")
		defer unsetEnv(t, "TEST_ENV")

		if err := fs.Parse([]string{}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}

		got, err := ResolveEnum(fs, "test-flag", "TEST_ENV", "ONLY_LOCAL", allowedValues, false)
		if err != nil {
			t.Errorf("ResolveEnum() error = %v", err)
		}
		if got != "default" {
			t.Errorf("ResolveEnum() = %v, want %v", got, "default")
		}
	})

	t.Run("Invalid value returns error", func(t *testing.T) {
		fs := flag.NewFlagSet("test", flag.ContinueOnError)
		fs.String("test-flag", "", "test flag")
		setEnv(t, "TEST_ENV", "INVALID")
		defer unsetEnv(t, "TEST_ENV")

		if err := fs.Parse([]string{}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}

		got, err := ResolveEnum(fs, "test-flag", "TEST_ENV", "INVALID", allowedValues, true)
		if err == nil {
			t.Error("ResolveEnum() error = nil, want error")
		}
		if got != "INVALID" {
			t.Errorf("ResolveEnum() = %v, want %v", got, "INVALID")
		}
	})
}

func TestResolveHostPort(t *testing.T) {
	t.Run("Valid CLI value", func(t *testing.T) {
		fs := flag.NewFlagSet("test", flag.ContinueOnError)
		fs.String("test-flag", "", "test flag")
		setEnv(t, "TEST_ENV", "localhost:6379")
		defer unsetEnv(t, "TEST_ENV")

		if err := fs.Parse([]string{"--test-flag", "example.com:8080"}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}

		host, port, err := ResolveHostPort(fs, "test-flag", "TEST_ENV", "localhost:6379")
		if err != nil {
			t.Errorf("ResolveHostPort() error = %v", err)
		}
		if host != "example.com" || port != 8080 {
			t.Errorf("ResolveHostPort() = (%q, %d), want (%q, %d)", host, port, "example.com", 8080)
		}
	})

	t.Run("Invalid format returns error", func(t *testing.T) {
		fs := flag.NewFlagSet("test", flag.ContinueOnError)
		fs.String("test-flag", "", "test flag")
		setEnv(t, "TEST_ENV", "invalid")
		defer unsetEnv(t, "TEST_ENV")

		if err := fs.Parse([]string{}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}

		_, _, err := ResolveHostPort(fs, "test-flag", "TEST_ENV", "invalid")
		if err == nil {
			t.Error("ResolveHostPort() error = nil, want error")
		}
	})
}

func TestResolvePort(t *testing.T) {
	t.Run("Valid CLI port", func(t *testing.T) {
		fs := flag.NewFlagSet("test", flag.ContinueOnError)
		fs.String("test-flag", "", "test flag")
		setEnv(t, "TEST_ENV", "8080")
		defer unsetEnv(t, "TEST_ENV")

		if err := fs.Parse([]string{"--test-flag", "9090"}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}

		got, err := ResolvePort(fs, "test-flag", "TEST_ENV", 8080)
		if err != nil {
			t.Errorf("ResolvePort() error = %v", err)
		}
		if got != 9090 {
			t.Errorf("ResolvePort() = %v, want %v", got, 9090)
		}
	})

	t.Run("Invalid port returns error", func(t *testing.T) {
		fs := flag.NewFlagSet("test", flag.ContinueOnError)
		fs.String("test-flag", "", "test flag")
		setEnv(t, "TEST_ENV", "99999")
		defer unsetEnv(t, "TEST_ENV")

		if err := fs.Parse([]string{}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}

		got, err := ResolvePort(fs, "test-flag", "TEST_ENV", 99999)
		if err == nil {
			t.Error("ResolvePort() error = nil, want error")
		}
		if got != 99999 {
			t.Errorf("ResolvePort() = %v, want %v", got, 99999)
		}
	})
}
