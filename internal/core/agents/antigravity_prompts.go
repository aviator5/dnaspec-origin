package agents

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/aviator5/dnaspec/internal/core/config"
	"github.com/aviator5/dnaspec/internal/core/files"
)

// GenerateAntigravityPrompt generates an Antigravity workflow file
func GenerateAntigravityPrompt(sourceName string, prompt config.ProjectPrompt, sourceDir string) error {
	// Generate filename: dnaspec-<source-name>-<prompt-name>.md
	filename := fmt.Sprintf("dnaspec-%s-%s.md", sourceName, prompt.Name)
	outputPath := filepath.Join(".agent", "workflows", filename)

	// Create directory if needed
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	// Read prompt content
	promptPath := filepath.Join(sourceDir, prompt.File)
	promptContent, err := os.ReadFile(promptPath)
	if err != nil {
		return fmt.Errorf("failed to read prompt file %s: %w", promptPath, err)
	}

	// Generate frontmatter and content
	content := generateAntigravityPromptContent(prompt, string(promptContent))

	// Write atomically
	return writeFileAtomic(outputPath, []byte(content))
}

// generateAntigravityPromptContent creates the full content of an Antigravity workflow file
func generateAntigravityPromptContent(prompt config.ProjectPrompt, promptContent string) string {
	var sb strings.Builder

	// Frontmatter
	sb.WriteString("---\n")
	sb.WriteString(fmt.Sprintf("description: %s\n", prompt.Description))
	sb.WriteString("---\n")

	// Managed block with prompt content
	sb.WriteString(files.ManagedBlockStart)
	sb.WriteString("\n")
	sb.WriteString(strings.TrimSpace(promptContent))
	sb.WriteString("\n")
	sb.WriteString(files.ManagedBlockEnd)
	sb.WriteString("\n")

	return sb.String()
}
