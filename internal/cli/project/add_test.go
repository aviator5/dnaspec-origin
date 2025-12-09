package project

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/aviator5/dnaspec/internal/core/config"
	"github.com/aviator5/dnaspec/internal/core/source"
)

func TestConvertToRelativePath_InsideProject(t *testing.T) {
	// Create temp project directory and change to it
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}

	// Create a local source inside project
	sourceDir := filepath.Join(tmpDir, "local-dna")
	if err := os.MkdirAll(sourceDir, 0755); err != nil {
		t.Fatal(err)
	}

	sourceInfo := &source.SourceInfo{
		SourceType: config.SourceTypeLocalPath,
		Path:       sourceDir,
	}

	result, err := convertToRelativePath(sourceInfo)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	// Should return relative path
	if filepath.IsAbs(result) {
		t.Errorf("expected relative path, got absolute: %s", result)
	}

	expected := "local-dna"
	if result != expected {
		t.Errorf("expected path %q, got %q", expected, result)
	}
}

func TestConvertToRelativePath_OutsideProject_NonInteractive(t *testing.T) {
	// Create temp directories
	tmpDir := t.TempDir()
	projectDir := filepath.Join(tmpDir, "project")
	outsideSource := filepath.Join(tmpDir, "outside-dna")

	if err := os.MkdirAll(projectDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(outsideSource, 0755); err != nil {
		t.Fatal(err)
	}

	// Change to project directory
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	if err := os.Chdir(projectDir); err != nil {
		t.Fatal(err)
	}

	sourceInfo := &source.SourceInfo{
		SourceType: config.SourceTypeLocalPath,
		Path:       outsideSource,
	}

	// convertToRelativePath should handle outside paths gracefully
	result, err := convertToRelativePath(sourceInfo)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	// Should return absolute path when outside project
	if !filepath.IsAbs(result) {
		t.Errorf("expected absolute path for outside source, got relative: %s", result)
	}
}

func TestConvertToRelativePath_GitSource(t *testing.T) {
	sourceInfo := &source.SourceInfo{
		SourceType: config.SourceTypeGitRepo,
		URL:        "https://github.com/example/dna",
		Path:       "",
	}

	result, err := convertToRelativePath(sourceInfo)
	if err != nil {
		t.Fatalf("expected no error for git source, got: %v", err)
	}

	// Git sources should keep empty path
	if result != "" {
		t.Errorf("expected empty path for git source, got: %s", result)
	}
}

func TestCheckLocalPathBeforeLoad_InsideProject(t *testing.T) {
	// Create temp project directory and change to it
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}

	// Create a local source inside project
	sourceDir := filepath.Join(tmpDir, "local-dna")
	if err := os.MkdirAll(sourceDir, 0755); err != nil {
		t.Fatal(err)
	}

	flags := addFlags{}
	args := []string{sourceDir}

	// Should not error for inside paths
	err := checkLocalPathBeforeLoad(flags, args)
	if err != nil {
		t.Errorf("expected no error for inside path, got: %v", err)
	}
}

func TestCheckLocalPathBeforeLoad_OutsideProject_NonInteractive(t *testing.T) {
	// Create temp directories
	tmpDir := t.TempDir()
	projectDir := filepath.Join(tmpDir, "project")
	outsideSource := filepath.Join(tmpDir, "outside-dna")

	if err := os.MkdirAll(projectDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(outsideSource, 0755); err != nil {
		t.Fatal(err)
	}

	// Change to project directory
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	if err := os.Chdir(projectDir); err != nil {
		t.Fatal(err)
	}

	flags := addFlags{all: true} // Non-interactive mode
	args := []string{outsideSource}

	// Non-interactive mode should auto-accept
	err := checkLocalPathBeforeLoad(flags, args)
	if err != nil {
		t.Errorf("expected no error in non-interactive mode, got: %v", err)
	}
}

func TestCheckLocalPathBeforeLoad_GitRepo(t *testing.T) {
	flags := addFlags{gitRepo: "https://github.com/example/dna"}
	args := []string{}

	// Should skip check for git repos
	err := checkLocalPathBeforeLoad(flags, args)
	if err != nil {
		t.Errorf("expected no error for git repo, got: %v", err)
	}
}
