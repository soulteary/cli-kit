package validator

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
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
// - Path traversal detection (..), including after normalization to resist bypasses
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

	// Security: after normalization, ensure no ".." segment remains (guards against
	// encoding/unicode bypasses or platform quirks that might bypass the string check)
	if opts.CheckTraversal {
		cleaned := filepath.Clean(absPath)
		if containsTraversalSegment(cleaned) {
			return "", fmt.Errorf("path cannot contain path traversal characters (..)")
		}
		absPath = cleaned
	}

	// Check directory restrictions: path must be exactly allowedDir or under it (no prefix bypass)
	if len(opts.AllowedDirs) > 0 {
		allowed := false
		sep := string(filepath.Separator)
		for _, allowedDir := range opts.AllowedDirs {
			allowedAbsDir, err := filepath.Abs(allowedDir)
			if err != nil {
				continue
			}
			allowedAbsDir = filepath.Clean(allowedAbsDir)
			if absPath == allowedAbsDir {
				allowed = true
				break
			}
			prefix := allowedAbsDir + sep
			if strings.HasPrefix(absPath, prefix) {
				allowed = true
				break
			}
		}
		if !allowed {
			// Do not include AllowedDirs in error to avoid leaking allowed paths to callers (e.g. API responses)
			return "", fmt.Errorf("path is not under allowed directories")
		}
	}

	return absPath, nil
}

// containsTraversalSegment returns true if path contains ".." as a path segment.
func containsTraversalSegment(path string) bool {
	for _, part := range strings.Split(path, string(filepath.Separator)) {
		if part == ".." {
			return true
		}
	}
	return false
}

// ErrFileNotFound is returned when a file does not exist
var ErrFileNotFound = fmt.Errorf("file not found")

// ErrNotAFile is returned when the path is not a regular file
var ErrNotAFile = fmt.Errorf("path is not a file")

// ErrDirNotFound is returned when a directory does not exist
var ErrDirNotFound = fmt.Errorf("directory not found")

// ErrNotADirectory is returned when the path is not a directory
var ErrNotADirectory = fmt.Errorf("path is not a directory")

// ErrFileNotReadable is returned when a file cannot be read
var ErrFileNotReadable = fmt.Errorf("file is not readable")

// ErrDirNotWritable is returned when a directory is not writable
var ErrDirNotWritable = fmt.Errorf("directory is not writable")

// ValidateFileExists validates that a file exists at the given path
//
// Parameters:
//   - path: The file path to validate
//
// Returns:
//   - error: Returns ErrFileNotFound if the file doesn't exist, ErrNotAFile if the path is a directory, nil otherwise
func ValidateFileExists(path string) error {
	if path == "" {
		return fmt.Errorf("path cannot be empty")
	}

	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("%w: %s", ErrFileNotFound, path)
		}
		return fmt.Errorf("unable to access path: %w", err)
	}

	if info.IsDir() {
		return fmt.Errorf("%w: %s is a directory", ErrNotAFile, path)
	}

	return nil
}

// ValidateFileReadable validates that a file exists and is readable
//
// Parameters:
//   - path: The file path to validate
//
// Returns:
//   - error: Returns error if file doesn't exist or can't be read, nil otherwise
func ValidateFileReadable(path string) error {
	// First check if file exists
	if err := ValidateFileExists(path); err != nil {
		return err
	}

	// Try to open the file for reading
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrFileNotReadable, path)
	}
	_ = f.Close()

	return nil
}

// ValidateDirExists validates that a directory exists at the given path
//
// Parameters:
//   - path: The directory path to validate
//
// Returns:
//   - error: Returns ErrDirNotFound if the directory doesn't exist, ErrNotADirectory if the path is a file, nil otherwise
func ValidateDirExists(path string) error {
	if path == "" {
		return fmt.Errorf("path cannot be empty")
	}

	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("%w: %s", ErrDirNotFound, path)
		}
		return fmt.Errorf("unable to access path: %w", err)
	}

	if !info.IsDir() {
		return fmt.Errorf("%w: %s is a file", ErrNotADirectory, path)
	}

	return nil
}

// ValidateDirWritable validates that a directory exists and is writable
// It creates a temporary file to verify write permissions
//
// Parameters:
//   - path: The directory path to validate
//
// Returns:
//   - error: Returns error if directory doesn't exist or is not writable, nil otherwise
func ValidateDirWritable(path string) error {
	// First check if directory exists
	if err := ValidateDirExists(path); err != nil {
		return err
	}

	// Try to create a temporary file to verify write permissions
	testFile := filepath.Join(path, ".write_test_"+randomSuffix())
	f, err := os.Create(testFile)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrDirNotWritable, path)
	}
	_ = f.Close()
	_ = os.Remove(testFile)

	return nil
}

// randomSuffix generates a random suffix for write-test filenames to avoid predictability and races.
func randomSuffix() string {
	b := make([]byte, 4)
	if _, err := rand.Read(b); err != nil {
		return fmt.Sprintf("%d", os.Getpid())
	}
	return hex.EncodeToString(b)
}
