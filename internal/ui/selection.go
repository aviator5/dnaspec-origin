package ui

import (
	"fmt"

	"github.com/charmbracelet/huh"

	"github.com/aviator5/dnaspec/internal/core/config"
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

// testMockSelection allows tests to bypass interactive UI
// When set, this function will be called instead of showing the UI
var testMockSelection func(available []config.ManifestGuideline, existing []string, orphaned []config.ProjectGuideline) ([]string, error)

// SetTestMockSelection sets a mock function for testing
// This allows tests to bypass interactive UI
func SetTestMockSelection(fn func(available []config.ManifestGuideline, existing []string, orphaned []config.ProjectGuideline) ([]string, error)) {
	testMockSelection = fn
}

// SelectGuidelinesWithStatus presents an interactive multi-select form for choosing guidelines
// with support for marking existing and orphaned guidelines.
//
// - available: Guidelines from the source manifest (new + existing + updated)
// - existing: Guidelines currently in project config that are also in manifest
// - orphaned: Guidelines in project config but missing from manifest
//
// Returns selected guideline names (not including orphaned ones)
func SelectGuidelinesWithStatus(
	available []config.ManifestGuideline,
	existing []string,
	orphaned []config.ProjectGuideline,
) ([]string, error) {
	// Allow tests to mock the selection
	if testMockSelection != nil {
		return testMockSelection(available, existing, orphaned)
	}

	if len(available) == 0 && len(orphaned) == 0 {
		return nil, fmt.Errorf("no guidelines available")
	}

	// Build maps for quick lookup
	existingMap := make(map[string]bool)
	for _, name := range existing {
		existingMap[name] = true
	}

	// Build options list with available guidelines
	var options []huh.Option[string]
	for _, g := range available {
		label := fmt.Sprintf("%s - %s", g.Name, g.Description)
		options = append(options, huh.NewOption(label, g.Name))
	}

	// Add orphaned guidelines at the end with warning icon
	for _, g := range orphaned {
		label := fmt.Sprintf("%s - %s ⚠️", g.Name, g.Description)
		options = append(options, huh.NewOption(label, g.Name))
	}

	// Pre-select existing and orphaned guidelines
	var preSelected []string
	preSelected = append(preSelected, existing...)
	for _, g := range orphaned {
		preSelected = append(preSelected, g.Name)
	}

	// Create multi-select form
	selected := preSelected // Initialize with pre-selected values
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Select guidelines to keep or add:").
				Description("Use space to select/deselect, enter to confirm. ⚠️ = missing from source").
				Options(options...).
				Value(&selected),
		),
	)

	if err := form.Run(); err != nil {
		return nil, err
	}

	// Filter out orphaned guidelines from the result
	// (user can deselect them to remove, but they're not in the source manifest)
	orphanedMap := make(map[string]bool)
	for _, g := range orphaned {
		orphanedMap[g.Name] = true
	}

	var result []string
	for _, name := range selected {
		if !orphanedMap[name] {
			result = append(result, name)
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
