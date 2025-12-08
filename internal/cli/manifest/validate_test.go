package manifest

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateCmd_Success(t *testing.T) {
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

	// Create a valid manifest
	manifestContent := `version: 1
guidelines:
  - name: test-guideline
    file: guidelines/test.md
    description: Test guideline
    applicable_scenarios:
      - Testing
prompts:
  - name: test-prompt
    file: prompts/test.md
    description: Test prompt
`
	err = os.WriteFile(manifestFileName, []byte(manifestContent), 0644)
	require.NoError(t, err)

	// Create the referenced files
	err = os.MkdirAll("guidelines", 0755)
	require.NoError(t, err)
	err = os.WriteFile("guidelines/test.md", []byte("# Test Guideline"), 0644)
	require.NoError(t, err)

	err = os.MkdirAll("prompts", 0755)
	require.NoError(t, err)
	err = os.WriteFile("prompts/test.md", []byte("# Test Prompt"), 0644)
	require.NoError(t, err)

	// Run validate
	err = runValidate()
	assert.NoError(t, err, "Valid manifest should pass validation")
}

func TestValidateCmd_ManifestNotFound(t *testing.T) {
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

	// Run validate without a manifest file
	err = runValidate()
	assert.Error(t, err, "Should error when manifest file doesn't exist")
}

func TestValidateCmd_InvalidYAML(t *testing.T) {
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

	// Create an invalid YAML manifest
	err = os.WriteFile(manifestFileName, []byte("invalid: yaml: : :"), 0644)
	require.NoError(t, err)

	// Run validate
	err = runValidate()
	assert.Error(t, err, "Should error on invalid YAML")
}

func TestValidateCmd_MissingFiles(t *testing.T) {
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

	// Create a manifest with missing file references
	manifestContent := `version: 1
guidelines:
  - name: test-guideline
    file: guidelines/missing.md
    description: Test guideline
    applicable_scenarios:
      - Testing
prompts: []
`
	err = os.WriteFile(manifestFileName, []byte(manifestContent), 0644)
	require.NoError(t, err)

	// Run validate (files don't exist)
	err = runValidate()
	assert.Error(t, err, "Should error when referenced files don't exist")
}

func TestValidateCmd_InvalidNaming(t *testing.T) {
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

	// Create a manifest with invalid naming
	manifestContent := `version: 1
guidelines:
  - name: Test_Guideline
    file: guidelines/test.md
    description: Test guideline
    applicable_scenarios:
      - Testing
prompts: []
`
	err = os.WriteFile(manifestFileName, []byte(manifestContent), 0644)
	require.NoError(t, err)

	// Create the file
	err = os.MkdirAll("guidelines", 0755)
	require.NoError(t, err)
	err = os.WriteFile("guidelines/test.md", []byte("# Test"), 0644)
	require.NoError(t, err)

	// Run validate
	err = runValidate()
	assert.Error(t, err, "Should error on invalid naming convention")
}

func TestValidateCmd_EmptyApplicableScenarios(t *testing.T) {
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

	// Create a manifest with empty applicable_scenarios
	manifestContent := `version: 1
guidelines:
  - name: test-guideline
    file: guidelines/test.md
    description: Test guideline
    applicable_scenarios: []
prompts: []
`
	err = os.WriteFile(manifestFileName, []byte(manifestContent), 0644)
	require.NoError(t, err)

	// Create the file
	err = os.MkdirAll("guidelines", 0755)
	require.NoError(t, err)
	err = os.WriteFile("guidelines/test.md", []byte("# Test"), 0644)
	require.NoError(t, err)

	// Run validate
	err = runValidate()
	assert.Error(t, err, "Should error on empty applicable_scenarios")
}

func TestValidateCmd_UndefinedPromptReference(t *testing.T) {
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

	// Create a manifest with undefined prompt reference
	manifestContent := `version: 1
guidelines:
  - name: test-guideline
    file: guidelines/test.md
    description: Test guideline
    applicable_scenarios:
      - Testing
    prompts:
      - nonexistent-prompt
prompts: []
`
	err = os.WriteFile(manifestFileName, []byte(manifestContent), 0644)
	require.NoError(t, err)

	// Create the file
	err = os.MkdirAll("guidelines", 0755)
	require.NoError(t, err)
	err = os.WriteFile("guidelines/test.md", []byte("# Test"), 0644)
	require.NoError(t, err)

	// Run validate
	err = runValidate()
	assert.Error(t, err, "Should error on undefined prompt reference")
}

func TestValidateCmd_PathTraversal(t *testing.T) {
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

	// Create a manifest with path traversal attempt
	manifestContent := `version: 1
guidelines:
  - name: test-guideline
    file: guidelines/../../etc/passwd
    description: Test guideline
    applicable_scenarios:
      - Testing
prompts: []
`
	err = os.WriteFile(manifestFileName, []byte(manifestContent), 0644)
	require.NoError(t, err)

	// Run validate
	err = runValidate()
	assert.Error(t, err, "Should error on path traversal attempt")
}

func TestValidateCmd_ComplexValid(t *testing.T) {
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

	// Create a complex valid manifest
	manifestContent := `version: 1
guidelines:
  - name: go-style
    file: guidelines/go-style.md
    description: Go coding style guidelines
    applicable_scenarios:
      - Writing Go code
      - Code reviews
    prompts:
      - code-review
  - name: rest-api
    file: guidelines/rest-api.md
    description: REST API design
    applicable_scenarios:
      - Designing APIs
prompts:
  - name: code-review
    file: prompts/code-review.md
    description: Code review checklist
`
	err = os.WriteFile(manifestFileName, []byte(manifestContent), 0644)
	require.NoError(t, err)

	// Create the referenced files
	err = os.MkdirAll("guidelines", 0755)
	require.NoError(t, err)
	err = os.WriteFile("guidelines/go-style.md", []byte("# Go Style"), 0644)
	require.NoError(t, err)
	err = os.WriteFile("guidelines/rest-api.md", []byte("# REST API"), 0644)
	require.NoError(t, err)

	err = os.MkdirAll("prompts", 0755)
	require.NoError(t, err)
	err = os.WriteFile("prompts/code-review.md", []byte("# Code Review"), 0644)
	require.NoError(t, err)

	// Run validate
	err = runValidate()
	assert.NoError(t, err, "Complex valid manifest should pass validation")
}

func TestNewValidateCmd(t *testing.T) {
	cmd := NewValidateCmd()
	assert.NotNil(t, cmd)
	assert.Equal(t, "validate", cmd.Use)
	assert.NotEmpty(t, cmd.Short)
	assert.NotEmpty(t, cmd.Long)
	assert.NotEmpty(t, cmd.Example)
}
