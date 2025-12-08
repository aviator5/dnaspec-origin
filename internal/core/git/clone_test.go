package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestCloneRepo_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Create a local git repository
	repoDir := t.TempDir()

	// Initialize repo
	runGit(t, repoDir, "init")
	runGit(t, repoDir, "config", "user.email", "test@example.com")
	runGit(t, repoDir, "config", "user.name", "Test User")
	runGit(t, repoDir, "branch", "-m", "main")

	// Create a file
	err := os.WriteFile(filepath.Join(repoDir, "test.txt"), []byte("content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}

	// Commit
	runGit(t, repoDir, "add", ".")
	runGit(t, repoDir, "commit", "-m", "Initial commit")

	// Get head hash
	headHash := getGitHead(t, repoDir)

	// Create a tag
	runGit(t, repoDir, "tag", "v1.0.0")

	t.Run("clone public repository", func(t *testing.T) {
		destDir := t.TempDir()

		// Use local repo
		url := "file://" + repoDir
		commit, err := CloneRepo(url, "", destDir)

		if err != nil {
			t.Fatalf("CloneRepo() error = %v", err)
		}

		if commit != headHash {
			t.Errorf("Expected commit hash %s, got %s", headHash, commit)
		}

		// Verify .git directory exists
		gitDir := filepath.Join(destDir, ".git")
		if _, err := os.Stat(gitDir); os.IsNotExist(err) {
			t.Error("Expected .git directory to exist")
		}
	})

	t.Run("clone with specific ref", func(t *testing.T) {
		destDir := t.TempDir()

		// Use local repo
		url := "file://" + repoDir
		commit, err := CloneRepo(url, "main", destDir)

		if err != nil {
			t.Fatalf("CloneRepo() error = %v", err)
		}

		if commit != headHash {
			t.Errorf("Expected commit hash %s, got %s", headHash, commit)
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

		url := "file://" + repoDir
		commit, err := CloneRepo(url, "", destDir)

		if err != nil {
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
			// file:// is now allowed, so removed from here
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
			"file:///local/path",
		}

		for _, url := range validURLs {
			err := ValidateGitURL(url)
			if err != nil {
				t.Errorf("Expected valid URL %s to pass validation, got error: %v", url, err)
			}
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
