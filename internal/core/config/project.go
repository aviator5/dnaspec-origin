package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// ProjectConfig represents the dnaspec.yaml structure
type ProjectConfig struct {
	Version int             `yaml:"version"`
	Agents  []string        `yaml:"agents,omitempty"`
	Sources []ProjectSource `yaml:"sources,omitempty"`
}

// ProjectSource represents a DNA source in the project configuration
type ProjectSource struct {
	Name       string              `yaml:"name"`
	Type       string              `yaml:"type"` // "git-repo" or "local-path"
	URL        string              `yaml:"url,omitempty"`
	Path       string              `yaml:"path,omitempty"`
	Ref        string              `yaml:"ref,omitempty"`
	Commit     string              `yaml:"commit,omitempty"`
	Guidelines []ProjectGuideline  `yaml:"guidelines,omitempty"`
	Prompts    []ProjectPrompt     `yaml:"prompts,omitempty"`
}

// ProjectGuideline represents a guideline in the project configuration
type ProjectGuideline struct {
	Name                string   `yaml:"name"`
	File                string   `yaml:"file"`
	Description         string   `yaml:"description"`
	ApplicableScenarios []string `yaml:"applicable_scenarios,omitempty"`
	Prompts             []string `yaml:"prompts,omitempty"`
}

// ProjectPrompt represents a prompt in the project configuration
type ProjectPrompt struct {
	Name        string `yaml:"name"`
	File        string `yaml:"file"`
	Description string `yaml:"description"`
}

// LoadProjectConfig loads and parses a project config file from the given path
func LoadProjectConfig(path string) (*ProjectConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config ProjectConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// SaveProjectConfig writes a project config to the given path
func SaveProjectConfig(path string, config *ProjectConfig) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// AtomicWriteProjectConfig writes a project config atomically using a temp file
func AtomicWriteProjectConfig(path string, config *ProjectConfig) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	// Write to temp file
	tmpFile := path + ".tmp"
	if err := os.WriteFile(tmpFile, data, 0644); err != nil {
		return err
	}

	// Atomic rename
	return os.Rename(tmpFile, path)
}
