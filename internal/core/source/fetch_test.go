package source

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFetchLocalSource(t *testing.T) {
	t.Run("fetch valid local source", func(t *testing.T) {
		// Use the test fixture
		testdataPath := filepath.Join("testdata", "valid-repo")

		// Get absolute path
		absPath, err := filepath.Abs(testdataPath)
		if err != nil {
			t.Fatalf("Failed to get absolute path: %v", err)
		}

		info, err := FetchLocalSource(absPath)
		if err != nil {
			t.Fatalf("FetchLocalSource() error = %v", err)
		}

		// Verify source info
		if info.SourceType != "local-path" {
			t.Errorf("SourceType = %s, want local-path", info.SourceType)
		}

		if info.SourceDir != absPath {
			t.Errorf("SourceDir = %s, want %s", info.SourceDir, absPath)
		}

		if info.Path != absPath {
			t.Errorf("Path = %s, want %s", info.Path, absPath)
		}

		// Verify manifest was loaded
		if info.Manifest == nil {
			t.Fatal("Manifest is nil")
		}

		if len(info.Manifest.Guidelines) != 2 {
			t.Errorf("Expected 2 guidelines, got %d", len(info.Manifest.Guidelines))
		}

		if len(info.Manifest.Prompts) != 1 {
			t.Errorf("Expected 1 prompt, got %d", len(info.Manifest.Prompts))
		}
	})

	t.Run("error on nonexistent path", func(t *testing.T) {
		_, err := FetchLocalSource("/nonexistent/path")
		if err == nil {
			t.Error("Expected error for nonexistent path, got nil")
		}
	})

	t.Run("error on file instead of directory", func(t *testing.T) {
		// Create a temp file
		tmpFile, _ := os.CreateTemp("", "testfile")
		defer os.Remove(tmpFile.Name())
		tmpFile.Close()

		_, err := FetchLocalSource(tmpFile.Name())
		if err == nil {
			t.Error("Expected error for file instead of directory, got nil")
		}
	})

	t.Run("error on missing manifest", func(t *testing.T) {
		// Create a directory without a manifest
		tmpDir := t.TempDir()

		_, err := FetchLocalSource(tmpDir)
		if err == nil {
			t.Error("Expected error for missing manifest, got nil")
		}
	})

	t.Run("error on invalid manifest", func(t *testing.T) {
		tmpDir := t.TempDir()

		// Create invalid manifest
		manifestPath := filepath.Join(tmpDir, "dnaspec-manifest.yaml")
		os.WriteFile(manifestPath, []byte("invalid: [[["), 0644)

		_, err := FetchLocalSource(tmpDir)
		if err == nil {
			t.Error("Expected error for invalid manifest, got nil")
		}
	})
}

func TestFetchLocalSource_RelativePath(t *testing.T) {
	// Test with relative path
	info, err := FetchLocalSource("testdata/valid-repo")
	if err != nil {
		t.Fatalf("FetchLocalSource() with relative path error = %v", err)
	}

	// Should convert to absolute path
	if !filepath.IsAbs(info.SourceDir) {
		t.Error("SourceDir should be absolute path")
	}

	if !filepath.IsAbs(info.Path) {
		t.Error("Path should be absolute path")
	}
}

func TestFetchLocalSource_Validation(t *testing.T) {
	t.Run("manifest validation runs", func(t *testing.T) {
		tmpDir := t.TempDir()

		// Create manifest with missing required fields
		manifestPath := filepath.Join(tmpDir, "dnaspec-manifest.yaml")
		invalidManifest := `version: 1
guidelines:
  - name: invalid-guideline
    file: guidelines/test.md
    description: Test
    applicable_scenarios:
      - "test"
`
		os.WriteFile(manifestPath, []byte(invalidManifest), 0644)

		_, err := FetchLocalSource(tmpDir)
		if err == nil {
			t.Error("Expected validation error for missing guideline file, got nil")
		}

		// Error should mention validation
		if err != nil && err.Error() != "" {
			t.Logf("Validation error: %s", err.Error())
		}
	})
}
