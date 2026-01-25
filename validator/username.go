package validator

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

var (
	// usernameRegexDefault matches alphanumeric usernames with underscore/hyphen (3-32 chars)
	usernameRegexDefault = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_-]{2,31}$`)
	// usernameRegexSimple matches simple alphanumeric usernames (3-32 chars)
	usernameRegexSimple = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9]{2,31}$`)
	// usernameRegexRelaxed matches relaxed usernames allowing dots (3-64 chars)
	usernameRegexRelaxed = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9._-]{2,63}$`)
)

// ErrInvalidUsername is returned when a username is invalid
var ErrInvalidUsername = fmt.Errorf("invalid username format")

// UsernameStyle represents username validation style
type UsernameStyle string

const (
	// UsernameStyleDefault allows alphanumeric, underscore, hyphen (3-32 chars, starts with letter)
	UsernameStyleDefault UsernameStyle = "default"
	// UsernameStyleSimple allows only alphanumeric (3-32 chars, starts with letter)
	UsernameStyleSimple UsernameStyle = "simple"
	// UsernameStyleRelaxed allows alphanumeric, underscore, hyphen, dot (3-64 chars, starts with letter)
	UsernameStyleRelaxed UsernameStyle = "relaxed"
	// UsernameStyleCustom uses custom regex pattern
	UsernameStyleCustom UsernameStyle = "custom"
)

// UsernameOptions configures username validation behavior
type UsernameOptions struct {
	// AllowEmpty allows empty usernames (default: false)
	AllowEmpty bool
	// Style specifies validation style (default: UsernameStyleDefault)
	Style UsernameStyle
	// MinLength minimum username length (default: 3, only for custom style)
	MinLength int
	// MaxLength maximum username length (default: 32, only for custom style)
	MaxLength int
	// CustomPattern custom regex pattern (only for UsernameStyleCustom)
	CustomPattern *regexp.Regexp
	// ReservedNames list of reserved usernames that are not allowed
	ReservedNames []string
}

// defaultUsernameOptions returns default username validation options
func defaultUsernameOptions() *UsernameOptions {
	return &UsernameOptions{
		AllowEmpty:    false,
		Style:         UsernameStyleDefault,
		MinLength:     3,
		MaxLength:     32,
		CustomPattern: nil,
		ReservedNames: nil,
	}
}

// ValidateUsername validates a username string
//
// This function performs validation on usernames:
// - Must start with a letter
// - Configurable allowed characters based on style
// - Length restrictions
// - Optional reserved name checking
//
// Styles:
//   - Default: alphanumeric + underscore + hyphen (3-32 chars)
//   - Simple: alphanumeric only (3-32 chars)
//   - Relaxed: alphanumeric + underscore + hyphen + dot (3-64 chars)
//   - Custom: user-provided regex pattern
//
// Parameters:
//   - username: Username string to validate
//   - opts: Optional validation options (nil uses defaults)
//
// Returns:
//   - error: Returns error if username is invalid; otherwise returns nil
func ValidateUsername(username string, opts *UsernameOptions) error {
	// Use default options if not provided
	if opts == nil {
		opts = defaultUsernameOptions()
	}

	// Trim whitespace
	username = strings.TrimSpace(username)

	// Check empty
	if username == "" {
		if opts.AllowEmpty {
			return nil
		}
		return fmt.Errorf("%w: username cannot be empty", ErrInvalidUsername)
	}

	// Check reserved names
	if len(opts.ReservedNames) > 0 {
		usernameLower := strings.ToLower(username)
		for _, reserved := range opts.ReservedNames {
			if strings.EqualFold(username, reserved) || usernameLower == strings.ToLower(reserved) {
				return fmt.Errorf("%w: %q is a reserved username", ErrInvalidUsername, username)
			}
		}
	}

	// Validate based on style
	switch opts.Style {
	case UsernameStyleSimple:
		if !usernameRegexSimple.MatchString(username) {
			return fmt.Errorf("%w: must be 3-32 alphanumeric characters starting with a letter", ErrInvalidUsername)
		}
	case UsernameStyleRelaxed:
		if !usernameRegexRelaxed.MatchString(username) {
			return fmt.Errorf("%w: must be 3-64 characters (letters, numbers, underscore, hyphen, dot) starting with a letter", ErrInvalidUsername)
		}
		// Additional check: no consecutive dots
		if strings.Contains(username, "..") {
			return fmt.Errorf("%w: username cannot contain consecutive dots", ErrInvalidUsername)
		}
		// Additional check: cannot end with dot
		if strings.HasSuffix(username, ".") {
			return fmt.Errorf("%w: username cannot end with a dot", ErrInvalidUsername)
		}
	case UsernameStyleCustom:
		if opts.CustomPattern == nil {
			return fmt.Errorf("%w: custom pattern is required for custom style", ErrInvalidUsername)
		}
		// Check length first
		if len(username) < opts.MinLength {
			return fmt.Errorf("%w: username must be at least %d characters", ErrInvalidUsername, opts.MinLength)
		}
		if len(username) > opts.MaxLength {
			return fmt.Errorf("%w: username must be at most %d characters", ErrInvalidUsername, opts.MaxLength)
		}
		if !opts.CustomPattern.MatchString(username) {
			return fmt.Errorf("%w: username does not match required pattern", ErrInvalidUsername)
		}
	case UsernameStyleDefault:
		fallthrough
	default:
		if !usernameRegexDefault.MatchString(username) {
			return fmt.Errorf("%w: must be 3-32 characters (letters, numbers, underscore, hyphen) starting with a letter", ErrInvalidUsername)
		}
	}

	return nil
}

// ValidateUsernameSimple validates a username with simple alphanumeric style
// Convenience function for ValidateUsername with Style: UsernameStyleSimple
func ValidateUsernameSimple(username string) error {
	return ValidateUsername(username, &UsernameOptions{Style: UsernameStyleSimple})
}

// ValidateUsernameRelaxed validates a username with relaxed style
// Convenience function for ValidateUsername with Style: UsernameStyleRelaxed
func ValidateUsernameRelaxed(username string) error {
	return ValidateUsername(username, &UsernameOptions{Style: UsernameStyleRelaxed})
}

// ValidateUsernameWithReserved validates a username and checks against reserved names
// Convenience function for ValidateUsername with ReservedNames option
func ValidateUsernameWithReserved(username string, reservedNames []string) error {
	return ValidateUsername(username, &UsernameOptions{ReservedNames: reservedNames})
}

// NormalizeUsername normalizes a username to lowercase
func NormalizeUsername(username string) string {
	return strings.ToLower(strings.TrimSpace(username))
}

// IsValidUsernameChar checks if a character is valid for username (default style)
func IsValidUsernameChar(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' || r == '-'
}
