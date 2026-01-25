package validator

import (
	"errors"
	"testing"
)

func TestValidatePhone(t *testing.T) {
	tests := []struct {
		name      string
		phone     string
		opts      *PhoneOptions
		wantErr   bool
		errSubstr string
	}{
		// Default options (any region)
		{"valid CN phone", "13800138000", nil, false, ""},
		{"valid CN phone 15x", "15000150000", nil, false, ""},
		{"valid CN phone 18x", "18800188000", nil, false, ""},
		{"valid US phone with plus", "+12025551234", nil, false, ""},
		{"valid US phone without plus", "12025551234", nil, false, ""},
		{"valid UK phone with plus", "+447911123456", nil, false, ""},
		{"valid UK phone without plus", "447911123456", nil, false, ""},
		{"valid international phone", "+8613800138000", nil, false, ""},
		{"valid international without plus", "8613800138000", nil, false, ""},

		// Invalid phones
		{"empty phone", "", nil, true, "empty"},
		{"whitespace only", "   ", nil, true, "empty"},
		{"contains letters", "138abc12345", nil, true, "does not match"},
		{"contains special chars", "138-001-38000", nil, true, "does not match"},
		{"too short", "12345", nil, true, "does not match"},
		{"starts with 0", "013800138000", nil, true, "does not match"},

		// With AllowEmpty option
		{"empty allowed", "", &PhoneOptions{AllowEmpty: true}, false, ""},
		{"whitespace allowed when empty allowed", "   ", &PhoneOptions{AllowEmpty: true}, false, ""},

		// CN region specific
		{"valid CN for CN region", "13800138000", &PhoneOptions{Region: PhoneRegionCN}, false, ""},
		{"invalid CN format for CN region", "+8613800138000", &PhoneOptions{Region: PhoneRegionCN}, true, "Chinese mainland"},
		{"US phone for CN region", "+12025551234", &PhoneOptions{Region: PhoneRegionCN}, true, "Chinese mainland"},

		// US region specific
		{"valid US for US region", "+12025551234", &PhoneOptions{Region: PhoneRegionUS}, false, ""},
		{"valid US without plus for US region", "12025551234", &PhoneOptions{Region: PhoneRegionUS}, false, ""},
		{"CN phone for US region", "13800138000", &PhoneOptions{Region: PhoneRegionUS}, true, "US format"},

		// UK region specific
		{"valid UK for UK region", "+447911123456", &PhoneOptions{Region: PhoneRegionUK}, false, ""},
		{"valid UK without plus for UK region", "447911123456", &PhoneOptions{Region: PhoneRegionUK}, false, ""},
		{"CN phone for UK region", "13800138000", &PhoneOptions{Region: PhoneRegionUK}, true, "UK format"},

		// International region
		{"valid international", "+8613800138000", &PhoneOptions{Region: PhoneRegionInternational}, false, ""},
		{"short international", "1234567", &PhoneOptions{Region: PhoneRegionInternational}, false, ""},
		{"too short for international", "123456", &PhoneOptions{Region: PhoneRegionInternational}, true, "international format"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePhone(tt.phone, tt.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePhone(%q, %+v) error = %v, wantErr %v", tt.phone, tt.opts, err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errSubstr != "" && err != nil {
				if !contains(err.Error(), tt.errSubstr) {
					t.Errorf("ValidatePhone(%q, %+v) error = %v, want error containing %q", tt.phone, tt.opts, err, tt.errSubstr)
				}
			}
		})
	}
}

func TestValidatePhoneCN(t *testing.T) {
	tests := []struct {
		name    string
		phone   string
		wantErr bool
	}{
		{"valid 13x", "13800138000", false},
		{"valid 15x", "15000150000", false},
		{"valid 17x", "17600176000", false},
		{"valid 18x", "18800188000", false},
		{"valid 19x", "19900199000", false},
		{"invalid 12x", "12000120000", true},
		{"invalid 10x", "10000100000", true},
		{"too short", "1380013800", true},
		{"too long", "138001380001", true},
		{"with country code", "+8613800138000", true},
		{"empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePhoneCN(tt.phone)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePhoneCN(%q) error = %v, wantErr %v", tt.phone, err, tt.wantErr)
			}
		})
	}
}

func TestValidatePhoneUS(t *testing.T) {
	tests := []struct {
		name    string
		phone   string
		wantErr bool
	}{
		{"valid with plus", "+12025551234", false},
		{"valid without plus", "12025551234", false},
		{"valid different area code", "+14155551234", false},
		{"invalid area code starts with 0", "+10025551234", true},
		{"invalid area code starts with 1", "+11025551234", true},
		{"too short", "+1202555123", true},
		{"too long", "+120255512345", true},
		{"empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePhoneUS(tt.phone)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePhoneUS(%q) error = %v, wantErr %v", tt.phone, err, tt.wantErr)
			}
		})
	}
}

func TestValidatePhoneUK(t *testing.T) {
	tests := []struct {
		name    string
		phone   string
		wantErr bool
	}{
		{"valid with plus 10 digits", "+447911123456", false},
		{"valid without plus 10 digits", "447911123456", false},
		{"valid 9 digits after 44", "+44791112345", false},
		{"invalid starts with 0", "+440911123456", true},
		{"too short 8 digits", "+4479111234", true},
		{"too long 12 digits", "+44791112345678", true},
		{"empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePhoneUK(tt.phone)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePhoneUK(%q) error = %v, wantErr %v", tt.phone, err, tt.wantErr)
			}
		})
	}
}

func TestValidatePhoneInternational(t *testing.T) {
	tests := []struct {
		name    string
		phone   string
		wantErr bool
	}{
		{"valid with plus", "+8613800138000", false},
		{"valid without plus", "8613800138000", false},
		{"valid 7 digits", "1234567", false},
		{"valid 15 digits", "123456789012345", false},
		{"too short 6 digits", "123456", true},
		{"too long 16 digits", "1234567890123456", true},
		{"starts with 0", "0123456789", true},
		{"empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePhoneInternational(tt.phone)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePhoneInternational(%q) error = %v, wantErr %v", tt.phone, err, tt.wantErr)
			}
		})
	}
}

func TestValidatePhone_TrimSpace(t *testing.T) {
	tests := []struct {
		name    string
		phone   string
		wantErr bool
	}{
		{"leading spaces", "  13800138000", false},
		{"trailing spaces", "13800138000  ", false},
		{"both spaces", "  13800138000  ", false},
		{"tabs", "\t13800138000\t", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePhone(tt.phone, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePhone(%q, nil) error = %v, wantErr %v", tt.phone, err, tt.wantErr)
			}
		})
	}
}

func TestErrInvalidPhone(t *testing.T) {
	err := ValidatePhone("invalid", nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, ErrInvalidPhone) {
		t.Errorf("error should wrap ErrInvalidPhone, got %v", err)
	}
}

func TestIsValidPhoneAny(t *testing.T) {
	tests := []struct {
		phone string
		want  bool
	}{
		{"13800138000", true},
		{"+12025551234", true},
		{"+447911123456", true},
		{"+8613800138000", true},
		{"abc", false},
		{"123", false},
	}

	for _, tt := range tests {
		t.Run(tt.phone, func(t *testing.T) {
			got := isValidPhoneAny(tt.phone)
			if got != tt.want {
				t.Errorf("isValidPhoneAny(%q) = %v, want %v", tt.phone, got, tt.want)
			}
		})
	}
}
