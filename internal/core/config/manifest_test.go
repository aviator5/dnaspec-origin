package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestManifest_UnmarshalYAML(t *testing.T) {
	tests := []struct {
		name    string
		yaml    string
		want    Manifest
		wantErr bool
	}{
		{
			name: "valid manifest",
			yaml: `version: 1
guidelines:
  - name: go-style
    file: guidelines/go-style.md
    description: Go style guide
    applicable_scenarios:
      - Writing Go code
    prompts:
      - code-review
prompts:
  - name: code-review
    file: prompts/code-review.md
    description: Code review prompt`,
			want: Manifest{
				Version: 1,
				Guidelines: []ManifestGuideline{
					{
						Name:                "go-style",
						File:                "guidelines/go-style.md",
						Description:         "Go style guide",
						ApplicableScenarios: []string{"Writing Go code"},
						Prompts:             []string{"code-review"},
					},
				},
				Prompts: []ManifestPrompt{
					{
						Name:        "code-review",
						File:        "prompts/code-review.md",
						Description: "Code review prompt",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "guideline without prompts",
			yaml: `version: 1
guidelines:
  - name: rest-api
    file: guidelines/rest-api.md
    description: REST API design
    applicable_scenarios:
      - Designing APIs
prompts: []`,
			want: Manifest{
				Version: 1,
				Guidelines: []ManifestGuideline{
					{
						Name:                "rest-api",
						File:                "guidelines/rest-api.md",
						Description:         "REST API design",
						ApplicableScenarios: []string{"Designing APIs"},
						Prompts:             nil,
					},
				},
				Prompts: []ManifestPrompt{},
			},
			wantErr: false,
		},
		{
			name: "empty manifest",
			yaml: `version: 1
guidelines: []
prompts: []`,
			want: Manifest{
				Version:    1,
				Guidelines: []ManifestGuideline{},
				Prompts:    []ManifestPrompt{},
			},
			wantErr: false,
		},
		{
			name:    "invalid yaml",
			yaml:    `invalid: yaml: structure:`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got Manifest
			err := yaml.Unmarshal([]byte(tt.yaml), &got)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want.Version, got.Version)
			assert.Equal(t, tt.want.Guidelines, got.Guidelines)
			assert.Equal(t, tt.want.Prompts, got.Prompts)
		})
	}
}

func TestManifest_MarshalYAML(t *testing.T) {
	manifest := Manifest{
		Version: 1,
		Guidelines: []ManifestGuideline{
			{
				Name:                "go-style",
				File:                "guidelines/go-style.md",
				Description:         "Go style guide",
				ApplicableScenarios: []string{"Writing Go code"},
				Prompts:             []string{"code-review"},
			},
		},
		Prompts: []ManifestPrompt{
			{
				Name:        "code-review",
				File:        "prompts/code-review.md",
				Description: "Code review prompt",
			},
		},
	}

	data, err := yaml.Marshal(&manifest)
	require.NoError(t, err)

	// Unmarshal to verify
	var got Manifest
	err = yaml.Unmarshal(data, &got)
	require.NoError(t, err)

	assert.Equal(t, manifest.Version, got.Version)
	assert.Equal(t, manifest.Guidelines, got.Guidelines)
	assert.Equal(t, manifest.Prompts, got.Prompts)
}

func TestLoadManifest(t *testing.T) {
	tests := []struct {
		name    string
		content string
		wantErr bool
	}{
		{
			name: "valid manifest file",
			content: `version: 1
guidelines:
  - name: test
    file: guidelines/test.md
    description: Test
    applicable_scenarios:
      - Testing
prompts: []`,
			wantErr: false,
		},
		{
			name:    "invalid yaml",
			content: `invalid: yaml: : :`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp file
			tmpDir := t.TempDir()
			manifestPath := filepath.Join(tmpDir, "dnaspec-manifest.yaml")
			err := os.WriteFile(manifestPath, []byte(tt.content), 0644)
			require.NoError(t, err)

			_, err = LoadManifest(manifestPath)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestLoadManifest_FileNotFound(t *testing.T) {
	_, err := LoadManifest("/nonexistent/path/manifest.yaml")
	assert.Error(t, err)
}

func TestSaveManifest(t *testing.T) {
	manifest := &Manifest{
		Version: 1,
		Guidelines: []ManifestGuideline{
			{
				Name:                "test",
				File:                "guidelines/test.md",
				Description:         "Test guideline",
				ApplicableScenarios: []string{"Testing"},
			},
		},
		Prompts: []ManifestPrompt{},
	}

	tmpDir := t.TempDir()
	manifestPath := filepath.Join(tmpDir, "dnaspec-manifest.yaml")

	// Save manifest
	err := SaveManifest(manifestPath, manifest)
	require.NoError(t, err)

	// Load and verify
	loaded, err := LoadManifest(manifestPath)
	require.NoError(t, err)

	assert.Equal(t, manifest.Version, loaded.Version)
	assert.Equal(t, manifest.Guidelines, loaded.Guidelines)
	assert.Equal(t, manifest.Prompts, loaded.Prompts)
}

func TestCreateExampleManifest(t *testing.T) {
	tmpDir := t.TempDir()
	manifestPath := filepath.Join(tmpDir, "dnaspec-manifest.yaml")

	err := CreateExampleManifest(manifestPath)
	require.NoError(t, err)

	// Verify file exists and can be loaded
	manifest, err := LoadManifest(manifestPath)
	require.NoError(t, err)

	assert.Equal(t, 1, manifest.Version)
	assert.NotEmpty(t, manifest.Guidelines, "Example manifest should have guidelines")
	assert.NotEmpty(t, manifest.Prompts, "Example manifest should have prompts")
}
