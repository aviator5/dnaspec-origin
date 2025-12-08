package git

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// CloneRepo clones a git repository to the specified directory
// Uses shallow clone for efficiency and enforces a timeout
// Returns the commit hash of the cloned repository
func CloneRepo(url, ref, destDir string) (commit string, err error) {
	// Validate URL first
	if err := ValidateGitURL(url); err != nil {
		return "", err
	}

	// Create timeout context (5 minutes)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Build git clone command with shallow clone
	args := []string{"clone", "--depth=1", "--single-branch"}

	// Add branch/tag if specified
	if ref != "" {
		args = append(args, "--branch", ref)
	}

	args = append(args, url, destDir)

	// Run git clone
	cmd := exec.CommandContext(ctx, "git", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("git clone timed out after 5 minutes")
		}
		return "", fmt.Errorf("git clone failed: %w\nOutput: %s", err, string(output))
	}

	// Get commit hash
	cmd = exec.CommandContext(ctx, "git", "-C", destDir, "rev-parse", "HEAD")
	output, err = cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get commit hash: %w", err)
	}

	commit = strings.TrimSpace(string(output))
	return commit, nil
}
