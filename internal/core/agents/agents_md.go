package agents

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/aviator5/dnaspec/internal/core/config"
	"github.com/aviator5/dnaspec/internal/core/files"
)

// GenerateAgentsMD generates or updates AGENTS.md with DNA guideline instructions
func GenerateAgentsMD(cfg *config.ProjectConfig) error {
	content := generateAgentsMDContent(cfg)

	// Read existing file if it exists
	existingContent, err := os.ReadFile("AGENTS.md")
	var finalContent string

	switch {
	case err == nil:
		// File exists, replace or append managed block
		finalContent = files.ReplaceManagedBlock(string(existingContent), content)
	case os.IsNotExist(err):
		// File doesn't exist, create new with header
		finalContent = files.CreateFileWithManagedBlock(content)
	default:
		return fmt.Errorf("failed to read AGENTS.md: %w", err)
	}

	// Write atomically
	return writeFileAtomic("AGENTS.md", []byte(finalContent))
}

// generateAgentsMDContent creates the managed block content for AGENTS.md
func generateAgentsMDContent(cfg *config.ProjectConfig) string {
	var sb strings.Builder

	sb.WriteString("## DNASpec Instructions\n\n")
	sb.WriteString("The project MUST follow shared DNA (Development Norms & Architecture) ")
	sb.WriteString("guidelines stored in the `@/dnaspec` directory. DNA contains reusable ")
	sb.WriteString("patterns and best practices applicable across different projects.\n\n")
	sb.WriteString("These instructions are for AI assistants working in this project.\n\n")

	if len(cfg.Sources) == 0 {
		sb.WriteString("No DNA sources configured yet. Run 'dnaspec add' to add guidelines.\n\n")
		sb.WriteString("Keep this managed block so 'dnaspec update-agents' can refresh the instructions.\n")
		return sb.String()
	}

	sb.WriteString("When working on the codebase, open and refer to the following DNA guidelines as needed:\n")

	for i := range cfg.Sources {
		source := &cfg.Sources[i]
		for _, guideline := range source.Guidelines {
			// Format: @/dnaspec/<source-name>/<file>
			path := fmt.Sprintf("@/dnaspec/%s/%s", source.Name, guideline.File)
			sb.WriteString(fmt.Sprintf("- `%s` for\n", path))

			// Add applicable scenarios as bullet points
			if len(guideline.ApplicableScenarios) > 0 {
				for _, scenario := range guideline.ApplicableScenarios {
					sb.WriteString(fmt.Sprintf("   * %s\n", scenario))
				}
			} else {
				// Fallback if scenarios somehow missing (should be prevented by validation)
				sb.WriteString(fmt.Sprintf("   * %s\n", guideline.Description))
			}
		}
	}

	sb.WriteString("\nKeep this managed block so 'dnaspec update-agents' can refresh the instructions.\n")

	return sb.String()
}

// writeFileAtomic writes content to file atomically using temp file + rename
func writeFileAtomic(path string, content []byte) error {
	dir := filepath.Dir(path)
	tmpFile, err := os.CreateTemp(dir, ".dnaspec-tmp-*")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()

	// Clean up temp file on error
	defer func() {
		if tmpFile != nil {
			_ = tmpFile.Close()
			_ = os.Remove(tmpPath)
		}
	}()

	// Write content
	if _, err := tmpFile.Write(content); err != nil {
		return fmt.Errorf("failed to write temp file: %w", err)
	}

	// Close before rename
	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("failed to close temp file: %w", err)
	}
	tmpFile = nil

	// Atomic rename
	if err := os.Rename(tmpPath, path); err != nil {
		return fmt.Errorf("failed to rename temp file: %w", err)
	}

	return nil
}
