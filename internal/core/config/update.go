package config

import (
	"fmt"
)

// AddSource adds a new source to the project configuration
// Validates that the source name doesn't already exist
// Extracts prompts referenced by selected guidelines
func AddSource(cfg *ProjectConfig, source ProjectSource) error {
	// Check for duplicate source names
	for _, existing := range cfg.Sources {
		if existing.Name == source.Name {
			return fmt.Errorf("source with name '%s' already exists, use --name to specify a different name", source.Name)
		}
	}

	// Append source to config
	cfg.Sources = append(cfg.Sources, source)

	return nil
}

// ExtractReferencedPrompts filters prompts to only those referenced by selected guidelines
func ExtractReferencedPrompts(selectedGuidelines []ManifestGuideline, allPrompts []ManifestPrompt) []ProjectPrompt {
	// Build set of referenced prompt names
	referenced := make(map[string]bool)
	for _, g := range selectedGuidelines {
		for _, pName := range g.Prompts {
			referenced[pName] = true
		}
	}

	// Filter prompts to referenced ones
	var result []ProjectPrompt
	for _, p := range allPrompts {
		if referenced[p.Name] {
			result = append(result, ProjectPrompt{
				Name:        p.Name,
				File:        p.File,
				Description: p.Description,
			})
		}
	}

	return result
}

// ManifestGuidelinesToProject converts manifest guidelines to project guidelines
func ManifestGuidelinesToProject(guidelines []ManifestGuideline) []ProjectGuideline {
	result := make([]ProjectGuideline, len(guidelines))
	for i, g := range guidelines {
		result[i] = ProjectGuideline{
			Name:                g.Name,
			File:                g.File,
			Description:         g.Description,
			ApplicableScenarios: g.ApplicableScenarios,
			Prompts:             g.Prompts,
		}
	}
	return result
}
