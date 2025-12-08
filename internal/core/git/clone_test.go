package git

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCloneRepo_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Disable git terminal prompts to prevent hanging on authentication
	t.Setenv("GIT_TERMINAL_PROMPT", "0")

	t.Run("clone public repository", func(t *testing.T) {
		destDir := t.TempDir()

		// Use a small, stable public repository
		url := "https://github.com/aviator5/dnaspec-test-repo.git"
		commit, err := CloneRepo(url, "", destDir)

		if err != nil {
			// If the test repo doesn't exist, skip this test
			if os.Getenv("CI") != "true" {
				t.Skipf("Test repository not available: %v", err)
			}
			t.Fatalf("CloneRepo() error = %v", err)
		}

		if commit == "" {
			t.Error("Expected non-empty commit hash")
		}

		// Verify .git directory exists
		gitDir := filepath.Join(destDir, ".git")
		if _, err := os.Stat(gitDir); os.IsNotExist(err) {
			t.Error("Expected .git directory to exist")
		}
	})

	t.Run("clone with specific ref", func(t *testing.T) {
		destDir := t.TempDir()

		// Use a repository with tags
		url := "https://github.com/aviator5/dnaspec-test-repo.git"
		commit, err := CloneRepo(url, "main", destDir)

		if err != nil {
			if os.Getenv("CI") != "true" {
				t.Skipf("Test repository not available: %v", err)
			}
			t.Fatalf("CloneRepo() error = %v", err)
		}

		if commit == "" {
			t.Error("Expected non-empty commit hash")
		}
	})

	t.Run("invalid repository URL", func(t *testing.T) {
		destDir := t.TempDir()

		url := "https://github.com/nonexistent/repo-that-does-not-exist-12345.git"
		_, err := CloneRepo(url, "", destDir)

		if err == nil {
			t.Error("Expected error for invalid repository URL, got nil")
		}
	})

	t.Run("invalid URL scheme rejected", func(t *testing.T) {
		destDir := t.TempDir()

		url := "git://github.com/test/repo.git"
		_, err := CloneRepo(url, "", destDir)

		if err == nil {
			t.Error("Expected error for git:// URL scheme, got nil")
		}
	})

	t.Run("clone produces commit hash", func(t *testing.T) {
		destDir := t.TempDir()

		url := "https://github.com/aviator5/dnaspec-test-repo.git"
		commit, err := CloneRepo(url, "", destDir)

		if err != nil {
			if os.Getenv("CI") != "true" {
				t.Skipf("Test repository not available: %v", err)
			}
			t.Fatalf("CloneRepo() error = %v", err)
		}

		// Commit hash should be 40 characters (SHA-1)
		if len(commit) != 40 {
			t.Errorf("Expected commit hash length 40, got %d: %s", len(commit), commit)
		}
	})
}

func TestCloneRepo_Unit(t *testing.T) {
	t.Run("validate URL before cloning", func(t *testing.T) {
		destDir := t.TempDir()

		// Invalid URL should be rejected before attempting clone
		invalidURLs := []string{
			"git://github.com/test/repo.git",
			"file:///local/path",
			"ftp://server.com/repo.git",
		}

		for _, url := range invalidURLs {
			_, err := CloneRepo(url, "", destDir)
			if err == nil {
				t.Errorf("Expected error for invalid URL %s, got nil", url)
			}
		}
	})

	t.Run("accepts valid URL schemes", func(t *testing.T) {
		// Note: These tests will fail at clone time, but should pass URL validation
		validURLs := []string{
			"https://github.com/test/repo.git",
			"git@github.com:test/repo.git",
		}

		for _, url := range validURLs {
			err := ValidateGitURL(url)
			if err != nil {
				t.Errorf("Expected valid URL %s to pass validation, got error: %v", url, err)
			}
		}
	})
}
