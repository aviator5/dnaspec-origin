package agents

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/aviator5/dnaspec/internal/core/config"
	"github.com/aviator5/dnaspec/internal/core/files"
)

// GenerationSummary contains counts of generated files
type GenerationSummary struct {
	AgentsMD          bool
	ClaudeMD          bool
	ClaudeCommands    int
	CopilotPrompts    int
	AntigravityPrompts int
	WindsurfWorkflows int
	CursorCommands    int
	Errors            []error
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
	// Check if Antigravity is selected
	hasAntigravity := contains(agents, "antigravity")
	// Check if Windsurf is selected
	hasWindsurf := contains(agents, "windsurf")
	// Check if Cursor is selected
	hasCursor := contains(agents, "cursor")

	// Generate CLAUDE.md if Claude Code is selected
	if hasClaudeCode {
		if err := GenerateClaudeMD(cfg); err != nil {
			summary.Errors = append(summary.Errors, fmt.Errorf("failed to generate CLAUDE.md: %w", err))
		} else {
			summary.ClaudeMD = true
		}
	}

	// Generate prompt files for each source
	for i := range cfg.Sources {
		source := &cfg.Sources[i]
		sourceDir := filepath.Join("dnaspec", source.Name)

		for _, prompt := range source.Prompts {
			generatePromptFiles(source.Name, prompt, sourceDir, summary,
				hasClaudeCode, hasCopilot, hasAntigravity, hasWindsurf, hasCursor)
		}
	}

	// Return error if there were any failures
	if len(summary.Errors) > 0 {
		return summary, fmt.Errorf("generation completed with %d errors", len(summary.Errors))
	}

	return summary, nil
}

// generatePromptFiles generates prompt files for a single prompt across all selected agents
func generatePromptFiles(sourceName string, prompt config.ProjectPrompt, sourceDir string,
	summary *GenerationSummary, hasClaudeCode, hasCopilot, hasAntigravity, hasWindsurf, hasCursor bool) {
	// Generate Claude command if Claude Code is selected
	if hasClaudeCode {
		if err := GenerateClaudeCommand(sourceName, prompt, sourceDir); err != nil {
			summary.Errors = append(summary.Errors,
				fmt.Errorf("failed to generate Claude command for %s/%s: %w",
					sourceName, prompt.Name, err))
		} else {
			summary.ClaudeCommands++
		}
	}

	// Generate Copilot prompt if GitHub Copilot is selected
	if hasCopilot {
		if err := GenerateCopilotPrompt(sourceName, prompt, sourceDir); err != nil {
			summary.Errors = append(summary.Errors,
				fmt.Errorf("failed to generate Copilot prompt for %s/%s: %w",
					sourceName, prompt.Name, err))
		} else {
			summary.CopilotPrompts++
		}
	}

	// Generate Antigravity prompt if Antigravity is selected
	if hasAntigravity {
		if err := GenerateAntigravityPrompt(sourceName, prompt, sourceDir); err != nil {
			summary.Errors = append(summary.Errors,
				fmt.Errorf("failed to generate Antigravity prompt for %s/%s: %w",
					sourceName, prompt.Name, err))
		} else {
			summary.AntigravityPrompts++
		}
	}

	// Generate Windsurf workflow if Windsurf is selected
	if hasWindsurf {
		if err := GenerateWindsurfPrompt(sourceName, prompt, sourceDir); err != nil {
			summary.Errors = append(summary.Errors,
				fmt.Errorf("failed to generate Windsurf workflow for %s/%s: %w",
					sourceName, prompt.Name, err))
		} else {
			summary.WindsurfWorkflows++
		}
	}

	// Generate Cursor command if Cursor is selected
	if hasCursor {
		if err := GenerateCursorCommand(sourceName, prompt, sourceDir); err != nil {
			summary.Errors = append(summary.Errors,
				fmt.Errorf("failed to generate Cursor command for %s/%s: %w",
					sourceName, prompt.Name, err))
		} else {
			summary.CursorCommands++
		}
	}
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

// CleanupSummary contains information about cleanup actions
type CleanupSummary struct {
	AgentsMDCleaned bool
	ClaudeMDCleaned bool
}

// CleanupAgentFiles removes DNASPEC blocks from AGENTS.md and CLAUDE.md if they exist
// Returns a summary of what was cleaned up
func CleanupAgentFiles() (*CleanupSummary, error) {
	summary := &CleanupSummary{}

	// Clean up AGENTS.md
	if err := cleanupFile("AGENTS.md"); err == nil {
		summary.AgentsMDCleaned = true
	} else if !os.IsNotExist(err) {
		return summary, fmt.Errorf("failed to cleanup AGENTS.md: %w", err)
	}

	// Clean up CLAUDE.md
	if err := cleanupFile("CLAUDE.md"); err == nil {
		summary.ClaudeMDCleaned = true
	} else if !os.IsNotExist(err) {
		return summary, fmt.Errorf("failed to cleanup CLAUDE.md: %w", err)
	}

	return summary, nil
}

// cleanupFile removes DNASPEC block from a file if it exists
// Returns nil if block was removed, os.ErrNotExist if file doesn't exist,
// or other error if something went wrong
func cleanupFile(path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	cleanedContent, removed := files.RemoveManagedBlock(string(content))
	if !removed {
		// No block to remove, nothing to do
		return os.ErrNotExist
	}

	// Write the cleaned content back
	return writeFileAtomic(path, []byte(cleanedContent))
}
