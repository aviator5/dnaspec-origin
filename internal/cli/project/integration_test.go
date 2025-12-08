package project

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/aviator5/dnaspec/internal/core/config"
)

func TestAddCommand_Integration(t *testing.T) {
	// Skip if we're in a short test run
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("add local source with --all flag", func(t *testing.T) {
		// Create a temporary project directory
		projectDir := t.TempDir()
		origDir, _ := os.Getwd()
		defer os.Chdir(origDir)

		// Change to project directory
		if err := os.Chdir(projectDir); err != nil {
			t.Fatalf("Failed to change directory: %v", err)
		}

		// Initialize project
		if err := runInit(); err != nil {
			t.Fatalf("runInit() error = %v", err)
		}

		// Get path to test fixture
		testdataPath, _ := filepath.Abs(filepath.Join(origDir, "../../core/source/testdata/valid-repo"))

		// Run add command with --all flag
		flags := addFlags{
			all: true,
		}
		args := []string{testdataPath}

		if err := runAdd(flags, args); err != nil {
			t.Fatalf("runAdd() error = %v", err)
		}

		// Verify config was updated
		cfg, err := config.LoadProjectConfig(projectConfigFileName)
		if err != nil {
			t.Fatalf("Failed to load config: %v", err)
		}

		if len(cfg.Sources) != 1 {
			t.Fatalf("Expected 1 source, got %d", len(cfg.Sources))
		}

		source := cfg.Sources[0]
		if source.Name != "valid-repo" {
			t.Errorf("Source name = %s, want valid-repo", source.Name)
		}

		if source.Type != "local-path" {
			t.Errorf("Source type = %s, want local-path", source.Type)
		}

		if len(source.Guidelines) != 2 {
			t.Errorf("Expected 2 guidelines, got %d", len(source.Guidelines))
		}

		// Verify files were copied
		guidelineFile := filepath.Join(projectDir, "dnaspec", "valid-repo", "guidelines", "test-guideline.md")
		if _, err := os.Stat(guidelineFile); os.IsNotExist(err) {
			t.Error("Guideline file was not copied")
		}

		promptFile := filepath.Join(projectDir, "dnaspec", "valid-repo", "prompts", "test-prompt.md")
		if _, err := os.Stat(promptFile); os.IsNotExist(err) {
			t.Error("Prompt file was not copied")
		}
	})

	t.Run("add with specific guidelines", func(t *testing.T) {
		projectDir := t.TempDir()
		origDir, _ := os.Getwd()
		defer os.Chdir(origDir)

		os.Chdir(projectDir)

		// Initialize project
		runInit()

		// Get path to test fixture
		testdataPath, _ := filepath.Abs(filepath.Join(origDir, "../../core/source/testdata/valid-repo"))

		// Run add command with specific guideline
		flags := addFlags{
			guidelines: []string{"test-guideline"},
		}
		args := []string{testdataPath}

		if err := runAdd(flags, args); err != nil {
			t.Fatalf("runAdd() error = %v", err)
		}

		// Verify only one guideline was added
		cfg, err := config.LoadProjectConfig(projectConfigFileName)
		if err != nil {
			t.Fatalf("Failed to load config: %v", err)
		}

		source := cfg.Sources[0]
		if len(source.Guidelines) != 1 {
			t.Errorf("Expected 1 guideline, got %d", len(source.Guidelines))
		}

		if source.Guidelines[0].Name != "test-guideline" {
			t.Errorf("Guideline name = %s, want test-guideline", source.Guidelines[0].Name)
		}

		// Verify the referenced prompt was included
		if len(source.Prompts) != 1 {
			t.Errorf("Expected 1 prompt, got %d", len(source.Prompts))
		}
	})

	t.Run("add with custom source name", func(t *testing.T) {
		projectDir := t.TempDir()
		origDir, _ := os.Getwd()
		defer os.Chdir(origDir)

		os.Chdir(projectDir)
		runInit()

		testdataPath, _ := filepath.Abs(filepath.Join(origDir, "../../core/source/testdata/valid-repo"))

		flags := addFlags{
			all:  true,
			name: "custom-name",
		}
		args := []string{testdataPath}

		if err := runAdd(flags, args); err != nil {
			t.Fatalf("runAdd() error = %v", err)
		}

		cfg, err := config.LoadProjectConfig(projectConfigFileName)
		if err != nil {
			t.Fatalf("Failed to load config: %v", err)
		}

		if cfg.Sources[0].Name != "custom-name" {
			t.Errorf("Source name = %s, want custom-name", cfg.Sources[0].Name)
		}

		// Verify files are under custom name
		destDir := filepath.Join(projectDir, "dnaspec", "custom-name")
		if _, err := os.Stat(destDir); os.IsNotExist(err) {
			t.Error("Custom source directory was not created")
		}
	})

	t.Run("error when project not initialized", func(t *testing.T) {
		projectDir := t.TempDir()
		origDir, _ := os.Getwd()
		defer os.Chdir(origDir)

		os.Chdir(projectDir)
		// Don't run init

		testdataPath, _ := filepath.Abs(filepath.Join(origDir, "../../core/source/testdata/valid-repo"))

		flags := addFlags{all: true}
		args := []string{testdataPath}

		err := runAdd(flags, args)
		if err == nil {
			t.Error("Expected error when project not initialized, got nil")
		}
	})

	t.Run("error on duplicate source name", func(t *testing.T) {
		projectDir := t.TempDir()
		origDir, _ := os.Getwd()
		defer os.Chdir(origDir)

		os.Chdir(projectDir)
		runInit()

		testdataPath, _ := filepath.Abs(filepath.Join(origDir, "../../core/source/testdata/valid-repo"))

		flags := addFlags{all: true}
		args := []string{testdataPath}

		// Add source first time
		if err := runAdd(flags, args); err != nil {
			t.Fatalf("First runAdd() error = %v", err)
		}

		// Try to add again with same name
		err := runAdd(flags, args)
		if err == nil {
			t.Error("Expected error for duplicate source name, got nil")
		}
	})

	t.Run("dry-run mode", func(t *testing.T) {
		projectDir := t.TempDir()
		origDir, _ := os.Getwd()
		defer os.Chdir(origDir)

		os.Chdir(projectDir)
		runInit()

		testdataPath, _ := filepath.Abs(filepath.Join(origDir, "../../core/source/testdata/valid-repo"))

		flags := addFlags{
			all:    true,
			dryRun: true,
		}
		args := []string{testdataPath}

		if err := runAdd(flags, args); err != nil {
			t.Fatalf("runAdd() with dry-run error = %v", err)
		}

		// Verify config was NOT updated
		cfg, err := config.LoadProjectConfig(projectConfigFileName)
		if err != nil {
			t.Fatalf("Failed to load config: %v", err)
		}

		if len(cfg.Sources) != 0 {
			t.Errorf("Expected 0 sources in dry-run mode, got %d", len(cfg.Sources))
		}

		// Verify files were NOT copied
		destDir := filepath.Join(projectDir, "dnaspec", "valid-repo")
		if _, err := os.Stat(destDir); !os.IsNotExist(err) {
			t.Error("Files should not be copied in dry-run mode")
		}
	})
}
