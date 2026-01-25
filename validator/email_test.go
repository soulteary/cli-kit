package validator

import (
	"errors"
	"testing"
)

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name      string
		email     string
		opts      *EmailOptions
		wantErr   bool
		errSubstr string
	}{
		// Valid emails
		{"valid simple", "test@example.com", nil, false, ""},
		{"valid with numbers", "test123@example.com", nil, false, ""},
		{"valid with underscore", "test_user@example.com", nil, false, ""},
		{"valid with dot", "test.user@example.com", nil, false, ""},
		{"valid with hyphen", "test-user@example.com", nil, false, ""},
		{"valid uppercase", "TEST@EXAMPLE.COM", nil, false, ""},
		{"valid subdomain", "test@mail.example.com", nil, false, ""},
		{"valid multiple subdomains", "test@mail.sub.example.com", nil, false, ""},
		{"valid multiple dots in local", "first.middle.last@example.com", nil, false, ""},
		{"valid numbers in domain", "user123@example123.com", nil, false, ""},
		{"valid mixed special chars", "user_name-test@example-site.com", nil, false, ""},

		// Invalid emails
		{"empty", "", nil, true, "empty"},
		{"whitespace only", "   ", nil, true, "empty"},
		{"missing @", "testexample.com", nil, true, "does not match"},
		{"multiple @", "test@@example.com", nil, true, "does not match"},
		{"missing domain", "test@", nil, true, "does not match"},
		{"missing local part", "@example.com", nil, true, "does not match"},
		{"local starts with dot", ".test@example.com", nil, true, "does not match"},
		{"local ends with dot", "test.@example.com", nil, true, "does not match"},
		{"consecutive dots", "test..user@example.com", nil, true, "consecutive dots"},
		{"domain starts with dot", "test@.example.com", nil, true, "does not match"},
		{"domain ends with dot", "test@example.com.", nil, true, "does not match"},
		{"missing TLD", "test@example", nil, true, "does not match"},
		{"TLD too short", "test@example.c", nil, true, "does not match"},
		{"contains space", "test user@example.com", nil, true, "does not match"},
		{"only @", "@", nil, true, "does not match"},
		{"consecutive dots in domain", "test@example..com", nil, true, "does not match"},

		// With AllowEmpty option
		{"empty allowed", "", &EmailOptions{AllowEmpty: true}, false, ""},
		{"whitespace allowed when empty allowed", "   ", &EmailOptions{AllowEmpty: true}, false, ""},

		// With AllowedDomains option
		{"valid with allowed domain", "test@example.com", &EmailOptions{AllowedDomains: []string{"example.com"}}, false, ""},
		{"valid with allowed subdomain", "test@mail.example.com", &EmailOptions{AllowedDomains: []string{"example.com"}}, false, ""},
		{"invalid with allowed domain", "test@other.com", &EmailOptions{AllowedDomains: []string{"example.com"}}, true, "not in allowed"},
		{"valid with multiple allowed domains", "test@other.com", &EmailOptions{AllowedDomains: []string{"example.com", "other.com"}}, false, ""},

		// With BlockedDomains option
		{"valid with blocked domain not matching", "test@example.com", &EmailOptions{BlockedDomains: []string{"blocked.com"}}, false, ""},
		{"invalid with blocked domain", "test@blocked.com", &EmailOptions{BlockedDomains: []string{"blocked.com"}}, true, "blocked"},
		{"invalid with blocked subdomain", "test@mail.blocked.com", &EmailOptions{BlockedDomains: []string{"blocked.com"}}, true, "blocked"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEmail(tt.email, tt.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateEmail(%q, %+v) error = %v, wantErr %v", tt.email, tt.opts, err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errSubstr != "" && err != nil {
				if !contains(err.Error(), tt.errSubstr) {
					t.Errorf("ValidateEmail(%q, %+v) error = %v, want error containing %q", tt.email, tt.opts, err, tt.errSubstr)
				}
			}
		})
	}
}

func TestValidateEmailSimple(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{"valid", "test@example.com", false},
		{"empty", "", true},
		{"invalid", "invalid", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEmailSimple(tt.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateEmailSimple(%q) error = %v, wantErr %v", tt.email, err, tt.wantErr)
			}
		})
	}
}

func TestValidateEmailWithDomains(t *testing.T) {
	tests := []struct {
		name           string
		email          string
		allowedDomains []string
		wantErr        bool
	}{
		{"valid with matching domain", "test@example.com", []string{"example.com"}, false},
		{"valid with subdomain", "test@mail.example.com", []string{"example.com"}, false},
		{"invalid with non-matching domain", "test@other.com", []string{"example.com"}, true},
		{"valid with multiple allowed", "test@other.com", []string{"example.com", "other.com"}, false},
		{"case insensitive domain", "test@EXAMPLE.COM", []string{"example.com"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEmailWithDomains(tt.email, tt.allowedDomains)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateEmailWithDomains(%q, %v) error = %v, wantErr %v", tt.email, tt.allowedDomains, err, tt.wantErr)
			}
		})
	}
}

func TestExtractEmailDomain(t *testing.T) {
	tests := []struct {
		name   string
		email  string
		expect string
	}{
		{"valid email", "test@example.com", "example.com"},
		{"valid with subdomain", "test@mail.example.com", "mail.example.com"},
		{"uppercase", "test@EXAMPLE.COM", "EXAMPLE.COM"},
		{"empty", "", ""},
		{"no @", "testexample.com", ""},
		{"multiple @", "test@@example.com", ""},
		{"only @", "@", ""},
		{"missing domain", "test@", ""},
		{"with spaces", "  test@example.com  ", "example.com"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExtractEmailDomain(tt.email)
			if got != tt.expect {
				t.Errorf("ExtractEmailDomain(%q) = %q, want %q", tt.email, got, tt.expect)
			}
		})
	}
}

func TestValidateEmail_TrimSpace(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{"leading spaces", "  test@example.com", false},
		{"trailing spaces", "test@example.com  ", false},
		{"both spaces", "  test@example.com  ", false},
		{"tabs", "\ttest@example.com\t", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEmail(tt.email, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateEmail(%q, nil) error = %v, wantErr %v", tt.email, err, tt.wantErr)
			}
		})
	}
}

func TestValidateEmail_CaseSensitivity(t *testing.T) {
	// Email validation should be case-insensitive
	emails := []string{
		"Test@Example.com",
		"TEST@EXAMPLE.COM",
		"test@example.com",
		"TeSt@ExAmPlE.CoM",
	}

	for _, email := range emails {
		t.Run(email, func(t *testing.T) {
			err := ValidateEmail(email, nil)
			if err != nil {
				t.Errorf("ValidateEmail(%q, nil) should pass, got error: %v", email, err)
			}
		})
	}
}

func TestErrInvalidEmail(t *testing.T) {
	err := ValidateEmail("invalid", nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, ErrInvalidEmail) {
		t.Errorf("error should wrap ErrInvalidEmail, got %v", err)
	}
}

func TestValidateEmail_DomainValidation(t *testing.T) {
	// Test domain-specific edge cases
	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{"valid long TLD", "test@example.museum", false},
		{"valid short TLD", "test@example.io", false},
		{"valid numeric subdomain", "test@123.example.com", false},
		{"valid hyphen in domain", "test@my-example.com", false},
		{"domain with many parts", "test@a.b.c.d.example.com", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEmail(tt.email, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateEmail(%q, nil) error = %v, wantErr %v", tt.email, err, tt.wantErr)
			}
		})
	}
}
