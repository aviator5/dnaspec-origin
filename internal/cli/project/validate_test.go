package project

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/aviator5/dnaspec/internal/core/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateCommand_Success(t *testing.T) {
	// Create temp directory for test
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(tmpDir)

	// Create valid configuration
	cfg := &config.ProjectConfig{
		Version: 1,
		Agents:  []string{"claude-code", "github-copilot"},
		Sources: []config.ProjectSource{
			{
				Name:   "test-source",
				Type:   "git-repo",
				URL:    "https://github.com/test/repo",
				Ref:    "main",
				Commit: "abc123",
				Guidelines: []config.ProjectGuideline{
					{
						Name:        "test-guideline",
						File:        "guidelines/test.md",
						Description: "Test guideline",
					},
				},
				Prompts: []config.ProjectPrompt{
					{
						Name:        "test-prompt",
						File:        "prompts/test.md",
						Description: "Test prompt",
					},
				},
			},
		},
	}

	// Write config
	err := config.SaveProjectConfig("dnaspec.yaml", cfg)
	require.NoError(t, err)

	// Create referenced files
	err = os.MkdirAll(filepath.Join("dnaspec", "test-source", "guidelines"), 0755)
	require.NoError(t, err)
	err = os.MkdirAll(filepath.Join("dnaspec", "test-source", "prompts"), 0755)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join("dnaspec", "test-source", "guidelines", "test.md"), []byte("# Test"), 0644)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join("dnaspec", "test-source", "prompts", "test.md"), []byte("# Test"), 0644)
	require.NoError(t, err)

	// Run validate
	err = runValidate()
	assert.NoError(t, err)
}

func TestValidateCommand_MissingConfig(t *testing.T) {
	// Create temp directory for test
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(tmpDir)

	// Run validate without config file
	err := runValidate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "project configuration not found")
}

func TestValidateCommand_UnsupportedVersion(t *testing.T) {
	// Create temp directory for test
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(tmpDir)

	// Create config with unsupported version
	cfg := &config.ProjectConfig{
		Version: 99,
		Sources: []config.ProjectSource{},
	}

	err := config.SaveProjectConfig("dnaspec.yaml", cfg)
	require.NoError(t, err)

	// Run validate
	err = runValidate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "validation failed")
}

func TestValidateCommand_MissingSourceFields(t *testing.T) {
	// Create temp directory for test
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(tmpDir)

	// Create config with missing required fields
	cfg := &config.ProjectConfig{
		Version: 1,
		Sources: []config.ProjectSource{
			{
				Name: "test-source",
				Type: "git-repo",
				// Missing URL and Commit
			},
		},
	}

	err := config.SaveProjectConfig("dnaspec.yaml", cfg)
	require.NoError(t, err)

	// Run validate
	err = runValidate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "validation failed")
}

func TestValidateCommand_MissingGuidelineFile(t *testing.T) {
	// Create temp directory for test
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(tmpDir)

	// Create config with guideline file reference
	cfg := &config.ProjectConfig{
		Version: 1,
		Sources: []config.ProjectSource{
			{
				Name:   "test-source",
				Type:   "git-repo",
				URL:    "https://github.com/test/repo",
				Commit: "abc123",
				Guidelines: []config.ProjectGuideline{
					{
						Name:        "test-guideline",
						File:        "guidelines/missing.md",
						Description: "Test guideline",
					},
				},
			},
		},
	}

	err := config.SaveProjectConfig("dnaspec.yaml", cfg)
	require.NoError(t, err)

	// Don't create the file - test missing file detection

	// Run validate
	err = runValidate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "validation failed")
}

func TestValidateCommand_InvalidAgentID(t *testing.T) {
	// Create temp directory for test
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(tmpDir)

	// Create config with invalid agent ID
	cfg := &config.ProjectConfig{
		Version: 1,
		Agents:  []string{"invalid-agent"},
		Sources: []config.ProjectSource{},
	}

	err := config.SaveProjectConfig("dnaspec.yaml", cfg)
	require.NoError(t, err)

	// Run validate
	err = runValidate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "validation failed")
}

func TestValidateCommand_DuplicateSourceNames(t *testing.T) {
	// Create temp directory for test
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(tmpDir)

	// Create config with duplicate source names
	cfg := &config.ProjectConfig{
		Version: 1,
		Sources: []config.ProjectSource{
			{
				Name:   "duplicate",
				Type:   "git-repo",
				URL:    "https://github.com/test/repo1",
				Commit: "abc123",
			},
			{
				Name:   "duplicate",
				Type:   "git-repo",
				URL:    "https://github.com/test/repo2",
				Commit: "def456",
			},
		},
	}

	err := config.SaveProjectConfig("dnaspec.yaml", cfg)
	require.NoError(t, err)

	// Run validate
	err = runValidate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "validation failed")
}

func TestValidateCommand_MultipleErrors(t *testing.T) {
	// Create temp directory for test
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(tmpDir)

	// Create config with multiple errors
	cfg := &config.ProjectConfig{
		Version: 1,
		Agents:  []string{"invalid-agent"},
		Sources: []config.ProjectSource{
			{
				Name: "test-source",
				Type: "git-repo",
				// Missing URL and Commit
				Guidelines: []config.ProjectGuideline{
					{
						Name:        "test-guideline",
						File:        "guidelines/missing.md",
						Description: "Test guideline",
					},
				},
			},
		},
	}

	err := config.SaveProjectConfig("dnaspec.yaml", cfg)
	require.NoError(t, err)

	// Run validate
	err = runValidate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "validation failed")
}
