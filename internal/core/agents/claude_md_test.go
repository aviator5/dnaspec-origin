package agents

import (
	"os"
	"testing"

	"github.com/aviator5/dnaspec/internal/core/config"
	"github.com/aviator5/dnaspec/internal/core/files"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateClaudeMD(t *testing.T) {
	// Create temp directory for test
	tempDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	err := os.Chdir(tempDir)
	require.NoError(t, err)

	config := &config.ProjectConfig{
		Version: 1,
		Sources: []config.ProjectSource{
			{
				Name: "test-source",
				Guidelines: []config.ProjectGuideline{
					{
						Name:        "test-guideline",
						File:        "guidelines/test.md",
						Description: "Test guideline",
						ApplicableScenarios: []string{
							"testing code",
						},
					},
				},
			},
		},
	}

	t.Run("create new CLAUDE.md", func(t *testing.T) {
		err := GenerateClaudeMD(config)
		require.NoError(t, err)

		content, err := os.ReadFile("CLAUDE.md")
		require.NoError(t, err)

		contentStr := string(content)
		assert.Contains(t, contentStr, "# DNASpec Agent Instructions")
		assert.Contains(t, contentStr, files.ManagedBlockStart)
		assert.Contains(t, contentStr, files.ManagedBlockEnd)
		assert.Contains(t, contentStr, "@/dnaspec/test-source/guidelines/test.md")
		assert.Contains(t, contentStr, "testing code")
	})

	t.Run("update existing CLAUDE.md preserving user content", func(t *testing.T) {
		// Create file with user content
		userContent := `# My Claude Instructions

Custom Claude setup.

<!-- DNASPEC:START -->
Old content
<!-- DNASPEC:END -->

Additional notes.`

		err := os.WriteFile("CLAUDE.md", []byte(userContent), 0644)
		require.NoError(t, err)

		// Update
		err = GenerateClaudeMD(config)
		require.NoError(t, err)

		content, err := os.ReadFile("CLAUDE.md")
		require.NoError(t, err)

		contentStr := string(content)

		// User content should be preserved
		assert.Contains(t, contentStr, "# My Claude Instructions")
		assert.Contains(t, contentStr, "Custom Claude setup.")
		assert.Contains(t, contentStr, "Additional notes.")

		// New content should be present
		assert.Contains(t, contentStr, "@/dnaspec/test-source/guidelines/test.md")

		// Old content should be replaced
		assert.NotContains(t, contentStr, "Old content")
	})

	t.Run("same content as AGENTS.md when both created fresh", func(t *testing.T) {
		// Remove existing files to ensure fresh start
		os.Remove("AGENTS.md")
		os.Remove("CLAUDE.md")

		// Generate both files fresh
		err := GenerateAgentsMD(config)
		require.NoError(t, err)

		err = GenerateClaudeMD(config)
		require.NoError(t, err)

		// Read both
		agentsContent, err := os.ReadFile("AGENTS.md")
		require.NoError(t, err)

		claudeContent, err := os.ReadFile("CLAUDE.md")
		require.NoError(t, err)

		// Should have identical content when created fresh
		assert.Equal(t, string(agentsContent), string(claudeContent))
	})
}
