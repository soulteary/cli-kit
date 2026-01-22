package env

import (
	"os"
	"strconv"
	"strings"
	"time"
)

// Get retrieves an environment variable value, returning defaultValue if not set
func Get(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetTrimmed retrieves a trimmed environment variable value, returning defaultValue if not set or empty
func GetTrimmed(key, defaultValue string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return defaultValue
}

// GetInt retrieves an environment variable as an integer, returning defaultValue if not set or invalid
func GetInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// GetDuration retrieves an environment variable as a duration, returning defaultValue if not set or invalid
func GetDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

// GetBool retrieves an environment variable as a boolean, returning defaultValue if not set or invalid
func GetBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

// GetInt64 retrieves an environment variable as an int64, returning defaultValue if not set or invalid
func GetInt64(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// GetUint retrieves an environment variable as a uint, returning defaultValue if not set or invalid
func GetUint(key string, defaultValue uint) uint {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.ParseUint(value, 10, 0); err == nil {
			return uint(intValue)
		}
	}
	return defaultValue
}

// GetUint64 retrieves an environment variable as a uint64, returning defaultValue if not set or invalid
func GetUint64(key string, defaultValue uint64) uint64 {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.ParseUint(value, 10, 64); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// GetFloat64 retrieves an environment variable as a float64, returning defaultValue if not set or invalid
func GetFloat64(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			return floatValue
		}
	}
	return defaultValue
}

// GetStringSlice retrieves a delimited environment variable as a string slice.
// Returns defaultValue if not set or no valid items found.
func GetStringSlice(key string, defaultValue []string, sep string) []string {
	if sep == "" {
		sep = ","
	}

	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	parts := strings.Split(value, sep)
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		item := strings.TrimSpace(part)
		if item == "" {
			continue
		}
		result = append(result, item)
	}

	if len(result) == 0 {
		return defaultValue
	}

	return result
}
