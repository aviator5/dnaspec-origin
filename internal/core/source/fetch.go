package source

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/aviator5/dnaspec/internal/core/config"
	"github.com/aviator5/dnaspec/internal/core/git"
	"github.com/aviator5/dnaspec/internal/core/validate"
)

// SourceInfo contains information about a fetched source
type SourceInfo struct {
	Manifest   *config.Manifest
	SourceDir  string
	SourceType string // "git-repo" or "local-path"
	URL        string
	Path       string
	Ref        string
	Commit     string
}

// FetchGitSource clones a git repository and parses its manifest
// Returns source info and a cleanup function
func FetchGitSource(url, ref string) (*SourceInfo, func(), error) {
	// Create temp directory for cloning
	tempDir, cleanup, err := git.CreateTempCloneDir()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create temp directory: %w", err)
	}

	// Clone the repository
	commit, err := git.CloneRepo(url, ref, tempDir)
	if err != nil {
		cleanup()
		return nil, nil, fmt.Errorf("failed to clone repository: %w", err)
	}

	// Parse manifest
	manifestPath := filepath.Join(tempDir, "dnaspec-manifest.yaml")
	manifest, err := config.LoadManifest(manifestPath)
	if err != nil {
		cleanup()
		return nil, nil, fmt.Errorf("failed to load manifest from repository: %w", err)
	}

	// Validate manifest
	validationErrors := validate.ValidateManifest(manifest, tempDir)
	if len(validationErrors) > 0 {
		cleanup()
		return nil, nil, fmt.Errorf("manifest validation failed: %s", validationErrors.Error())
	}

	info := &SourceInfo{
		Manifest:   manifest,
		SourceDir:  tempDir,
		SourceType: "git-repo",
		URL:        url,
		Ref:        ref,
		Commit:     commit,
	}

	return info, cleanup, nil
}

// FetchLocalSource reads a manifest from a local directory
// Does not return a cleanup function (no temporary files created)
func FetchLocalSource(path string) (*SourceInfo, error) {
	// Validate path exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("path does not exist: %s", path)
	}

	// Check if path is a directory
	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("failed to stat path: %w", err)
	}
	if !fileInfo.IsDir() {
		return nil, fmt.Errorf("path is not a directory: %s", path)
	}

	// Convert to absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Parse manifest
	manifestPath := filepath.Join(absPath, "dnaspec-manifest.yaml")
	manifest, err := config.LoadManifest(manifestPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load manifest from directory: %w", err)
	}

	// Validate manifest
	validationErrors := validate.ValidateManifest(manifest, absPath)
	if len(validationErrors) > 0 {
		return nil, fmt.Errorf("manifest validation failed: %s", validationErrors.Error())
	}

	info := &SourceInfo{
		Manifest:   manifest,
		SourceDir:  absPath,
		SourceType: "local-path",
		Path:       absPath,
	}

	return info, nil
}
