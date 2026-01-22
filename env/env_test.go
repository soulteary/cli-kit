package env

import (
	"os"
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	// Test with environment variable set
	os.Setenv("TEST_KEY", "test_value")
	defer os.Unsetenv("TEST_KEY")

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
	os.Setenv("TEST_TRIMMED", "  value  ")
	defer os.Unsetenv("TEST_TRIMMED")

	if got := GetTrimmed("TEST_TRIMMED", "default"); got != "value" {
		t.Errorf("GetTrimmed() = %v, want %v", got, "value")
	}

	// Test with whitespace only (should return default)
	os.Setenv("TEST_TRIMMED_EMPTY", "   ")
	defer os.Unsetenv("TEST_TRIMMED_EMPTY")

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
	os.Setenv("TEST_INT", "42")
	defer os.Unsetenv("TEST_INT")

	if got := GetInt("TEST_INT", 0); got != 42 {
		t.Errorf("GetInt() = %v, want %v", got, 42)
	}

	// Test with invalid integer (should return default)
	os.Setenv("TEST_INVALID", "not_a_number")
	defer os.Unsetenv("TEST_INVALID")

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
	os.Setenv("TEST_INT64", "922337203685477580")
	defer os.Unsetenv("TEST_INT64")

	if got := GetInt64("TEST_INT64", 0); got != 922337203685477580 {
		t.Errorf("GetInt64() = %v, want %v", got, int64(922337203685477580))
	}

	// Test with invalid int64 (should return default)
	os.Setenv("TEST_INVALID_INT64", "not_a_number")
	defer os.Unsetenv("TEST_INVALID_INT64")

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
	os.Setenv("TEST_UINT", "42")
	defer os.Unsetenv("TEST_UINT")

	if got := GetUint("TEST_UINT", 0); got != uint(42) {
		t.Errorf("GetUint() = %v, want %v", got, uint(42))
	}

	// Test with invalid uint (should return default)
	os.Setenv("TEST_INVALID_UINT", "not_a_number")
	defer os.Unsetenv("TEST_INVALID_UINT")

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
	os.Setenv("TEST_UINT64", "184467440737095516")
	defer os.Unsetenv("TEST_UINT64")

	if got := GetUint64("TEST_UINT64", 0); got != 184467440737095516 {
		t.Errorf("GetUint64() = %v, want %v", got, uint64(184467440737095516))
	}

	// Test with invalid uint64 (should return default)
	os.Setenv("TEST_INVALID_UINT64", "not_a_number")
	defer os.Unsetenv("TEST_INVALID_UINT64")

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
	os.Setenv("TEST_DURATION", "5m")
	defer os.Unsetenv("TEST_DURATION")

	expected := 5 * time.Minute
	if got := GetDuration("TEST_DURATION", time.Second); got != expected {
		t.Errorf("GetDuration() = %v, want %v", got, expected)
	}

	// Test with invalid duration (should return default)
	os.Setenv("TEST_INVALID_DURATION", "invalid")
	defer os.Unsetenv("TEST_INVALID_DURATION")

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
	os.Setenv("TEST_FLOAT64", "3.14")
	defer os.Unsetenv("TEST_FLOAT64")

	if got := GetFloat64("TEST_FLOAT64", 0); got != 3.14 {
		t.Errorf("GetFloat64() = %v, want %v", got, 3.14)
	}

	// Test with invalid float (should return default)
	os.Setenv("TEST_INVALID_FLOAT64", "not_a_number")
	defer os.Unsetenv("TEST_INVALID_FLOAT64")

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
	os.Setenv("TEST_BOOL", "true")
	defer os.Unsetenv("TEST_BOOL")

	if got := GetBool("TEST_BOOL", false); got != true {
		t.Errorf("GetBool() = %v, want %v", got, true)
	}

	// Test with false
	os.Setenv("TEST_BOOL_FALSE", "false")
	defer os.Unsetenv("TEST_BOOL_FALSE")

	if got := GetBool("TEST_BOOL_FALSE", true); got != false {
		t.Errorf("GetBool() = %v, want %v", got, false)
	}

	// Test with invalid boolean (should return default)
	os.Setenv("TEST_INVALID_BOOL", "maybe")
	defer os.Unsetenv("TEST_INVALID_BOOL")

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
	os.Setenv("TEST_STRING_SLICE", "a, b, , c")
	defer os.Unsetenv("TEST_STRING_SLICE")

	got := GetStringSlice("TEST_STRING_SLICE", defaultValue, ",")
	if len(got) != 3 || got[0] != "a" || got[1] != "b" || got[2] != "c" {
		t.Errorf("GetStringSlice() = %v, want %v", got, []string{"a", "b", "c"})
	}

	// Test with custom separator
	os.Setenv("TEST_STRING_SLICE_CUSTOM", "one|two| three ")
	defer os.Unsetenv("TEST_STRING_SLICE_CUSTOM")

	got = GetStringSlice("TEST_STRING_SLICE_CUSTOM", defaultValue, "|")
	if len(got) != 3 || got[0] != "one" || got[1] != "two" || got[2] != "three" {
		t.Errorf("GetStringSlice() = %v, want %v", got, []string{"one", "two", "three"})
	}

	// Test with empty value (should return default)
	os.Setenv("TEST_STRING_SLICE_EMPTY", " , , ")
	defer os.Unsetenv("TEST_STRING_SLICE_EMPTY")

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
	os.Setenv("TEST_STRING_SLICE_EMPTY_SEP", "x,y,z")
	defer os.Unsetenv("TEST_STRING_SLICE_EMPTY_SEP")

	got = GetStringSlice("TEST_STRING_SLICE_EMPTY_SEP", defaultValue, "")
	if len(got) != 3 || got[0] != "x" || got[1] != "y" || got[2] != "z" {
		t.Errorf("GetStringSlice() with empty sep = %v, want %v", got, []string{"x", "y", "z"})
	}
}
