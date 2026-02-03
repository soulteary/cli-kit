package validator

import (
	"net"
	"testing"
)

func TestValidateURL(t *testing.T) {
	// optsNoResolve disables DNS resolution so tests don't require network
	optsNoResolve := &URLOptions{ResolveHostTimeout: 0}
	tests := []struct {
		name      string
		urlStr    string
		opts      *URLOptions
		wantErr   bool
		errSubstr string
	}{
		// Default options (secure); no resolution to avoid network in unit tests
		{"valid https", "https://example.com", optsNoResolve, false, ""},
		{"valid http", "http://example.com", optsNoResolve, false, ""},
		{"empty URL", "", nil, true, "empty"},
		{"invalid format", "not-a-url", nil, true, "invalid"},
		{"localhost blocked", "http://localhost:8080", nil, true, "localhost"},
		{"127.0.0.1 blocked", "http://127.0.0.1:8080", nil, true, "localhost"},
		{"private IP blocked", "http://192.168.1.1:8080", nil, true, "private"},
		{"file scheme blocked", "file:///etc/passwd", nil, true, "not allowed"},
		{"ftp scheme blocked", "ftp://example.com", nil, true, "not allowed"},

		// With options
		{"localhost allowed", "http://localhost:8080", &URLOptions{AllowLocalhost: true, ResolveHostTimeout: 0}, false, ""},
		{"private IP allowed", "http://192.168.1.1:8080", &URLOptions{AllowPrivateIP: true}, false, ""},
		{"custom schemes", "ftp://example.com", &URLOptions{AllowedSchemes: []string{"ftp", "ftps"}, ResolveHostTimeout: 0}, false, ""},
		{"custom schemes blocked", "http://example.com", &URLOptions{AllowedSchemes: []string{"ftp"}, ResolveHostTimeout: 0}, true, "not allowed"},
		{"localhost + private allowed", "http://192.168.1.1:8080", &URLOptions{AllowLocalhost: true, AllowPrivateIP: true}, false, ""},
		{"empty schemes allowed", "http://example.com", &URLOptions{AllowedSchemes: []string{}, ResolveHostTimeout: 0}, false, ""},
		{"no host", "http://", nil, true, "host"},
		{"IPv6 loopback blocked", "http://[::1]:8080", nil, true, "localhost"},
		{"IPv6 loopback allowed", "http://[::1]:8080", &URLOptions{AllowLocalhost: true}, false, ""},
		{"IPv6 private blocked", "http://[fe80::1]:8080", nil, true, "private"},
		{"IPv6 private allowed", "http://[fe80::1]:8080", &URLOptions{AllowPrivateIP: true}, false, ""},
		{"loopback IP blocked", "http://127.1.1.1:8080", nil, true, "loopback"},
		{"loopback IP allowed", "http://127.1.1.1:8080", &URLOptions{AllowLocalhost: true}, false, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateURL(tt.urlStr, tt.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateURL(%q, %+v) error = %v, wantErr %v", tt.urlStr, tt.opts, err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errSubstr != "" && err != nil {
				if !contains(err.Error(), tt.errSubstr) {
					t.Errorf("ValidateURL(%q, %+v) error = %v, want error containing %q", tt.urlStr, tt.opts, err, tt.errSubstr)
				}
			}
		})
	}
}

func TestIsPrivateIP(t *testing.T) {
	tests := []struct {
		name string
		ip   string
		want bool
	}{
		{"10.0.0.1", "10.0.0.1", true},
		{"10.255.255.255", "10.255.255.255", true},
		{"172.16.0.1", "172.16.0.1", true},
		{"172.31.255.255", "172.31.255.255", true},
		{"172.15.255.255", "172.15.255.255", false}, // Below range
		{"172.32.0.1", "172.32.0.1", false},         // Above range
		{"192.168.1.1", "192.168.1.1", true},
		{"192.168.255.255", "192.168.255.255", true},
		{"192.169.0.1", "192.169.0.1", false}, // Not 192.168.x.x
		{"127.0.0.1", "127.0.0.1", true},
		{"127.255.255.255", "127.255.255.255", true},
		{"public IP", "8.8.8.8", false},
		{"public IP 2", "1.1.1.1", false},
		{"invalid IP", "invalid", false},
		// IPv6 tests
		{"IPv6 loopback", "::1", true},
		{"IPv6 link-local unicast", "fe80::1", true},
		{"IPv6 link-local multicast", "ff02::1", true},
		{"IPv6 public", "2001:4860:4860::8888", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ip := net.ParseIP(tt.ip)
			if ip == nil && tt.want {
				t.Skipf("Skipping %s: invalid IP format", tt.ip)
				return
			}
			if ip == nil {
				return // Invalid IP, expected false
			}
			got := isPrivateIP(ip)
			if got != tt.want {
				t.Errorf("isPrivateIP(%q) = %v, want %v", tt.ip, got, tt.want)
			}
		})
	}
}
