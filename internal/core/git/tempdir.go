package git

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
)

// CreateTempCloneDir creates a temporary directory for git cloning
// Returns the directory path and a cleanup function
func CreateTempCloneDir() (path string, cleanup func(), err error) {
	// Generate unique ID using random bytes
	randomBytes := make([]byte, 8)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", nil, fmt.Errorf("failed to generate random ID: %w", err)
	}
	randomID := hex.EncodeToString(randomBytes)

	// Create unique temp directory
	tempDir := filepath.Join(
		os.TempDir(),
		"dnaspec",
		fmt.Sprintf("%d-%s", os.Getpid(), randomID),
	)

	if err := os.MkdirAll(tempDir, 0o755); err != nil {
		return "", nil, fmt.Errorf("failed to create temp directory: %w", err)
	}

	// Cleanup function
	cleanupFn := func() {
		_ = os.RemoveAll(tempDir)
	}

	return tempDir, cleanupFn, nil
}
