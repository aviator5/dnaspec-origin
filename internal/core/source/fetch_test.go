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

func TestFetchGitSource_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("fetch from git repository", func(t *testing.T) {
		// Use a small test repository
		url := "https://github.com/aviator5/dnaspec-test-repo.git"

		info, cleanup, err := FetchGitSource(url, "")
		if err != nil {
			// Skip if test repo not available
			if os.Getenv("CI") != "true" {
				t.Skipf("Test repository not available: %v", err)
			}
			t.Fatalf("FetchGitSource() error = %v", err)
		}
		defer cleanup()

		// Verify source info
		if info.SourceType != "git-repo" {
			t.Errorf("SourceType = %s, want git-repo", info.SourceType)
		}

		if info.URL != url {
			t.Errorf("URL = %s, want %s", info.URL, url)
		}

		if info.Commit == "" {
			t.Error("Expected non-empty commit hash")
		}

		// Verify source directory exists and is temporary
		if _, err := os.Stat(info.SourceDir); os.IsNotExist(err) {
			t.Error("Source directory does not exist")
		}

		// Verify cleanup works
		tempDir := info.SourceDir
		cleanup()

		// After cleanup, temp dir should be gone
		if _, err := os.Stat(tempDir); !os.IsNotExist(err) {
			t.Error("Temporary directory was not cleaned up")
		}
	})

	t.Run("fetch with specific ref", func(t *testing.T) {
		url := "https://github.com/aviator5/dnaspec-test-repo.git"

		info, cleanup, err := FetchGitSource(url, "main")
		if err != nil {
			if os.Getenv("CI") != "true" {
				t.Skipf("Test repository not available: %v", err)
			}
			t.Fatalf("FetchGitSource() error = %v", err)
		}
		defer cleanup()

		if info.Ref != "main" {
			t.Errorf("Ref = %s, want main", info.Ref)
		}
	})

	t.Run("error on invalid repository", func(t *testing.T) {
		url := "https://github.com/nonexistent/repo-does-not-exist-12345.git"

		_, cleanup, err := FetchGitSource(url, "")
		if cleanup != nil {
			defer cleanup()
		}

		if err == nil {
			t.Error("Expected error for invalid repository, got nil")
		}
	})

	t.Run("error on missing manifest in repo", func(t *testing.T) {
		// This would need a test repo without a manifest
		// For now, we can test the error path with a non-DNA repo
		url := "https://github.com/golang/example.git"

		_, cleanup, err := FetchGitSource(url, "")
		if cleanup != nil {
			defer cleanup()
		}

		if err == nil {
			t.Error("Expected error for repository without manifest, got nil")
		}
	})
}
