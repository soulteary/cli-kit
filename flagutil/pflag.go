package flagutil

import (
	"strconv"
	"time"

	"github.com/spf13/pflag"
)

// HasFlagPflag checks if a command-line flag is set in the given pflag.FlagSet.
// It returns true if the flag was explicitly set (e.g. --port 8080 or -p 8080).
func HasFlagPflag(fs *pflag.FlagSet, name string) bool {
	if fs == nil || name == "" {
		return false
	}
	f := fs.Lookup(name)
	return f != nil && f.Changed
}

// GetFlagValuePflag returns the string value for a flag if it was set.
func GetFlagValuePflag(fs *pflag.FlagSet, name string) (string, bool) {
	if fs == nil || name == "" {
		return "", false
	}
	f := fs.Lookup(name)
	if f == nil || !f.Changed {
		return "", false
	}
	return f.Value.String(), true
}

// GetStringPflag returns flag value or defaultValue when not set.
func GetStringPflag(fs *pflag.FlagSet, name, defaultValue string) string {
	value, ok := GetFlagValuePflag(fs, name)
	if !ok {
		return defaultValue
	}
	return value
}

// GetIntPflag returns flag value as int or defaultValue when not set/invalid.
func GetIntPflag(fs *pflag.FlagSet, name string, defaultValue int) int {
	value, ok := GetFlagValuePflag(fs, name)
	if !ok {
		return defaultValue
	}
	if parsed, err := strconv.Atoi(value); err == nil {
		return parsed
	}
	return defaultValue
}

// GetInt64Pflag returns flag value as int64 or defaultValue when not set/invalid.
func GetInt64Pflag(fs *pflag.FlagSet, name string, defaultValue int64) int64 {
	value, ok := GetFlagValuePflag(fs, name)
	if !ok {
		return defaultValue
	}
	if parsed, err := strconv.ParseInt(value, 10, 64); err == nil {
		return parsed
	}
	return defaultValue
}

// GetBoolPflag returns flag value as bool or defaultValue when not set/invalid.
func GetBoolPflag(fs *pflag.FlagSet, name string, defaultValue bool) bool {
	value, ok := GetFlagValuePflag(fs, name)
	if !ok {
		return defaultValue
	}
	if parsed, err := strconv.ParseBool(value); err == nil {
		return parsed
	}
	return defaultValue
}

// GetDurationPflag returns flag value as time.Duration or defaultValue when not set/invalid.
func GetDurationPflag(fs *pflag.FlagSet, name string, defaultValue time.Duration) time.Duration {
	value, ok := GetFlagValuePflag(fs, name)
	if !ok {
		return defaultValue
	}
	if parsed, err := time.ParseDuration(value); err == nil {
		return parsed
	}
	return defaultValue
}
