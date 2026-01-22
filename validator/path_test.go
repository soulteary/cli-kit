package validator

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidatePath(t *testing.T) {
	// Get current working directory for tests
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}

	tests := []struct {
		name      string
		path      string
		opts      *PathOptions
		wantErr   bool
		errSubstr string
	}{
		// Default options
		{"valid relative path", "test.txt", nil, false, ""},
		{"valid absolute path", cwd, nil, false, ""},
		{"empty path", "", nil, true, "empty"},
		{"path with ..", "../test.txt", nil, true, "traversal"},
		{"path with multiple ..", "../../etc/passwd", nil, true, "traversal"},

		// With directory restrictions
		{"path in allowed dir", cwd, &PathOptions{AllowedDirs: []string{cwd}}, false, ""},
		{"path outside allowed dir", "/tmp", &PathOptions{AllowedDirs: []string{cwd}}, true, "allowed"},
		{"path in one of allowed dirs", cwd, &PathOptions{AllowedDirs: []string{"/tmp", cwd}}, false, ""},

		// With traversal check disabled
		{"traversal check disabled", "../test.txt", &PathOptions{CheckTraversal: false}, false, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidatePath(tt.path, tt.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePath(%q, %+v) error = %v, wantErr %v", tt.path, tt.opts, err, tt.wantErr)
				return
			}
			if tt.wantErr {
				if tt.errSubstr != "" && err != nil && !contains(err.Error(), tt.errSubstr) {
					t.Errorf("ValidatePath(%q, %+v) error = %v, want error containing %q", tt.path, tt.opts, err, tt.errSubstr)
				}
				return
			}
			// Check that returned path is absolute
			if !filepath.IsAbs(got) {
				t.Errorf("ValidatePath(%q, %+v) = %q, want absolute path", tt.path, tt.opts, got)
			}
		})
	}
}

func TestValidatePath_EdgeCases(t *testing.T) {
	// Test with non-existent path (should still work, just normalize)
	path, err := ValidatePath("./nonexistent.txt", nil)
	if err != nil {
		t.Errorf("ValidatePath() with non-existent path error = %v", err)
	}
	if !filepath.IsAbs(path) {
		t.Errorf("ValidatePath() = %q, want absolute path", path)
	}

	// Test with allowed dir that fails to convert to absolute (should skip it)
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}

	// This test case covers the error handling in AllowedDirs processing
	// where filepath.Abs might fail for an allowed dir
	path, err = ValidatePath(cwd, &PathOptions{
		AllowedDirs: []string{cwd, "/nonexistent/dir/that/might/fail"},
	})
	if err != nil {
		t.Errorf("ValidatePath() with potentially invalid allowed dir error = %v", err)
	}
	if !filepath.IsAbs(path) {
		t.Errorf("ValidatePath() = %q, want absolute path", path)
	}
}
