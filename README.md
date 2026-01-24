# cli-kit

A comprehensive Go library for building command-line applications. This package provides utilities for environment variable management, command-line flag handling, configuration resolution, validation, and testing support.

## Features

- **Environment Variable Management**: Safe and flexible environment variable operations with type conversion
- **Flag Utilities**: Enhanced command-line flag handling with type-safe getters
- **Configuration Resolution**: Priority-based configuration resolution (CLI flags > environment variables > defaults)
- **Validators**: Comprehensive validation for URLs, paths, ports, host:port, and enums with SSRF protection
- **Test Utilities**: Helper functions for testing CLI applications and configuration resolution

## Installation

```bash
go get github.com/soulteary/cli-kit
```

## Usage

### Environment Variables

```go
import "github.com/soulteary/cli-kit/env"

// Check if environment variable exists
if env.Has("PORT") {
    // Variable is set
}

// Get with default value
port := env.Get("PORT", "8080")

// Get typed values
portInt := env.GetInt("PORT", 8080)
timeout := env.GetDuration("TIMEOUT", 5*time.Second)
enabled := env.GetBool("ENABLED", false)

// Get trimmed string (removes leading/trailing whitespace)
value := env.GetTrimmed("CONFIG_PATH", "")

// Lookup (distinguish between not set and empty)
value, ok := env.Lookup("API_KEY")
```

### Flag Utilities

```go
import "github.com/soulteary/cli-kit/flagutil"

fs := flag.NewFlagSet("app", flag.ContinueOnError)
port := fs.Int("port", 8080, "Server port")

// Check if flag was set
if flagutil.HasFlag(fs, "port") {
    // Flag was provided
}

// Get flag value with type conversion
portValue := flagutil.GetInt(fs, "port", 8080)
timeout := flagutil.GetDuration(fs, "timeout", 5*time.Second)

// Check if flag exists in command-line arguments
if flagutil.HasFlagInOSArgs("verbose") {
    // -verbose or --verbose was provided
}

// Read password from file (with security checks)
password, err := flagutil.ReadPasswordFromFile("/path/to/password.txt")
```

### Configuration Resolution

```go
import "github.com/soulteary/cli-kit/configutil"

fs := flag.NewFlagSet("app", flag.ContinueOnError)
portFlag := fs.Int("port", 0, "Server port")

// Resolve with priority: CLI flag > ENV > default
port := configutil.ResolveInt(fs, "port", "PORT", 8080, false)

// Resolve string with validation
url, err := configutil.ResolveStringWithValidation(
    fs, "url", "API_URL", "https://api.example.com",
    true, // trimmed
    func(s string) error {
        return validator.ValidateURL(s, nil)
    },
)

// Resolve enum value
mode, err := configutil.ResolveEnum(
    fs, "mode", "APP_MODE", "production",
    []string{"development", "production", "staging"},
    false, // case-insensitive
)

// Resolve host:port with validation
host, port, err := configutil.ResolveHostPort(
    fs, "addr", "SERVER_ADDR", "localhost:8080",
)
```

### Validators

```go
import "github.com/soulteary/cli-kit/validator"

// Validate URL (with SSRF protection)
err := validator.ValidateURL("https://api.example.com", nil)

// With custom options
opts := &validator.URLOptions{
    AllowedSchemes: []string{"http", "https", "ws", "wss"},
    AllowLocalhost: true,
    AllowPrivateIP: false,
}
err := validator.ValidateURL("http://localhost:8080", opts)

// Validate port
err := validator.ValidatePort(8080) // Valid: 1-65535

// Validate host:port
host, port, err := validator.ValidateHostPort("localhost:8080")

// Validate path
err := validator.ValidatePath("/var/log/app.log")

// Validate enum
err := validator.ValidateEnum("production", 
    []string{"development", "production", "staging"},
    false, // case-insensitive
)
```

### Test Utilities

```go
import (
    "github.com/soulteary/cli-kit/testutil"
    "testing"
)

// Environment variable management in tests
func TestMyFunction(t *testing.T) {
    envMgr := testutil.NewEnvManager()
    defer envMgr.Cleanup()
    
    envMgr.Set("PORT", "8080")
    envMgr.SetMultiple(map[string]string{
        "HOST": "localhost",
        "DEBUG": "true",
    })
    
    // Your test code here
}

// Configuration resolution tests
func TestConfigResolution(t *testing.T) {
    cases := []testutil.ConfigTestCase{
        {
            Name:     "CLI flag takes priority",
            CLIArgs:  []string{"-port", "9090"},
            EnvVars:  map[string]string{"PORT": "8080"},
            Expected: 9090,
        },
        {
            Name:     "ENV used when no CLI flag",
            CLIArgs:  []string{},
            EnvVars:  map[string]string{"PORT": "8080"},
            Expected: 8080,
        },
        {
            Name:     "Default used when neither set",
            CLIArgs:  []string{},
            EnvVars:  map[string]string{},
            Expected: 3000,
        },
    }
    
    resolver := func(fs *flag.FlagSet, envVars map[string]string) (interface{}, error) {
        // Set up flags and parse
        portFlag := fs.Int("port", 0, "Port")
        fs.Parse(tc.CLIArgs)
        
        // Resolve config
        return configutil.ResolveInt(fs, "port", "PORT", 3000, false), nil
    }
    
    testutil.RunConfigTests(t, cases, resolver)
}
```

## Project Structure

```
cli-kit/
├── env/              # Environment variable utilities
├── flagutil/         # Command-line flag utilities
├── configutil/       # Configuration resolution with priority
├── validator/        # Input validation (URL, path, port, etc.)
└── testutil/         # Testing utilities
```

## Security Features

- **SSRF Protection**: URL validator blocks private IPs and localhost by default
- **Path Security**: Path validator prevents directory traversal attacks
- **Safe File Reading**: Password file reading with path validation

## License

See LICENSE file for details.
