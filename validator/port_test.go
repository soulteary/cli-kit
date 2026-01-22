package validator

import (
	"testing"
)

func TestValidatePort(t *testing.T) {
	tests := []struct {
		name    string
		port    int
		wantErr bool
	}{
		{"valid port 1", 1, false},
		{"valid port 8080", 8080, false},
		{"valid port 65535", 65535, false},
		{"invalid port 0", 0, true},
		{"invalid port -1", -1, true},
		{"invalid port 65536", 65536, true},
		{"invalid port 99999", 99999, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePort(tt.port)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePort(%d) error = %v, wantErr %v", tt.port, err, tt.wantErr)
			}
		})
	}
}

func TestValidatePortString(t *testing.T) {
	tests := []struct {
		name    string
		portStr string
		want    int
		wantErr bool
	}{
		{"valid port string 1", "1", 1, false},
		{"valid port string 8080", "8080", 8080, false},
		{"valid port string 65535", "65535", 65535, false},
		{"empty string", "", 0, true},
		{"invalid format", "abc", 0, true},
		{"invalid format with letters", "80abc", 0, true},
		{"invalid port 0", "0", 0, true},
		{"invalid port -1", "-1", 0, true},
		{"invalid port 65536", "65536", 0, true},
		{"whitespace", " 8080 ", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidatePortString(tt.portStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePortString(%q) error = %v, wantErr %v", tt.portStr, err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("ValidatePortString(%q) = %d, want %d", tt.portStr, got, tt.want)
			}
		})
	}
}
