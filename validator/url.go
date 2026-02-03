package validator

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"strings"
	"time"
)

// URLOptions configures URL validation behavior
type URLOptions struct {
	// AllowedSchemes specifies allowed URL schemes (default: ["http", "https"])
	AllowedSchemes []string
	// AllowLocalhost allows localhost and 127.0.0.1 (default: false)
	AllowLocalhost bool
	// AllowPrivateIP allows private IP addresses (default: false)
	AllowPrivateIP bool
	// ResolveHostTimeout enables DNS resolution for hostnames and sets timeout; 0 disables resolution (default: 5s)
	ResolveHostTimeout time.Duration
}

// defaultURLOptions returns default URL validation options
func defaultURLOptions() *URLOptions {
	return &URLOptions{
		AllowedSchemes:     []string{"http", "https"},
		AllowLocalhost:     false,
		AllowPrivateIP:     false,
		ResolveHostTimeout: 5 * time.Second,
	}
}

// ValidateURL validates a URL string with SSRF protection
//
// This function performs strict validation on URLs, including:
// - Scheme validation (default: only http and https)
// - Host validation
// - Private IP address blocking (default: blocked)
// - Localhost blocking (default: blocked)
//
// Parameters:
//   - urlStr: URL string to validate
//   - opts: Optional validation options (nil uses secure defaults)
//
// Returns:
//   - error: Returns error if URL is invalid or has security risks; otherwise returns nil
func ValidateURL(urlStr string, opts *URLOptions) error {
	if urlStr == "" {
		return fmt.Errorf("URL cannot be empty")
	}

	// Use default options if not provided
	if opts == nil {
		opts = defaultURLOptions()
	}

	// Parse URL
	u, err := url.ParseRequestURI(urlStr)
	if err != nil {
		return fmt.Errorf("invalid URL format: %w", err)
	}

	// Validate scheme
	if len(opts.AllowedSchemes) > 0 {
		schemeAllowed := false
		for _, allowedScheme := range opts.AllowedSchemes {
			if strings.EqualFold(u.Scheme, allowedScheme) {
				schemeAllowed = true
				break
			}
		}
		if !schemeAllowed {
			return fmt.Errorf("scheme %q is not allowed, allowed schemes: %v", u.Scheme, opts.AllowedSchemes)
		}
	}

	// Validate host
	host := u.Hostname()
	if host == "" {
		return fmt.Errorf("URL must contain a valid host")
	}

	// Check localhost
	if !opts.AllowLocalhost {
		hostLower := strings.ToLower(host)
		if hostLower == "localhost" || hostLower == "127.0.0.1" || hostLower == "::1" {
			return fmt.Errorf("access to localhost is not allowed")
		}
	}

	// Parse IP address or resolve hostname for SSRF protection
	ip := net.ParseIP(host)
	if ip != nil {
		if err := checkIPAllowed(ip, opts); err != nil {
			return err
		}
		return nil
	}

	// Host is a hostname: resolve and check all resolved IPs when timeout is set
	if opts.ResolveHostTimeout > 0 {
		ctx, cancel := context.WithTimeout(context.Background(), opts.ResolveHostTimeout)
		defer cancel()
		resolver := &net.Resolver{}
		addrs, err := resolver.LookupIPAddr(ctx, host)
		if err != nil {
			return fmt.Errorf("failed to resolve host %q: %w", host, err)
		}
		if len(addrs) == 0 {
			return fmt.Errorf("host %q resolved to no addresses", host)
		}
		for _, ipAddr := range addrs {
			if err := checkIPAllowed(ipAddr.IP, opts); err != nil {
				return err
			}
		}
	}

	return nil
}

// checkIPAllowed returns an error if the IP is not allowed by opts (loopback/private checks).
func checkIPAllowed(ip net.IP, opts *URLOptions) error {
	if !opts.AllowLocalhost && ip.IsLoopback() {
		return fmt.Errorf("access to loopback address is not allowed: %s", ip.String())
	}
	if !opts.AllowPrivateIP && isPrivateIP(ip) {
		if opts.AllowLocalhost && ip.IsLoopback() {
			return nil
		}
		return fmt.Errorf("access to private IP address is not allowed: %s", ip.String())
	}
	return nil
}

// isPrivateIP checks if IP is a private IP address
//
// Private IP address ranges:
// - 10.0.0.0/8 (10.0.0.0 to 10.255.255.255)
// - 172.16.0.0/12 (172.16.0.0 to 172.31.255.255)
// - 192.168.0.0/16 (192.168.0.0 to 192.168.255.255)
// - 127.0.0.0/8 (127.0.0.0 to 127.255.255.255) - loopback address
func isPrivateIP(ip net.IP) bool {
	if ip4 := ip.To4(); ip4 != nil {
		return ip4[0] == 10 ||
			(ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31) ||
			(ip4[0] == 192 && ip4[1] == 168) ||
			ip4[0] == 127
	}
	// IPv6 private address check
	if ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
		return true
	}
	return false
}
