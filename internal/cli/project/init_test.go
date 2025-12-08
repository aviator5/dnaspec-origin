package project

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRunInit(t *testing.T) {
	// Save current directory
	origDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current directory: %v", err)
	}
	defer func() {
		err := os.Chdir(origDir)
		if err != nil {
			// t.Error cannot be used in main test default cleanup easily if main test is ending,
			// but we can log it if strictly needed. Since this is non-critical cleanup of chdir:
			_ = err
		}
	}()

	t.Run("success", func(t *testing.T) {
		tmpDir := t.TempDir()
		if err := os.Chdir(tmpDir); err != nil {
			t.Fatalf("failed to change directory: %v", err)
		}

		if err := runInit(); err != nil {
			t.Errorf("runInit() error = %v", err)
		}

		// Verify file was created
		configPath := filepath.Join(tmpDir, projectConfigFileName)
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			t.Error("runInit() did not create config file")
		}
	})

	t.Run("error_when_file_exists", func(t *testing.T) {
		tmpDir := t.TempDir()
		if err := os.Chdir(tmpDir); err != nil {
			t.Fatalf("failed to change directory: %v", err)
		}

		// Create existing file
		configPath := filepath.Join(tmpDir, projectConfigFileName)
		if err := os.WriteFile(configPath, []byte("existing"), 0644); err != nil {
			t.Fatalf("failed to create existing file: %v", err)
		}

		// Should fail
		if err := runInit(); err == nil {
			t.Error("runInit() expected error when file exists, got nil")
		}
	})
}

func TestNewInitCmd(t *testing.T) {
	cmd := NewInitCmd()
	if cmd == nil {
		t.Fatal("NewInitCmd() returned nil")
	}

	if cmd.Use != "init" {
		t.Errorf("Use = %v, want %v", cmd.Use, "init")
	}

	if cmd.Short == "" {
		t.Error("Short description is empty")
	}

	if cmd.RunE == nil {
		t.Error("RunE is nil")
	}
}
