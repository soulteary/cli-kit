package flagutil

import (
	"os"
	"testing"
	"time"

	"github.com/spf13/pflag"
)

func TestHasFlagPflag(t *testing.T) {
	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	fs.IntP("port", "p", 8080, "port")
	fs.BoolP("verbose", "v", false, "verbose")

	os.Args = []string{"test", "--port", "9090"}
	if err := fs.Parse(os.Args[1:]); err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if !HasFlagPflag(fs, "port") {
		t.Error("HasFlagPflag(port) should be true when set")
	}
	if HasFlagPflag(fs, "verbose") {
		t.Error("HasFlagPflag(verbose) should be false when not set")
	}
	if HasFlagPflag(fs, "missing") {
		t.Error("HasFlagPflag(missing) should be false for unknown flag")
	}
	if HasFlagPflag(nil, "port") {
		t.Error("HasFlagPflag(nil, ...) should be false")
	}
}

func TestGetFlagValuePflag(t *testing.T) {
	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	fs.StringP("name", "n", "default", "name")
	if err := fs.Parse([]string{"--name", "foo"}); err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	val, ok := GetFlagValuePflag(fs, "name")
	if !ok || val != "foo" {
		t.Errorf("GetFlagValuePflag(name) = %q, %v; want \"foo\", true", val, ok)
	}
	_, ok = GetFlagValuePflag(fs, "other")
	if ok {
		t.Error("GetFlagValuePflag(other) should be false")
	}
}

func TestGetIntPflag_GetBoolPflag(t *testing.T) {
	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	fs.Int("port", 8080, "port")
	fs.Bool("enable", false, "enable")
	if err := fs.Parse([]string{"--port", "3000", "--enable"}); err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if got := GetIntPflag(fs, "port", 8080); got != 3000 {
		t.Errorf("GetIntPflag(port) = %d; want 3000", got)
	}
	if got := GetBoolPflag(fs, "enable", false); got != true {
		t.Errorf("GetBoolPflag(enable) = %v; want true", got)
	}
	if got := GetIntPflag(fs, "missing", 42); got != 42 {
		t.Errorf("GetIntPflag(missing) = %d; want 42", got)
	}
}

func TestGetDurationPflag(t *testing.T) {
	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	fs.Duration("timeout", 5*time.Second, "timeout")
	if err := fs.Parse([]string{"--timeout", "10s"}); err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if got := GetDurationPflag(fs, "timeout", time.Second); got != 10*time.Second {
		t.Errorf("GetDurationPflag(timeout) = %v; want 10s", got)
	}
	if got := GetDurationPflag(fs, "missing", time.Minute); got != time.Minute {
		t.Errorf("GetDurationPflag(missing) = %v; want 1m", got)
	}
	if got := GetDurationPflag(nil, "timeout", time.Second); got != time.Second {
		t.Errorf("GetDurationPflag(nil fs) = %v; want default", got)
	}
}

func TestGetStringPflag(t *testing.T) {
	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	fs.StringP("name", "n", "default", "name")
	if err := fs.Parse([]string{"--name", "foo"}); err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if got := GetStringPflag(fs, "name", "fallback"); got != "foo" {
		t.Errorf("GetStringPflag(name) = %q; want foo", got)
	}
	if got := GetStringPflag(fs, "other", "fallback"); got != "fallback" {
		t.Errorf("GetStringPflag(other) = %q; want fallback", got)
	}
	if got := GetStringPflag(nil, "name", "fallback"); got != "fallback" {
		t.Errorf("GetStringPflag(nil fs) = %q; want fallback", got)
	}
	val, ok := GetFlagValuePflag(fs, "")
	if ok {
		t.Error("GetFlagValuePflag(empty name) should be false")
	}
	if val != "" {
		t.Errorf("GetFlagValuePflag(empty name) = %q; want empty", val)
	}
}

func TestGetInt64Pflag(t *testing.T) {
	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	fs.Int64("limit", 100, "limit")
	if err := fs.Parse([]string{"--limit", "9223372036854775807"}); err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if got := GetInt64Pflag(fs, "limit", 0); got != 9223372036854775807 {
		t.Errorf("GetInt64Pflag(limit) = %d; want 9223372036854775807", got)
	}
	if got := GetInt64Pflag(fs, "missing", 42); got != 42 {
		t.Errorf("GetInt64Pflag(missing) = %d; want 42", got)
	}
	fs2 := pflag.NewFlagSet("test2", pflag.ContinueOnError)
	fs2.String("bad", "", "bad")
	if err := fs2.Parse([]string{"--bad", "not-a-number"}); err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	if got := GetInt64Pflag(fs2, "bad", 10); got != 10 {
		t.Errorf("GetInt64Pflag(invalid) = %d; want 10", got)
	}
}
