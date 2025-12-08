package config

import (
	"slices"
	"testing"
)

func TestCompareGuidelines_NoChanges(t *testing.T) {
	current := []ProjectGuideline{
		{
			Name:                "guideline-1",
			File:                "guidelines/g1.md",
			Description:         "Guideline 1",
			ApplicableScenarios: []string{"scenario-1"},
			Prompts:             []string{"prompt-1"},
		},
	}

	manifest := []ManifestGuideline{
		{
			Name:                "guideline-1",
			File:                "guidelines/g1.md",
			Description:         "Guideline 1",
			ApplicableScenarios: []string{"scenario-1"},
			Prompts:             []string{"prompt-1"},
		},
	}

	result := CompareGuidelines(current, manifest)

	if len(result.Updated) != 0 {
		t.Errorf("Expected 0 updated, got %d", len(result.Updated))
	}
	if len(result.New) != 0 {
		t.Errorf("Expected 0 new, got %d", len(result.New))
	}
	if len(result.Removed) != 0 {
		t.Errorf("Expected 0 removed, got %d", len(result.Removed))
	}
	if len(result.Unchanged) != 1 {
		t.Errorf("Expected 1 unchanged, got %d", len(result.Unchanged))
	}
}

func TestCompareGuidelines_Updated(t *testing.T) {
	current := []ProjectGuideline{
		{
			Name:                "guideline-1",
			File:                "guidelines/g1.md",
			Description:         "Old Description",
			ApplicableScenarios: []string{"scenario-1"},
			Prompts:             []string{"prompt-1"},
		},
	}

	manifest := []ManifestGuideline{
		{
			Name:                "guideline-1",
			File:                "guidelines/g1.md",
			Description:         "New Description",
			ApplicableScenarios: []string{"scenario-1"},
			Prompts:             []string{"prompt-1"},
		},
	}

	result := CompareGuidelines(current, manifest)

	if len(result.Updated) != 1 {
		t.Errorf("Expected 1 updated, got %d", len(result.Updated))
	}
	if len(result.Unchanged) != 0 {
		t.Errorf("Expected 0 unchanged, got %d", len(result.Unchanged))
	}
}

func TestCompareGuidelines_New(t *testing.T) {
	current := []ProjectGuideline{}

	manifest := []ManifestGuideline{
		{
			Name:                "guideline-1",
			File:                "guidelines/g1.md",
			Description:         "Guideline 1",
			ApplicableScenarios: []string{"scenario-1"},
			Prompts:             []string{"prompt-1"},
		},
	}

	result := CompareGuidelines(current, manifest)

	if len(result.New) != 1 {
		t.Errorf("Expected 1 new, got %d", len(result.New))
	}
	if !slices.Contains(result.New, "guideline-1") {
		t.Errorf("Expected guideline-1 in new list")
	}
}

func TestCompareGuidelines_Removed(t *testing.T) {
	current := []ProjectGuideline{
		{
			Name:                "guideline-1",
			File:                "guidelines/g1.md",
			Description:         "Guideline 1",
			ApplicableScenarios: []string{"scenario-1"},
			Prompts:             []string{"prompt-1"},
		},
	}

	manifest := []ManifestGuideline{}

	result := CompareGuidelines(current, manifest)

	if len(result.Removed) != 1 {
		t.Errorf("Expected 1 removed, got %d", len(result.Removed))
	}
	if !slices.Contains(result.Removed, "guideline-1") {
		t.Errorf("Expected guideline-1 in removed list")
	}
}

func TestCompareGuidelines_Mixed(t *testing.T) {
	current := []ProjectGuideline{
		{
			Name:        "unchanged-guideline",
			Description: "Unchanged",
		},
		{
			Name:        "updated-guideline",
			Description: "Old Description",
		},
		{
			Name:        "removed-guideline",
			Description: "Removed",
		},
	}

	manifest := []ManifestGuideline{
		{
			Name:        "unchanged-guideline",
			Description: "Unchanged",
		},
		{
			Name:        "updated-guideline",
			Description: "New Description",
		},
		{
			Name:        "new-guideline",
			Description: "New",
		},
	}

	result := CompareGuidelines(current, manifest)

	if len(result.Updated) != 1 || !slices.Contains(result.Updated, "updated-guideline") {
		t.Errorf("Expected updated-guideline in updated list")
	}
	if len(result.New) != 1 || !slices.Contains(result.New, "new-guideline") {
		t.Errorf("Expected new-guideline in new list")
	}
	if len(result.Removed) != 1 || !slices.Contains(result.Removed, "removed-guideline") {
		t.Errorf("Expected removed-guideline in removed list")
	}
	if len(result.Unchanged) != 1 || !slices.Contains(result.Unchanged, "unchanged-guideline") {
		t.Errorf("Expected unchanged-guideline in unchanged list")
	}
}

func TestHasChanges_DescriptionChanged(t *testing.T) {
	current := ProjectGuideline{
		Name:        "test",
		Description: "Old",
	}
	manifest := ManifestGuideline{
		Name:        "test",
		Description: "New",
	}

	if !hasChanges(current, manifest) {
		t.Error("Expected hasChanges to return true for description change")
	}
}

func TestHasChanges_ScenariosChanged(t *testing.T) {
	current := ProjectGuideline{
		Name:                "test",
		Description:         "Same",
		ApplicableScenarios: []string{"scenario-1"},
	}
	manifest := ManifestGuideline{
		Name:                "test",
		Description:         "Same",
		ApplicableScenarios: []string{"scenario-1", "scenario-2"},
	}

	if !hasChanges(current, manifest) {
		t.Error("Expected hasChanges to return true for scenarios change")
	}
}

func TestHasChanges_PromptsChanged(t *testing.T) {
	current := ProjectGuideline{
		Name:        "test",
		Description: "Same",
		Prompts:     []string{"prompt-1"},
	}
	manifest := ManifestGuideline{
		Name:        "test",
		Description: "Same",
		Prompts:     []string{"prompt-1", "prompt-2"},
	}

	if !hasChanges(current, manifest) {
		t.Error("Expected hasChanges to return true for prompts change")
	}
}

func TestHasChanges_NoChanges(t *testing.T) {
	current := ProjectGuideline{
		Name:                "test",
		Description:         "Same",
		ApplicableScenarios: []string{"scenario-1"},
		Prompts:             []string{"prompt-1"},
	}
	manifest := ManifestGuideline{
		Name:                "test",
		Description:         "Same",
		ApplicableScenarios: []string{"scenario-1"},
		Prompts:             []string{"prompt-1"},
	}

	if hasChanges(current, manifest) {
		t.Error("Expected hasChanges to return false for no changes")
	}
}

func TestFindSourceByName_Found(t *testing.T) {
	cfg := &ProjectConfig{
		Sources: []ProjectSource{
			{Name: "source-1"},
			{Name: "source-2"},
		},
	}

	result := FindSourceByName(cfg, "source-2")
	if result == nil {
		t.Fatal("Expected to find source-2")
	}
	if result.Name != "source-2" {
		t.Errorf("Expected source-2, got %s", result.Name)
	}
}

func TestFindSourceByName_NotFound(t *testing.T) {
	cfg := &ProjectConfig{
		Sources: []ProjectSource{
			{Name: "source-1"},
		},
	}

	result := FindSourceByName(cfg, "nonexistent")
	if result != nil {
		t.Error("Expected nil for nonexistent source")
	}
}

func TestUpdateSourceInConfig(t *testing.T) {
	cfg := &ProjectConfig{
		Sources: []ProjectSource{
			{Name: "source-1", Commit: "old"},
			{Name: "source-2", Commit: "old"},
		},
	}

	updated := ProjectSource{Name: "source-1", Commit: "new"}
	err := UpdateSourceInConfig(cfg, "source-1", updated)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if cfg.Sources[0].Commit != "new" {
		t.Errorf("Expected source-1 commit to be 'new', got %s", cfg.Sources[0].Commit)
	}
	if cfg.Sources[1].Commit != "old" {
		t.Errorf("Expected source-2 commit to remain 'old', got %s", cfg.Sources[1].Commit)
	}
}

func TestUpdateSourceInConfig_NotFound(t *testing.T) {
	cfg := &ProjectConfig{
		Sources: []ProjectSource{
			{Name: "source-1"},
		},
	}

	updated := ProjectSource{Name: "nonexistent"}
	err := UpdateSourceInConfig(cfg, "nonexistent", updated)

	if err != nil {
		t.Errorf("Expected no error for nonexistent source, got %v", err)
	}
}
