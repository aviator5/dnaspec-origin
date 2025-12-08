package project

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/aviator5/dnaspec/internal/core/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListCommand_FullConfiguration(t *testing.T) {
	// Create temp directory
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(tmpDir)

	// Create full configuration
	cfg := &config.ProjectConfig{
		Version: 1,
		Agents:  []string{"claude-code", "github-copilot"},
		Sources: []config.ProjectSource{
			{
				Name:   "company-dna",
				Type:   "git-repo",
				URL:    "https://github.com/company/dna",
				Ref:    "v1.0.0",
				Commit: "abc123def456",
				Guidelines: []config.ProjectGuideline{
					{
						Name:                "go-style",
						File:                "guidelines/go-style.md",
						Description:         "Go code style conventions",
						ApplicableScenarios: []string{"writing Go code"},
						Prompts:             []string{"go-review"},
					},
				},
				Prompts: []config.ProjectPrompt{
					{
						Name:        "go-review",
						File:        "prompts/go-review.md",
						Description: "Review Go code",
					},
				},
			},
		},
	}

	err := config.SaveProjectConfig(projectConfigFileName, cfg)
	require.NoError(t, err)

	// Run list command
	cmd := NewListCmd()
	err = cmd.Execute()
	assert.NoError(t, err)
}

func TestListCommand_EmptyAgents(t *testing.T) {
	// Create temp directory
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(tmpDir)

	// Create configuration with no agents
	cfg := &config.ProjectConfig{
		Version: 1,
		Agents:  []string{},
		Sources: []config.ProjectSource{
			{
				Name: "test-source",
				Type: config.SourceTypeLocalPath,
				Path: "/test/path",
			},
		},
	}

	err := config.SaveProjectConfig(projectConfigFileName, cfg)
	require.NoError(t, err)

	// Run list command
	cmd := NewListCmd()
	err = cmd.Execute()
	assert.NoError(t, err)
}

func TestListCommand_EmptySources(t *testing.T) {
	// Create temp directory
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(tmpDir)

	// Create configuration with no sources
	cfg := &config.ProjectConfig{
		Version: 1,
		Agents:  []string{"claude-code"},
		Sources: []config.ProjectSource{},
	}

	err := config.SaveProjectConfig(projectConfigFileName, cfg)
	require.NoError(t, err)

	// Run list command
	cmd := NewListCmd()
	err = cmd.Execute()
	assert.NoError(t, err)
}

func TestListCommand_SourceWithNoGuidelines(t *testing.T) {
	// Create temp directory
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(tmpDir)

	// Create configuration with source that has no guidelines
	cfg := &config.ProjectConfig{
		Version: 1,
		Sources: []config.ProjectSource{
			{
				Name:       "test-source",
				Type:       config.SourceTypeLocalPath,
				Path:       "/test/path",
				Guidelines: []config.ProjectGuideline{},
				Prompts:    []config.ProjectPrompt{},
			},
		},
	}

	err := config.SaveProjectConfig(projectConfigFileName, cfg)
	require.NoError(t, err)

	// Run list command
	cmd := NewListCmd()
	err = cmd.Execute()
	assert.NoError(t, err)
}

func TestListCommand_GitRepoSource(t *testing.T) {
	// Create temp directory
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(tmpDir)

	// Create configuration with git-repo source
	cfg := &config.ProjectConfig{
		Version: 1,
		Sources: []config.ProjectSource{
			{
				Name:   "git-source",
				Type:   "git-repo",
				URL:    "https://github.com/test/repo",
				Ref:    "main",
				Commit: "1234567890abcdef",
				Guidelines: []config.ProjectGuideline{
					{
						Name:        "test-guideline",
						File:        "guidelines/test.md",
						Description: "Test guideline",
					},
				},
			},
		},
	}

	err := config.SaveProjectConfig(projectConfigFileName, cfg)
	require.NoError(t, err)

	// Run list command
	cmd := NewListCmd()
	err = cmd.Execute()
	assert.NoError(t, err)
}

func TestListCommand_LocalDirSource(t *testing.T) {
	// Create temp directory
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(tmpDir)

	// Create configuration with local-dir source
	cfg := &config.ProjectConfig{
		Version: 1,
		Sources: []config.ProjectSource{
			{
				Name: "local-source",
				Type: config.SourceTypeLocalPath,
				Path: "/path/to/local/dna",
				Prompts: []config.ProjectPrompt{
					{
						Name:        "local-prompt",
						File:        "prompts/local.md",
						Description: "Local prompt",
					},
				},
			},
		},
	}

	err := config.SaveProjectConfig(projectConfigFileName, cfg)
	require.NoError(t, err)

	// Run list command
	cmd := NewListCmd()
	err = cmd.Execute()
	assert.NoError(t, err)
}

func TestListCommand_MixedSourceTypes(t *testing.T) {
	// Create temp directory
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(tmpDir)

	// Create configuration with both git-repo and local-dir sources
	cfg := &config.ProjectConfig{
		Version: 1,
		Agents:  []string{"claude-code"},
		Sources: []config.ProjectSource{
			{
				Name:   "git-source",
				Type:   "git-repo",
				URL:    "https://github.com/test/repo",
				Ref:    "v1.0.0",
				Commit: "abc123",
				Guidelines: []config.ProjectGuideline{
					{
						Name:        "git-guideline",
						File:        "guidelines/git.md",
						Description: "Git guideline",
					},
				},
			},
			{
				Name: "local-source",
				Type: config.SourceTypeLocalPath,
				Path: "/path/to/local",
				Prompts: []config.ProjectPrompt{
					{
						Name:        "local-prompt",
						File:        "prompts/local.md",
						Description: "Local prompt",
					},
				},
			},
		},
	}

	err := config.SaveProjectConfig(projectConfigFileName, cfg)
	require.NoError(t, err)

	// Run list command
	cmd := NewListCmd()
	err = cmd.Execute()
	assert.NoError(t, err)
}

func TestListCommand_MissingConfigFile(t *testing.T) {
	// Create temp directory
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(tmpDir)

	// Don't create config file

	// Run list command - should error
	cmd := NewListCmd()
	err := cmd.Execute()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestListCommand_Integration(t *testing.T) {
	// Create temp directory
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(tmpDir)

	// Create complete test scenario with multiple sources
	cfg := &config.ProjectConfig{
		Version: 1,
		Agents:  []string{"claude-code", "github-copilot"},
		Sources: []config.ProjectSource{
			{
				Name:   "company-dna",
				Type:   "git-repo",
				URL:    "https://github.com/company/dna",
				Ref:    "v1.2.0",
				Commit: "abc123def456789",
				Guidelines: []config.ProjectGuideline{
					{
						Name:                "go-style",
						File:                "guidelines/go-style.md",
						Description:         "Go code style conventions",
						ApplicableScenarios: []string{"writing Go code", "refactoring Go code"},
						Prompts:             []string{"go-review"},
					},
					{
						Name:        "rest-api",
						File:        "guidelines/rest-api.md",
						Description: "REST API design guidelines",
					},
				},
				Prompts: []config.ProjectPrompt{
					{
						Name:        "go-review",
						File:        "prompts/go-review.md",
						Description: "Review Go code against go-style guideline",
					},
				},
			},
			{
				Name: "api-patterns",
				Type: config.SourceTypeLocalPath,
				Path: "/Users/me/dna/api",
				Guidelines: []config.ProjectGuideline{
					{
						Name:        "rest-best-practices",
						File:        "guidelines/rest-best-practices.md",
						Description: "REST API best practices",
					},
				},
				Prompts: []config.ProjectPrompt{
					{
						Name:        "api-review",
						File:        "prompts/api-review.md",
						Description: "Review API design",
					},
				},
			},
		},
	}

	configPath := filepath.Join(tmpDir, projectConfigFileName)
	err := config.SaveProjectConfig(configPath, cfg)
	require.NoError(t, err)

	// Verify config was saved
	loadedCfg, err := config.LoadProjectConfig(configPath)
	require.NoError(t, err)
	assert.Equal(t, 2, len(loadedCfg.Sources))
	assert.Equal(t, 2, len(loadedCfg.Agents))

	// Run list command
	cmd := NewListCmd()
	err = cmd.Execute()
	assert.NoError(t, err)
}
