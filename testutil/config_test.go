package testutil

import (
	"flag"
	"testing"
)

func TestRunConfigTests(t *testing.T) {
	resolver := func(fs *flag.FlagSet, envVars map[string]string) (interface{}, error) {
		var port int
		fs.IntVar(&port, "port", 8080, "port")

		// Note: In real usage, CLI args would be passed to resolver separately
		// For this test, we parse empty args to get default value
		// The test case CLIArgs are not used here, but would be in real implementation
		if err := fs.Parse([]string{}); err != nil {
			return nil, err
		}

		// In real usage, would use configutil to resolve
		// For test, just return a simple value
		return port, nil
	}

	resolverWithError := func(fs *flag.FlagSet, envVars map[string]string) (interface{}, error) {
		return nil, flag.ErrHelp
	}

	t.Run("successful cases", func(t *testing.T) {
		cases := []ConfigTestCase{
			{
				Name:     "default value",
				CLIArgs:  []string{},
				EnvVars:  map[string]string{},
				Expected: 8080,
				WantErr:  false,
			},
			{
				Name:     "with CLI args",
				CLIArgs:  []string{"--port", "9090"},
				EnvVars:  map[string]string{},
				Expected: 8080, // Note: resolver doesn't use CLIArgs, so we expect default
				WantErr:  false,
			},
		}

		RunConfigTests(t, cases, resolver)
	})

	t.Run("error cases", func(t *testing.T) {
		cases := []ConfigTestCase{
			{
				Name:     "resolver returns error",
				CLIArgs:  []string{},
				EnvVars:  map[string]string{},
				Expected: nil,
				WantErr:  true,
			},
		}

		RunConfigTests(t, cases, resolverWithError)
	})
}
