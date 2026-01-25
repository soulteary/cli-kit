package validator

import (
	"errors"
	"regexp"
	"testing"
)

func TestValidateUsername(t *testing.T) {
	tests := []struct {
		name      string
		username  string
		opts      *UsernameOptions
		wantErr   bool
		errSubstr string
	}{
		// Default style - valid
		{"valid simple", "john", nil, false, ""},
		{"valid with numbers", "john123", nil, false, ""},
		{"valid with underscore", "john_doe", nil, false, ""},
		{"valid with hyphen", "john-doe", nil, false, ""},
		{"valid 3 chars", "abc", nil, false, ""},
		{"valid 32 chars", "abcdefghijklmnopqrstuvwxyz123456", nil, false, ""},
		{"valid uppercase", "JohnDoe", nil, false, ""},
		{"valid mixed case", "JohnDoe123", nil, false, ""},

		// Default style - invalid
		{"empty", "", nil, true, "empty"},
		{"whitespace only", "   ", nil, true, "empty"},
		{"too short 2 chars", "ab", nil, true, "3-32"},
		{"too long 33 chars", "abcdefghijklmnopqrstuvwxyz1234567", nil, true, "3-32"},
		{"starts with number", "1john", nil, true, "starting with a letter"},
		{"starts with underscore", "_john", nil, true, "starting with a letter"},
		{"starts with hyphen", "-john", nil, true, "starting with a letter"},
		{"contains dot", "john.doe", nil, true, "3-32"},
		{"contains space", "john doe", nil, true, "3-32"},
		{"contains special char", "john@doe", nil, true, "3-32"},

		// With AllowEmpty option
		{"empty allowed", "", &UsernameOptions{AllowEmpty: true}, false, ""},
		{"whitespace allowed when empty allowed", "   ", &UsernameOptions{AllowEmpty: true}, false, ""},

		// Simple style
		{"valid simple style", "johndoe", &UsernameOptions{Style: UsernameStyleSimple}, false, ""},
		{"simple with underscore invalid", "john_doe", &UsernameOptions{Style: UsernameStyleSimple}, true, "alphanumeric"},
		{"simple with hyphen invalid", "john-doe", &UsernameOptions{Style: UsernameStyleSimple}, true, "alphanumeric"},

		// Relaxed style
		{"valid relaxed with dot", "john.doe", &UsernameOptions{Style: UsernameStyleRelaxed}, false, ""},
		{"valid relaxed 64 chars", "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijkl", &UsernameOptions{Style: UsernameStyleRelaxed}, false, ""},
		{"relaxed consecutive dots", "john..doe", &UsernameOptions{Style: UsernameStyleRelaxed}, true, "consecutive dots"},
		{"relaxed ends with dot", "john.", &UsernameOptions{Style: UsernameStyleRelaxed}, true, "end with a dot"},
		{"relaxed too long 65 chars", "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklm", &UsernameOptions{Style: UsernameStyleRelaxed}, true, "3-64"},

		// Custom style
		{"custom valid", "test123", &UsernameOptions{
			Style:         UsernameStyleCustom,
			CustomPattern: regexp.MustCompile(`^[a-z]+\d+$`),
			MinLength:     3,
			MaxLength:     10,
		}, false, ""},
		{"custom invalid pattern", "TEST123", &UsernameOptions{
			Style:         UsernameStyleCustom,
			CustomPattern: regexp.MustCompile(`^[a-z]+\d+$`),
			MinLength:     3,
			MaxLength:     10,
		}, true, "does not match"},
		{"custom too short", "ab", &UsernameOptions{
			Style:         UsernameStyleCustom,
			CustomPattern: regexp.MustCompile(`^[a-z]+$`),
			MinLength:     3,
			MaxLength:     10,
		}, true, "at least 3"},
		{"custom too long", "abcdefghijk", &UsernameOptions{
			Style:         UsernameStyleCustom,
			CustomPattern: regexp.MustCompile(`^[a-z]+$`),
			MinLength:     3,
			MaxLength:     10,
		}, true, "at most 10"},
		{"custom nil pattern", "test", &UsernameOptions{Style: UsernameStyleCustom}, true, "custom pattern is required"},

		// Reserved names
		{"reserved name admin", "admin", &UsernameOptions{ReservedNames: []string{"admin", "root", "system"}}, true, "reserved"},
		{"reserved name case insensitive", "ADMIN", &UsernameOptions{ReservedNames: []string{"admin", "root", "system"}}, true, "reserved"},
		{"non-reserved name", "john", &UsernameOptions{ReservedNames: []string{"admin", "root", "system"}}, false, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUsername(tt.username, tt.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateUsername(%q, %+v) error = %v, wantErr %v", tt.username, tt.opts, err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errSubstr != "" && err != nil {
				if !contains(err.Error(), tt.errSubstr) {
					t.Errorf("ValidateUsername(%q, %+v) error = %v, want error containing %q", tt.username, tt.opts, err, tt.errSubstr)
				}
			}
		})
	}
}

func TestValidateUsernameSimple(t *testing.T) {
	tests := []struct {
		name     string
		username string
		wantErr  bool
	}{
		{"valid alphanumeric", "johndoe", false},
		{"valid with numbers", "john123", false},
		{"valid uppercase", "JohnDoe", false},
		{"invalid with underscore", "john_doe", true},
		{"invalid with hyphen", "john-doe", true},
		{"invalid with dot", "john.doe", true},
		{"empty", "", true},
		{"too short", "ab", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUsernameSimple(tt.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateUsernameSimple(%q) error = %v, wantErr %v", tt.username, err, tt.wantErr)
			}
		})
	}
}

func TestValidateUsernameRelaxed(t *testing.T) {
	tests := []struct {
		name     string
		username string
		wantErr  bool
	}{
		{"valid with dot", "john.doe", false},
		{"valid with underscore", "john_doe", false},
		{"valid with hyphen", "john-doe", false},
		{"valid mixed", "john.doe_test-123", false},
		{"invalid consecutive dots", "john..doe", true},
		{"invalid ends with dot", "john.doe.", true},
		{"empty", "", true},
		{"too short", "ab", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUsernameRelaxed(tt.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateUsernameRelaxed(%q) error = %v, wantErr %v", tt.username, err, tt.wantErr)
			}
		})
	}
}

func TestValidateUsernameWithReserved(t *testing.T) {
	reserved := []string{"admin", "root", "system", "moderator"}

	tests := []struct {
		name     string
		username string
		wantErr  bool
	}{
		{"reserved admin", "admin", true},
		{"reserved root", "root", true},
		{"reserved case insensitive", "ADMIN", true},
		{"reserved mixed case", "AdMiN", true},
		{"not reserved", "john", false},
		{"similar to reserved", "administrator", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUsernameWithReserved(tt.username, reserved)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateUsernameWithReserved(%q, %v) error = %v, wantErr %v", tt.username, reserved, err, tt.wantErr)
			}
		})
	}
}

func TestNormalizeUsername(t *testing.T) {
	tests := []struct {
		name     string
		username string
		expect   string
	}{
		{"lowercase", "john", "john"},
		{"uppercase", "JOHN", "john"},
		{"mixed case", "JohnDoe", "johndoe"},
		{"with spaces", "  John  ", "john"},
		{"with tabs", "\tJohn\t", "john"},
		{"empty", "", ""},
		{"only spaces", "   ", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NormalizeUsername(tt.username)
			if got != tt.expect {
				t.Errorf("NormalizeUsername(%q) = %q, want %q", tt.username, got, tt.expect)
			}
		})
	}
}

func TestIsValidUsernameChar(t *testing.T) {
	tests := []struct {
		name string
		r    rune
		want bool
	}{
		{"lowercase letter", 'a', true},
		{"uppercase letter", 'Z', true},
		{"digit", '5', true},
		{"underscore", '_', true},
		{"hyphen", '-', true},
		{"dot", '.', false},
		{"at sign", '@', false},
		{"space", ' ', false},
		{"chinese char", '中', true}, // unicode letter
		{"emoji", '😀', false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsValidUsernameChar(tt.r)
			if got != tt.want {
				t.Errorf("IsValidUsernameChar(%q) = %v, want %v", tt.r, got, tt.want)
			}
		})
	}
}

func TestValidateUsername_TrimSpace(t *testing.T) {
	tests := []struct {
		name     string
		username string
		wantErr  bool
	}{
		{"leading spaces", "  john", false},
		{"trailing spaces", "john  ", false},
		{"both spaces", "  john  ", false},
		{"tabs", "\tjohn\t", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUsername(tt.username, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateUsername(%q, nil) error = %v, wantErr %v", tt.username, err, tt.wantErr)
			}
		})
	}
}

func TestErrInvalidUsername(t *testing.T) {
	err := ValidateUsername("", nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, ErrInvalidUsername) {
		t.Errorf("error should wrap ErrInvalidUsername, got %v", err)
	}
}

func TestValidateUsername_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		username string
		wantErr  bool
	}{
		{"all underscores after first char", "a___", false},
		{"all hyphens after first char", "a---", false},
		{"alternating", "a_b-c_d", false},
		{"numbers only after first char", "a123456789", false},
		{"max length boundary", "abcdefghijklmnopqrstuvwxyz12345", false},  // 31 chars
		{"over max length", "abcdefghijklmnopqrstuvwxyz123456", false},     // 32 chars (exactly at max)
		{"one over max length", "abcdefghijklmnopqrstuvwxyz1234567", true}, // 33 chars
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUsername(tt.username, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateUsername(%q, nil) error = %v, wantErr %v", tt.username, err, tt.wantErr)
			}
		})
	}
}
