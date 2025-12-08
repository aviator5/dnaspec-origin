package files

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/aviator5/dnaspec/internal/core/config"
)

func TestCopyGuidelineFiles(t *testing.T) {
	t.Run("copy guidelines and prompts", func(t *testing.T) {
		// Create source directory structure
		sourceDir := t.TempDir()
		destDir := t.TempDir()

		// Create source files
		guidelineContent := []byte("# Guideline content")
		promptContent := []byte("# Prompt content")

		os.MkdirAll(filepath.Join(sourceDir, "guidelines"), 0755)
		os.MkdirAll(filepath.Join(sourceDir, "prompts"), 0755)

		os.WriteFile(filepath.Join(sourceDir, "guidelines", "test.md"), guidelineContent, 0644)
		os.WriteFile(filepath.Join(sourceDir, "prompts", "review.md"), promptContent, 0644)

		// Define what to copy
		guidelines := []config.ManifestGuideline{
			{
				Name: "test-guideline",
				File: "guidelines/test.md",
			},
		}

		prompts := []config.ManifestPrompt{
			{
				Name: "review",
				File: "prompts/review.md",
			},
		}

		// Copy files
		err := CopyGuidelineFiles(sourceDir, destDir, guidelines, prompts)
		if err != nil {
			t.Fatalf("CopyGuidelineFiles() error = %v", err)
		}

		// Verify files were copied
		guidelineDest := filepath.Join(destDir, "guidelines", "test.md")
		if _, err := os.Stat(guidelineDest); os.IsNotExist(err) {
			t.Error("Guideline file was not copied")
		}

		promptDest := filepath.Join(destDir, "prompts", "review.md")
		if _, err := os.Stat(promptDest); os.IsNotExist(err) {
			t.Error("Prompt file was not copied")
		}

		// Verify content
		copiedGuideline, _ := os.ReadFile(guidelineDest)
		if string(copiedGuideline) != string(guidelineContent) {
			t.Error("Guideline content does not match")
		}

		copiedPrompt, _ := os.ReadFile(promptDest)
		if string(copiedPrompt) != string(promptContent) {
			t.Error("Prompt content does not match")
		}
	})

	t.Run("preserve directory structure", func(t *testing.T) {
		sourceDir := t.TempDir()
		destDir := t.TempDir()

		// Create nested directory structure
		nestedPath := filepath.Join(sourceDir, "guidelines", "subdir")
		os.MkdirAll(nestedPath, 0755)
		os.WriteFile(filepath.Join(nestedPath, "nested.md"), []byte("nested"), 0644)

		guidelines := []config.ManifestGuideline{
			{
				Name: "nested",
				File: "guidelines/subdir/nested.md",
			},
		}

		err := CopyGuidelineFiles(sourceDir, destDir, guidelines, []config.ManifestPrompt{})
		if err != nil {
			t.Fatalf("CopyGuidelineFiles() error = %v", err)
		}

		// Verify nested structure was preserved
		destPath := filepath.Join(destDir, "guidelines", "subdir", "nested.md")
		if _, err := os.Stat(destPath); os.IsNotExist(err) {
			t.Error("Nested directory structure was not preserved")
		}
	})

	t.Run("error on missing source file", func(t *testing.T) {
		sourceDir := t.TempDir()
		destDir := t.TempDir()

		guidelines := []config.ManifestGuideline{
			{
				Name: "missing",
				File: "guidelines/missing.md",
			},
		}

		err := CopyGuidelineFiles(sourceDir, destDir, guidelines, []config.ManifestPrompt{})
		if err == nil {
			t.Error("Expected error for missing source file, got nil")
		}
	})

	t.Run("copy multiple files", func(t *testing.T) {
		sourceDir := t.TempDir()
		destDir := t.TempDir()

		os.MkdirAll(filepath.Join(sourceDir, "guidelines"), 0755)
		os.WriteFile(filepath.Join(sourceDir, "guidelines", "g1.md"), []byte("g1"), 0644)
		os.WriteFile(filepath.Join(sourceDir, "guidelines", "g2.md"), []byte("g2"), 0644)
		os.WriteFile(filepath.Join(sourceDir, "guidelines", "g3.md"), []byte("g3"), 0644)

		guidelines := []config.ManifestGuideline{
			{Name: "g1", File: "guidelines/g1.md"},
			{Name: "g2", File: "guidelines/g2.md"},
			{Name: "g3", File: "guidelines/g3.md"},
		}

		err := CopyGuidelineFiles(sourceDir, destDir, guidelines, []config.ManifestPrompt{})
		if err != nil {
			t.Fatalf("CopyGuidelineFiles() error = %v", err)
		}

		// Verify all files were copied
		for _, g := range guidelines {
			path := filepath.Join(destDir, g.File)
			if _, err := os.Stat(path); os.IsNotExist(err) {
				t.Errorf("File %s was not copied", g.File)
			}
		}
	})

	t.Run("copy empty lists", func(t *testing.T) {
		sourceDir := t.TempDir()
		destDir := t.TempDir()

		err := CopyGuidelineFiles(sourceDir, destDir, []config.ManifestGuideline{}, []config.ManifestPrompt{})
		if err != nil {
			t.Errorf("CopyGuidelineFiles() with empty lists error = %v, want nil", err)
		}
	})

	t.Run("rollback on partial failure", func(t *testing.T) {
		sourceDir := t.TempDir()
		destDir := t.TempDir()

		// Create source files
		os.MkdirAll(filepath.Join(sourceDir, "guidelines"), 0755)
		os.WriteFile(filepath.Join(sourceDir, "guidelines", "g1.md"), []byte("g1"), 0644)
		os.WriteFile(filepath.Join(sourceDir, "guidelines", "g2.md"), []byte("g2"), 0644)
		// Intentionally skip creating g3.md to trigger error

		guidelines := []config.ManifestGuideline{
			{Name: "g1", File: "guidelines/g1.md"},
			{Name: "g2", File: "guidelines/g2.md"},
			{Name: "g3", File: "guidelines/g3.md"}, // This will fail
		}

		err := CopyGuidelineFiles(sourceDir, destDir, guidelines, []config.ManifestPrompt{})
		if err == nil {
			t.Fatal("Expected error for missing file, got nil")
		}

		// Verify rollback: previously copied files should be removed
		g1Path := filepath.Join(destDir, "guidelines", "g1.md")
		if _, err := os.Stat(g1Path); !os.IsNotExist(err) {
			t.Error("First file was not rolled back after failure")
		}

		g2Path := filepath.Join(destDir, "guidelines", "g2.md")
		if _, err := os.Stat(g2Path); !os.IsNotExist(err) {
			t.Error("Second file was not rolled back after failure")
		}
	})

	t.Run("rollback on prompt copy failure", func(t *testing.T) {
		sourceDir := t.TempDir()
		destDir := t.TempDir()

		// Create only guideline, not prompt
		os.MkdirAll(filepath.Join(sourceDir, "guidelines"), 0755)
		os.WriteFile(filepath.Join(sourceDir, "guidelines", "g1.md"), []byte("g1"), 0644)

		guidelines := []config.ManifestGuideline{
			{Name: "g1", File: "guidelines/g1.md"},
		}

		prompts := []config.ManifestPrompt{
			{Name: "p1", File: "prompts/missing.md"}, // This will fail
		}

		err := CopyGuidelineFiles(sourceDir, destDir, guidelines, prompts)
		if err == nil {
			t.Fatal("Expected error for missing prompt file, got nil")
		}

		// Verify rollback: guideline file should be removed
		g1Path := filepath.Join(destDir, "guidelines", "g1.md")
		if _, err := os.Stat(g1Path); !os.IsNotExist(err) {
			t.Error("Guideline file was not rolled back after prompt copy failure")
		}
	})
}

func TestCopyFile(t *testing.T) {
	t.Run("copy file successfully", func(t *testing.T) {
		tmpDir := t.TempDir()

		srcPath := filepath.Join(tmpDir, "source.txt")
		dstPath := filepath.Join(tmpDir, "dest.txt")

		content := []byte("test content")
		os.WriteFile(srcPath, content, 0644)

		err := copyFile(srcPath, dstPath)
		if err != nil {
			t.Fatalf("copyFile() error = %v", err)
		}

		// Verify file was copied
		copiedContent, err := os.ReadFile(dstPath)
		if err != nil {
			t.Fatalf("Failed to read copied file: %v", err)
		}

		if string(copiedContent) != string(content) {
			t.Error("Copied content does not match original")
		}
	})

	t.Run("create parent directories", func(t *testing.T) {
		tmpDir := t.TempDir()

		srcPath := filepath.Join(tmpDir, "source.txt")
		dstPath := filepath.Join(tmpDir, "nested", "dirs", "dest.txt")

		os.WriteFile(srcPath, []byte("content"), 0644)

		err := copyFile(srcPath, dstPath)
		if err != nil {
			t.Fatalf("copyFile() error = %v", err)
		}

		// Verify parent directories were created
		if _, err := os.Stat(filepath.Dir(dstPath)); os.IsNotExist(err) {
			t.Error("Parent directories were not created")
		}

		// Verify file exists
		if _, err := os.Stat(dstPath); os.IsNotExist(err) {
			t.Error("Destination file was not created")
		}
	})

	t.Run("error on missing source", func(t *testing.T) {
		tmpDir := t.TempDir()

		srcPath := filepath.Join(tmpDir, "nonexistent.txt")
		dstPath := filepath.Join(tmpDir, "dest.txt")

		err := copyFile(srcPath, dstPath)
		if err == nil {
			t.Error("Expected error for missing source file, got nil")
		}
	})

	t.Run("overwrite existing file", func(t *testing.T) {
		tmpDir := t.TempDir()

		srcPath := filepath.Join(tmpDir, "source.txt")
		dstPath := filepath.Join(tmpDir, "dest.txt")

		// Create both files
		os.WriteFile(srcPath, []byte("new content"), 0644)
		os.WriteFile(dstPath, []byte("old content"), 0644)

		err := copyFile(srcPath, dstPath)
		if err != nil {
			t.Fatalf("copyFile() error = %v", err)
		}

		// Verify content was overwritten
		content, _ := os.ReadFile(dstPath)
		if string(content) != "new content" {
			t.Error("File was not overwritten with new content")
		}
	})
}
