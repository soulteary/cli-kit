package validator

import (
	"errors"
	"testing"
)

func TestValidatePositive(t *testing.T) {
	tests := []struct {
		name    string
		value   int
		wantErr bool
	}{
		{"positive value", 1, false},
		{"large positive", 1000000, false},
		{"zero is not positive", 0, true},
		{"negative is not positive", -1, true},
		{"large negative", -1000000, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePositive(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePositive(%d) error = %v, wantErr %v", tt.value, err, tt.wantErr)
			}
			if tt.wantErr && !errors.Is(err, ErrNotPositive) {
				t.Errorf("ValidatePositive(%d) error should wrap ErrNotPositive", tt.value)
			}
		})
	}
}

func TestValidatePositiveInt64(t *testing.T) {
	tests := []struct {
		name    string
		value   int64
		wantErr bool
	}{
		{"positive value", 1, false},
		{"large positive", 9223372036854775807, false},
		{"zero is not positive", 0, true},
		{"negative is not positive", -1, true},
		{"large negative", -9223372036854775807, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePositiveInt64(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePositiveInt64(%d) error = %v, wantErr %v", tt.value, err, tt.wantErr)
			}
			if tt.wantErr && !errors.Is(err, ErrNotPositive) {
				t.Errorf("ValidatePositiveInt64(%d) error should wrap ErrNotPositive", tt.value)
			}
		})
	}
}

func TestValidateNonNegative(t *testing.T) {
	tests := []struct {
		name    string
		value   int
		wantErr bool
	}{
		{"positive value", 1, false},
		{"zero is allowed", 0, false},
		{"large positive", 1000000, false},
		{"negative is not allowed", -1, true},
		{"large negative", -1000000, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateNonNegative(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateNonNegative(%d) error = %v, wantErr %v", tt.value, err, tt.wantErr)
			}
			if tt.wantErr && !errors.Is(err, ErrNegative) {
				t.Errorf("ValidateNonNegative(%d) error should wrap ErrNegative", tt.value)
			}
		})
	}
}

func TestValidateNonNegativeInt64(t *testing.T) {
	tests := []struct {
		name    string
		value   int64
		wantErr bool
	}{
		{"positive value", 1, false},
		{"zero is allowed", 0, false},
		{"large positive", 9223372036854775807, false},
		{"negative is not allowed", -1, true},
		{"large negative", -9223372036854775807, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateNonNegativeInt64(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateNonNegativeInt64(%d) error = %v, wantErr %v", tt.value, err, tt.wantErr)
			}
			if tt.wantErr && !errors.Is(err, ErrNegative) {
				t.Errorf("ValidateNonNegativeInt64(%d) error should wrap ErrNegative", tt.value)
			}
		})
	}
}

func TestValidateInRange(t *testing.T) {
	tests := []struct {
		name    string
		value   int
		min     int
		max     int
		wantErr bool
	}{
		{"within range", 50, 0, 100, false},
		{"at minimum", 0, 0, 100, false},
		{"at maximum", 100, 0, 100, false},
		{"below minimum", -1, 0, 100, true},
		{"above maximum", 101, 0, 100, true},
		{"negative range", -50, -100, -1, false},
		{"below negative range", -101, -100, -1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateInRange(tt.value, tt.min, tt.max)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateInRange(%d, %d, %d) error = %v, wantErr %v", tt.value, tt.min, tt.max, err, tt.wantErr)
			}
		})
	}
}

func TestValidateInRangeInt64(t *testing.T) {
	tests := []struct {
		name    string
		value   int64
		min     int64
		max     int64
		wantErr bool
	}{
		{"within range", 50, 0, 100, false},
		{"at minimum", 0, 0, 100, false},
		{"at maximum", 100, 0, 100, false},
		{"below minimum", -1, 0, 100, true},
		{"above maximum", 101, 0, 100, true},
		{"large values", 9223372036854775806, 9223372036854775805, 9223372036854775807, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateInRangeInt64(tt.value, tt.min, tt.max)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateInRangeInt64(%d, %d, %d) error = %v, wantErr %v", tt.value, tt.min, tt.max, err, tt.wantErr)
			}
		})
	}
}
