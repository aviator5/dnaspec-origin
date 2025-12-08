package files

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/aviator5/dnaspec/internal/core/config"
)

// CopyGuidelineFiles copies guideline and prompt files from source to destination
// Preserves the relative path structure from the manifest
// If an error occurs during copying, attempts to rollback by removing any files that were created
func CopyGuidelineFiles(sourceDir, destDir string, guidelines []config.ManifestGuideline, prompts []config.ManifestPrompt) error {
	var copiedFiles []string

	// Helper to rollback on error
	rollback := func() {
		for _, file := range copiedFiles {
			_ = os.Remove(file) // Best effort cleanup, ignore errors
		}
	}

	// Copy guidelines
	for _, g := range guidelines {
		src := filepath.Join(sourceDir, g.File)
		dst := filepath.Join(destDir, g.File)
		if err := copyFile(src, dst); err != nil {
			rollback()
			return fmt.Errorf("failed to copy guideline %s: %w", g.File, err)
		}
		copiedFiles = append(copiedFiles, dst)
	}

	// Copy prompts
	for _, p := range prompts {
		src := filepath.Join(sourceDir, p.File)
		dst := filepath.Join(destDir, p.File)
		if err := copyFile(src, dst); err != nil {
			rollback()
			return fmt.Errorf("failed to copy prompt %s: %w", p.File, err)
		}
		copiedFiles = append(copiedFiles, dst)
	}

	return nil
}

// copyFile copies a single file from src to dst
// Creates parent directories as needed
func copyFile(src, dst string) error {
	// Create destination directory
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Read source file
	data, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("failed to read source file: %w", err)
	}

	// Write destination file
	if err := os.WriteFile(dst, data, 0o644); err != nil {
		return fmt.Errorf("failed to write destination file: %w", err)
	}

	return nil
}
