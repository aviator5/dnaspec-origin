package agents

// Agent represents an AI agent that can consume DNA guidelines
type Agent struct {
	ID          string
	DisplayName string
	Description string
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
