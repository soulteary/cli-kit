package validator

import (
	"testing"
)

func TestValidateHostPort(t *testing.T) {
	tests := []struct {
		name      string
		addr      string
		wantHost  string
		wantPort  int
		wantErr   bool
		errSubstr string
	}{
		{"valid localhost", "localhost:6379", "localhost", 6379, false, ""},
		{"valid IP", "192.168.1.1:8080", "192.168.1.1", 8080, false, ""},
		{"valid IPv6", "[::1]:8080", "::1", 8080, false, ""},
		{"valid domain", "example.com:443", "example.com", 443, false, ""},
		{"empty address", "", "", 0, true, "empty"},
		{"missing port", "localhost", "", 0, true, "invalid"},
		{"invalid port format", "localhost:abc", "", 0, true, "invalid"},
		{"port out of range", "localhost:65536", "", 0, true, "invalid"},
		{"port zero", "localhost:0", "", 0, true, "invalid"},
		{"negative port", "localhost:-1", "", 0, true, "invalid"},
		{"missing host", ":8080", "", 0, true, "empty"},
		{"multiple colons", "host:port:extra", "", 0, true, "invalid"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHost, gotPort, err := ValidateHostPort(tt.addr)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateHostPort(%q) error = %v, wantErr %v", tt.addr, err, tt.wantErr)
				return
			}
			if tt.wantErr {
				if tt.errSubstr != "" && err != nil && !contains(err.Error(), tt.errSubstr) {
					t.Errorf("ValidateHostPort(%q) error = %v, want error containing %q", tt.addr, err, tt.errSubstr)
				}
				return
			}
			if gotHost != tt.wantHost {
				t.Errorf("ValidateHostPort(%q) host = %v, want %v", tt.addr, gotHost, tt.wantHost)
			}
			if gotPort != tt.wantPort {
				t.Errorf("ValidateHostPort(%q) port = %v, want %v", tt.addr, gotPort, tt.wantPort)
			}
		})
	}
}

func TestParseHostPort(t *testing.T) {
	// ParseHostPort is an alias, test that it works the same
	host, port, err := ParseHostPort("localhost:6379")
	if err != nil {
		t.Errorf("ParseHostPort() error = %v", err)
	}
	if host != "localhost" || port != 6379 {
		t.Errorf("ParseHostPort() = (%q, %d), want (%q, %d)", host, port, "localhost", 6379)
	}
}

func TestValidateHostPortWithDefaults(t *testing.T) {
	tests := []struct {
		name        string
		addr        string
		defaultHost string
		defaultPort int
		wantHost    string
		wantPort    int
		wantErr     bool
	}{
		{"empty with defaults", "", "localhost", 6379, "localhost", 6379, false},
		{"host only", "example.com", "localhost", 6379, "example.com", 6379, false},
		{"full address", "example.com:8080", "localhost", 6379, "example.com", 8080, false},
		{"invalid default port", "", "localhost", 0, "", 0, true},
		{"invalid port in addr", "example.com:99999", "localhost", 6379, "", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHost, gotPort, err := ValidateHostPortWithDefaults(tt.addr, tt.defaultHost, tt.defaultPort)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateHostPortWithDefaults() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if gotHost != tt.wantHost {
					t.Errorf("ValidateHostPortWithDefaults() host = %v, want %v", gotHost, tt.wantHost)
				}
				if gotPort != tt.wantPort {
					t.Errorf("ValidateHostPortWithDefaults() port = %v, want %v", gotPort, tt.wantPort)
				}
			}
		})
	}
}
