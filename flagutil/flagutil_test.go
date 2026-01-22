package flagutil

import (
	"flag"
	"os"
	"testing"
	"time"
)

func TestHasFlag(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	var testFlag string
	var otherFlag string
	fs.StringVar(&testFlag, "test-flag", "", "test flag")
	fs.StringVar(&otherFlag, "other-flag", "", "other flag")

	// Test with flag not set
	if HasFlag(fs, "test-flag") {
		t.Error("HasFlag() should return false when flag is not set")
	}

	// Parse with flag set
	os.Args = []string{"test", "--test-flag", "value"}
	if err := fs.Parse(os.Args[1:]); err != nil {
		t.Fatalf("fs.Parse() failed: %v", err)
	}

	// Test with flag set
	if !HasFlag(fs, "test-flag") {
		t.Error("HasFlag() should return true when flag is set")
	}

	// Test with non-existent flag
	if HasFlag(fs, "nonexistent-flag") {
		t.Error("HasFlag() should return false for non-existent flag")
	}

	// Test with other flag set, querying unset flag (covers f.Name == name false branch)
	fs2 := flag.NewFlagSet("test2", flag.ContinueOnError)
	var flag1, flag2 string
	fs2.StringVar(&flag1, "flag1", "", "flag1")
	fs2.StringVar(&flag2, "flag2", "", "flag2")
	if err := fs2.Parse([]string{"--flag1", "value1"}); err != nil {
		t.Fatalf("fs2.Parse() failed: %v", err)
	}
	if HasFlag(fs2, "flag2") {
		t.Error("HasFlag() should return false when other flag is set but queried flag is not")
	}
}

func TestHasFlagInArgs(t *testing.T) {
	args := []string{
		"--test-flag",
		"value",
		"--with-value=foo",
		"-short",
		"-short=value",
	}

	if !HasFlagInArgs(args, "test-flag") {
		t.Error("HasFlagInArgs() should detect --test-flag")
	}
	if !HasFlagInArgs(args, "with-value") {
		t.Error("HasFlagInArgs() should detect --with-value=foo")
	}
	if !HasFlagInArgs(args, "short") {
		t.Error("HasFlagInArgs() should detect -short and -short=value")
	}
	if HasFlagInArgs(args, "missing") {
		t.Error("HasFlagInArgs() should return false for missing flag")
	}
	if HasFlagInArgs(args, "") {
		t.Error("HasFlagInArgs() should return false for empty flag name")
	}
}

func TestHasFlagInOSArgs(t *testing.T) {
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	os.Args = []string{"test", "--test-flag", "value"}
	if !HasFlagInOSArgs("test-flag") {
		t.Error("HasFlagInOSArgs() should detect flag in os.Args")
	}
	if HasFlagInOSArgs("missing") {
		t.Error("HasFlagInOSArgs() should return false for missing flag")
	}
}

func TestGetFlagValue(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	var testFlag string
	fs.StringVar(&testFlag, "test-flag", "", "test flag")

	if err := fs.Parse([]string{"--test-flag", "value"}); err != nil {
		t.Fatalf("fs.Parse() failed: %v", err)
	}

	value, ok := GetFlagValue(fs, "test-flag")
	if !ok || value != "value" {
		t.Errorf("GetFlagValue() = (%v, %v), want (%v, %v)", value, ok, "value", true)
	}

	_, ok = GetFlagValue(fs, "missing")
	if ok {
		t.Error("GetFlagValue() should return false for missing flag")
	}

	// Test with nil FlagSet
	_, ok = GetFlagValue(nil, "test-flag")
	if ok {
		t.Error("GetFlagValue() should return false for nil FlagSet")
	}

	// Test with empty name
	_, ok = GetFlagValue(fs, "")
	if ok {
		t.Error("GetFlagValue() should return false for empty name")
	}

	// Test with flag defined but not set (found != nil but HasFlag returns false)
	fs2 := flag.NewFlagSet("test2", flag.ContinueOnError)
	var unsetFlag string
	fs2.StringVar(&unsetFlag, "unset-flag", "default", "unset flag")
	if err := fs2.Parse([]string{}); err != nil {
		t.Fatalf("fs2.Parse() failed: %v", err)
	}
	_, ok = GetFlagValue(fs2, "unset-flag")
	if ok {
		t.Error("GetFlagValue() should return false for flag that exists but was not set")
	}
}

func TestGetString(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	fs.String("name", "flag-default", "name")

	if err := fs.Parse([]string{"--name", "alice"}); err != nil {
		t.Fatalf("fs.Parse() failed: %v", err)
	}
	if got := GetString(fs, "name", "fallback"); got != "alice" {
		t.Errorf("GetString() = %v, want %v", got, "alice")
	}

	fs = flag.NewFlagSet("test", flag.ContinueOnError)
	fs.String("name", "flag-default", "name")
	if err := fs.Parse([]string{}); err != nil {
		t.Fatalf("fs.Parse() failed: %v", err)
	}
	if got := GetString(fs, "name", "fallback"); got != "fallback" {
		t.Errorf("GetString() = %v, want %v", got, "fallback")
	}
}

func TestGetInt(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	fs.String("port", "", "port")
	fs.String("bad-port", "", "bad port")

	if err := fs.Parse([]string{"--port", "8080", "--bad-port", "not-a-number"}); err != nil {
		t.Fatalf("fs.Parse() failed: %v", err)
	}

	if got := GetInt(fs, "port", 80); got != 8080 {
		t.Errorf("GetInt() = %v, want %v", got, 8080)
	}
	if got := GetInt(fs, "bad-port", 80); got != 80 {
		t.Errorf("GetInt() = %v, want %v", got, 80)
	}

	fs = flag.NewFlagSet("test", flag.ContinueOnError)
	fs.Int("port", 9090, "port")
	if err := fs.Parse([]string{}); err != nil {
		t.Fatalf("fs.Parse() failed: %v", err)
	}
	if got := GetInt(fs, "port", 80); got != 80 {
		t.Errorf("GetInt() = %v, want %v", got, 80)
	}
}

func TestGetInt64(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	fs.String("limit", "", "limit")
	fs.String("bad-limit", "", "bad limit")
	if err := fs.Parse([]string{"--limit", "922337203685477580", "--bad-limit", "not-a-number"}); err != nil {
		t.Fatalf("fs.Parse() failed: %v", err)
	}

	if got := GetInt64(fs, "limit", 10); got != 922337203685477580 {
		t.Errorf("GetInt64() = %v, want %v", got, int64(922337203685477580))
	}
	if got := GetInt64(fs, "missing", 10); got != 10 {
		t.Errorf("GetInt64() = %v, want %v", got, int64(10))
	}
	// Test parse failure branch
	if got := GetInt64(fs, "bad-limit", 10); got != 10 {
		t.Errorf("GetInt64() with invalid value = %v, want %v", got, int64(10))
	}
}

func TestGetUint(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	fs.String("size", "", "size")
	fs.String("bad-size", "", "bad size")
	if err := fs.Parse([]string{"--size", "42", "--bad-size", "not-a-number"}); err != nil {
		t.Fatalf("fs.Parse() failed: %v", err)
	}

	if got := GetUint(fs, "size", 1); got != uint(42) {
		t.Errorf("GetUint() = %v, want %v", got, uint(42))
	}
	if got := GetUint(fs, "missing", 1); got != uint(1) {
		t.Errorf("GetUint() = %v, want %v", got, uint(1))
	}
	// Test parse failure branch
	if got := GetUint(fs, "bad-size", 1); got != uint(1) {
		t.Errorf("GetUint() with invalid value = %v, want %v", got, uint(1))
	}
}

func TestGetUint64(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	fs.String("total", "", "total")
	fs.String("bad-total", "", "bad total")
	if err := fs.Parse([]string{"--total", "184467440737095516", "--bad-total", "not-a-number"}); err != nil {
		t.Fatalf("fs.Parse() failed: %v", err)
	}

	if got := GetUint64(fs, "total", 5); got != uint64(184467440737095516) {
		t.Errorf("GetUint64() = %v, want %v", got, uint64(184467440737095516))
	}
	if got := GetUint64(fs, "missing", 5); got != uint64(5) {
		t.Errorf("GetUint64() = %v, want %v", got, uint64(5))
	}
	// Test parse failure branch
	if got := GetUint64(fs, "bad-total", 5); got != uint64(5) {
		t.Errorf("GetUint64() with invalid value = %v, want %v", got, uint64(5))
	}
}

func TestGetFloat64(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	fs.String("ratio", "", "ratio")
	fs.String("bad-ratio", "", "bad ratio")
	if err := fs.Parse([]string{"--ratio", "0.75", "--bad-ratio", "not-a-number"}); err != nil {
		t.Fatalf("fs.Parse() failed: %v", err)
	}

	if got := GetFloat64(fs, "ratio", 0.5); got != 0.75 {
		t.Errorf("GetFloat64() = %v, want %v", got, 0.75)
	}
	if got := GetFloat64(fs, "missing", 0.5); got != 0.5 {
		t.Errorf("GetFloat64() = %v, want %v", got, 0.5)
	}
	// Test parse failure branch
	if got := GetFloat64(fs, "bad-ratio", 0.5); got != 0.5 {
		t.Errorf("GetFloat64() with invalid value = %v, want %v", got, 0.5)
	}
}

func TestGetBool(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	fs.String("enabled", "", "enabled")
	fs.String("bad-bool", "", "bad bool")
	if err := fs.Parse([]string{"--enabled", "true", "--bad-bool", "maybe"}); err != nil {
		t.Fatalf("fs.Parse() failed: %v", err)
	}

	if got := GetBool(fs, "enabled", false); got != true {
		t.Errorf("GetBool() = %v, want %v", got, true)
	}
	if got := GetBool(fs, "missing", true); got != true {
		t.Errorf("GetBool() = %v, want %v", got, true)
	}
	// Test parse failure branch
	if got := GetBool(fs, "bad-bool", false); got != false {
		t.Errorf("GetBool() with invalid value = %v, want %v", got, false)
	}
}

func TestGetDuration(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	fs.String("timeout", "", "timeout")
	fs.String("bad-timeout", "", "bad timeout")
	if err := fs.Parse([]string{"--timeout", "5s", "--bad-timeout", "invalid-duration"}); err != nil {
		t.Fatalf("fs.Parse() failed: %v", err)
	}

	if got := GetDuration(fs, "timeout", 3*time.Second); got != 5*time.Second {
		t.Errorf("GetDuration() = %v, want %v", got, 5*time.Second)
	}
	if got := GetDuration(fs, "missing", 3*time.Second); got != 3*time.Second {
		t.Errorf("GetDuration() = %v, want %v", got, 3*time.Second)
	}
	// Test parse failure branch
	if got := GetDuration(fs, "bad-timeout", 3*time.Second); got != 3*time.Second {
		t.Errorf("GetDuration() with invalid value = %v, want %v", got, 3*time.Second)
	}
}

func TestReadPasswordFromFile(t *testing.T) {
	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "test-password-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer func() {
		if err := os.Remove(tmpFile.Name()); err != nil {
			t.Logf("Failed to remove temp file: %v", err)
		}
	}()

	// Write password to file
	testPassword := "  my-secret-password  \n"
	if _, err := tmpFile.WriteString(testPassword); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	// Read password from file
	password, err := ReadPasswordFromFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("ReadPasswordFromFile() error = %v", err)
	}

	// Password should be trimmed
	expected := "my-secret-password"
	if password != expected {
		t.Errorf("ReadPasswordFromFile() = %v, want %v", password, expected)
	}
}

func TestReadPasswordFromFile_Nonexistent(t *testing.T) {
	_, err := ReadPasswordFromFile("/nonexistent/path/to/file.txt")
	if err == nil {
		t.Error("ReadPasswordFromFile() should return error for nonexistent file")
	}
}

func TestReadPasswordFromFile_AbsError(t *testing.T) {
	// Test filepath.Abs error branch by using an invalid path
	// On some systems, very long paths or special characters might cause Abs to fail
	// We'll use a path that might cause issues on certain systems
	originalWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer func() {
		if err := os.Chdir(originalWD); err != nil {
			t.Logf("Failed to restore working directory: %v", err)
		}
	}()

	// Create a temp directory and then remove it to trigger potential Abs issues
	tmpDir, err := os.MkdirTemp("", "test-abs-error-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Logf("Failed to remove temp dir: %v", err)
		}
	}()

	// Change to the temp directory
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	// Remove the directory while we're in it (on some systems this might cause issues)
	if err := os.RemoveAll(tmpDir); err != nil {
		t.Logf("Failed to remove temp dir: %v", err)
	}

	// Try to read a file with a relative path that might cause Abs to fail
	// Note: This test might not always trigger the error branch depending on OS behavior
	// but it's the best we can do to test the error path
	_, err = ReadPasswordFromFile("nonexistent.txt")
	// We expect an error (either from Abs or ReadFile), but the exact error depends on OS
	if err == nil {
		t.Error("ReadPasswordFromFile() should return error when file doesn't exist")
	}

	// Restore working directory (already handled by defer)
}

func TestReadPasswordFromFile_EmptyFile(t *testing.T) {
	// Create a temporary file with empty content
	tmpFile, err := os.CreateTemp("", "test-password-empty-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer func() {
		if err := os.Remove(tmpFile.Name()); err != nil {
			t.Logf("Failed to remove temp file: %v", err)
		}
	}()

	if err := tmpFile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	// Read password from empty file
	password, err := ReadPasswordFromFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("ReadPasswordFromFile() error = %v", err)
	}

	// Password should be empty string after trimming
	if password != "" {
		t.Errorf("ReadPasswordFromFile() = %q, want empty string", password)
	}
}
