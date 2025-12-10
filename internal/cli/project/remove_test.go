package project

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/aviator5/dnaspec/internal/core/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRemoveCommand_Success(t *testing.T) {
	// Create temp directory
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(tmpDir)

	// Create configuration
	cfg := &config.ProjectConfig{
		Version: 1,
		Sources: []config.ProjectSource{
			{
				Name: "test-source",
				Type: "git-repo",
				URL:  "https://github.com/test/repo",
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

	// Create source directory with files
	sourceDir := filepath.Join("dnaspec", "test-source")
	err = os.MkdirAll(sourceDir, 0755)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(sourceDir, "test-guideline.md"), []byte("test"), 0644)
	require.NoError(t, err)

	// Create prompts directory
	promptsDir := filepath.Join(sourceDir, "prompts")
	err = os.MkdirAll(promptsDir, 0755)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(promptsDir, "test-prompt.md"), []byte("prompt"), 0644)
	require.NoError(t, err)

	// Create generated agent files
	claudeDir := filepath.Join(".claude", "commands", "dnaspec")
	err = os.MkdirAll(claudeDir, 0755)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(claudeDir, "test-source-guideline.md"), []byte("claude"), 0644)
	require.NoError(t, err)

	copilotDir := filepath.Join(".github", "prompts")
	err = os.MkdirAll(copilotDir, 0755)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(copilotDir, "dnaspec-test-source-prompt.prompt.md"), []byte("copilot"), 0644)
	require.NoError(t, err)

	// Run remove with --force
	err = runRemove("test-source", true)
	assert.NoError(t, err)

	// Verify source directory was deleted
	_, err = os.Stat(sourceDir)
	assert.True(t, os.IsNotExist(err))

	// Verify generated files were deleted
	_, err = os.Stat(filepath.Join(claudeDir, "test-source-guideline.md"))
	assert.True(t, os.IsNotExist(err))
	_, err = os.Stat(filepath.Join(copilotDir, "dnaspec-test-source-prompt.prompt.md"))
	assert.True(t, os.IsNotExist(err))

	// Verify config was updated
	updatedCfg, err := config.LoadProjectConfig(projectConfigFileName)
	require.NoError(t, err)
	assert.Equal(t, 0, len(updatedCfg.Sources))
}

func TestRemoveCommand_SourceNotFound(t *testing.T) {
	// Create temp directory
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(tmpDir)

	// Create configuration with different source
	cfg := &config.ProjectConfig{
		Version: 1,
		Sources: []config.ProjectSource{
			{
				Name: "existing-source",
				Type: "git-repo",
				URL:  "https://github.com/test/repo",
			},
		},
	}
	err := config.SaveProjectConfig(projectConfigFileName, cfg)
	require.NoError(t, err)

	// Try to remove non-existent source
	err = runRemove("non-existent", true)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "source not found")

	// Verify config unchanged
	updatedCfg, err := config.LoadProjectConfig(projectConfigFileName)
	require.NoError(t, err)
	assert.Equal(t, 1, len(updatedCfg.Sources))
}

func TestRemoveCommand_MissingConfigFile(t *testing.T) {
	// Create temp directory
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(tmpDir)

	// Don't create config file

	// Try to remove source
	err := runRemove("test-source", true)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestRemoveCommand_IdempotentDirectoryDeletion(t *testing.T) {
	// Create temp directory
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(tmpDir)

	// Create configuration
	cfg := &config.ProjectConfig{
		Version: 1,
		Sources: []config.ProjectSource{
			{
				Name: "test-source",
				Type: "git-repo",
				URL:  "https://github.com/test/repo",
			},
		},
	}
	err := config.SaveProjectConfig(projectConfigFileName, cfg)
	require.NoError(t, err)

	// Don't create source directory - should still succeed

	// Run remove with --force
	err = runRemove("test-source", true)
	assert.NoError(t, err)

	// Verify config was updated
	updatedCfg, err := config.LoadProjectConfig(projectConfigFileName)
	require.NoError(t, err)
	assert.Equal(t, 0, len(updatedCfg.Sources))
}

func TestRemoveCommand_MultipleSourcesRemoveOne(t *testing.T) {
	// Create temp directory
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(tmpDir)

	// Create configuration with multiple sources
	cfg := &config.ProjectConfig{
		Version: 1,
		Sources: []config.ProjectSource{
			{
				Name: "source-1",
				Type: "git-repo",
				URL:  "https://github.com/test/repo1",
			},
			{
				Name: "source-2",
				Type: "git-repo",
				URL:  "https://github.com/test/repo2",
			},
			{
				Name: "source-3",
				Type: config.SourceTypeLocalPath,
				Path: "/test/path",
			},
		},
	}
	err := config.SaveProjectConfig(projectConfigFileName, cfg)
	require.NoError(t, err)

	// Create source directories
	for _, source := range cfg.Sources {
		sourceDir := filepath.Join("dnaspec", source.Name)
		err = os.MkdirAll(sourceDir, 0755)
		require.NoError(t, err)
	}

	// Remove middle source
	err = runRemove("source-2", true)
	assert.NoError(t, err)

	// Verify only source-2 was removed
	updatedCfg, err := config.LoadProjectConfig(projectConfigFileName)
	require.NoError(t, err)
	assert.Equal(t, 2, len(updatedCfg.Sources))
	assert.Equal(t, "source-1", updatedCfg.Sources[0].Name)
	assert.Equal(t, "source-3", updatedCfg.Sources[1].Name)

	// Verify source-2 directory was deleted
	_, err = os.Stat(filepath.Join("dnaspec", "source-2"))
	assert.True(t, os.IsNotExist(err))

	// Verify other directories still exist
	_, err = os.Stat(filepath.Join("dnaspec", "source-1"))
	assert.NoError(t, err)
	_, err = os.Stat(filepath.Join("dnaspec", "source-3"))
	assert.NoError(t, err)
}

func TestRemoveCommand_NoGeneratedFiles(t *testing.T) {
	// Create temp directory
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(tmpDir)

	// Create configuration
	cfg := &config.ProjectConfig{
		Version: 1,
		Sources: []config.ProjectSource{
			{
				Name: "test-source",
				Type: "git-repo",
				URL:  "https://github.com/test/repo",
			},
		},
	}
	err := config.SaveProjectConfig(projectConfigFileName, cfg)
	require.NoError(t, err)

	// Create source directory but no generated agent files
	sourceDir := filepath.Join("dnaspec", "test-source")
	err = os.MkdirAll(sourceDir, 0755)
	require.NoError(t, err)

	// Run remove with --force
	err = runRemove("test-source", true)
	assert.NoError(t, err)

	// Verify source directory was deleted
	_, err = os.Stat(sourceDir)
	assert.True(t, os.IsNotExist(err))

	// Verify config was updated
	updatedCfg, err := config.LoadProjectConfig(projectConfigFileName)
	require.NoError(t, err)
	assert.Equal(t, 0, len(updatedCfg.Sources))
}

func TestRemoveCommand_OnlyClaudeFiles(t *testing.T) {
	// Create temp directory
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(tmpDir)

	// Create configuration
	cfg := &config.ProjectConfig{
		Version: 1,
		Sources: []config.ProjectSource{
			{
				Name: "test-source",
				Type: "git-repo",
				URL:  "https://github.com/test/repo",
			},
		},
	}
	err := config.SaveProjectConfig(projectConfigFileName, cfg)
	require.NoError(t, err)

	// Create only Claude command files
	claudeDir := filepath.Join(".claude", "commands", "dnaspec")
	err = os.MkdirAll(claudeDir, 0755)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(claudeDir, "test-source-guideline1.md"), []byte("claude1"), 0644)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(claudeDir, "test-source-guideline2.md"), []byte("claude2"), 0644)
	require.NoError(t, err)

	// Run remove with --force
	err = runRemove("test-source", true)
	assert.NoError(t, err)

	// Verify generated files were deleted
	files, _ := filepath.Glob(filepath.Join(claudeDir, "test-source-*.md"))
	assert.Equal(t, 0, len(files))
}

func TestRemoveCommand_OnlyCopilotFiles(t *testing.T) {
	// Create temp directory
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(tmpDir)

	// Create configuration
	cfg := &config.ProjectConfig{
		Version: 1,
		Sources: []config.ProjectSource{
			{
				Name: "test-source",
				Type: "git-repo",
				URL:  "https://github.com/test/repo",
			},
		},
	}
	err := config.SaveProjectConfig(projectConfigFileName, cfg)
	require.NoError(t, err)

	// Create only Copilot prompt files
	copilotDir := filepath.Join(".github", "prompts")
	err = os.MkdirAll(copilotDir, 0755)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(copilotDir, "dnaspec-test-source-prompt1.prompt.md"), []byte("copilot1"), 0644)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(copilotDir, "dnaspec-test-source-prompt2.prompt.md"), []byte("copilot2"), 0644)
	require.NoError(t, err)

	// Run remove with --force
	err = runRemove("test-source", true)
	assert.NoError(t, err)

	// Verify generated files were deleted
	files, _ := filepath.Glob(filepath.Join(copilotDir, "dnaspec-test-source-*.prompt.md"))
	assert.Equal(t, 0, len(files))
}

func TestRemoveCommand_PreservesOtherSourceFiles(t *testing.T) {
	// Create temp directory
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(tmpDir)

	// Create configuration with multiple sources
	cfg := &config.ProjectConfig{
		Version: 1,
		Sources: []config.ProjectSource{
			{
				Name: "source-to-remove",
				Type: "git-repo",
				URL:  "https://github.com/test/repo1",
			},
			{
				Name: "source-to-keep",
				Type: "git-repo",
				URL:  "https://github.com/test/repo2",
			},
		},
	}
	err := config.SaveProjectConfig(projectConfigFileName, cfg)
	require.NoError(t, err)

	// Create generated files for both sources
	claudeDir := filepath.Join(".claude", "commands", "dnaspec")
	err = os.MkdirAll(claudeDir, 0755)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(claudeDir, "source-to-remove-guideline.md"), []byte("to-remove"), 0644)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(claudeDir, "source-to-keep-guideline.md"), []byte("to-keep"), 0644)
	require.NoError(t, err)

	copilotDir := filepath.Join(".github", "prompts")
	err = os.MkdirAll(copilotDir, 0755)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(copilotDir, "dnaspec-source-to-remove-prompt.prompt.md"), []byte("to-remove"), 0644)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(copilotDir, "dnaspec-source-to-keep-prompt.prompt.md"), []byte("to-keep"), 0644)
	require.NoError(t, err)

	// Run remove with --force
	err = runRemove("source-to-remove", true)
	assert.NoError(t, err)

	// Verify only source-to-remove files were deleted
	_, err = os.Stat(filepath.Join(claudeDir, "source-to-remove-guideline.md"))
	assert.True(t, os.IsNotExist(err))
	_, err = os.Stat(filepath.Join(copilotDir, "dnaspec-source-to-remove-prompt.prompt.md"))
	assert.True(t, os.IsNotExist(err))

	// Verify source-to-keep files still exist
	_, err = os.Stat(filepath.Join(claudeDir, "source-to-keep-guideline.md"))
	assert.NoError(t, err)
	_, err = os.Stat(filepath.Join(copilotDir, "dnaspec-source-to-keep-prompt.prompt.md"))
	assert.NoError(t, err)
}

func TestRemoveCommand_AllAgentFiles(t *testing.T) {
	// Create temp directory
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(tmpDir)

	// Create configuration
	cfg := &config.ProjectConfig{
		Version: 1,
		Sources: []config.ProjectSource{
			{
				Name: "test-source",
				Type: "git-repo",
				URL:  "https://github.com/test/repo",
			},
		},
	}
	err := config.SaveProjectConfig(projectConfigFileName, cfg)
	require.NoError(t, err)

	// Create generated files for all agents
	agentFiles := []struct {
		dir  string
		file string
	}{
		{".agent/workflows", "dnaspec-test-source-workflow.md"},
		{".claude/commands/dnaspec", "test-source-command.md"},
		{".cursor/commands", "dnaspec-test-source-cursor.md"},
		{".github/prompts", "dnaspec-test-source-prompt.prompt.md"},
		{".windsurf/workflows", "dnaspec-test-source-windsurf.md"},
	}

	for _, af := range agentFiles {
		err = os.MkdirAll(af.dir, 0755)
		require.NoError(t, err)
		err = os.WriteFile(filepath.Join(af.dir, af.file), []byte("content"), 0644)
		require.NoError(t, err)
	}

	// Run remove with --force
	err = runRemove("test-source", true)
	assert.NoError(t, err)

	// Verify all generated files were deleted
	for _, af := range agentFiles {
		_, err = os.Stat(filepath.Join(af.dir, af.file))
		assert.True(t, os.IsNotExist(err), "File %s should be deleted", filepath.Join(af.dir, af.file))
	}
}

func TestRemoveCommand_CursorFiles(t *testing.T) {
	// Create temp directory
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(tmpDir)

	// Create configuration
	cfg := &config.ProjectConfig{
		Version: 1,
		Sources: []config.ProjectSource{
			{
				Name: "test-source",
				Type: "git-repo",
				URL:  "https://github.com/test/repo",
			},
		},
	}
	err := config.SaveProjectConfig(projectConfigFileName, cfg)
	require.NoError(t, err)

	// Create only Cursor command files
	cursorDir := filepath.Join(".cursor", "commands")
	err = os.MkdirAll(cursorDir, 0755)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(cursorDir, "dnaspec-test-source-cmd1.md"), []byte("cursor1"), 0644)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(cursorDir, "dnaspec-test-source-cmd2.md"), []byte("cursor2"), 0644)
	require.NoError(t, err)

	// Run remove with --force
	err = runRemove("test-source", true)
	assert.NoError(t, err)

	// Verify generated files were deleted
	files, _ := filepath.Glob(filepath.Join(cursorDir, "dnaspec-test-source-*.md"))
	assert.Equal(t, 0, len(files))
}

func TestRemoveCommand_WindsurfFiles(t *testing.T) {
	// Create temp directory
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(tmpDir)

	// Create configuration
	cfg := &config.ProjectConfig{
		Version: 1,
		Sources: []config.ProjectSource{
			{
				Name: "test-source",
				Type: "git-repo",
				URL:  "https://github.com/test/repo",
			},
		},
	}
	err := config.SaveProjectConfig(projectConfigFileName, cfg)
	require.NoError(t, err)

	// Create only Windsurf workflow files
	windsurfDir := filepath.Join(".windsurf", "workflows")
	err = os.MkdirAll(windsurfDir, 0755)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(windsurfDir, "dnaspec-test-source-flow1.md"), []byte("windsurf1"), 0644)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(windsurfDir, "dnaspec-test-source-flow2.md"), []byte("windsurf2"), 0644)
	require.NoError(t, err)

	// Run remove with --force
	err = runRemove("test-source", true)
	assert.NoError(t, err)

	// Verify generated files were deleted
	files, _ := filepath.Glob(filepath.Join(windsurfDir, "dnaspec-test-source-*.md"))
	assert.Equal(t, 0, len(files))
}

func TestRemoveCommand_AntigravityFiles(t *testing.T) {
	// Create temp directory
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(tmpDir)

	// Create configuration
	cfg := &config.ProjectConfig{
		Version: 1,
		Sources: []config.ProjectSource{
			{
				Name: "test-source",
				Type: "git-repo",
				URL:  "https://github.com/test/repo",
			},
		},
	}
	err := config.SaveProjectConfig(projectConfigFileName, cfg)
	require.NoError(t, err)

	// Create only Antigravity workflow files
	antigravityDir := filepath.Join(".agent", "workflows")
	err = os.MkdirAll(antigravityDir, 0755)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(antigravityDir, "dnaspec-test-source-workflow1.md"), []byte("antigravity1"), 0644)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(antigravityDir, "dnaspec-test-source-workflow2.md"), []byte("antigravity2"), 0644)
	require.NoError(t, err)

	// Run remove with --force
	err = runRemove("test-source", true)
	assert.NoError(t, err)

	// Verify generated files were deleted
	files, _ := filepath.Glob(filepath.Join(antigravityDir, "dnaspec-test-source-*.md"))
	assert.Equal(t, 0, len(files))
}
