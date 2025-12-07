package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Manifest represents the dnaspec-manifest.yaml structure
type Manifest struct {
	Version    int                 `yaml:"version"`
	Guidelines []ManifestGuideline `yaml:"guidelines"`
	Prompts    []ManifestPrompt    `yaml:"prompts"`
}

// ManifestGuideline represents a single guideline entry
type ManifestGuideline struct {
	Name                string   `yaml:"name"`
	File                string   `yaml:"file"`
	Description         string   `yaml:"description"`
	ApplicableScenarios []string `yaml:"applicable_scenarios"`
	Prompts             []string `yaml:"prompts,omitempty"`
}

// ManifestPrompt represents a single prompt entry
type ManifestPrompt struct {
	Name        string `yaml:"name"`
	File        string `yaml:"file"`
	Description string `yaml:"description"`
}

// LoadManifest loads and parses a manifest file from the given path
func LoadManifest(path string) (*Manifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var manifest Manifest
	if err := yaml.Unmarshal(data, &manifest); err != nil {
		return nil, err
	}

	return &manifest, nil
}

// SaveManifest writes a manifest to the given path
func SaveManifest(path string, manifest *Manifest) error {
	data, err := yaml.Marshal(manifest)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
