package agents

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/aviator5/dnaspec/internal/core/config"
	"github.com/aviator5/dnaspec/internal/core/files"
)

// GenerateCopilotPrompt generates a GitHub Copilot prompt file
func GenerateCopilotPrompt(sourceName string, prompt config.ProjectPrompt, sourceDir string) error {
	// Generate filename: dnaspec-<source-name>-<prompt-name>.prompt.md
	filename := fmt.Sprintf("dnaspec-%s-%s.prompt.md", sourceName, prompt.Name)
	outputPath := filepath.Join(".github", "prompts", filename)

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
	content := generateCopilotPromptContent(prompt, string(promptContent))

	// Write atomically
	return writeFileAtomic(outputPath, []byte(content))
}

// generateCopilotPromptContent creates the full content of a Copilot prompt file
func generateCopilotPromptContent(prompt config.ProjectPrompt, promptContent string) string {
	var sb strings.Builder

	// Frontmatter
	sb.WriteString("---\n")
	sb.WriteString(fmt.Sprintf("description: %s\n", prompt.Description))
	sb.WriteString("---\n\n")

	// $ARGUMENTS placeholder
	sb.WriteString("$ARGUMENTS\n\n")

	// Managed block with prompt content
	sb.WriteString(files.ManagedBlockStart)
	sb.WriteString("\n")
	sb.WriteString(strings.TrimSpace(promptContent))
	sb.WriteString("\n")
	sb.WriteString(files.ManagedBlockEnd)
	sb.WriteString("\n")

	return sb.String()
}
