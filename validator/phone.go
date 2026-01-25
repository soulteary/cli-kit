package validator

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	// phoneRegexCN matches Chinese mainland phone numbers (11 digits starting with 1[3-9])
	phoneRegexCN = regexp.MustCompile(`^1[3-9]\d{9}$`)
	// phoneRegexUS matches US phone numbers (+1 followed by 10 digits, area code 2-9, exchange 2-9)
	phoneRegexUS = regexp.MustCompile(`^\+?1[2-9]\d{2}[2-9]\d{6}$`)
	// phoneRegexUK matches UK phone numbers (+44 followed by 9-10 digits starting with 1-9)
	phoneRegexUK = regexp.MustCompile(`^\+?44[1-9]\d{8,9}$`)
	// phoneRegexInternational matches international phone numbers (general format, 7-15 digits, may include +)
	phoneRegexInternational = regexp.MustCompile(`^\+?[1-9]\d{6,14}$`)
)

// ErrInvalidPhone is returned when a phone number is invalid
var ErrInvalidPhone = fmt.Errorf("invalid phone number format")

// PhoneRegion represents a phone number region/country
type PhoneRegion string

const (
	// PhoneRegionAny allows any valid phone format
	PhoneRegionAny PhoneRegion = "any"
	// PhoneRegionCN requires Chinese mainland phone format
	PhoneRegionCN PhoneRegion = "cn"
	// PhoneRegionUS requires US phone format
	PhoneRegionUS PhoneRegion = "us"
	// PhoneRegionUK requires UK phone format
	PhoneRegionUK PhoneRegion = "uk"
	// PhoneRegionInternational allows international phone format
	PhoneRegionInternational PhoneRegion = "international"
)

// PhoneOptions configures phone validation behavior
type PhoneOptions struct {
	// AllowEmpty allows empty phone numbers (default: false)
	AllowEmpty bool
	// Region specifies required phone format (default: PhoneRegionAny)
	Region PhoneRegion
}

// defaultPhoneOptions returns default phone validation options
func defaultPhoneOptions() *PhoneOptions {
	return &PhoneOptions{
		AllowEmpty: false,
		Region:     PhoneRegionAny,
	}
}

// ValidatePhone validates a phone number string
//
// This function performs validation on phone numbers, supporting:
// - Chinese mainland format (11 digits starting with 1[3-9])
// - US format (+1 followed by 10 digits)
// - UK format (+44 followed by 9-10 digits)
// - International format (7-15 digits with optional +)
//
// Parameters:
//   - phone: Phone number string to validate
//   - opts: Optional validation options (nil uses defaults)
//
// Returns:
//   - error: Returns error if phone number is invalid; otherwise returns nil
func ValidatePhone(phone string, opts *PhoneOptions) error {
	// Use default options if not provided
	if opts == nil {
		opts = defaultPhoneOptions()
	}

	// Trim whitespace
	phone = strings.TrimSpace(phone)

	// Check empty
	if phone == "" {
		if opts.AllowEmpty {
			return nil
		}
		return fmt.Errorf("%w: phone number cannot be empty", ErrInvalidPhone)
	}

	// Validate based on region
	switch opts.Region {
	case PhoneRegionCN:
		if !phoneRegexCN.MatchString(phone) {
			return fmt.Errorf("%w: expected Chinese mainland format (e.g., 13800138000)", ErrInvalidPhone)
		}
	case PhoneRegionUS:
		if !phoneRegexUS.MatchString(phone) {
			return fmt.Errorf("%w: expected US format (e.g., +12025551234)", ErrInvalidPhone)
		}
	case PhoneRegionUK:
		if !phoneRegexUK.MatchString(phone) {
			return fmt.Errorf("%w: expected UK format (e.g., +447911123456)", ErrInvalidPhone)
		}
	case PhoneRegionInternational:
		if !phoneRegexInternational.MatchString(phone) {
			return fmt.Errorf("%w: expected international format (7-15 digits)", ErrInvalidPhone)
		}
	case PhoneRegionAny:
		fallthrough
	default:
		// Try all formats
		if !isValidPhoneAny(phone) {
			return fmt.Errorf("%w: %q does not match any known phone format", ErrInvalidPhone, phone)
		}
	}

	return nil
}

// isValidPhoneAny checks if phone matches any supported format
func isValidPhoneAny(phone string) bool {
	return phoneRegexCN.MatchString(phone) ||
		phoneRegexUS.MatchString(phone) ||
		phoneRegexUK.MatchString(phone) ||
		phoneRegexInternational.MatchString(phone)
}

// ValidatePhoneCN validates a Chinese mainland phone number
// Convenience function for ValidatePhone with Region: PhoneRegionCN
func ValidatePhoneCN(phone string) error {
	return ValidatePhone(phone, &PhoneOptions{Region: PhoneRegionCN})
}

// ValidatePhoneUS validates a US phone number
// Convenience function for ValidatePhone with Region: PhoneRegionUS
func ValidatePhoneUS(phone string) error {
	return ValidatePhone(phone, &PhoneOptions{Region: PhoneRegionUS})
}

// ValidatePhoneUK validates a UK phone number
// Convenience function for ValidatePhone with Region: PhoneRegionUK
func ValidatePhoneUK(phone string) error {
	return ValidatePhone(phone, &PhoneOptions{Region: PhoneRegionUK})
}

// ValidatePhoneInternational validates an international phone number
// Convenience function for ValidatePhone with Region: PhoneRegionInternational
func ValidatePhoneInternational(phone string) error {
	return ValidatePhone(phone, &PhoneOptions{Region: PhoneRegionInternational})
}
