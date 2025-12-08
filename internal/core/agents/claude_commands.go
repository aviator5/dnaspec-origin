package agents

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/aviator5/dnaspec/internal/core/config"
	"github.com/aviator5/dnaspec/internal/core/files"
)

// GenerateClaudeCommand generates a Claude slash command for a prompt
func GenerateClaudeCommand(sourceName string, prompt config.ProjectPrompt, sourceDir string) error {
	// Generate filename: <source-name>-<prompt-name>.md
	filename := fmt.Sprintf("%s-%s.md", sourceName, prompt.Name)
	outputPath := filepath.Join(".claude", "commands", "dnaspec", filename)

	// Create directory if needed
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	// Read prompt content
	promptPath := filepath.Join(sourceDir, prompt.File)
	promptContent, err := os.ReadFile(promptPath)
	if err != nil {
		return fmt.Errorf("failed to read prompt file %s: %w", promptPath, err)
	}

	// Generate frontmatter and content
	content := generateClaudeCommandContent(sourceName, prompt, string(promptContent))

	// Write atomically
	return writeFileAtomic(outputPath, []byte(content))
}

// generateClaudeCommandContent creates the full content of a Claude command file
func generateClaudeCommandContent(sourceName string, prompt config.ProjectPrompt, promptContent string) string {
	var sb strings.Builder

	// Frontmatter
	sb.WriteString("---\n")
	sb.WriteString(fmt.Sprintf("name: DNASpec: %s %s\n", formatSourceName(sourceName), formatPromptName(prompt.Name)))
	sb.WriteString(fmt.Sprintf("description: %s\n", prompt.Description))
	sb.WriteString("category: DNASpec\n")
	sb.WriteString(fmt.Sprintf("tags: [dnaspec, \"%s-%s\"]\n", sourceName, prompt.Name))
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

// formatSourceName converts source name to title case for display
func formatSourceName(name string) string {
	parts := strings.Split(name, "-")
	for i, part := range parts {
		if len(part) > 0 {
			parts[i] = strings.ToUpper(part[:1]) + part[1:]
		}
	}
	return strings.Join(parts, " ")
}

// formatPromptName converts prompt name to title case for display
func formatPromptName(name string) string {
	parts := strings.Split(name, "-")
	for i, part := range parts {
		if len(part) > 0 {
			parts[i] = strings.ToUpper(part[:1]) + part[1:]
		}
	}
	return strings.Join(parts, " ")
}
