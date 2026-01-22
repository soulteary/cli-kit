package env

import (
	"os"
	"testing"
	"time"
)

// setEnv sets an environment variable and panics on error (should never happen in tests)
func setEnv(t *testing.T, key, value string) {
	t.Helper()
	if err := os.Setenv(key, value); err != nil {
		t.Fatalf("os.Setenv(%q, %q) failed: %v", key, value, err)
	}
}

// unsetEnv unsets an environment variable and panics on error (should never happen in tests)
func unsetEnv(t *testing.T, key string) {
	t.Helper()
	if err := os.Unsetenv(key); err != nil {
		t.Fatalf("os.Unsetenv(%q) failed: %v", key, err)
	}
}

func TestGet(t *testing.T) {
	// Test with environment variable set
	setEnv(t, "TEST_KEY", "test_value")
	defer unsetEnv(t, "TEST_KEY")

	if got := Get("TEST_KEY", "default"); got != "test_value" {
		t.Errorf("Get() = %v, want %v", got, "test_value")
	}

	// Test with default value
	if got := Get("NONEXISTENT_KEY", "default"); got != "default" {
		t.Errorf("Get() = %v, want %v", got, "default")
	}
}

func TestGetTrimmed(t *testing.T) {
	// Test with trimmed value
	setEnv(t, "TEST_TRIMMED", "  value  ")
	defer unsetEnv(t, "TEST_TRIMMED")

	if got := GetTrimmed("TEST_TRIMMED", "default"); got != "value" {
		t.Errorf("GetTrimmed() = %v, want %v", got, "value")
	}

	// Test with whitespace only (should return default)
	setEnv(t, "TEST_TRIMMED_EMPTY", "   ")
	defer unsetEnv(t, "TEST_TRIMMED_EMPTY")

	if got := GetTrimmed("TEST_TRIMMED_EMPTY", "default"); got != "default" {
		t.Errorf("GetTrimmed() = %v, want %v", got, "default")
	}

	// Test with default value
	if got := GetTrimmed("NONEXISTENT_TRIMMED", "default"); got != "default" {
		t.Errorf("GetTrimmed() = %v, want %v", got, "default")
	}
}

func TestGetInt(t *testing.T) {
	// Test with valid integer
	setEnv(t, "TEST_INT", "42")
	defer unsetEnv(t, "TEST_INT")

	if got := GetInt("TEST_INT", 0); got != 42 {
		t.Errorf("GetInt() = %v, want %v", got, 42)
	}

	// Test with invalid integer (should return default)
	setEnv(t, "TEST_INVALID", "not_a_number")
	defer unsetEnv(t, "TEST_INVALID")

	if got := GetInt("TEST_INVALID", 10); got != 10 {
		t.Errorf("GetInt() = %v, want %v", got, 10)
	}

	// Test with default value
	if got := GetInt("NONEXISTENT_INT", 5); got != 5 {
		t.Errorf("GetInt() = %v, want %v", got, 5)
	}
}

func TestGetInt64(t *testing.T) {
	// Test with valid int64
	setEnv(t, "TEST_INT64", "922337203685477580")
	defer unsetEnv(t, "TEST_INT64")

	if got := GetInt64("TEST_INT64", 0); got != 922337203685477580 {
		t.Errorf("GetInt64() = %v, want %v", got, int64(922337203685477580))
	}

	// Test with invalid int64 (should return default)
	setEnv(t, "TEST_INVALID_INT64", "not_a_number")
	defer unsetEnv(t, "TEST_INVALID_INT64")

	if got := GetInt64("TEST_INVALID_INT64", 10); got != 10 {
		t.Errorf("GetInt64() = %v, want %v", got, int64(10))
	}

	// Test with default value
	if got := GetInt64("NONEXISTENT_INT64", 5); got != 5 {
		t.Errorf("GetInt64() = %v, want %v", got, int64(5))
	}
}

func TestGetUint(t *testing.T) {
	// Test with valid uint
	setEnv(t, "TEST_UINT", "42")
	defer unsetEnv(t, "TEST_UINT")

	if got := GetUint("TEST_UINT", 0); got != uint(42) {
		t.Errorf("GetUint() = %v, want %v", got, uint(42))
	}

	// Test with invalid uint (should return default)
	setEnv(t, "TEST_INVALID_UINT", "not_a_number")
	defer unsetEnv(t, "TEST_INVALID_UINT")

	if got := GetUint("TEST_INVALID_UINT", 10); got != uint(10) {
		t.Errorf("GetUint() = %v, want %v", got, uint(10))
	}

	// Test with default value
	if got := GetUint("NONEXISTENT_UINT", 5); got != uint(5) {
		t.Errorf("GetUint() = %v, want %v", got, uint(5))
	}
}

func TestGetUint64(t *testing.T) {
	// Test with valid uint64
	setEnv(t, "TEST_UINT64", "184467440737095516")
	defer unsetEnv(t, "TEST_UINT64")

	if got := GetUint64("TEST_UINT64", 0); got != 184467440737095516 {
		t.Errorf("GetUint64() = %v, want %v", got, uint64(184467440737095516))
	}

	// Test with invalid uint64 (should return default)
	setEnv(t, "TEST_INVALID_UINT64", "not_a_number")
	defer unsetEnv(t, "TEST_INVALID_UINT64")

	if got := GetUint64("TEST_INVALID_UINT64", 10); got != 10 {
		t.Errorf("GetUint64() = %v, want %v", got, uint64(10))
	}

	// Test with default value
	if got := GetUint64("NONEXISTENT_UINT64", 5); got != 5 {
		t.Errorf("GetUint64() = %v, want %v", got, uint64(5))
	}
}

func TestGetDuration(t *testing.T) {
	// Test with valid duration
	setEnv(t, "TEST_DURATION", "5m")
	defer unsetEnv(t, "TEST_DURATION")

	expected := 5 * time.Minute
	if got := GetDuration("TEST_DURATION", time.Second); got != expected {
		t.Errorf("GetDuration() = %v, want %v", got, expected)
	}

	// Test with invalid duration (should return default)
	setEnv(t, "TEST_INVALID_DURATION", "invalid")
	defer unsetEnv(t, "TEST_INVALID_DURATION")

	defaultDuration := 10 * time.Second
	if got := GetDuration("TEST_INVALID_DURATION", defaultDuration); got != defaultDuration {
		t.Errorf("GetDuration() = %v, want %v", got, defaultDuration)
	}

	// Test with default value
	defaultVal := 1 * time.Hour
	if got := GetDuration("NONEXISTENT_DURATION", defaultVal); got != defaultVal {
		t.Errorf("GetDuration() = %v, want %v", got, defaultVal)
	}
}

func TestGetFloat64(t *testing.T) {
	// Test with valid float
	setEnv(t, "TEST_FLOAT64", "3.14")
	defer unsetEnv(t, "TEST_FLOAT64")

	if got := GetFloat64("TEST_FLOAT64", 0); got != 3.14 {
		t.Errorf("GetFloat64() = %v, want %v", got, 3.14)
	}

	// Test with invalid float (should return default)
	setEnv(t, "TEST_INVALID_FLOAT64", "not_a_number")
	defer unsetEnv(t, "TEST_INVALID_FLOAT64")

	if got := GetFloat64("TEST_INVALID_FLOAT64", 2.71); got != 2.71 {
		t.Errorf("GetFloat64() = %v, want %v", got, 2.71)
	}

	// Test with default value
	if got := GetFloat64("NONEXISTENT_FLOAT64", 1.23); got != 1.23 {
		t.Errorf("GetFloat64() = %v, want %v", got, 1.23)
	}
}

func TestGetBool(t *testing.T) {
	// Test with true
	setEnv(t, "TEST_BOOL", "true")
	defer unsetEnv(t, "TEST_BOOL")

	if got := GetBool("TEST_BOOL", false); got != true {
		t.Errorf("GetBool() = %v, want %v", got, true)
	}

	// Test with false
	setEnv(t, "TEST_BOOL_FALSE", "false")
	defer unsetEnv(t, "TEST_BOOL_FALSE")

	if got := GetBool("TEST_BOOL_FALSE", true); got != false {
		t.Errorf("GetBool() = %v, want %v", got, false)
	}

	// Test with invalid boolean (should return default)
	setEnv(t, "TEST_INVALID_BOOL", "maybe")
	defer unsetEnv(t, "TEST_INVALID_BOOL")

	if got := GetBool("TEST_INVALID_BOOL", false); got != false {
		t.Errorf("GetBool() = %v, want %v", got, false)
	}

	// Test with default value
	if got := GetBool("NONEXISTENT_BOOL", true); got != true {
		t.Errorf("GetBool() = %v, want %v", got, true)
	}
}

func TestGetStringSlice(t *testing.T) {
	defaultValue := []string{"default"}

	// Test with comma separated values
	setEnv(t, "TEST_STRING_SLICE", "a, b, , c")
	defer unsetEnv(t, "TEST_STRING_SLICE")

	got := GetStringSlice("TEST_STRING_SLICE", defaultValue, ",")
	if len(got) != 3 || got[0] != "a" || got[1] != "b" || got[2] != "c" {
		t.Errorf("GetStringSlice() = %v, want %v", got, []string{"a", "b", "c"})
	}

	// Test with custom separator
	setEnv(t, "TEST_STRING_SLICE_CUSTOM", "one|two| three ")
	defer unsetEnv(t, "TEST_STRING_SLICE_CUSTOM")

	got = GetStringSlice("TEST_STRING_SLICE_CUSTOM", defaultValue, "|")
	if len(got) != 3 || got[0] != "one" || got[1] != "two" || got[2] != "three" {
		t.Errorf("GetStringSlice() = %v, want %v", got, []string{"one", "two", "three"})
	}

	// Test with empty value (should return default)
	setEnv(t, "TEST_STRING_SLICE_EMPTY", " , , ")
	defer unsetEnv(t, "TEST_STRING_SLICE_EMPTY")

	got = GetStringSlice("TEST_STRING_SLICE_EMPTY", defaultValue, ",")
	if len(got) != 1 || got[0] != "default" {
		t.Errorf("GetStringSlice() = %v, want %v", got, defaultValue)
	}

	// Test with default value
	got = GetStringSlice("NONEXISTENT_STRING_SLICE", defaultValue, ",")
	if len(got) != 1 || got[0] != "default" {
		t.Errorf("GetStringSlice() = %v, want %v", got, defaultValue)
	}

	// Test with empty separator (should default to comma)
	setEnv(t, "TEST_STRING_SLICE_EMPTY_SEP", "x,y,z")
	defer unsetEnv(t, "TEST_STRING_SLICE_EMPTY_SEP")

	got = GetStringSlice("TEST_STRING_SLICE_EMPTY_SEP", defaultValue, "")
	if len(got) != 3 || got[0] != "x" || got[1] != "y" || got[2] != "z" {
		t.Errorf("GetStringSlice() with empty sep = %v, want %v", got, []string{"x", "y", "z"})
	}
}

func TestLookup(t *testing.T) {
	// Test with environment variable set
	setEnv(t, "TEST_LOOKUP", "test_value")
	defer unsetEnv(t, "TEST_LOOKUP")

	value, ok := Lookup("TEST_LOOKUP")
	if !ok || value != "test_value" {
		t.Errorf("Lookup() = (%v, %v), want (%v, %v)", value, ok, "test_value", true)
	}

	// Test with empty string (should return true, empty string)
	setEnv(t, "TEST_LOOKUP_EMPTY", "")
	defer unsetEnv(t, "TEST_LOOKUP_EMPTY")

	value, ok = Lookup("TEST_LOOKUP_EMPTY")
	if !ok || value != "" {
		t.Errorf("Lookup() with empty value = (%v, %v), want (%v, %v)", value, ok, "", true)
	}

	// Test with non-existent variable
	value, ok = Lookup("NONEXISTENT_LOOKUP")
	if ok || value != "" {
		t.Errorf("Lookup() with non-existent = (%v, %v), want (%v, %v)", value, ok, "", false)
	}
}

func TestHas(t *testing.T) {
	// Test with environment variable set
	setEnv(t, "TEST_HAS", "test_value")
	defer unsetEnv(t, "TEST_HAS")

	if !Has("TEST_HAS") {
		t.Error("Has() should return true when variable is set")
	}

	// Test with empty string (should return true)
	setEnv(t, "TEST_HAS_EMPTY", "")
	defer unsetEnv(t, "TEST_HAS_EMPTY")

	if !Has("TEST_HAS_EMPTY") {
		t.Error("Has() should return true when variable is set to empty string")
	}

	// Test with non-existent variable
	if Has("NONEXISTENT_HAS") {
		t.Error("Has() should return false when variable does not exist")
	}
}
