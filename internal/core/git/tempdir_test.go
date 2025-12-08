package git

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCreateTempCloneDir(t *testing.T) {
	path, cleanup, err := CreateTempCloneDir()
	if err != nil {
		t.Fatalf("CreateTempCloneDir() error = %v", err)
	}
	defer cleanup()

	// Verify path is under temp directory
	if !strings.HasPrefix(path, os.TempDir()) {
		t.Errorf("Path not under temp directory: %s", path)
	}

	// Verify path contains dnaspec
	if !strings.Contains(path, "dnaspec") {
		t.Errorf("Path does not contain 'dnaspec': %s", path)
	}

	// Verify directory was created
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Error("Directory was not created")
	}

	// Verify cleanup works
	cleanup()
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Error("Cleanup did not remove directory")
	}
}

func TestCreateTempCloneDir_Uniqueness(t *testing.T) {
	// Create two temp directories
	path1, cleanup1, err := CreateTempCloneDir()
	if err != nil {
		t.Fatalf("CreateTempCloneDir() error = %v", err)
	}
	defer cleanup1()

	path2, cleanup2, err := CreateTempCloneDir()
	if err != nil {
		t.Fatalf("CreateTempCloneDir() error = %v", err)
	}
	defer cleanup2()

	// Verify they are different
	if path1 == path2 {
		t.Errorf("Paths are not unique: %s == %s", path1, path2)
	}

	// Verify both exist
	if _, err := os.Stat(path1); os.IsNotExist(err) {
		t.Error("First directory was not created")
	}
	if _, err := os.Stat(path2); os.IsNotExist(err) {
		t.Error("Second directory was not created")
	}
}

func TestCreateTempCloneDir_Structure(t *testing.T) {
	path, cleanup, err := CreateTempCloneDir()
	if err != nil {
		t.Fatalf("CreateTempCloneDir() error = %v", err)
	}
	defer cleanup()

	// Verify path has the expected structure: /tmp/dnaspec/PID-RANDOMID
	parts := strings.Split(path, string(filepath.Separator))

	// Find "dnaspec" in the path
	foundDnaspec := false
	for _, part := range parts {
		if part == "dnaspec" {
			foundDnaspec = true
			break
		}
	}

	if !foundDnaspec {
		t.Errorf("Path does not contain 'dnaspec' directory: %s", path)
	}

	// Verify last component contains PID and random ID (contains hyphen)
	lastPart := filepath.Base(path)
	if !strings.Contains(lastPart, "-") {
		t.Errorf("Last path component does not contain hyphen (PID-RANDOMID): %s", lastPart)
	}
}
