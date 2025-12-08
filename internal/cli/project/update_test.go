package project

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/aviator5/dnaspec/internal/core/config"
)

func TestUpdateCommand_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("update local source with changes", func(t *testing.T) {
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

		// Get path to test fixtures
		testdataPath, _ := filepath.Abs(filepath.Join(origDir, "../../core/source/testdata/valid-repo"))

		// Add initial source
		addFlags := addFlags{
			all: true,
		}
		if err := runAdd(addFlags, []string{testdataPath}); err != nil {
			t.Fatalf("runAdd() error = %v", err)
		}

		// Verify initial state
		cfg, err := config.LoadProjectConfig(projectConfigFileName)
		if err != nil {
			t.Fatalf("Failed to load config: %v", err)
		}

		if len(cfg.Sources[0].Guidelines) != 2 {
			t.Fatalf("Expected 2 guidelines initially, got %d", len(cfg.Sources[0].Guidelines))
		}

		// Update the source path to point to updated fixture
		updatedPath, _ := filepath.Abs(filepath.Join(origDir, "../../core/source/testdata/valid-repo-updated"))
		cfg.Sources[0].Path = updatedPath
		if err := config.SaveProjectConfig(projectConfigFileName, cfg); err != nil {
			t.Fatalf("Failed to update config path: %v", err)
		}

		// Run update command with --add-new=all
		updateFlags := updateFlags{
			addNew: "all",
		}
		if err := updateSingleSource(cfg, "valid-repo", updateFlags); err != nil {
			t.Fatalf("updateSingleSource() error = %v", err)
		}

		// Verify updated state
		cfg, err = config.LoadProjectConfig(projectConfigFileName)
		if err != nil {
			t.Fatalf("Failed to load config after update: %v", err)
		}

		source := cfg.Sources[0]

		// Should now have 3 guidelines (2 original + 1 new)
		if len(source.Guidelines) != 3 {
			t.Errorf("Expected 3 guidelines after update, got %d", len(source.Guidelines))
		}

		// Check that test-guideline description was updated
		var testGuideline *config.ProjectGuideline
		for i := range source.Guidelines {
			if source.Guidelines[i].Name == "test-guideline" {
				testGuideline = &source.Guidelines[i]
				break
			}
		}

		if testGuideline == nil {
			t.Fatal("test-guideline not found after update")
		}

		if testGuideline.Description != "Updated test guideline description" {
			t.Errorf("Description not updated, got: %s", testGuideline.Description)
		}

		if len(testGuideline.ApplicableScenarios) != 2 {
			t.Errorf("Expected 2 scenarios after update, got %d", len(testGuideline.ApplicableScenarios))
		}

		// Verify new guideline was added
		var newGuideline *config.ProjectGuideline
		for i := range source.Guidelines {
			if source.Guidelines[i].Name == "new-guideline" {
				newGuideline = &source.Guidelines[i]
				break
			}
		}

		if newGuideline == nil {
			t.Error("new-guideline not found after update")
		}

		// Verify new guideline file was copied
		newGuidelineFile := filepath.Join(projectDir, "dnaspec", "valid-repo", "guidelines", "new-guideline.md")
		if _, err := os.Stat(newGuidelineFile); os.IsNotExist(err) {
			t.Error("New guideline file was not copied")
		}
	})

	t.Run("update with --add-new=none", func(t *testing.T) {
		projectDir := t.TempDir()
		origDir, _ := os.Getwd()
		defer os.Chdir(origDir)

		os.Chdir(projectDir)
		runInit()

		testdataPath, _ := filepath.Abs(filepath.Join(origDir, "../../core/source/testdata/valid-repo"))

		// Add initial source
		addFlags := addFlags{
			all: true,
		}
		runAdd(addFlags, []string{testdataPath})

		// Update config to point to updated fixture
		cfg, _ := config.LoadProjectConfig(projectConfigFileName)
		updatedPath, _ := filepath.Abs(filepath.Join(origDir, "../../core/source/testdata/valid-repo-updated"))
		cfg.Sources[0].Path = updatedPath
		config.SaveProjectConfig(projectConfigFileName, cfg)

		// Run update with --add-new=none
		updateFlags := updateFlags{
			addNew: "none",
		}
		if err := updateSingleSource(cfg, "valid-repo", updateFlags); err != nil {
			t.Fatalf("updateSingleSource() error = %v", err)
		}

		// Verify state
		cfg, _ = config.LoadProjectConfig(projectConfigFileName)
		source := cfg.Sources[0]

		// Should still have 2 guidelines (new one not added)
		if len(source.Guidelines) != 2 {
			t.Errorf("Expected 2 guidelines (new not added), got %d", len(source.Guidelines))
		}

		// Verify new guideline was NOT added
		for _, g := range source.Guidelines {
			if g.Name == "new-guideline" {
				t.Error("new-guideline should not have been added with --add-new=none")
			}
		}
	})

	t.Run("update with dry-run", func(t *testing.T) {
		projectDir := t.TempDir()
		origDir, _ := os.Getwd()
		defer os.Chdir(origDir)

		os.Chdir(projectDir)
		runInit()

		testdataPath, _ := filepath.Abs(filepath.Join(origDir, "../../core/source/testdata/valid-repo"))

		// Add initial source
		addFlags := addFlags{
			all: true,
		}
		runAdd(addFlags, []string{testdataPath})

		// Get initial config state
		cfgBefore, _ := config.LoadProjectConfig(projectConfigFileName)
		initialGuidelineCount := len(cfgBefore.Sources[0].Guidelines)

		// Update config to point to updated fixture
		updatedPath, _ := filepath.Abs(filepath.Join(origDir, "../../core/source/testdata/valid-repo-updated"))
		cfgBefore.Sources[0].Path = updatedPath
		config.SaveProjectConfig(projectConfigFileName, cfgBefore)

		// Run update with dry-run
		updateFlags := updateFlags{
			dryRun: true,
			addNew: "all",
		}
		if err := updateSingleSource(cfgBefore, "valid-repo", updateFlags); err != nil {
			t.Fatalf("updateSingleSource() error = %v", err)
		}

		// Verify config was NOT changed
		cfgAfter, _ := config.LoadProjectConfig(projectConfigFileName)

		if len(cfgAfter.Sources[0].Guidelines) != initialGuidelineCount {
			t.Errorf("Dry-run should not modify config. Expected %d guidelines, got %d",
				initialGuidelineCount, len(cfgAfter.Sources[0].Guidelines))
		}

		// Verify new guideline file was NOT copied
		newGuidelineFile := filepath.Join(projectDir, "dnaspec", "valid-repo", "guidelines", "new-guideline.md")
		if _, err := os.Stat(newGuidelineFile); err == nil {
			t.Error("Dry-run should not copy new files")
		}
	})

	t.Run("update all sources", func(t *testing.T) {
		projectDir := t.TempDir()
		origDir, _ := os.Getwd()
		defer os.Chdir(origDir)

		os.Chdir(projectDir)
		runInit()

		testdataPath, _ := filepath.Abs(filepath.Join(origDir, "../../core/source/testdata/valid-repo"))

		// Add two sources
		addFlags1 := addFlags{
			all:  true,
			name: "source-1",
		}
		runAdd(addFlags1, []string{testdataPath})

		addFlags2 := addFlags{
			all:  true,
			name: "source-2",
		}
		runAdd(addFlags2, []string{testdataPath})

		// Verify we have 2 sources
		cfg, _ := config.LoadProjectConfig(projectConfigFileName)
		if len(cfg.Sources) != 2 {
			t.Fatalf("Expected 2 sources, got %d", len(cfg.Sources))
		}

		// Update both to point to updated fixture
		updatedPath, _ := filepath.Abs(filepath.Join(origDir, "../../core/source/testdata/valid-repo-updated"))
		cfg.Sources[0].Path = updatedPath
		cfg.Sources[1].Path = updatedPath
		config.SaveProjectConfig(projectConfigFileName, cfg)

		// Run update --all
		updateFlags := updateFlags{
			all:    true,
			addNew: "all",
		}

		// Reload config
		cfg, _ = config.LoadProjectConfig(projectConfigFileName)
		if err := updateAllSources(cfg, updateFlags); err != nil {
			t.Fatalf("updateAllSources() error = %v", err)
		}

		// Verify both sources were updated
		cfg, _ = config.LoadProjectConfig(projectConfigFileName)

		for _, source := range cfg.Sources {
			if len(source.Guidelines) != 3 {
				t.Errorf("Source %s: expected 3 guidelines after update, got %d",
					source.Name, len(source.Guidelines))
			}
		}
	})

	t.Run("update with no changes", func(t *testing.T) {
		projectDir := t.TempDir()
		origDir, _ := os.Getwd()
		defer os.Chdir(origDir)

		os.Chdir(projectDir)
		runInit()

		testdataPath, _ := filepath.Abs(filepath.Join(origDir, "../../core/source/testdata/valid-repo"))

		// Add initial source
		addFlags := addFlags{
			all: true,
		}
		runAdd(addFlags, []string{testdataPath})

		// Run update without changing the source
		cfg, _ := config.LoadProjectConfig(projectConfigFileName)
		updateFlags := updateFlags{}

		// This should succeed and report no changes
		if err := updateSingleSource(cfg, "valid-repo", updateFlags); err != nil {
			t.Fatalf("updateSingleSource() should succeed with no changes: %v", err)
		}

		// Verify config unchanged
		cfgAfter, _ := config.LoadProjectConfig(projectConfigFileName)
		if len(cfgAfter.Sources[0].Guidelines) != 2 {
			t.Errorf("Expected 2 guidelines (unchanged), got %d", len(cfgAfter.Sources[0].Guidelines))
		}
	})
}

func TestUpdateCommand_ErrorCases(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("source not found", func(t *testing.T) {
		projectDir := t.TempDir()
		origDir, _ := os.Getwd()
		defer os.Chdir(origDir)

		os.Chdir(projectDir)
		runInit()

		cfg, _ := config.LoadProjectConfig(projectConfigFileName)
		updateFlags := updateFlags{}

		// Try to update non-existent source
		err := updateSingleSource(cfg, "nonexistent", updateFlags)
		if err == nil {
			t.Error("Expected error for nonexistent source")
		}
	})

	t.Run("project not initialized", func(t *testing.T) {
		projectDir := t.TempDir()
		origDir, _ := os.Getwd()
		defer os.Chdir(origDir)

		os.Chdir(projectDir)

		// Don't initialize project
		updateFlags := updateFlags{}

		err := runUpdate(updateFlags, []string{"some-source"})
		if err == nil {
			t.Error("Expected error when project not initialized")
		}
	})

	t.Run("invalid add-new value", func(t *testing.T) {
		updateFlags := updateFlags{
			addNew: "invalid",
		}

		err := runUpdate(updateFlags, []string{"some-source"})
		if err == nil {
			t.Error("Expected error for invalid --add-new value")
		}
	})

	t.Run("both source name and --all", func(t *testing.T) {
		updateFlags := updateFlags{
			all: true,
		}

		err := runUpdate(updateFlags, []string{"some-source"})
		if err == nil {
			t.Error("Expected error when both source name and --all are specified")
		}
	})

	t.Run("neither source name nor --all", func(t *testing.T) {
		updateFlags := updateFlags{}

		err := runUpdate(updateFlags, []string{})
		if err == nil {
			t.Error("Expected error when neither source name nor --all are specified")
		}
	})
}
