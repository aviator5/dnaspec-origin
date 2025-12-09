package paths

import (
	"fmt"
	"path/filepath"
	"strings"
)

// MakeRelative converts an absolute path to relative from project root.
// Returns error if path is not under project root.
func MakeRelative(projectRoot, absPath string) (string, error) {
	// Clean both paths
	cleanRoot := filepath.Clean(projectRoot)
	cleanPath := filepath.Clean(absPath)

	// Resolve symlinks to get real paths
	realRoot, err := filepath.EvalSymlinks(cleanRoot)
	if err != nil {
		// If project root symlink can't be resolved, use cleaned path
		realRoot = cleanRoot
	}

	realPath, err := filepath.EvalSymlinks(cleanPath)
	if err != nil {
		// If path doesn't exist yet or symlink can't be resolved, use cleaned path
		realPath = cleanPath
	}

	// Calculate relative path
	relPath, err := filepath.Rel(realRoot, realPath)
	if err != nil {
		return "", fmt.Errorf("cannot make relative path: %w", err)
	}

	// Validate doesn't escape (no leading ..)
	if strings.HasPrefix(relPath, "..") {
		return "", fmt.Errorf("path is outside project root")
	}

	// Normalize: remove ./ prefix if present
	relPath = strings.TrimPrefix(relPath, "./")

	// Handle special case where path equals project root
	if relPath == "." {
		return ".", nil
	}

	return relPath, nil
}

// ResolveRelative converts a relative path to absolute based on project root.
// Returns error if resolved path escapes project root.
func ResolveRelative(projectRoot, relPath string) (string, error) {
	// Validate input is relative
	if filepath.IsAbs(relPath) {
		return "", fmt.Errorf("expected relative path, got absolute: %s", relPath)
	}

	// Join with project root
	absPath := filepath.Join(projectRoot, relPath)

	// Clean the path
	cleanPath := filepath.Clean(absPath)

	// Resolve symlinks if path exists
	realPath, err := filepath.EvalSymlinks(cleanPath)
	if err != nil {
		// Path doesn't exist yet (ok for some operations)
		realPath = cleanPath
	}

	// Validate within project root
	realRoot, err := filepath.EvalSymlinks(filepath.Clean(projectRoot))
	if err != nil {
		realRoot = filepath.Clean(projectRoot)
	}

	// Check if realPath is within realRoot
	if !isWithinPath(realRoot, realPath) {
		return "", fmt.Errorf("path escapes project root")
	}

	return realPath, nil
}

// ValidateLocalPath ensures a path is safe and within project root.
// Works with both absolute and relative paths.
func ValidateLocalPath(projectRoot, path string) error {
	var absPath string
	var err error

	if filepath.IsAbs(path) {
		// Absolute path
		absPath = path
	} else {
		// Relative path - resolve it
		absPath, err = ResolveRelative(projectRoot, path)
		if err != nil {
			return err
		}
	}

	// Clean and resolve symlinks
	cleanPath := filepath.Clean(absPath)
	realPath, err := filepath.EvalSymlinks(cleanPath)
	if err != nil {
		// Path doesn't exist - check the clean path at least
		realPath = cleanPath
	}

	// Validate within project root
	realRoot, err := filepath.EvalSymlinks(filepath.Clean(projectRoot))
	if err != nil {
		realRoot = filepath.Clean(projectRoot)
	}

	if !isWithinPath(realRoot, realPath) {
		return fmt.Errorf("path is outside project directory")
	}

	return nil
}

// IsWithinProject checks if a path (after resolving symlinks) is within project root.
func IsWithinProject(projectRoot, path string) (bool, error) {
	// Clean paths
	cleanRoot := filepath.Clean(projectRoot)
	cleanPath := filepath.Clean(path)

	// Resolve symlinks
	realRoot, err := filepath.EvalSymlinks(cleanRoot)
	if err != nil {
		realRoot = cleanRoot
	}

	realPath, err := filepath.EvalSymlinks(cleanPath)
	if err != nil {
		// If path doesn't exist, we can't determine if it's within project
		return false, fmt.Errorf("cannot resolve path: %w", err)
	}

	return isWithinPath(realRoot, realPath), nil
}

// isWithinPath checks if childPath is within or equal to parentPath.
// Both paths should be cleaned and resolved before calling this function.
func isWithinPath(parentPath, childPath string) bool {
	// Handle exact match
	if parentPath == childPath {
		return true
	}

	// Ensure parent path ends with separator for proper prefix checking
	if !strings.HasSuffix(parentPath, string(filepath.Separator)) {
		parentPath += string(filepath.Separator)
	}

	// Check if child starts with parent
	return strings.HasPrefix(childPath, parentPath)
}
