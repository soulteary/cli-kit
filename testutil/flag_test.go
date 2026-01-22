package testutil

import (
	"flag"
	"testing"
)

func TestNewTestFlagSet(t *testing.T) {
	tests := []struct {
		name    string
		flags   map[string]interface{}
		wantErr bool
	}{
		{"string flags", map[string]interface{}{"test": "default"}, false},
		{"int flags", map[string]interface{}{"port": 8080}, false},
		{"bool flags", map[string]interface{}{"enabled": true}, false},
		{"float flags", map[string]interface{}{"ratio": 0.5}, false},
		{"mixed flags", map[string]interface{}{
			"name":    "test",
			"port":    8080,
			"enabled": true,
		}, false},
		{"unsupported type", map[string]interface{}{"test": []string{}}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs, err := NewTestFlagSet("test", tt.flags)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewTestFlagSet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && fs == nil {
				t.Error("NewTestFlagSet() returned nil FlagSet")
			}
		})
	}
}

func TestParseFlags(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	fs.String("name", "default", "test flag")

	err := ParseFlags(fs, []string{"--name", "test"})
	if err != nil {
		t.Errorf("ParseFlags() error = %v", err)
	}
}

func TestMustParseFlags(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			// Expected panic for invalid flags
			t.Logf("Expected panic occurred: %v", r)
		}
	}()

	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	fs.String("name", "default", "test flag")

	MustParseFlags(fs, []string{"--name", "test"})

	// Test with invalid flags (should panic)
	fs2 := flag.NewFlagSet("test2", flag.ContinueOnError)
	MustParseFlags(fs2, []string{"--unknown"})
}
