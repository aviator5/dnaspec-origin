package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadProjectConfig(t *testing.T) {
	tests := []struct {
		name    string
		content string
		wantErr bool
	}{
		{
			name: "valid config",
			content: `version: 1
agents:
  - claude-code
sources:
  - name: test-source
    type: git-repo
    url: https://github.com/test/repo
    ref: v1.0.0
    commit: abc123
    guidelines:
      - name: test-guideline
        file: guidelines/test.md
        description: Test guideline
`,
			wantErr: false,
		},
		{
			name:    "empty config",
			content: `version: 1`,
			wantErr: false,
		},
		{
			name:    "invalid yaml",
			content: `invalid: [[[`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpFile := filepath.Join(t.TempDir(), "dnaspec.yaml")
			if err := os.WriteFile(tmpFile, []byte(tt.content), 0644); err != nil {
				t.Fatalf("failed to write test file: %v", err)
			}

			config, err := LoadProjectConfig(tmpFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadProjectConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && config == nil {
				t.Error("LoadProjectConfig() returned nil config")
			}
		})
	}
}

func TestSaveProjectConfig(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "dnaspec.yaml")

	config := &ProjectConfig{
		Version: 1,
		Agents:  []string{"claude-code"},
		Sources: []ProjectSource{
			{
				Name:   "test-source",
				Type:   "git-repo",
				URL:    "https://github.com/test/repo",
				Ref:    "v1.0.0",
				Commit: "abc123",
			},
		},
	}

	if err := SaveProjectConfig(tmpFile, config); err != nil {
		t.Fatalf("SaveProjectConfig() error = %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(tmpFile); os.IsNotExist(err) {
		t.Error("SaveProjectConfig() did not create file")
	}

	// Load back and verify
	loaded, err := LoadProjectConfig(tmpFile)
	if err != nil {
		t.Fatalf("LoadProjectConfig() error = %v", err)
	}

	if loaded.Version != config.Version {
		t.Errorf("Version mismatch: got %d, want %d", loaded.Version, config.Version)
	}

	if len(loaded.Agents) != len(config.Agents) {
		t.Errorf("Agents length mismatch: got %d, want %d", len(loaded.Agents), len(config.Agents))
	}

	if len(loaded.Sources) != len(config.Sources) {
		t.Errorf("Sources length mismatch: got %d, want %d", len(loaded.Sources), len(config.Sources))
	}
}

func TestAtomicWriteProjectConfig(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "dnaspec.yaml")

	config := &ProjectConfig{
		Version: 1,
		Agents:  []string{"claude-code"},
	}

	if err := AtomicWriteProjectConfig(tmpFile, config); err != nil {
		t.Fatalf("AtomicWriteProjectConfig() error = %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(tmpFile); os.IsNotExist(err) {
		t.Error("AtomicWriteProjectConfig() did not create file")
	}

	// Verify temp file was cleaned up
	tmpTempFile := tmpFile + ".tmp"
	if _, err := os.Stat(tmpTempFile); !os.IsNotExist(err) {
		t.Error("AtomicWriteProjectConfig() did not clean up temp file")
	}

	// Load back and verify
	loaded, err := LoadProjectConfig(tmpFile)
	if err != nil {
		t.Fatalf("LoadProjectConfig() error = %v", err)
	}

	if loaded.Version != config.Version {
		t.Errorf("Version mismatch: got %d, want %d", loaded.Version, config.Version)
	}
}
