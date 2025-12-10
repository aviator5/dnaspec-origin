package agents

import (
	"fmt"
	"path/filepath"
)

// Agent represents an AI agent that can consume DNA guidelines
type Agent struct {
	ID          string
	DisplayName string
	Description string
}

// AgentFilePattern defines the file pattern for a specific agent
type AgentFilePattern struct {
	AgentID       string
	PatternFormat string // Format string with sourceName placeholder
	DisplayFormat string // Format string for displaying to users
}

// Phase1Agents are the initially supported agents
var Phase1Agents = []Agent{
	{
		ID:          "antigravity",
		DisplayName: "Antigravity",
		Description: "AI development assistant",
	},
	{
		ID:          "claude-code",
		DisplayName: "Claude Code",
		Description: "Anthropic's AI assistant with slash commands",
	},
	{
		ID:          "cursor",
		DisplayName: "Cursor",
		Description: "AI-first code editor",
	},
	{
		ID:          "github-copilot",
		DisplayName: "GitHub Copilot",
		Description: "GitHub's AI pair programmer",
	},
	{
		ID:          "windsurf",
		DisplayName: "Windsurf",
		Description: "AI-powered code editor",
	},
}

// AgentFilePatterns defines file patterns for all supported agents
var AgentFilePatterns = []AgentFilePattern{
	{
		AgentID:       "antigravity",
		PatternFormat: ".agent/workflows/dnaspec-%s-*.md",
		DisplayFormat: ".agent/workflows/dnaspec-%s-*.md",
	},
	{
		AgentID:       "claude-code",
		PatternFormat: ".claude/commands/dnaspec/%s-*.md",
		DisplayFormat: ".claude/commands/dnaspec/%s-*.md",
	},
	{
		AgentID:       "cursor",
		PatternFormat: ".cursor/commands/dnaspec-%s-*.md",
		DisplayFormat: ".cursor/commands/dnaspec-%s-*.md",
	},
	{
		AgentID:       "github-copilot",
		PatternFormat: ".github/prompts/dnaspec-%s-*.prompt.md",
		DisplayFormat: ".github/prompts/dnaspec-%s-*.prompt.md",
	},
	{
		AgentID:       "windsurf",
		PatternFormat: ".windsurf/workflows/dnaspec-%s-*.md",
		DisplayFormat: ".windsurf/workflows/dnaspec-%s-*.md",
	},
}

// GetFilePatternForSource returns the file glob pattern for a specific source
func (afp AgentFilePattern) GetFilePatternForSource(sourceName string) string {
	return filepath.FromSlash(fmt.Sprintf(afp.PatternFormat, sourceName))
}

// GetDisplayPatternForSource returns the display string for a specific source
func (afp AgentFilePattern) GetDisplayPatternForSource(sourceName string) string {
	return fmt.Sprintf(afp.DisplayFormat, sourceName)
}

// GetAvailableAgents returns the list of supported agents
func GetAvailableAgents() []Agent {
	return Phase1Agents
}

// IsValidAgent checks if the given agent ID is supported
func IsValidAgent(id string) bool {
	for _, agent := range Phase1Agents {
		if agent.ID == id {
			return true
		}
	}
	return false
}

// GetAgent returns the agent with the given ID, or nil if not found
func GetAgent(id string) *Agent {
	for _, agent := range Phase1Agents {
		if agent.ID == id {
			return &agent
		}
	}
	return nil
}
