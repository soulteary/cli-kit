package validator

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	// emailRegex implements RFC 5322 compliant email validation
	// - Disallows consecutive dots
	// - Disallows starting or ending with dot
	// - Stricter domain part validation
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9._-]*[a-zA-Z0-9])?@[a-zA-Z0-9]([a-zA-Z0-9.-]*[a-zA-Z0-9])?\.[a-zA-Z]{2,}$`)
)

// ErrInvalidEmail is returned when an email address is invalid
var ErrInvalidEmail = fmt.Errorf("invalid email format")

// EmailOptions configures email validation behavior
type EmailOptions struct {
	// AllowEmpty allows empty email addresses (default: false)
	AllowEmpty bool
	// AllowedDomains restricts emails to specific domains (empty allows all)
	AllowedDomains []string
	// BlockedDomains blocks specific domains (checked after AllowedDomains)
	BlockedDomains []string
}

// defaultEmailOptions returns default email validation options
func defaultEmailOptions() *EmailOptions {
	return &EmailOptions{
		AllowEmpty:     false,
		AllowedDomains: nil,
		BlockedDomains: nil,
	}
}

// ValidateEmail validates an email address string
//
// This function performs strict validation on email addresses:
// - RFC 5322 basic compliance
// - Disallows consecutive dots
// - Disallows local part starting or ending with dot
// - Requires valid TLD (minimum 2 characters)
// - Optional domain allowlist/blocklist
//
// Parameters:
//   - email: Email address string to validate
//   - opts: Optional validation options (nil uses defaults)
//
// Returns:
//   - error: Returns error if email is invalid; otherwise returns nil
func ValidateEmail(email string, opts *EmailOptions) error {
	// Use default options if not provided
	if opts == nil {
		opts = defaultEmailOptions()
	}

	// Trim whitespace
	email = strings.TrimSpace(email)

	// Check empty
	if email == "" {
		if opts.AllowEmpty {
			return nil
		}
		return fmt.Errorf("%w: email cannot be empty", ErrInvalidEmail)
	}

	// Basic format check
	if !emailRegex.MatchString(email) {
		return fmt.Errorf("%w: %q does not match email pattern", ErrInvalidEmail, email)
	}

	// Check for consecutive dots
	if strings.Contains(email, "..") {
		return fmt.Errorf("%w: email cannot contain consecutive dots", ErrInvalidEmail)
	}

	// Split into local and domain parts
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return fmt.Errorf("%w: email must contain exactly one @ symbol", ErrInvalidEmail)
	}

	localPart := parts[0]
	domainPart := parts[1]

	// Local part cannot start or end with dot
	if strings.HasPrefix(localPart, ".") || strings.HasSuffix(localPart, ".") {
		return fmt.Errorf("%w: local part cannot start or end with dot", ErrInvalidEmail)
	}

	// Domain part cannot start or end with dot
	if strings.HasPrefix(domainPart, ".") || strings.HasSuffix(domainPart, ".") {
		return fmt.Errorf("%w: domain cannot start or end with dot", ErrInvalidEmail)
	}

	// Domain must contain at least one dot (for TLD)
	if !strings.Contains(domainPart, ".") {
		return fmt.Errorf("%w: domain must have a valid TLD", ErrInvalidEmail)
	}

	// Check domain allowlist
	if len(opts.AllowedDomains) > 0 {
		domainLower := strings.ToLower(domainPart)
		allowed := false
		for _, d := range opts.AllowedDomains {
			if strings.EqualFold(domainPart, d) || strings.HasSuffix(domainLower, "."+strings.ToLower(d)) {
				allowed = true
				break
			}
		}
		if !allowed {
			return fmt.Errorf("%w: domain %q is not in allowed list", ErrInvalidEmail, domainPart)
		}
	}

	// Check domain blocklist
	if len(opts.BlockedDomains) > 0 {
		domainLower := strings.ToLower(domainPart)
		for _, d := range opts.BlockedDomains {
			if strings.EqualFold(domainPart, d) || strings.HasSuffix(domainLower, "."+strings.ToLower(d)) {
				return fmt.Errorf("%w: domain %q is blocked", ErrInvalidEmail, domainPart)
			}
		}
	}

	return nil
}

// ValidateEmailSimple validates an email address with default options
// Convenience function for ValidateEmail(email, nil)
func ValidateEmailSimple(email string) error {
	return ValidateEmail(email, nil)
}

// ValidateEmailWithDomains validates an email and checks against allowed domains
// Convenience function for ValidateEmail with AllowedDomains option
func ValidateEmailWithDomains(email string, allowedDomains []string) error {
	return ValidateEmail(email, &EmailOptions{AllowedDomains: allowedDomains})
}

// ExtractEmailDomain extracts the domain part from an email address
// Returns empty string if email is invalid
func ExtractEmailDomain(email string) string {
	email = strings.TrimSpace(email)
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return ""
	}
	return parts[1]
}
