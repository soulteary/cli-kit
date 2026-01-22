package validator

import (
	"fmt"
	"path/filepath"
	"strings"
)

// PathOptions configures path validation behavior
type PathOptions struct {
	// AllowRelative allows relative paths (default: true)
	AllowRelative bool
	// AllowedDirs restricts paths to specific directories (default: empty, no restriction)
	AllowedDirs []string
	// CheckTraversal checks for path traversal attacks (default: true)
	CheckTraversal bool
}

// defaultPathOptions returns default path validation options
func defaultPathOptions() *PathOptions {
	return &PathOptions{
		AllowRelative:  true,
		AllowedDirs:    nil,
		CheckTraversal: true,
	}
}

// ValidatePath validates a file path to prevent path traversal attacks
//
// This function validates file paths, including:
// - Path traversal detection (..)
// - Absolute vs relative path handling
// - Optional directory restrictions
//
// Parameters:
//   - path: File path to validate
//   - opts: Optional validation options (nil uses defaults)
//
// Returns:
//   - string: Normalized absolute path
//   - error: Returns error if path is invalid or has security risks; otherwise returns nil
func ValidatePath(path string, opts *PathOptions) (string, error) {
	if path == "" {
		return "", fmt.Errorf("path cannot be empty")
	}

	// Use default options if not provided
	if opts == nil {
		opts = defaultPathOptions()
	}

	// Check for path traversal in original path before converting to absolute
	if opts.CheckTraversal {
		if strings.Contains(path, "..") {
			return "", fmt.Errorf("path cannot contain path traversal characters (..)")
		}
	}

	// Convert to absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("unable to parse path: %w", err)
	}

	// Check directory restrictions
	if len(opts.AllowedDirs) > 0 {
		allowed := false
		for _, allowedDir := range opts.AllowedDirs {
			allowedAbsDir, err := filepath.Abs(allowedDir)
			if err != nil {
				continue
			}
			if strings.HasPrefix(absPath, allowedAbsDir) {
				allowed = true
				break
			}
		}
		if !allowed {
			return "", fmt.Errorf("path must be under allowed directories: %v", opts.AllowedDirs)
		}
	}

	return absPath, nil
}
