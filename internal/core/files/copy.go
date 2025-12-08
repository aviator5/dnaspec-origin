package files

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/aviator5/dnaspec/internal/core/config"
)

// CopyGuidelineFiles copies guideline and prompt files from source to destination
// Preserves the relative path structure from the manifest
func CopyGuidelineFiles(sourceDir, destDir string, guidelines []config.ManifestGuideline, prompts []config.ManifestPrompt) error {
	// Copy guidelines
	for _, g := range guidelines {
		src := filepath.Join(sourceDir, g.File)
		dst := filepath.Join(destDir, g.File)
		if err := copyFile(src, dst); err != nil {
			return fmt.Errorf("failed to copy guideline %s: %w", g.File, err)
		}
	}

	// Copy prompts
	for _, p := range prompts {
		src := filepath.Join(sourceDir, p.File)
		dst := filepath.Join(destDir, p.File)
		if err := copyFile(src, dst); err != nil {
			return fmt.Errorf("failed to copy prompt %s: %w", p.File, err)
		}
	}

	return nil
}

// copyFile copies a single file from src to dst
// Creates parent directories as needed
func copyFile(src, dst string) error {
	// Create destination directory
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Read source file
	data, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("failed to read source file: %w", err)
	}

	// Write destination file
	if err := os.WriteFile(dst, data, 0644); err != nil {
		return fmt.Errorf("failed to write destination file: %w", err)
	}

	return nil
}
