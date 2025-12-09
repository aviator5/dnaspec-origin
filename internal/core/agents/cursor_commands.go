package agents

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/aviator5/dnaspec/internal/core/config"
	"github.com/aviator5/dnaspec/internal/core/files"
)

// GenerateCursorCommand generates a Cursor command file
func GenerateCursorCommand(sourceName string, prompt config.ProjectPrompt, sourceDir string) error {
	// Generate filename: dnaspec-<source-name>-<prompt-name>.md
	filename := fmt.Sprintf("dnaspec-%s-%s.md", sourceName, prompt.Name)
	outputPath := filepath.Join(".cursor", "commands", filename)

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
	content := generateCursorCommandContent(sourceName, prompt, string(promptContent))

	// Write atomically
	return writeFileAtomic(outputPath, []byte(content))
}

// generateCursorCommandContent creates the full content of a Cursor command file
func generateCursorCommandContent(sourceName string, prompt config.ProjectPrompt, promptContent string) string {
	var sb strings.Builder

	// Command name: /dnaspec-<source-name>-<prompt-name>
	commandName := fmt.Sprintf("/dnaspec-%s-%s", sourceName, prompt.Name)
	// ID: dnaspec-<source-name>-<prompt-name>
	commandID := fmt.Sprintf("dnaspec-%s-%s", sourceName, prompt.Name)

	// Frontmatter
	sb.WriteString("---\n")
	sb.WriteString(fmt.Sprintf("name: %s\n", commandName))
	sb.WriteString(fmt.Sprintf("id: %s\n", commandID))
	sb.WriteString("category: DNASpec\n")
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
