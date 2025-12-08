package agents

import (
	"fmt"
	"os"

	"github.com/aviator5/dnaspec/internal/core/config"
	"github.com/aviator5/dnaspec/internal/core/files"
)

// GenerateClaudeMD generates or updates CLAUDE.md with DNA guideline instructions
// This file has the same content as AGENTS.md but specifically for Claude Code
func GenerateClaudeMD(cfg *config.ProjectConfig) error {
	// Reuse the same content generation as AGENTS.md
	content := generateAgentsMDContent(cfg)

	// Read existing file if it exists
	existingContent, err := os.ReadFile("CLAUDE.md")
	var finalContent string

	switch {
	case err == nil:
		// File exists, replace or append managed block
		finalContent = files.ReplaceManagedBlock(string(existingContent), content)
	case os.IsNotExist(err):
		// File doesn't exist, create new with header
		finalContent = files.CreateFileWithManagedBlock(content)
	default:
		return fmt.Errorf("failed to read CLAUDE.md: %w", err)
	}

	// Write atomically
	return writeFileAtomic("CLAUDE.md", []byte(finalContent))
}
