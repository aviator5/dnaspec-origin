package agents

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAvailableAgents(t *testing.T) {
	agents := GetAvailableAgents()

	assert.Len(t, agents, 2, "Should return 2 phase 1 agents")
	assert.Equal(t, "claude-code", agents[0].ID)
	assert.Equal(t, "github-copilot", agents[1].ID)
}

func TestIsValidAgent(t *testing.T) {
	tests := []struct {
		name     string
		agentID  string
		expected bool
	}{
		{
			name:     "valid claude-code",
			agentID:  "claude-code",
			expected: true,
		},
		{
			name:     "valid github-copilot",
			agentID:  "github-copilot",
			expected: true,
		},
		{
			name:     "invalid agent",
			agentID:  "invalid-agent",
			expected: false,
		},
		{
			name:     "empty string",
			agentID:  "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidAgent(tt.agentID)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetAgent(t *testing.T) {
	tests := []struct {
		name         string
		agentID      string
		expectNil    bool
		expectedID   string
		expectedName string
	}{
		{
			name:         "get claude-code",
			agentID:      "claude-code",
			expectNil:    false,
			expectedID:   "claude-code",
			expectedName: "Claude Code",
		},
		{
			name:         "get github-copilot",
			agentID:      "github-copilot",
			expectNil:    false,
			expectedID:   "github-copilot",
			expectedName: "GitHub Copilot",
		},
		{
			name:      "get invalid agent",
			agentID:   "invalid",
			expectNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agent := GetAgent(tt.agentID)

			if tt.expectNil {
				assert.Nil(t, agent)
			} else {
				assert.NotNil(t, agent)
				assert.Equal(t, tt.expectedID, agent.ID)
				assert.Equal(t, tt.expectedName, agent.DisplayName)
			}
		})
	}
}
