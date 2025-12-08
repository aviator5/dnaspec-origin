package config

import (
	"slices"
)

// GuidelineComparison categorizes guidelines into updated, new, removed, and unchanged
type GuidelineComparison struct {
	Updated   []string // Guidelines that exist in both but have changes
	New       []string // Guidelines in manifest but not in config
	Removed   []string // Guidelines in config but not in manifest
	Unchanged []string // Guidelines with no changes
}

// CompareGuidelines compares current project guidelines with latest manifest guidelines
// Returns a categorization of all guidelines into updated, new, removed, and unchanged
func CompareGuidelines(
	currentGuidelines []ProjectGuideline,
	manifestGuidelines []ManifestGuideline,
) GuidelineComparison {
	// Build maps for O(1) lookup
	currentMap := make(map[string]ProjectGuideline)
	manifestMap := make(map[string]ManifestGuideline)

	for _, g := range currentGuidelines {
		currentMap[g.Name] = g
	}
	for _, g := range manifestGuidelines {
		manifestMap[g.Name] = g
	}

	var result GuidelineComparison

	// Check each current guideline
	for name, current := range currentMap {
		if manifest, exists := manifestMap[name]; exists {
			if hasChanges(current, manifest) {
				result.Updated = append(result.Updated, name)
			} else {
				result.Unchanged = append(result.Unchanged, name)
			}
		} else {
			result.Removed = append(result.Removed, name)
		}
	}

	// Find new guidelines
	for name := range manifestMap {
		if _, exists := currentMap[name]; !exists {
			result.New = append(result.New, name)
		}
	}

	return result
}

// hasChanges detects if guideline metadata has changed
// Note: File content changes will be copied regardless of metadata changes
func hasChanges(current ProjectGuideline, manifest ManifestGuideline) bool {
	if current.Description != manifest.Description {
		return true
	}
	if !slices.Equal(current.ApplicableScenarios, manifest.ApplicableScenarios) {
		return true
	}
	if !slices.Equal(current.Prompts, manifest.Prompts) {
		return true
	}
	return false
}

// FindSourceByName finds a source in the project config by name
// Returns nil if source not found
func FindSourceByName(cfg *ProjectConfig, name string) *ProjectSource {
	for i := range cfg.Sources {
		if cfg.Sources[i].Name == name {
			return &cfg.Sources[i]
		}
	}
	return nil
}

// UpdateSourceInConfig updates a source in the project configuration
// Returns an error if the source is not found
func UpdateSourceInConfig(cfg *ProjectConfig, sourceName string, updatedSource ProjectSource) error {
	for i := range cfg.Sources {
		if cfg.Sources[i].Name == sourceName {
			cfg.Sources[i] = updatedSource
			return nil
		}
	}
	return nil
}
