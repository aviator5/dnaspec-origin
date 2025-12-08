package manifest

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/aviator5/dnaspec/internal/core/config"
)

func TestInitCmd_Success(t *testing.T) {
	// Create temp directory and change to it
	tmpDir := t.TempDir()
	originalDir, err := os.Getwd()
	require.NoError(t, err)
	defer func() {
		err := os.Chdir(originalDir)
		require.NoError(t, err)
	}()

	err = os.Chdir(tmpDir)
	require.NoError(t, err)

	// Run init
	err = runInit()
	assert.NoError(t, err)

	// Verify manifest file was created
	manifestPath := filepath.Join(tmpDir, manifestFileName)
	_, err = os.Stat(manifestPath)
	assert.NoError(t, err, "Manifest file should exist")

	// Verify manifest can be loaded and is valid
	manifest, err := config.LoadManifest(manifestPath)
	require.NoError(t, err)
	assert.Equal(t, 1, manifest.Version)
	assert.NotEmpty(t, manifest.Guidelines, "Should have example guidelines")
	assert.NotEmpty(t, manifest.Prompts, "Should have example prompts")
}

func TestInitCmd_FileAlreadyExists(t *testing.T) {
	// Create temp directory and change to it
	tmpDir := t.TempDir()
	originalDir, err := os.Getwd()
	require.NoError(t, err)
	defer func() {
		err := os.Chdir(originalDir)
		require.NoError(t, err)
	}()

	err = os.Chdir(tmpDir)
	require.NoError(t, err)

	// Create an existing manifest file
	manifestPath := filepath.Join(tmpDir, manifestFileName)
	err = os.WriteFile(manifestPath, []byte("existing content"), 0644)
	require.NoError(t, err)

	// Run init should fail
	err = runInit()
	assert.Error(t, err, "Should error when manifest already exists")

	// Verify original file is unchanged
	content, err := os.ReadFile(manifestPath)
	require.NoError(t, err)
	assert.Equal(t, "existing content", string(content), "Original file should not be modified")
}

func TestInitCmd_CreatesValidManifest(t *testing.T) {
	// Create temp directory and change to it
	tmpDir := t.TempDir()
	originalDir, err := os.Getwd()
	require.NoError(t, err)
	defer func() {
		err := os.Chdir(originalDir)
		require.NoError(t, err)
	}()

	err = os.Chdir(tmpDir)
	require.NoError(t, err)

	// Run init
	err = runInit()
	require.NoError(t, err)

	// Load and verify manifest structure
	manifest, err := config.LoadManifest(manifestFileName)
	require.NoError(t, err)

	// Verify guidelines structure
	assert.NotEmpty(t, manifest.Guidelines)
	for i, guideline := range manifest.Guidelines {
		assert.NotEmpty(t, guideline.Name, "Guideline %d should have a name", i)
		assert.NotEmpty(t, guideline.File, "Guideline %d should have a file", i)
		assert.NotEmpty(t, guideline.Description, "Guideline %d should have a description", i)
		assert.NotEmpty(t, guideline.ApplicableScenarios, "Guideline %d should have applicable_scenarios", i)
		assert.True(t, hasPrefix(guideline.File, "guidelines/"), "Guideline file should be in guidelines/ directory")
	}

	// Verify prompts structure
	assert.NotEmpty(t, manifest.Prompts)
	for i, prompt := range manifest.Prompts {
		assert.NotEmpty(t, prompt.Name, "Prompt %d should have a name", i)
		assert.NotEmpty(t, prompt.File, "Prompt %d should have a file", i)
		assert.NotEmpty(t, prompt.Description, "Prompt %d should have a description", i)
		assert.True(t, hasPrefix(prompt.File, "prompts/"), "Prompt file should be in prompts/ directory")
	}
}

func hasPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}

func TestNewInitCmd(t *testing.T) {
	cmd := NewInitCmd()
	assert.NotNil(t, cmd)
	assert.Equal(t, "init", cmd.Use)
	assert.NotEmpty(t, cmd.Short)
	assert.NotEmpty(t, cmd.Long)
	assert.NotEmpty(t, cmd.Example)
}
