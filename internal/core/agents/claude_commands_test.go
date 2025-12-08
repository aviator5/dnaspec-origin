package agents

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/aviator5/dnaspec/internal/core/config"
	"github.com/aviator5/dnaspec/internal/core/files"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateClaudeCommandContent(t *testing.T) {
	tests := []struct {
		name          string
		sourceName    string
		prompt        config.ProjectPrompt
		promptContent string
		contains      []string
	}{
		{
			name:       "basic command generation",
			sourceName: "company-dna",
			prompt: config.ProjectPrompt{
				Name:        "code-review",
				File:        "prompts/code-review.md",
				Description: "Review code against guidelines",
			},
			promptContent: "Review the code carefully.\n\nCheck for style issues.",
			contains: []string{
				"name: DNASpec: Company Dna Code Review",
				"description: Review code against guidelines",
				"category: DNASpec",
				"tags: [dnaspec, \"company-dna-code-review\"]",
				files.ManagedBlockStart,
				"Review the code carefully.",
				"Check for style issues.",
				files.ManagedBlockEnd,
			},
		},
		{
			name:       "multi-word source and prompt names",
			sourceName: "my-team-patterns",
			prompt: config.ProjectPrompt{
				Name:        "api-design-review",
				File:        "prompts/api-review.md",
				Description: "Review API design",
			},
			promptContent: "Check API endpoints.",
			contains: []string{
				"name: DNASpec: My Team Patterns Api Design Review",
				"description: Review API design",
				"tags: [dnaspec, \"my-team-patterns-api-design-review\"]",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content := generateClaudeCommandContent(tt.sourceName, tt.prompt, tt.promptContent)

			for _, substr := range tt.contains {
				assert.Contains(t, content, substr)
			}
		})
	}
}

func TestGenerateClaudeCommand(t *testing.T) {
	tempDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer func() {
		err := os.Chdir(originalDir)
		require.NoError(t, err)
	}()

	err := os.Chdir(tempDir)
	require.NoError(t, err)

	// Create source directory with prompt file
	sourceDir := filepath.Join(tempDir, "dnaspec", "test-source")
	err = os.MkdirAll(filepath.Join(sourceDir, "prompts"), 0755)
	require.NoError(t, err)

	promptContent := "This is a test prompt.\n\nCheck the code."
	err = os.WriteFile(filepath.Join(sourceDir, "prompts", "review.md"), []byte(promptContent), 0644)
	require.NoError(t, err)

	prompt := config.ProjectPrompt{
		Name:        "review",
		File:        "prompts/review.md",
		Description: "Review code",
	}

	t.Run("generate new command file", func(t *testing.T) {
		err := GenerateClaudeCommand("test-source", prompt, sourceDir)
		require.NoError(t, err)

		// Check file was created
		expectedPath := filepath.Join(".claude", "commands", "dnaspec", "test-source-review.md")
		content, err := os.ReadFile(expectedPath)
		require.NoError(t, err)

		contentStr := string(content)
		assert.Contains(t, contentStr, "name: DNASpec: Test Source Review")
		assert.Contains(t, contentStr, "description: Review code")
		assert.Contains(t, contentStr, "category: DNASpec")
		assert.Contains(t, contentStr, "This is a test prompt.")
		assert.Contains(t, contentStr, files.ManagedBlockStart)
		assert.Contains(t, contentStr, files.ManagedBlockEnd)
	})

	t.Run("overwrite existing command file", func(t *testing.T) {
		// Create file with old content
		commandPath := filepath.Join(".claude", "commands", "dnaspec", "test-source-review.md")
		oldContent := "---\nname: Old Command\n---\nOld content"
		err := os.WriteFile(commandPath, []byte(oldContent), 0644)
		require.NoError(t, err)

		// Generate new
		err = GenerateClaudeCommand("test-source", prompt, sourceDir)
		require.NoError(t, err)

		// Check it was overwritten
		content, err := os.ReadFile(commandPath)
		require.NoError(t, err)

		contentStr := string(content)
		assert.Contains(t, contentStr, "name: DNASpec: Test Source Review")
		assert.NotContains(t, contentStr, "Old Command")
	})

	t.Run("error on missing prompt file", func(t *testing.T) {
		missingPrompt := config.ProjectPrompt{
			Name:        "missing",
			File:        "prompts/missing.md",
			Description: "Missing prompt",
		}

		err := GenerateClaudeCommand("test-source", missingPrompt, sourceDir)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to read prompt file")
	})
}

func TestFormatSourceName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "company-dna",
			expected: "Company Dna",
		},
		{
			input:    "my-team-patterns",
			expected: "My Team Patterns",
		},
		{
			input:    "single",
			expected: "Single",
		},
		{
			input:    "api-v2-guidelines",
			expected: "Api V2 Guidelines",
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := formatSourceName(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFormatPromptName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "code-review",
			expected: "Code Review",
		},
		{
			input:    "api-design-check",
			expected: "Api Design Check",
		},
		{
			input:    "lint",
			expected: "Lint",
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := formatPromptName(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
