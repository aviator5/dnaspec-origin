package agents

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/aviator5/dnaspec/internal/core/config"
	"github.com/aviator5/dnaspec/internal/core/files"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateCopilotPromptContent(t *testing.T) {
	tests := []struct {
		name          string
		prompt        config.ProjectPrompt
		promptContent string
		contains      []string
	}{
		{
			name: "basic prompt generation",
			prompt: config.ProjectPrompt{
				Name:        "code-review",
				File:        "prompts/code-review.md",
				Description: "Review code against guidelines",
			},
			promptContent: "Review the code carefully.\n\nCheck for style issues.",
			contains: []string{
				"description: Review code against guidelines",
				"$ARGUMENTS",
				files.ManagedBlockStart,
				"Review the code carefully.",
				"Check for style issues.",
				files.ManagedBlockEnd,
			},
		},
		{
			name: "prompt with special characters in description",
			prompt: config.ProjectPrompt{
				Name:        "api-review",
				File:        "prompts/api-review.md",
				Description: "Review API design & implementation",
			},
			promptContent: "Check API endpoints.",
			contains: []string{
				"description: Review API design & implementation",
				"$ARGUMENTS",
				"Check API endpoints.",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content := generateCopilotPromptContent(tt.prompt, tt.promptContent)

			for _, substr := range tt.contains {
				assert.Contains(t, content, substr)
			}

			// Check structure: frontmatter, then $ARGUMENTS, then managed block
			assert.True(t, strings.Index(content, "---") < strings.Index(content, "$ARGUMENTS"))
			assert.True(t, strings.Index(content, "$ARGUMENTS") < strings.Index(content, files.ManagedBlockStart))
		})
	}
}

func TestGenerateCopilotPrompt(t *testing.T) {
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

	promptContent := "This is a test prompt for Copilot.\n\nCheck the code thoroughly."
	err = os.WriteFile(filepath.Join(sourceDir, "prompts", "review.md"), []byte(promptContent), 0644)
	require.NoError(t, err)

	prompt := config.ProjectPrompt{
		Name:        "review",
		File:        "prompts/review.md",
		Description: "Review code",
	}

	t.Run("generate new prompt file", func(t *testing.T) {
		err := GenerateCopilotPrompt("test-source", prompt, sourceDir)
		require.NoError(t, err)

		// Check file was created
		expectedPath := filepath.Join(".github", "prompts", "dnaspec-test-source-review.prompt.md")
		content, err := os.ReadFile(expectedPath)
		require.NoError(t, err)

		contentStr := string(content)
		assert.Contains(t, contentStr, "description: Review code")
		assert.Contains(t, contentStr, "$ARGUMENTS")
		assert.Contains(t, contentStr, "This is a test prompt for Copilot.")
		assert.Contains(t, contentStr, files.ManagedBlockStart)
		assert.Contains(t, contentStr, files.ManagedBlockEnd)
	})

	t.Run("overwrite existing prompt file", func(t *testing.T) {
		// Create file with old content
		promptPath := filepath.Join(".github", "prompts", "dnaspec-test-source-review.prompt.md")
		oldContent := "---\ndescription: Old\n---\nOld content"
		err := os.WriteFile(promptPath, []byte(oldContent), 0644)
		require.NoError(t, err)

		// Generate new
		err = GenerateCopilotPrompt("test-source", prompt, sourceDir)
		require.NoError(t, err)

		// Check it was overwritten
		content, err := os.ReadFile(promptPath)
		require.NoError(t, err)

		contentStr := string(content)
		assert.Contains(t, contentStr, "description: Review code")
		assert.NotContains(t, contentStr, "description: Old")
	})

	t.Run("error on missing prompt file", func(t *testing.T) {
		missingPrompt := config.ProjectPrompt{
			Name:        "missing",
			File:        "prompts/missing.md",
			Description: "Missing prompt",
		}

		err := GenerateCopilotPrompt("test-source", missingPrompt, sourceDir)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to read prompt file")
	})

	t.Run("filename format with source namespacing", func(t *testing.T) {
		err := GenerateCopilotPrompt("company-dna", prompt, sourceDir)
		require.NoError(t, err)

		// Check file was created with correct name
		expectedPath := filepath.Join(".github", "prompts", "dnaspec-company-dna-review.prompt.md")
		_, err = os.Stat(expectedPath)
		require.NoError(t, err, "file should exist at expected path")
	})
}
