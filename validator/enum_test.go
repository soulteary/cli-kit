package validator

import (
	"testing"
)

func TestValidateEnum(t *testing.T) {
	allowedValues := []string{"DEFAULT", "REMOTE_FIRST", "ONLY_LOCAL"}

	tests := []struct {
		name          string
		value         string
		allowedValues []string
		caseSensitive bool
		wantErr       bool
		errSubstr     string
	}{
		// Case sensitive
		{"valid case sensitive", "DEFAULT", allowedValues, true, false, ""},
		{"invalid case sensitive", "default", allowedValues, true, true, "not in allowed"},
		{"invalid value", "INVALID", allowedValues, true, true, "not in allowed"},

		// Case insensitive
		{"valid case insensitive", "DEFAULT", allowedValues, false, false, ""},
		{"valid case insensitive lower", "default", allowedValues, false, false, ""},
		{"valid case insensitive mixed", "ReMoTe_FiRsT", allowedValues, false, false, ""},
		{"invalid value case insensitive", "INVALID", allowedValues, false, true, "not in allowed"},

		// Edge cases
		{"empty value", "", allowedValues, true, true, "empty"},
		{"empty allowed list", "DEFAULT", []string{}, true, true, "empty"},
		{"single allowed value", "ONLY", []string{"ONLY"}, true, false, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEnum(tt.value, tt.allowedValues, tt.caseSensitive)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateEnum(%q, %v, %v) error = %v, wantErr %v", tt.value, tt.allowedValues, tt.caseSensitive, err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errSubstr != "" && err != nil {
				if !contains(err.Error(), tt.errSubstr) {
					t.Errorf("ValidateEnum(%q, %v, %v) error = %v, want error containing %q", tt.value, tt.allowedValues, tt.caseSensitive, err, tt.errSubstr)
				}
			}
		})
	}
}

func TestValidateEnumCaseInsensitive(t *testing.T) {
	allowedValues := []string{"DEFAULT", "ONLY_LOCAL"}

	err := ValidateEnumCaseInsensitive("default", allowedValues)
	if err != nil {
		t.Errorf("ValidateEnumCaseInsensitive() error = %v, want nil", err)
	}

	err = ValidateEnumCaseInsensitive("INVALID", allowedValues)
	if err == nil {
		t.Error("ValidateEnumCaseInsensitive() error = nil, want error")
	}
}

func TestValidateEnumCaseSensitive(t *testing.T) {
	allowedValues := []string{"DEFAULT", "ONLY_LOCAL"}

	err := ValidateEnumCaseSensitive("DEFAULT", allowedValues)
	if err != nil {
		t.Errorf("ValidateEnumCaseSensitive() error = %v, want nil", err)
	}

	err = ValidateEnumCaseSensitive("default", allowedValues)
	if err == nil {
		t.Error("ValidateEnumCaseSensitive() error = nil, want error")
	}
}
