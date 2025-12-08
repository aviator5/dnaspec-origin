package agents

import (
	"fmt"
	"path/filepath"

	"github.com/aviator5/dnaspec/internal/core/config"
)

// GenerationSummary contains counts of generated files
type GenerationSummary struct {
	AgentsMD        bool
	ClaudeMD        bool
	ClaudeCommands  int
	CopilotPrompts  int
	Errors          []error
}

// GenerateAgentFiles generates all agent integration files based on config and selected agents
func GenerateAgentFiles(cfg *config.ProjectConfig, agents []string) (*GenerationSummary, error) {
	summary := &GenerationSummary{
		Errors: []error{},
	}

	// Always generate AGENTS.md regardless of selected agents
	if err := GenerateAgentsMD(cfg); err != nil {
		summary.Errors = append(summary.Errors, fmt.Errorf("failed to generate AGENTS.md: %w", err))
	} else {
		summary.AgentsMD = true
	}

	// Check if Claude Code is selected
	hasClaudeCode := contains(agents, "claude-code")
	// Check if GitHub Copilot is selected
	hasCopilot := contains(agents, "github-copilot")

	// Generate CLAUDE.md if Claude Code is selected
	if hasClaudeCode {
		if err := GenerateClaudeMD(cfg); err != nil {
			summary.Errors = append(summary.Errors, fmt.Errorf("failed to generate CLAUDE.md: %w", err))
		} else {
			summary.ClaudeMD = true
		}
	}

	// Generate prompt files for each source
	for _, source := range cfg.Sources {
		sourceDir := filepath.Join("dnaspec", source.Name)

		for _, prompt := range source.Prompts {
			// Generate Claude command if Claude Code is selected
			if hasClaudeCode {
				if err := GenerateClaudeCommand(source.Name, prompt, sourceDir); err != nil {
					summary.Errors = append(summary.Errors, fmt.Errorf("failed to generate Claude command for %s/%s: %w", source.Name, prompt.Name, err))
				} else {
					summary.ClaudeCommands++
				}
			}

			// Generate Copilot prompt if GitHub Copilot is selected
			if hasCopilot {
				if err := GenerateCopilotPrompt(source.Name, prompt, sourceDir); err != nil {
					summary.Errors = append(summary.Errors, fmt.Errorf("failed to generate Copilot prompt for %s/%s: %w", source.Name, prompt.Name, err))
				} else {
					summary.CopilotPrompts++
				}
			}
		}
	}

	// Return error if there were any failures
	if len(summary.Errors) > 0 {
		return summary, fmt.Errorf("generation completed with %d errors", len(summary.Errors))
	}

	return summary, nil
}

// contains checks if a string slice contains a value
func contains(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}
