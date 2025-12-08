package ui

import (
	"fmt"

	"github.com/aviator5/dnaspec/internal/core/config"
	"github.com/charmbracelet/huh"
)

// SelectGuidelines presents an interactive multi-select form for choosing guidelines
// Returns the selected guidelines or an error
func SelectGuidelines(available []config.ManifestGuideline) ([]config.ManifestGuideline, error) {
	if len(available) == 0 {
		return nil, fmt.Errorf("no guidelines available")
	}

	// Build options
	var options []huh.Option[string]
	for _, g := range available {
		label := fmt.Sprintf("%s - %s", g.Name, g.Description)
		options = append(options, huh.NewOption(label, g.Name))
	}

	// Create multi-select form
	var selected []string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Select guidelines to add:").
				Description("Use space to select/deselect, enter to confirm").
				Options(options...).
				Value(&selected),
		),
	)

	if err := form.Run(); err != nil {
		return nil, err
	}

	// If nothing selected, return empty
	if len(selected) == 0 {
		return []config.ManifestGuideline{}, nil
	}

	// Filter to selected guidelines
	var result []config.ManifestGuideline
	for _, name := range selected {
		for _, g := range available {
			if g.Name == name {
				result = append(result, g)
				break
			}
		}
	}

	return result, nil
}

// SelectGuidelinesByName selects guidelines by their names
// Validates that all requested names exist in the available guidelines
func SelectGuidelinesByName(available []config.ManifestGuideline, names []string) ([]config.ManifestGuideline, error) {
	if len(names) == 0 {
		return nil, fmt.Errorf("no guideline names provided")
	}

	// Build map of available guidelines
	availableMap := make(map[string]config.ManifestGuideline)
	for _, g := range available {
		availableMap[g.Name] = g
	}

	// Validate all names exist and collect selected guidelines
	var result []config.ManifestGuideline
	var missing []string

	for _, name := range names {
		if g, ok := availableMap[name]; ok {
			result = append(result, g)
		} else {
			missing = append(missing, name)
		}
	}

	// Report missing guidelines
	if len(missing) > 0 {
		// Build list of available names
		var availableNames []string
		for _, g := range available {
			availableNames = append(availableNames, g.Name)
		}
		return nil, fmt.Errorf("guidelines not found: %v (available: %v)", missing, availableNames)
	}

	return result, nil
}
