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

func TestGenerateAgentsMDContent(t *testing.T) {
	tests := []struct {
		name     string
		config   *config.ProjectConfig
		contains []string
		notContains []string
	}{
		{
			name: "single source with one guideline",
			config: &config.ProjectConfig{
				Version: 1,
				Sources: []config.ProjectSource{
					{
						Name: "company-dna",
						Guidelines: []config.ProjectGuideline{
							{
								Name:        "go-style",
								File:        "guidelines/go-style.md",
								Description: "Go code style conventions",
								ApplicableScenarios: []string{
									"writing new Go code",
									"refactoring existing Go code",
								},
							},
						},
					},
				},
			},
			contains: []string{
				"DNASpec Instructions",
				"@/dnaspec/company-dna/guidelines/go-style.md",
				"writing new Go code",
				"refactoring existing Go code",
				"Keep this managed block",
			},
		},
		{
			name: "multiple sources with multiple guidelines",
			config: &config.ProjectConfig{
				Version: 1,
				Sources: []config.ProjectSource{
					{
						Name: "company-dna",
						Guidelines: []config.ProjectGuideline{
							{
								Name:        "go-style",
								File:        "guidelines/go-style.md",
								Description: "Go code style",
								ApplicableScenarios: []string{
									"writing new Go code",
								},
							},
						},
					},
					{
						Name: "team-patterns",
						Guidelines: []config.ProjectGuideline{
							{
								Name:        "rest-api",
								File:        "guidelines/rest-api.md",
								Description: "REST API design",
								ApplicableScenarios: []string{
									"designing API endpoints",
									"implementing HTTP handlers",
								},
							},
						},
					},
				},
			},
			contains: []string{
				"@/dnaspec/company-dna/guidelines/go-style.md",
				"@/dnaspec/team-patterns/guidelines/rest-api.md",
				"writing new Go code",
				"designing API endpoints",
				"implementing HTTP handlers",
			},
		},
		{
			name: "empty sources",
			config: &config.ProjectConfig{
				Version: 1,
				Sources: []config.ProjectSource{},
			},
			contains: []string{
				"DNASpec Instructions",
				"No DNA sources configured",
				"Run 'dnaspec add'",
			},
			notContains: []string{
				"@/dnaspec/",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content := generateAgentsMDContent(tt.config)

			for _, substr := range tt.contains {
				assert.Contains(t, content, substr, "should contain: %s", substr)
			}

			for _, substr := range tt.notContains {
				assert.NotContains(t, content, substr, "should not contain: %s", substr)
			}
		})
	}
}

func TestGenerateAgentsMD(t *testing.T) {
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

	t.Run("create new AGENTS.md", func(t *testing.T) {
		err := GenerateAgentsMD(config)
		require.NoError(t, err)

		content, err := os.ReadFile("AGENTS.md")
		require.NoError(t, err)

		contentStr := string(content)
		assert.Contains(t, contentStr, "# DNASpec Agent Instructions")
		assert.Contains(t, contentStr, files.ManagedBlockStart)
		assert.Contains(t, contentStr, files.ManagedBlockEnd)
		assert.Contains(t, contentStr, "@/dnaspec/test-source/guidelines/test.md")
		assert.Contains(t, contentStr, "testing code")
	})

	t.Run("update existing AGENTS.md preserving user content", func(t *testing.T) {
		// Create file with user content
		userContent := `# My Custom Header

User content before managed block.

<!-- DNASPEC:START -->
Old content to be replaced
<!-- DNASPEC:END -->

User content after managed block.`

		err := os.WriteFile("AGENTS.md", []byte(userContent), 0644)
		require.NoError(t, err)

		// Update
		err = GenerateAgentsMD(config)
		require.NoError(t, err)

		content, err := os.ReadFile("AGENTS.md")
		require.NoError(t, err)

		contentStr := string(content)

		// User content should be preserved
		assert.Contains(t, contentStr, "# My Custom Header")
		assert.Contains(t, contentStr, "User content before managed block.")
		assert.Contains(t, contentStr, "User content after managed block.")

		// New content should be present
		assert.Contains(t, contentStr, "@/dnaspec/test-source/guidelines/test.md")
		assert.Contains(t, contentStr, "testing code")

		// Old content should be replaced
		assert.NotContains(t, contentStr, "Old content to be replaced")
	})

	t.Run("append to existing AGENTS.md without managed block", func(t *testing.T) {
		// Create file without managed block
		existingContent := "# Existing Content\n\nSome notes.\n"
		err := os.WriteFile("AGENTS.md", []byte(existingContent), 0644)
		require.NoError(t, err)

		// Generate
		err = GenerateAgentsMD(config)
		require.NoError(t, err)

		content, err := os.ReadFile("AGENTS.md")
		require.NoError(t, err)

		contentStr := string(content)

		// Original content should be preserved
		assert.Contains(t, contentStr, "# Existing Content")
		assert.Contains(t, contentStr, "Some notes.")

		// Managed block should be appended
		assert.Contains(t, contentStr, files.ManagedBlockStart)
		assert.Contains(t, contentStr, "@/dnaspec/test-source/guidelines/test.md")
	})
}

func TestWriteFileAtomic(t *testing.T) {
	tempDir := t.TempDir()

	t.Run("write new file", func(t *testing.T) {
		path := filepath.Join(tempDir, "test.txt")
		content := []byte("test content")

		err := writeFileAtomic(path, content)
		require.NoError(t, err)

		readContent, err := os.ReadFile(path)
		require.NoError(t, err)
		assert.Equal(t, content, readContent)
	})

	t.Run("overwrite existing file", func(t *testing.T) {
		path := filepath.Join(tempDir, "existing.txt")

		// Write initial content
		err := os.WriteFile(path, []byte("old content"), 0644)
		require.NoError(t, err)

		// Overwrite atomically
		newContent := []byte("new content")
		err = writeFileAtomic(path, newContent)
		require.NoError(t, err)

		readContent, err := os.ReadFile(path)
		require.NoError(t, err)
		assert.Equal(t, newContent, readContent)
	})

	t.Run("no temp files left on success", func(t *testing.T) {
		path := filepath.Join(tempDir, "cleanup-test.txt")
		content := []byte("test content")

		err := writeFileAtomic(path, content)
		require.NoError(t, err)

		// Check no temp files remain
		entries, err := os.ReadDir(tempDir)
		require.NoError(t, err)

		for _, entry := range entries {
			assert.False(t, strings.HasPrefix(entry.Name(), ".dnaspec-tmp-"),
				"temp file should be cleaned up: %s", entry.Name())
		}
	})
}
