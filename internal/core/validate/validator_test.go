package validate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/aviator5/dnaspec/internal/core/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidator_Validate(t *testing.T) {
	tests := []struct {
		name         string
		manifest     *config.Manifest
		setupFiles   map[string]string // file path -> content
		wantErrCount int
		wantErrTypes []string // error types we expect
	}{
		{
			name: "valid manifest",
			manifest: &config.Manifest{
				Version: 1,
				Guidelines: []config.ManifestGuideline{
					{
						Name:                "go-style",
						File:                "guidelines/go-style.md",
						Description:         "Go style guide",
						ApplicableScenarios: []string{"Writing Go code"},
						Prompts:             []string{"code-review"},
					},
				},
				Prompts: []config.ManifestPrompt{
					{
						Name:        "code-review",
						File:        "prompts/code-review.md",
						Description: "Code review prompt",
					},
				},
			},
			setupFiles: map[string]string{
				"guidelines/go-style.md": "# Go Style Guide",
				"prompts/code-review.md": "# Code Review",
			},
			wantErrCount: 0,
		},
		{
			name: "missing guideline name",
			manifest: &config.Manifest{
				Version: 1,
				Guidelines: []config.ManifestGuideline{
					{
						File:                "guidelines/test.md",
						Description:         "Test",
						ApplicableScenarios: []string{"Testing"},
					},
				},
				Prompts: []config.ManifestPrompt{},
			},
			setupFiles: map[string]string{
				"guidelines/test.md": "# Test",
			},
			wantErrCount: 1,
			wantErrTypes: []string{"MissingRequiredField"},
		},
		{
			name: "missing guideline file",
			manifest: &config.Manifest{
				Version: 1,
				Guidelines: []config.ManifestGuideline{
					{
						Name:                "test",
						Description:         "Test",
						ApplicableScenarios: []string{"Testing"},
					},
				},
				Prompts: []config.ManifestPrompt{},
			},
			wantErrCount: 1,
			wantErrTypes: []string{"MissingRequiredField"},
		},
		{
			name: "missing prompt name",
			manifest: &config.Manifest{
				Version:    1,
				Guidelines: []config.ManifestGuideline{},
				Prompts: []config.ManifestPrompt{
					{
						File:        "prompts/test.md",
						Description: "Test",
					},
				},
			},
			setupFiles: map[string]string{
				"prompts/test.md": "# Test",
			},
			wantErrCount: 1,
			wantErrTypes: []string{"MissingRequiredField"},
		},
		{
			name: "guideline file does not exist",
			manifest: &config.Manifest{
				Version: 1,
				Guidelines: []config.ManifestGuideline{
					{
						Name:                "test",
						File:                "guidelines/nonexistent.md",
						Description:         "Test",
						ApplicableScenarios: []string{"Testing"},
					},
				},
				Prompts: []config.ManifestPrompt{},
			},
			wantErrCount: 1,
			wantErrTypes: []string{"FileNotFound"},
		},
		{
			name: "prompt file does not exist",
			manifest: &config.Manifest{
				Version:    1,
				Guidelines: []config.ManifestGuideline{},
				Prompts: []config.ManifestPrompt{
					{
						Name:        "test",
						File:        "prompts/nonexistent.md",
						Description: "Test",
					},
				},
			},
			wantErrCount: 1,
			wantErrTypes: []string{"FileNotFound"},
		},
		{
			name: "duplicate guideline names",
			manifest: &config.Manifest{
				Version: 1,
				Guidelines: []config.ManifestGuideline{
					{
						Name:                "test",
						File:                "guidelines/test1.md",
						Description:         "Test 1",
						ApplicableScenarios: []string{"Testing"},
					},
					{
						Name:                "test",
						File:                "guidelines/test2.md",
						Description:         "Test 2",
						ApplicableScenarios: []string{"Testing"},
					},
				},
				Prompts: []config.ManifestPrompt{},
			},
			setupFiles: map[string]string{
				"guidelines/test1.md": "# Test 1",
				"guidelines/test2.md": "# Test 2",
			},
			wantErrCount: 1,
			wantErrTypes: []string{"DuplicateName"},
		},
		{
			name: "duplicate prompt names",
			manifest: &config.Manifest{
				Version:    1,
				Guidelines: []config.ManifestGuideline{},
				Prompts: []config.ManifestPrompt{
					{
						Name:        "review",
						File:        "prompts/review1.md",
						Description: "Review 1",
					},
					{
						Name:        "review",
						File:        "prompts/review2.md",
						Description: "Review 2",
					},
				},
			},
			setupFiles: map[string]string{
				"prompts/review1.md": "# Review 1",
				"prompts/review2.md": "# Review 2",
			},
			wantErrCount: 1,
			wantErrTypes: []string{"DuplicateName"},
		},
		{
			name: "undefined prompt reference",
			manifest: &config.Manifest{
				Version: 1,
				Guidelines: []config.ManifestGuideline{
					{
						Name:                "test",
						File:                "guidelines/test.md",
						Description:         "Test",
						ApplicableScenarios: []string{"Testing"},
						Prompts:             []string{"nonexistent"},
					},
				},
				Prompts: []config.ManifestPrompt{},
			},
			setupFiles: map[string]string{
				"guidelines/test.md": "# Test",
			},
			wantErrCount: 1,
			wantErrTypes: []string{"UndefinedReference"},
		},
		{
			name: "guideline references multiple prompts",
			manifest: &config.Manifest{
				Version: 1,
				Guidelines: []config.ManifestGuideline{
					{
						Name:                "guide1",
						File:                "guidelines/guide1.md",
						Description:         "Guide 1",
						Prompts:             []string{"prompt1", "prompt2"},
						ApplicableScenarios: []string{"Testing"},
					},
				},
				Prompts: []config.ManifestPrompt{
					{
						Name:        "prompt1",
						File:        "prompts/prompt1.md",
						Description: "Prompt 1",
					},
					{
						Name:        "prompt2",
						File:        "prompts/prompt2.md",
						Description: "Prompt 2",
					},
				},
			},
			setupFiles: map[string]string{
				"guidelines/guide1.md": "# Guide 1",
				"prompts/prompt1.md":   "# Prompt 1",
				"prompts/prompt2.md":   "# Prompt 2",
			},
			wantErrCount: 0,
		},
		{
			name: "multiple validation errors",
			manifest: &config.Manifest{
				Version: 1,
				Guidelines: []config.ManifestGuideline{
					{
						Name:                "test",
						File:                "guidelines/missing.md",
						Description:         "Test",
						ApplicableScenarios: []string{"Testing"},
						Prompts:             []string{"nonexistent"},
					},
				},
				Prompts: []config.ManifestPrompt{
					{
						File:        "prompts/test.md",
						Description: "Missing name",
					},
				},
			},
			wantErrCount: 4, // missing file, undefined reference, missing prompt name, missing prompt file
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp directory
			tmpDir := t.TempDir()

			// Setup test files
			for path, content := range tt.setupFiles {
				fullPath := filepath.Join(tmpDir, path)
				err := os.MkdirAll(filepath.Dir(fullPath), 0755)
				require.NoError(t, err)
				err = os.WriteFile(fullPath, []byte(content), 0644)
				require.NoError(t, err)
			}

			// Validate
			errs := ValidateManifest(tt.manifest, tmpDir)

			// Check error count
			assert.Len(t, errs, tt.wantErrCount, "Unexpected error count. Errors: %v", errs)

			// Check error types if specified (by field prefix)
			for i, errType := range tt.wantErrTypes {
				if i >= len(errs) {
					break
				}
				validationErr := errs[i]
				var typeMatches bool
				switch errType {
				case "MissingRequiredField":
					typeMatches = strings.Contains(validationErr.Message, "missing required field")
				case "FileNotFound":
					typeMatches = strings.Contains(validationErr.Message, "file not found")
				case "DuplicateName":
					typeMatches = strings.Contains(validationErr.Message, "duplicate")
				case "UndefinedReference":
					typeMatches = strings.Contains(validationErr.Message, "non-existent")
				default:
					typeMatches = true // Unknown type, skip check
				}
				assert.True(t, typeMatches, "Error %d message '%s' doesn't match type %s", i, validationErr.Message, errType)
			}
		})
	}
}

func TestValidator_ValidateManifestFile(t *testing.T) {
	tests := []struct {
		name         string
		manifestYAML string
		setupFiles   map[string]string
		wantErrCount int
	}{
		{
			name: "valid manifest file",
			manifestYAML: `version: 1
guidelines:
  - name: test
    file: guidelines/test.md
    description: Test guideline
    applicable_scenarios:
      - Testing
prompts:
  - name: review
    file: prompts/review.md
    description: Review prompt`,
			setupFiles: map[string]string{
				"guidelines/test.md": "# Test",
				"prompts/review.md":  "# Review",
			},
			wantErrCount: 0,
		},
		{
			name: "invalid manifest file",
			manifestYAML: `version: 1
guidelines:
  - name: test
    file: guidelines/missing.md
    description: Test
    applicable_scenarios:
      - Testing
prompts: []`,
			wantErrCount: 1, // file not found
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp directory
			tmpDir := t.TempDir()

			// Write manifest file
			manifestPath := filepath.Join(tmpDir, "dnaspec-manifest.yaml")
			err := os.WriteFile(manifestPath, []byte(tt.manifestYAML), 0644)
			require.NoError(t, err)

			// Setup test files
			for path, content := range tt.setupFiles {
				fullPath := filepath.Join(tmpDir, path)
				err := os.MkdirAll(filepath.Dir(fullPath), 0755)
				require.NoError(t, err)
				err = os.WriteFile(fullPath, []byte(content), 0644)
				require.NoError(t, err)
			}

			// Load and validate
			manifest, err := config.LoadManifest(manifestPath)
			var errs ValidationErrors
			if err != nil {
				// If we can't load the manifest, that's a validation error
				errs.Add("manifest", err.Error())
			} else {
				errs = ValidateManifest(manifest, tmpDir)
			}

			// Check error count
			assert.Len(t, errs, tt.wantErrCount, "Unexpected error count. Errors: %v", errs)
		})
	}
}

func TestValidator_ValidateManifestFile_FileNotFound(t *testing.T) {
	tmpDir := t.TempDir()

	_, err := config.LoadManifest(filepath.Join(tmpDir, "nonexistent.yaml"))
	assert.Error(t, err)
}

func TestValidator_ValidateManifestFile_InvalidYAML(t *testing.T) {
	tmpDir := t.TempDir()
	manifestPath := filepath.Join(tmpDir, "dnaspec-manifest.yaml")

	// Write invalid YAML
	err := os.WriteFile(manifestPath, []byte("invalid: yaml: : :"), 0644)
	require.NoError(t, err)

	_, err = config.LoadManifest(manifestPath)
	assert.Error(t, err)
}

func TestValidator_PathSecurity(t *testing.T) {
	tests := []struct {
		name         string
		path         string
		expectedErr  bool
		errorMessage string
	}{
		{
			name:         "absolute path blocked",
			path:         "/etc/passwd",
			expectedErr:  true,
			errorMessage: "absolute paths not allowed",
		},
		{
			name:         "path traversal blocked",
			path:         "guidelines/../../../etc/passwd",
			expectedErr:  true,
			errorMessage: "path traversal not allowed",
		},
		{
			name:         "simple path traversal blocked",
			path:         "../etc/passwd",
			expectedErr:  true,
			errorMessage: "path traversal not allowed",
		},
		{
			name:         "wrong directory prefix blocked",
			path:         "other/file.md",
			expectedErr:  true,
			errorMessage: "path must be within",
		},
		{
			name:        "valid guideline path allowed",
			path:        "guidelines/test.md",
			expectedErr: false,
		},
		{
			name:        "valid prompt path allowed",
			path:        "prompts/test.md",
			expectedErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()

			// Create the file if it's expected to pass path validation
			if !tt.expectedErr {
				fullPath := filepath.Join(tmpDir, tt.path)
				err := os.MkdirAll(filepath.Dir(fullPath), 0755)
				require.NoError(t, err)
				err = os.WriteFile(fullPath, []byte("test"), 0644)
				require.NoError(t, err)
			}

			// Determine expected prefix based on path
			var expectedPrefix string
			if strings.HasPrefix(tt.path, "guidelines/") || strings.Contains(tt.name, "guideline") {
				expectedPrefix = "guidelines/"
			} else {
				expectedPrefix = "prompts/"
			}

			errs := validateFilePath(tt.path, "test.field", tmpDir, expectedPrefix)

			if tt.expectedErr {
				assert.NotEmpty(t, errs, "Expected validation error for path: %s", tt.path)
				if len(errs) > 0 && tt.errorMessage != "" {
					assert.Contains(t, errs[0].Message, tt.errorMessage)
				}
			} else {
				assert.Empty(t, errs, "Unexpected validation error for path: %s", tt.path)
			}
		})
	}
}

func TestValidator_NamingConventions(t *testing.T) {
	tests := []struct {
		name        string
		inputName   string
		expectError bool
	}{
		{
			name:        "valid spinal-case",
			inputName:   "go-style",
			expectError: false,
		},
		{
			name:        "valid single word",
			inputName:   "testing",
			expectError: false,
		},
		{
			name:        "valid with numbers",
			inputName:   "go-1-18-features",
			expectError: false,
		},
		{
			name:        "invalid uppercase",
			inputName:   "Go-Style",
			expectError: true,
		},
		{
			name:        "invalid camelCase",
			inputName:   "goStyle",
			expectError: true,
		},
		{
			name:        "invalid snake_case",
			inputName:   "go_style",
			expectError: true,
		},
		{
			name:        "invalid spaces",
			inputName:   "go style",
			expectError: true,
		},
		{
			name:        "invalid leading hyphen",
			inputName:   "-go-style",
			expectError: true,
		},
		{
			name:        "invalid trailing hyphen",
			inputName:   "go-style-",
			expectError: true,
		},
		{
			name:        "invalid double hyphen",
			inputName:   "go--style",
			expectError: true,
		},
		{
			name:        "invalid starting with number",
			inputName:   "1-go-style",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()

			// Create the file
			guidelinePath := "guidelines/test.md"
			fullPath := filepath.Join(tmpDir, guidelinePath)
			err := os.MkdirAll(filepath.Dir(fullPath), 0755)
			require.NoError(t, err)
			err = os.WriteFile(fullPath, []byte("test"), 0644)
			require.NoError(t, err)

			manifest := &config.Manifest{
				Version: 1,
				Guidelines: []config.ManifestGuideline{
					{
						Name:                tt.inputName,
						File:                guidelinePath,
						Description:         "Test guideline",
						ApplicableScenarios: []string{"Testing"},
					},
				},
			}

			errs := ValidateManifest(manifest, tmpDir)

			if tt.expectError {
				assert.NotEmpty(t, errs, "Expected validation error for name: %s", tt.inputName)
				// Check that at least one error mentions naming format
				hasNamingError := false
				for _, err := range errs {
					if strings.Contains(err.Message, "invalid naming format") {
						hasNamingError = true
						break
					}
				}
				assert.True(t, hasNamingError, "Expected naming format error for: %s", tt.inputName)
			} else {
				// Filter out only naming errors
				namingErrors := []ValidationError{}
				for _, err := range errs {
					if strings.Contains(err.Message, "invalid naming format") {
						namingErrors = append(namingErrors, err)
					}
				}
				assert.Empty(t, namingErrors, "Unexpected naming error for: %s", tt.inputName)
			}
		})
	}
}

func TestValidator_EmptyApplicableScenarios(t *testing.T) {
	tmpDir := t.TempDir()

	// Create the guideline file
	guidelinePath := "guidelines/test.md"
	fullPath := filepath.Join(tmpDir, guidelinePath)
	err := os.MkdirAll(filepath.Dir(fullPath), 0755)
	require.NoError(t, err)
	err = os.WriteFile(fullPath, []byte("test"), 0644)
	require.NoError(t, err)

	manifest := &config.Manifest{
		Version: 1,
		Guidelines: []config.ManifestGuideline{
			{
				Name:                "test-guideline",
				File:                guidelinePath,
				Description:         "Test guideline",
				ApplicableScenarios: []string{}, // Empty!
			},
		},
	}

	errs := ValidateManifest(manifest, tmpDir)

	assert.NotEmpty(t, errs, "Expected validation error for empty applicable_scenarios")

	// Check that error mentions applicable_scenarios
	hasApplicableError := false
	for _, err := range errs {
		if strings.Contains(err.Message, "applicable_scenarios") {
			hasApplicableError = true
			assert.Contains(t, err.Message, "AGENTS.md", "Error should mention AGENTS.md")
			break
		}
	}
	assert.True(t, hasApplicableError, "Expected applicable_scenarios error")
}

func TestValidator_MissingVersion(t *testing.T) {
	manifest := &config.Manifest{
		Version:    0, // Missing/zero version
		Guidelines: []config.ManifestGuideline{},
		Prompts:    []config.ManifestPrompt{},
	}

	errs := ValidateManifest(manifest, t.TempDir())

	assert.NotEmpty(t, errs, "Expected validation error for missing version")

	hasVersionError := false
	for _, err := range errs {
		if strings.Contains(err.Field, "version") {
			hasVersionError = true
			assert.Contains(t, err.Message, "missing required field")
			break
		}
	}
	assert.True(t, hasVersionError, "Expected version error")
}

func TestValidator_ComplexScenario(t *testing.T) {
	tmpDir := t.TempDir()

	// Setup files
	files := map[string]string{
		"guidelines/go-style.md":   "# Go Style",
		"guidelines/rest-api.md":   "# REST API",
		"prompts/code-review.md":   "# Code Review",
		"prompts/api-design.md":    "# API Design",
		"prompts/documentation.md": "# Documentation",
	}

	for path, content := range files {
		fullPath := filepath.Join(tmpDir, path)
		err := os.MkdirAll(filepath.Dir(fullPath), 0755)
		require.NoError(t, err)
		err = os.WriteFile(fullPath, []byte(content), 0644)
		require.NoError(t, err)
	}

	manifest := &config.Manifest{
		Version: 1,
		Guidelines: []config.ManifestGuideline{
			{
				Name:                "go-style",
				File:                "guidelines/go-style.md",
				Description:         "Go coding style guidelines",
				ApplicableScenarios: []string{"Writing Go code", "Code reviews"},
				Prompts:             []string{"code-review", "documentation"},
			},
			{
				Name:                "rest-api",
				File:                "guidelines/rest-api.md",
				Description:         "REST API design principles",
				ApplicableScenarios: []string{"Designing APIs", "API documentation"},
				Prompts:             []string{"api-design"},
			},
		},
		Prompts: []config.ManifestPrompt{
			{
				Name:        "code-review",
				File:        "prompts/code-review.md",
				Description: "Code review checklist",
			},
			{
				Name:        "api-design",
				File:        "prompts/api-design.md",
				Description: "API design template",
			},
			{
				Name:        "documentation",
				File:        "prompts/documentation.md",
				Description: "Documentation standards",
			},
		},
	}

	errs := ValidateManifest(manifest, tmpDir)
	assert.Empty(t, errs, "Valid complex manifest should have no errors")
}
