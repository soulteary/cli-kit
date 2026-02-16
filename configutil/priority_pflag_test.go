package configutil

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/spf13/pflag"
)

func setEnvPflag(t *testing.T, key, value string) {
	t.Helper()
	if err := os.Setenv(key, value); err != nil {
		t.Fatalf("os.Setenv(%q, %q) failed: %v", key, value, err)
	}
}

func unsetEnvPflag(t *testing.T, key string) {
	t.Helper()
	if err := os.Unsetenv(key); err != nil {
		t.Fatalf("os.Unsetenv(%q) failed: %v", key, err)
	}
}

func TestResolveStringPflag(t *testing.T) {
	t.Run("CLI has highest priority", func(t *testing.T) {
		fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
		fs.StringP("name", "n", "", "name")
		setEnvPflag(t, "TEST_NAME", "env_value")
		defer unsetEnvPflag(t, "TEST_NAME")
		if err := fs.Parse([]string{"--name", "cli_value"}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}
		got := ResolveStringPflag(fs, "name", "TEST_NAME", "default", false)
		if got != "cli_value" {
			t.Errorf("ResolveStringPflag() = %q, want cli_value", got)
		}
	})
	t.Run("ENV over default when CLI not set", func(t *testing.T) {
		fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
		fs.String("name", "", "name")
		setEnvPflag(t, "TEST_NAME", "env_value")
		defer unsetEnvPflag(t, "TEST_NAME")
		if err := fs.Parse([]string{}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}
		got := ResolveStringPflag(fs, "name", "TEST_NAME", "default", false)
		if got != "env_value" {
			t.Errorf("ResolveStringPflag() = %q, want env_value", got)
		}
	})
	t.Run("default when neither CLI nor ENV set", func(t *testing.T) {
		fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
		fs.String("name", "", "name")
		if err := fs.Parse([]string{}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}
		got := ResolveStringPflag(fs, "name", "TEST_NAME", "default", false)
		if got != "default" {
			t.Errorf("ResolveStringPflag() = %q, want default", got)
		}
	})
	t.Run("empty envKey ignores ENV", func(t *testing.T) {
		fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
		fs.String("name", "", "name")
		setEnvPflag(t, "TEST_NAME", "env_value")
		defer unsetEnvPflag(t, "TEST_NAME")
		if err := fs.Parse([]string{}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}
		got := ResolveStringPflag(fs, "name", "", "default", false)
		if got != "default" {
			t.Errorf("ResolveStringPflag(empty envKey) = %q, want default", got)
		}
	})
	t.Run("trimmed ENV value", func(t *testing.T) {
		fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
		fs.String("name", "", "name")
		setEnvPflag(t, "TEST_NAME", "  env_value  ")
		defer unsetEnvPflag(t, "TEST_NAME")
		if err := fs.Parse([]string{}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}
		got := ResolveStringPflag(fs, "name", "TEST_NAME", "default", true)
		if got != "env_value" {
			t.Errorf("ResolveStringPflag(trimmed) = %q, want env_value", got)
		}
	})
	t.Run("empty ENV value falls back to default", func(t *testing.T) {
		fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
		fs.String("name", "", "name")
		setEnvPflag(t, "TEST_NAME", "")
		defer unsetEnvPflag(t, "TEST_NAME")
		if err := fs.Parse([]string{}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}
		got := ResolveStringPflag(fs, "name", "TEST_NAME", "default", false)
		if got != "default" {
			t.Errorf("ResolveStringPflag(empty ENV) = %q, want default", got)
		}
	})
}

func TestResolveIntPflag(t *testing.T) {
	t.Run("CLI has highest priority", func(t *testing.T) {
		fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
		fs.IntP("port", "p", 8080, "port")
		setEnvPflag(t, "TEST_PORT", "9090")
		defer unsetEnvPflag(t, "TEST_PORT")
		if err := fs.Parse([]string{"--port", "3000"}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}
		got := ResolveIntPflag(fs, "port", "TEST_PORT", 8080, true)
		if got != 3000 {
			t.Errorf("ResolveIntPflag() = %d, want 3000", got)
		}
	})
	t.Run("ENV over default when CLI not set", func(t *testing.T) {
		fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
		fs.Int("port", 8080, "port")
		setEnvPflag(t, "TEST_PORT", "9090")
		defer unsetEnvPflag(t, "TEST_PORT")
		if err := fs.Parse([]string{}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}
		got := ResolveIntPflag(fs, "port", "TEST_PORT", 8080, true)
		if got != 9090 {
			t.Errorf("ResolveIntPflag() = %d, want 9090", got)
		}
	})
	t.Run("allowZero=false with ENV 0 uses default", func(t *testing.T) {
		fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
		fs.Int("port", 8080, "port")
		setEnvPflag(t, "TEST_PORT", "0")
		defer unsetEnvPflag(t, "TEST_PORT")
		if err := fs.Parse([]string{}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}
		got := ResolveIntPflag(fs, "port", "TEST_PORT", 42, false)
		if got != 42 {
			t.Errorf("ResolveIntPflag(allowZero=false, ENV=0) = %d, want 42", got)
		}
	})
	t.Run("allowZero=true with ENV 0 returns 0", func(t *testing.T) {
		fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
		fs.Int("port", 8080, "port")
		setEnvPflag(t, "TEST_PORT", "0")
		defer unsetEnvPflag(t, "TEST_PORT")
		if err := fs.Parse([]string{}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}
		got := ResolveIntPflag(fs, "port", "TEST_PORT", 42, true)
		if got != 0 {
			t.Errorf("ResolveIntPflag(allowZero=true, ENV=0) = %d, want 0", got)
		}
	})
}

func TestResolveBoolPflag(t *testing.T) {
	t.Run("empty envKey CLI set", func(t *testing.T) {
		fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
		fs.BoolP("on", "o", false, "switch")
		if err := fs.Parse([]string{"--on"}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}
		got := ResolveBoolPflag(fs, "on", "", false)
		if !got {
			t.Errorf("ResolveBoolPflag(empty envKey, CLI set) = %v, want true", got)
		}
	})
	t.Run("empty envKey not set uses default", func(t *testing.T) {
		fs := pflag.NewFlagSet("test2", pflag.ContinueOnError)
		fs.Bool("off", false, "off")
		if err := fs.Parse([]string{}); err != nil {
			t.Fatalf("fs2.Parse() failed: %v", err)
		}
		got := ResolveBoolPflag(fs, "off", "", true)
		if !got {
			t.Errorf("ResolveBoolPflag(empty envKey, not set) = %v, want true (default)", got)
		}
	})
	t.Run("ENV over default when CLI not set", func(t *testing.T) {
		fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
		fs.Bool("flag", false, "flag")
		setEnvPflag(t, "TEST_BOOL", "true")
		defer unsetEnvPflag(t, "TEST_BOOL")
		if err := fs.Parse([]string{}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}
		got := ResolveBoolPflag(fs, "flag", "TEST_BOOL", false)
		if !got {
			t.Errorf("ResolveBoolPflag() = %v, want true", got)
		}
	})
}

func TestResolvePortPflag_EmptyEnvKey(t *testing.T) {
	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	fs.IntP("port", "p", 8080, "port")
	if err := fs.Parse([]string{"--port", "3000"}); err != nil {
		t.Fatalf("fs.Parse() failed: %v", err)
	}
	port, err := ResolvePortPflag(fs, "port", "", 8080)
	if err != nil {
		t.Fatalf("ResolvePortPflag() err = %v", err)
	}
	if port != 3000 {
		t.Errorf("ResolvePortPflag() = %d, want 3000", port)
	}

	// Not set: use default
	fs2 := pflag.NewFlagSet("test2", pflag.ContinueOnError)
	fs2.Int("port", 5000, "port")
	if err := fs2.Parse([]string{}); err != nil {
		t.Fatalf("fs2.Parse() failed: %v", err)
	}
	port2, err2 := ResolvePortPflag(fs2, "port", "", 5005)
	if err2 != nil {
		t.Fatalf("ResolvePortPflag() err = %v", err2)
	}
	if port2 != 5005 {
		t.Errorf("ResolvePortPflag(not set) = %d, want 5005", port2)
	}
}

func TestResolveDurationPflag(t *testing.T) {
	t.Run("CLI has highest priority", func(t *testing.T) {
		fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
		fs.Duration("timeout", time.Second, "timeout")
		setEnvPflag(t, "TEST_TIMEOUT", "5m")
		defer unsetEnvPflag(t, "TEST_TIMEOUT")
		if err := fs.Parse([]string{"--timeout", "10s"}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}
		got := ResolveDurationPflag(fs, "timeout", "TEST_TIMEOUT", time.Minute)
		if got != 10*time.Second {
			t.Errorf("ResolveDurationPflag() = %v, want 10s", got)
		}
	})
	t.Run("ENV over default when CLI not set", func(t *testing.T) {
		fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
		fs.Duration("timeout", time.Second, "timeout")
		setEnvPflag(t, "TEST_TIMEOUT", "5m")
		defer unsetEnvPflag(t, "TEST_TIMEOUT")
		if err := fs.Parse([]string{}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}
		got := ResolveDurationPflag(fs, "timeout", "TEST_TIMEOUT", time.Minute)
		if got != 5*time.Minute {
			t.Errorf("ResolveDurationPflag() = %v, want 5m", got)
		}
	})
	t.Run("empty envKey uses default when not set", func(t *testing.T) {
		fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
		fs.Duration("timeout", time.Second, "timeout")
		setEnvPflag(t, "TEST_TIMEOUT", "5m")
		defer unsetEnvPflag(t, "TEST_TIMEOUT")
		if err := fs.Parse([]string{}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}
		got := ResolveDurationPflag(fs, "timeout", "", 2*time.Minute)
		if got != 2*time.Minute {
			t.Errorf("ResolveDurationPflag(empty envKey) = %v, want 2m", got)
		}
	})
}

func TestResolveIntAsStringPflag(t *testing.T) {
	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	fs.Int("port", 8080, "port")
	if err := fs.Parse([]string{"--port", "9090"}); err != nil {
		t.Fatalf("fs.Parse() failed: %v", err)
	}
	got := ResolveIntAsStringPflag(fs, "port", "", 8080, true)
	if got != "9090" {
		t.Errorf("ResolveIntAsStringPflag() = %q, want 9090", got)
	}
	fs2 := pflag.NewFlagSet("test2", pflag.ContinueOnError)
	fs2.Int("port", 8080, "port")
	if err := fs2.Parse([]string{}); err != nil {
		t.Fatalf("fs2.Parse() failed: %v", err)
	}
	got2 := ResolveIntAsStringPflag(fs2, "port", "", 8080, true)
	if got2 != "8080" {
		t.Errorf("ResolveIntAsStringPflag(not set) = %q, want 8080", got2)
	}
}

func TestResolveStringWithValidationPflag(t *testing.T) {
	validate := func(s string) error {
		if len(s) < 3 {
			return fmt.Errorf("too short")
		}
		return nil
	}
	t.Run("valid CLI value", func(t *testing.T) {
		fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
		fs.String("name", "", "name")
		setEnvPflag(t, "TEST_NAME", "env_ok")
		defer unsetEnvPflag(t, "TEST_NAME")
		if err := fs.Parse([]string{"--name", "cli_ok"}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}
		got, err := ResolveStringWithValidationPflag(fs, "name", "TEST_NAME", "default_ok", false, validate)
		if err != nil {
			t.Fatalf("ResolveStringWithValidationPflag() err = %v", err)
		}
		if got != "cli_ok" {
			t.Errorf("ResolveStringWithValidationPflag() = %q, want cli_ok", got)
		}
	})
	t.Run("invalid CLI falls back to ENV", func(t *testing.T) {
		fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
		fs.String("name", "", "name")
		setEnvPflag(t, "TEST_NAME", "env_ok")
		defer unsetEnvPflag(t, "TEST_NAME")
		if err := fs.Parse([]string{"--name", "ab"}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}
		got, err := ResolveStringWithValidationPflag(fs, "name", "TEST_NAME", "default_ok", false, validate)
		if err != nil {
			t.Fatalf("ResolveStringWithValidationPflag() err = %v", err)
		}
		if got != "env_ok" {
			t.Errorf("ResolveStringWithValidationPflag() = %q, want env_ok", got)
		}
	})
	t.Run("all invalid returns error", func(t *testing.T) {
		fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
		fs.String("name", "", "name")
		setEnvPflag(t, "TEST_NAME", "x")
		defer unsetEnvPflag(t, "TEST_NAME")
		if err := fs.Parse([]string{"--name", "y"}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}
		got, err := ResolveStringWithValidationPflag(fs, "name", "TEST_NAME", "z", false, validate)
		if err == nil {
			t.Error("ResolveStringWithValidationPflag() want error when all invalid")
		}
		if got != "z" {
			t.Errorf("ResolveStringWithValidationPflag() = %q, want default z", got)
		}
	})
}

func TestResolveIntWithValidationPflag(t *testing.T) {
	validate := func(i int) error {
		if i < 1 || i > 100 {
			return fmt.Errorf("must be 1-100")
		}
		return nil
	}
	t.Run("valid CLI value", func(t *testing.T) {
		fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
		fs.Int("n", 50, "n")
		setEnvPflag(t, "TEST_N", "50")
		defer unsetEnvPflag(t, "TEST_N")
		if err := fs.Parse([]string{"--n", "42"}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}
		got, err := ResolveIntWithValidationPflag(fs, "n", "TEST_N", 10, true, validate)
		if err != nil {
			t.Fatalf("ResolveIntWithValidationPflag() err = %v", err)
		}
		if got != 42 {
			t.Errorf("ResolveIntWithValidationPflag() = %d, want 42", got)
		}
	})
	t.Run("invalid CLI falls back to ENV", func(t *testing.T) {
		fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
		fs.Int("n", 50, "n")
		setEnvPflag(t, "TEST_N", "50")
		defer unsetEnvPflag(t, "TEST_N")
		if err := fs.Parse([]string{"--n", "200"}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}
		got, err := ResolveIntWithValidationPflag(fs, "n", "TEST_N", 10, true, validate)
		if err != nil {
			t.Fatalf("ResolveIntWithValidationPflag() err = %v", err)
		}
		if got != 50 {
			t.Errorf("ResolveIntWithValidationPflag() = %d, want 50", got)
		}
	})
	t.Run("all invalid returns error", func(t *testing.T) {
		fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
		fs.Int("n", 50, "n")
		setEnvPflag(t, "TEST_N", "200")
		defer unsetEnvPflag(t, "TEST_N")
		if err := fs.Parse([]string{"--n", "300"}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}
		got, err := ResolveIntWithValidationPflag(fs, "n", "TEST_N", 999, true, validate)
		if err == nil {
			t.Error("ResolveIntWithValidationPflag() want error when all invalid")
		}
		if got != 999 {
			t.Errorf("ResolveIntWithValidationPflag() = %d, want default 999", got)
		}
	})
}

func TestResolvePortPflag(t *testing.T) {
	t.Run("valid port", func(t *testing.T) {
		fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
		fs.Int("port", 8080, "port")
		if err := fs.Parse([]string{"--port", "3000"}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}
		port, err := ResolvePortPflag(fs, "port", "", 8080)
		if err != nil {
			t.Fatalf("ResolvePortPflag() err = %v", err)
		}
		if port != 3000 {
			t.Errorf("ResolvePortPflag() = %d, want 3000", port)
		}
	})
	t.Run("invalid CLI port falls back to default", func(t *testing.T) {
		fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
		fs.Int("port", 8080, "port")
		if err := fs.Parse([]string{"--port", "99999"}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}
		port, err := ResolvePortPflag(fs, "port", "", 8080)
		if err != nil {
			t.Fatalf("ResolvePortPflag() err = %v", err)
		}
		if port != 8080 {
			t.Errorf("ResolvePortPflag(invalid CLI) = %d, want default 8080", port)
		}
	})
	t.Run("invalid port when default also invalid returns error", func(t *testing.T) {
		fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
		fs.Int("port", 8080, "port")
		if err := fs.Parse([]string{"--port", "99999"}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}
		_, err := ResolvePortPflag(fs, "port", "", 99999)
		if err == nil {
			t.Error("ResolvePortPflag() want error when CLI and default both invalid")
		}
	})
	t.Run("port 0 invalid", func(t *testing.T) {
		fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
		fs.Int("port", 8080, "port")
		setEnvPflag(t, "TEST_PORT", "0")
		defer unsetEnvPflag(t, "TEST_PORT")
		if err := fs.Parse([]string{}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}
		port, err := ResolvePortPflag(fs, "port", "TEST_PORT", 8080)
		if err != nil {
			t.Fatalf("ResolvePortPflag() err = %v", err)
		}
		if port != 8080 {
			t.Errorf("ResolvePortPflag(ENV=0) = %d, want default 8080", port)
		}
	})
}

func TestResolveEnumPflag(t *testing.T) {
	t.Run("empty envKey valid CLI", func(t *testing.T) {
		fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
		fs.String("mode", "DEFAULT", "mode")
		if err := fs.Parse([]string{"--mode", "PRIVATE"}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}
		mode, err := ResolveEnumPflag(fs, "mode", "", "DEFAULT", []string{"DEFAULT", "PRIVATE"}, false)
		if err != nil {
			t.Fatalf("ResolveEnumPflag() err = %v", err)
		}
		if mode != "PRIVATE" {
			t.Errorf("ResolveEnumPflag() = %q, want PRIVATE", mode)
		}
	})
	t.Run("case insensitive", func(t *testing.T) {
		fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
		fs.String("mode", "DEFAULT", "mode")
		setEnvPflag(t, "TEST_MODE", "private")
		defer unsetEnvPflag(t, "TEST_MODE")
		if err := fs.Parse([]string{}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}
		mode, err := ResolveEnumPflag(fs, "mode", "TEST_MODE", "DEFAULT", []string{"DEFAULT", "PRIVATE"}, false)
		if err != nil {
			t.Fatalf("ResolveEnumPflag() err = %v", err)
		}
		if mode != "private" {
			t.Errorf("ResolveEnumPflag(caseInsensitive) = %q, want private", mode)
		}
	})
	t.Run("invalid CLI falls back to default", func(t *testing.T) {
		fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
		fs.String("mode", "DEFAULT", "mode")
		if err := fs.Parse([]string{"--mode", "INVALID"}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}
		mode, err := ResolveEnumPflag(fs, "mode", "", "DEFAULT", []string{"DEFAULT", "PRIVATE"}, true)
		if err != nil {
			t.Fatalf("ResolveEnumPflag() err = %v", err)
		}
		if mode != "DEFAULT" {
			t.Errorf("ResolveEnumPflag(invalid CLI) = %q, want default DEFAULT", mode)
		}
	})
	t.Run("invalid value when default also invalid returns error", func(t *testing.T) {
		fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
		fs.String("mode", "DEFAULT", "mode")
		if err := fs.Parse([]string{"--mode", "INVALID"}); err != nil {
			t.Fatalf("fs.Parse() failed: %v", err)
		}
		_, err := ResolveEnumPflag(fs, "mode", "", "INVALID", []string{"DEFAULT", "PRIVATE"}, true)
		if err == nil {
			t.Error("ResolveEnumPflag() want error when CLI and default both invalid")
		}
	})
}
