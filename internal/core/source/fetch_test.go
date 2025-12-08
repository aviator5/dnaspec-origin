package source

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
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
		defer func() {
			err := os.Remove(tmpFile.Name())
			require.NoError(t, err)
		}()
		err := tmpFile.Close()
		require.NoError(t, err)

		_, err = FetchLocalSource(tmpFile.Name())
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
		err := os.WriteFile(manifestPath, []byte("invalid: [[["), 0644)
		require.NoError(t, err)

		_, err = FetchLocalSource(tmpDir)
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
		err := os.WriteFile(manifestPath, []byte(invalidManifest), 0644)
		require.NoError(t, err)

		_, err = FetchLocalSource(tmpDir)
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

	// Create a local git repository
	repoDir := t.TempDir()

	// Initialize repo
	runGit(t, repoDir, "init")
	runGit(t, repoDir, "config", "user.email", "test@example.com")
	runGit(t, repoDir, "config", "user.name", "Test User")
	// Set default branch to main to avoid confusion
	runGit(t, repoDir, "branch", "-m", "main")

	// Create manifest
	manifestContent := `version: 1
guidelines: []
prompts: []
`
	err := os.WriteFile(filepath.Join(repoDir, "dnaspec-manifest.yaml"), []byte(manifestContent), 0644)
	require.NoError(t, err)

	// Commit
	runGit(t, repoDir, "add", ".")
	runGit(t, repoDir, "commit", "-m", "Initial commit")

	// Get head hash
	headHash := getGitHead(t, repoDir)

	t.Run("fetch from git repository", func(t *testing.T) {
		url := "file://" + repoDir

		info, cleanup, err := FetchGitSource(url, "")
		if err != nil {
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

		if info.Commit != headHash {
			t.Errorf("Commit = %s, want %s", info.Commit, headHash)
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
		url := "file://" + repoDir

		info, cleanup, err := FetchGitSource(url, "main")
		if err != nil {
			t.Fatalf("FetchGitSource() error = %v", err)
		}
		defer cleanup()

		if info.Ref != "main" {
			t.Errorf("Ref = %s, want main", info.Ref)
		}

		if info.Commit != headHash {
			t.Errorf("Commit = %s, want %s", info.Commit, headHash)
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
		// Use the same local repo but delete the manifest and commit
		// Or create a new one
		badRepoDir := t.TempDir()
		runGit(t, badRepoDir, "init")
		runGit(t, badRepoDir, "config", "user.email", "test@example.com")
		runGit(t, badRepoDir, "config", "user.name", "Test User")
		runGit(t, badRepoDir, "branch", "-m", "main")
		runGit(t, badRepoDir, "commit", "--allow-empty", "-m", "Empty commit")

		url := "file://" + badRepoDir
		_, cleanup, err := FetchGitSource(url, "")
		if cleanup != nil {
			defer cleanup()
		}

		if err == nil {
			t.Error("Expected error for repository without manifest, got nil")
		}
	})
}

func runGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v failed: %v\nOutput: %s", args, err, out)
	}
}

func getGitHead(t *testing.T, dir string) string {
	t.Helper()
	cmd := exec.Command("git", "rev-parse", "HEAD")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("git rev-parse HEAD failed: %v", err)
	}
	return strings.TrimSpace(string(out))
}
