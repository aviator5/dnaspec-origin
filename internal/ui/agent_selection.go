package ui

import (
	"fmt"

	"github.com/aviator5/dnaspec/internal/core/agents"
	"github.com/charmbracelet/huh"
)

// SelectAgents displays an interactive agent selection UI
// Returns selected agent IDs or error if cancelled
func SelectAgents(currentSelection []string) ([]string, error) {
	availableAgents := agents.GetAvailableAgents()

	// Build options for multi-select
	options := make([]huh.Option[string], len(availableAgents))
	for i, agent := range availableAgents {
		label := fmt.Sprintf("%s - %s", agent.DisplayName, agent.Description)
		options[i] = huh.NewOption(label, agent.ID)
	}

	// Prepare initial selection
	var selected []string
	if currentSelection != nil {
		selected = currentSelection
	}

	// Create multi-select form
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Select AI agents to integrate with:").
				Options(options...).
				Value(&selected),
		),
	)

	// Run form
	if err := form.Run(); err != nil {
		return nil, fmt.Errorf("agent selection cancelled: %w", err)
	}

	return selected, nil
}
