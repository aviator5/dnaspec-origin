package agents

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/aviator5/dnaspec/internal/core/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateAgentFiles(t *testing.T) {
	tempDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer func() {
		err := os.Chdir(originalDir)
		require.NoError(t, err)
	}()

	err := os.Chdir(tempDir)
	require.NoError(t, err)

	// Create test source with guidelines and prompts
	setupTestSource(t, "test-source")

	cfg := &config.ProjectConfig{
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
				Prompts: []config.ProjectPrompt{
					{
						Name:        "review",
						File:        "prompts/review.md",
						Description: "Review code",
					},
					{
						Name:        "lint",
						File:        "prompts/lint.md",
						Description: "Lint code",
					},
				},
			},
		},
	}

	t.Run("generate for Claude Code only", func(t *testing.T) {
		summary, err := GenerateAgentFiles(cfg, []string{"claude-code"})
		require.NoError(t, err)

		assert.True(t, summary.AgentsMD, "should generate AGENTS.md")
		assert.True(t, summary.ClaudeMD, "should generate CLAUDE.md")
		assert.Equal(t, 2, summary.ClaudeCommands, "should generate 2 Claude commands")
		assert.Equal(t, 0, summary.CopilotPrompts, "should not generate Copilot prompts")
		assert.Empty(t, summary.Errors, "should have no errors")

		// Verify files exist
		assert.FileExists(t, "AGENTS.md")
		assert.FileExists(t, "CLAUDE.md")
		assert.FileExists(t, ".claude/commands/dnaspec/test-source-review.md")
		assert.FileExists(t, ".claude/commands/dnaspec/test-source-lint.md")
	})

	t.Run("generate for Copilot only", func(t *testing.T) {
		// Clean up previous files
		err := os.RemoveAll(".claude")
		require.NoError(t, err)
		err = os.Remove("CLAUDE.md")
		require.NoError(t, err)

		summary, err := GenerateAgentFiles(cfg, []string{"github-copilot"})
		require.NoError(t, err)

		assert.True(t, summary.AgentsMD, "should generate AGENTS.md")
		assert.False(t, summary.ClaudeMD, "should not generate CLAUDE.md")
		assert.Equal(t, 0, summary.ClaudeCommands, "should not generate Claude commands")
		assert.Equal(t, 2, summary.CopilotPrompts, "should generate 2 Copilot prompts")
		assert.Empty(t, summary.Errors, "should have no errors")

		// Verify files exist
		assert.FileExists(t, "AGENTS.md")
		assert.NoFileExists(t, "CLAUDE.md")
		assert.FileExists(t, ".github/prompts/dnaspec-test-source-review.prompt.md")
		assert.FileExists(t, ".github/prompts/dnaspec-test-source-lint.prompt.md")
	})

	t.Run("generate for both agents", func(t *testing.T) {
		summary, err := GenerateAgentFiles(cfg, []string{"claude-code", "github-copilot"})
		require.NoError(t, err)

		assert.True(t, summary.AgentsMD)
		assert.True(t, summary.ClaudeMD)
		assert.Equal(t, 2, summary.ClaudeCommands)
		assert.Equal(t, 2, summary.CopilotPrompts)
		assert.Empty(t, summary.Errors)
	})

	t.Run("generate with no agents still creates AGENTS.md", func(t *testing.T) {
		summary, err := GenerateAgentFiles(cfg, []string{})
		require.NoError(t, err)

		assert.True(t, summary.AgentsMD, "should always generate AGENTS.md")
		assert.False(t, summary.ClaudeMD)
		assert.Equal(t, 0, summary.ClaudeCommands)
		assert.Equal(t, 0, summary.CopilotPrompts)
	})

	t.Run("handle errors for missing prompt files", func(t *testing.T) {
		badCfg := &config.ProjectConfig{
			Version: 1,
			Sources: []config.ProjectSource{
				{
					Name: "test-source",
					Guidelines: []config.ProjectGuideline{
						{
							Name:                "test",
							File:                "guidelines/test.md",
							Description:         "Test",
							ApplicableScenarios: []string{"testing"},
						},
					},
					Prompts: []config.ProjectPrompt{
						{
							Name:        "missing",
							File:        "prompts/missing.md",
							Description: "Missing prompt",
						},
					},
				},
			},
		}

		summary, err := GenerateAgentFiles(badCfg, []string{"claude-code"})

		assert.Error(t, err, "should return error")
		assert.True(t, summary.AgentsMD, "should still generate AGENTS.md")
		assert.True(t, summary.ClaudeMD, "should still generate CLAUDE.md")
		assert.Equal(t, 0, summary.ClaudeCommands, "should not count failed commands")
		assert.NotEmpty(t, summary.Errors, "should have error details")
	})
}

func TestGenerateAgentFilesWithMultipleSources(t *testing.T) {
	tempDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer func() {
		err := os.Chdir(originalDir)
		require.NoError(t, err)
	}()

	err := os.Chdir(tempDir)
	require.NoError(t, err)

	// Create two test sources
	setupTestSource(t, "source-a")
	setupTestSource(t, "source-b")

	cfg := &config.ProjectConfig{
		Version: 1,
		Sources: []config.ProjectSource{
			{
				Name: "source-a",
				Guidelines: []config.ProjectGuideline{
					{
						Name:                "guideline-a",
						File:                "guidelines/test.md",
						Description:         "Guideline A",
						ApplicableScenarios: []string{"scenario a"},
					},
				},
				Prompts: []config.ProjectPrompt{
					{
						Name:        "review",
						File:        "prompts/review.md",
						Description: "Review for source A",
					},
				},
			},
			{
				Name: "source-b",
				Guidelines: []config.ProjectGuideline{
					{
						Name:                "guideline-b",
						File:                "guidelines/test.md",
						Description:         "Guideline B",
						ApplicableScenarios: []string{"scenario b"},
					},
				},
				Prompts: []config.ProjectPrompt{
					{
						Name:        "review",
						File:        "prompts/review.md",
						Description: "Review for source B",
					},
				},
			},
		},
	}

	summary, err := GenerateAgentFiles(cfg, []string{"claude-code"})
	require.NoError(t, err)

	// Should generate separate files for each source due to namespacing
	assert.Equal(t, 2, summary.ClaudeCommands)
	assert.FileExists(t, ".claude/commands/dnaspec/source-a-review.md")
	assert.FileExists(t, ".claude/commands/dnaspec/source-b-review.md")

	// AGENTS.md should contain both guidelines
	agentsContent, err := os.ReadFile("AGENTS.md")
	require.NoError(t, err)
	agentsStr := string(agentsContent)
	assert.Contains(t, agentsStr, "source-a")
	assert.Contains(t, agentsStr, "source-b")
}

func TestContains(t *testing.T) {
	tests := []struct {
		name     string
		slice    []string
		value    string
		expected bool
	}{
		{
			name:     "value exists",
			slice:    []string{"a", "b", "c"},
			value:    "b",
			expected: true,
		},
		{
			name:     "value does not exist",
			slice:    []string{"a", "b", "c"},
			value:    "d",
			expected: false,
		},
		{
			name:     "empty slice",
			slice:    []string{},
			value:    "a",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := contains(tt.slice, tt.value)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Helper function to setup a test source directory with files
func setupTestSource(t *testing.T, sourceName string) {
	sourceDir := filepath.Join("dnaspec", sourceName)

	// Create directories
	err := os.MkdirAll(filepath.Join(sourceDir, "guidelines"), 0755)
	require.NoError(t, err)
	err = os.MkdirAll(filepath.Join(sourceDir, "prompts"), 0755)
	require.NoError(t, err)

	// Create test guideline file
	guidelineContent := "# Test Guideline\n\nThis is a test guideline."
	err = os.WriteFile(filepath.Join(sourceDir, "guidelines", "test.md"), []byte(guidelineContent), 0644)
	require.NoError(t, err)

	// Create test prompt files
	reviewContent := "Review the code against guidelines."
	err = os.WriteFile(filepath.Join(sourceDir, "prompts", "review.md"), []byte(reviewContent), 0644)
	require.NoError(t, err)

	lintContent := "Lint the code for issues."
	err = os.WriteFile(filepath.Join(sourceDir, "prompts", "lint.md"), []byte(lintContent), 0644)
	require.NoError(t, err)
}
