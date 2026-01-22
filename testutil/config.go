package testutil

import (
	"flag"
	"reflect"
	"testing"
)

// ConfigTestCase represents a test case for configuration resolution
type ConfigTestCase struct {
	// Name is the test case name
	Name string
	// CLIArgs are command-line arguments to parse
	CLIArgs []string
	// EnvVars are environment variables to set
	EnvVars map[string]string
	// Expected is the expected configuration value
	Expected interface{}
	// WantErr indicates whether an error is expected
	WantErr bool
}

// RunConfigTests runs a table-driven test for configuration resolution
//
// Parameters:
//   - t: Testing.T instance
//   - cases: Test cases to run
//   - resolver: Function that resolves configuration and returns the result
//     The resolver should set up FlagSet, parse CLI args, and resolve config
func RunConfigTests(t *testing.T, cases []ConfigTestCase, resolver func(*flag.FlagSet, map[string]string) (interface{}, error)) {
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			// Set up environment variables
			envMgr := NewEnvManager()
			defer envMgr.Cleanup()

			if err := envMgr.SetMultiple(tc.EnvVars); err != nil {
				t.Fatalf("Failed to set environment variables: %v", err)
			}

			// Create FlagSet and parse CLI args
			fs := flag.NewFlagSet("test", flag.ContinueOnError)

			// Resolve configuration
			got, err := resolver(fs, tc.EnvVars)

			// Check error expectation
			if (err != nil) != tc.WantErr {
				t.Errorf("resolver() error = %v, wantErr %v", err, tc.WantErr)
				return
			}

			// Check value if no error expected
			if !tc.WantErr {
				if !reflect.DeepEqual(got, tc.Expected) {
					t.Errorf("resolver() = %v, want %v", got, tc.Expected)
				}
			}
		})
	}
}
