package testutil

import (
	"flag"
	"fmt"
)

// NewTestFlagSet creates a new FlagSet with predefined flags for testing
//
// Parameters:
//   - name: Name of the FlagSet
//   - flags: Map of flag names to their default values and types
//     Supported types: string, int, bool, float64
//
// Returns:
//   - *flag.FlagSet: The created FlagSet
//   - error: Returns error if flag type is unsupported
func NewTestFlagSet(name string, flags map[string]interface{}) (*flag.FlagSet, error) {
	fs := flag.NewFlagSet(name, flag.ContinueOnError)

	for flagName, defaultValue := range flags {
		switch v := defaultValue.(type) {
		case string:
			fs.String(flagName, v, fmt.Sprintf("test flag %s", flagName))
		case int:
			fs.Int(flagName, v, fmt.Sprintf("test flag %s", flagName))
		case bool:
			fs.Bool(flagName, v, fmt.Sprintf("test flag %s", flagName))
		case float64:
			fs.Float64(flagName, v, fmt.Sprintf("test flag %s", flagName))
		default:
			return nil, fmt.Errorf("unsupported flag type for %q: %T", flagName, defaultValue)
		}
	}

	return fs, nil
}

// ParseFlags parses command-line arguments with a FlagSet
//
// Parameters:
//   - fs: FlagSet to parse
//   - args: Command-line arguments to parse
//
// Returns:
//   - error: Returns error if parsing fails
func ParseFlags(fs *flag.FlagSet, args []string) error {
	return fs.Parse(args)
}

// MustParseFlags parses command-line arguments and panics on error
//
// Parameters:
//   - fs: FlagSet to parse
//   - args: Command-line arguments to parse
func MustParseFlags(fs *flag.FlagSet, args []string) {
	if err := fs.Parse(args); err != nil {
		panic(fmt.Sprintf("failed to parse flags: %v", err))
	}
}
