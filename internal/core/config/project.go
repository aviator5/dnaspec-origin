package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/aviator5/dnaspec/internal/core/paths"
	"gopkg.in/yaml.v3"
)

// Source types
const (
	SourceTypeGitRepo   = "git-repo"
	SourceTypeLocalPath = "local-path"
)

// ProjectConfig represents the dnaspec.yaml structure
type ProjectConfig struct {
	Version int             `yaml:"version"`
	Agents  []string        `yaml:"agents,omitempty"`
	Sources []ProjectSource `yaml:"sources,omitempty"`
}

// ProjectSource represents a DNA source in the project configuration
type ProjectSource struct {
	Name       string             `yaml:"name"`
	Type       string             `yaml:"type"` // "git-repo" or "local-path"
	URL        string             `yaml:"url,omitempty"`
	Path       string             `yaml:"path,omitempty"`
	Ref        string             `yaml:"ref,omitempty"`
	Commit     string             `yaml:"commit,omitempty"`
	Guidelines []ProjectGuideline `yaml:"guidelines,omitempty"`
	Prompts    []ProjectPrompt    `yaml:"prompts,omitempty"`
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

	// Note: Absolute path validation is handled by 'dnaspec validate' command
	// We don't warn here to avoid noise in every command

	return &config, nil
}

// SaveProjectConfig writes a project config to the given path
func SaveProjectConfig(path string, config *ProjectConfig) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0o644)
}

// AtomicWriteProjectConfig writes a project config atomically using a temp file
func AtomicWriteProjectConfig(path string, config *ProjectConfig) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	// Write to temp file
	tmpFile := path + ".tmp"
	if err := os.WriteFile(tmpFile, data, 0o644); err != nil {
		return err
	}

	// Atomic rename
	return os.Rename(tmpFile, path)
}

// MigrateToRelativePaths converts absolute paths to relative paths in-place.
// Returns error if any conversion fails.
func (c *ProjectConfig) MigrateToRelativePaths(projectRoot string) error {
	for i := range c.Sources {
		source := &c.Sources[i]
		if source.Type == SourceTypeLocalPath && filepath.IsAbs(source.Path) {
			relPath, err := paths.MakeRelative(projectRoot, source.Path)
			if err != nil {
				return fmt.Errorf("source %s: %w", source.Name, err)
			}
			source.Path = relPath
		}
	}
	return nil
}
