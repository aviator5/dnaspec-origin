package config

import (
	"testing"
)

func TestAddSource(t *testing.T) {
	t.Run("add source to empty config", func(t *testing.T) {
		cfg := &ProjectConfig{
			Version: 1,
			Sources: []ProjectSource{},
		}

		newSource := ProjectSource{
			Name: "test-source",
			Type: "git-repo",
			URL:  "https://github.com/test/repo",
		}

		err := AddSource(cfg, newSource)
		if err != nil {
			t.Fatalf("AddSource() error = %v", err)
		}

		if len(cfg.Sources) != 1 {
			t.Errorf("Expected 1 source, got %d", len(cfg.Sources))
		}

		if cfg.Sources[0].Name != "test-source" {
			t.Errorf("Expected source name 'test-source', got %s", cfg.Sources[0].Name)
		}
	})

	t.Run("add multiple sources", func(t *testing.T) {
		cfg := &ProjectConfig{
			Version: 1,
			Sources: []ProjectSource{},
		}

		source1 := ProjectSource{Name: "source-1", Type: "git-repo"}
		source2 := ProjectSource{Name: "source-2", Type: "local-path"}

		if err := AddSource(cfg, source1); err != nil {
			t.Fatalf("AddSource() error = %v", err)
		}

		if err := AddSource(cfg, source2); err != nil {
			t.Fatalf("AddSource() error = %v", err)
		}

		if len(cfg.Sources) != 2 {
			t.Errorf("Expected 2 sources, got %d", len(cfg.Sources))
		}
	})

	t.Run("error on duplicate source name", func(t *testing.T) {
		cfg := &ProjectConfig{
			Version: 1,
			Sources: []ProjectSource{
				{Name: "existing-source", Type: "git-repo"},
			},
		}

		newSource := ProjectSource{
			Name: "existing-source",
			Type: "git-repo",
		}

		err := AddSource(cfg, newSource)
		if err == nil {
			t.Error("Expected error for duplicate source name, got nil")
		}
	})
}

func TestExtractReferencedPrompts(t *testing.T) {
	t.Run("extract prompts referenced by guidelines", func(t *testing.T) {
		guidelines := []ManifestGuideline{
			{
				Name:    "guideline-1",
				Prompts: []string{"prompt-1", "prompt-2"},
			},
			{
				Name:    "guideline-2",
				Prompts: []string{"prompt-2", "prompt-3"},
			},
		}

		allPrompts := []ManifestPrompt{
			{Name: "prompt-1", File: "prompts/p1.md", Description: "Prompt 1"},
			{Name: "prompt-2", File: "prompts/p2.md", Description: "Prompt 2"},
			{Name: "prompt-3", File: "prompts/p3.md", Description: "Prompt 3"},
			{Name: "prompt-4", File: "prompts/p4.md", Description: "Prompt 4"},
		}

		result := ExtractReferencedPrompts(guidelines, allPrompts)

		if len(result) != 3 {
			t.Errorf("Expected 3 prompts, got %d", len(result))
		}

		// Check that prompt-4 is not included
		for _, p := range result {
			if p.Name == "prompt-4" {
				t.Error("Unreferenced prompt-4 should not be included")
			}
		}

		// Check that all referenced prompts are included
		expectedNames := map[string]bool{"prompt-1": true, "prompt-2": true, "prompt-3": true}
		for _, p := range result {
			if !expectedNames[p.Name] {
				t.Errorf("Unexpected prompt: %s", p.Name)
			}
		}
	})

	t.Run("no prompts referenced", func(t *testing.T) {
		guidelines := []ManifestGuideline{
			{Name: "guideline-1", Prompts: []string{}},
		}

		allPrompts := []ManifestPrompt{
			{Name: "prompt-1", File: "prompts/p1.md"},
		}

		result := ExtractReferencedPrompts(guidelines, allPrompts)

		if len(result) != 0 {
			t.Errorf("Expected 0 prompts, got %d", len(result))
		}
	})

	t.Run("empty guidelines", func(t *testing.T) {
		guidelines := []ManifestGuideline{}
		allPrompts := []ManifestPrompt{
			{Name: "prompt-1", File: "prompts/p1.md"},
		}

		result := ExtractReferencedPrompts(guidelines, allPrompts)

		if len(result) != 0 {
			t.Errorf("Expected 0 prompts, got %d", len(result))
		}
	})
}

func TestManifestGuidelinesToProject(t *testing.T) {
	t.Run("convert manifest guidelines to project guidelines", func(t *testing.T) {
		manifestGuidelines := []ManifestGuideline{
			{
				Name:                "test-guideline",
				File:                "guidelines/test.md",
				Description:         "Test guideline",
				ApplicableScenarios: []string{"scenario-1", "scenario-2"},
				Prompts:             []string{"prompt-1"},
			},
		}

		result := ManifestGuidelinesToProject(manifestGuidelines)

		if len(result) != 1 {
			t.Fatalf("Expected 1 guideline, got %d", len(result))
		}

		g := result[0]
		if g.Name != "test-guideline" {
			t.Errorf("Name = %s, want test-guideline", g.Name)
		}
		if g.File != "guidelines/test.md" {
			t.Errorf("File = %s, want guidelines/test.md", g.File)
		}
		if g.Description != "Test guideline" {
			t.Errorf("Description = %s, want 'Test guideline'", g.Description)
		}
		if len(g.ApplicableScenarios) != 2 {
			t.Errorf("ApplicableScenarios length = %d, want 2", len(g.ApplicableScenarios))
		}
		if len(g.Prompts) != 1 {
			t.Errorf("Prompts length = %d, want 1", len(g.Prompts))
		}
	})

	t.Run("empty list", func(t *testing.T) {
		result := ManifestGuidelinesToProject([]ManifestGuideline{})

		if len(result) != 0 {
			t.Errorf("Expected 0 guidelines, got %d", len(result))
		}
	})
}

func TestUpdateAgents(t *testing.T) {
	t.Run("update agents and save config", func(t *testing.T) {
		tmpDir := t.TempDir()
		configPath := tmpDir + "/dnaspec.yaml"

		// Create initial config with no agents
		cfg := &ProjectConfig{
			Version: 1,
			Agents:  []string{},
			Sources: []ProjectSource{},
		}

		// Save it
		if err := SaveProjectConfig(configPath, cfg); err != nil {
			t.Fatalf("SaveProjectConfig() error = %v", err)
		}

		// Load it
		loaded, err := LoadProjectConfig(configPath)
		if err != nil {
			t.Fatalf("LoadProjectConfig() error = %v", err)
		}

		// Update agents
		newAgents := []string{"claude-code", "github-copilot"}
		UpdateAgents(loaded, newAgents)

		// Save the updated config
		if err := SaveProjectConfig(configPath, loaded); err != nil {
			t.Fatalf("SaveProjectConfig() error = %v", err)
		}

		// Load again to verify
		final, err := LoadProjectConfig(configPath)
		if err != nil {
			t.Fatalf("LoadProjectConfig() error = %v", err)
		}

		if len(final.Agents) != 2 {
			t.Errorf("Expected 2 agents, got %d", len(final.Agents))
		}

		if final.Agents[0] != "claude-code" || final.Agents[1] != "github-copilot" {
			t.Errorf("Agents = %v, want [claude-code github-copilot]", final.Agents)
		}
	})

	t.Run("update agents replaces previous selection", func(t *testing.T) {
		tmpDir := t.TempDir()
		configPath := tmpDir + "/dnaspec.yaml"

		// Create config with initial agent selection
		cfg := &ProjectConfig{
			Version: 1,
			Agents:  []string{"claude-code"},
			Sources: []ProjectSource{},
		}

		if err := SaveProjectConfig(configPath, cfg); err != nil {
			t.Fatalf("SaveProjectConfig() error = %v", err)
		}

		loaded, err := LoadProjectConfig(configPath)
		if err != nil {
			t.Fatalf("LoadProjectConfig() error = %v", err)
		}

		// Update to different agents
		UpdateAgents(loaded, []string{"github-copilot"})

		// Save the updated config
		if err := SaveProjectConfig(configPath, loaded); err != nil {
			t.Fatalf("SaveProjectConfig() error = %v", err)
		}

		// Load and verify
		final, err := LoadProjectConfig(configPath)
		if err != nil {
			t.Fatalf("LoadProjectConfig() error = %v", err)
		}

		if len(final.Agents) != 1 {
			t.Errorf("Expected 1 agent, got %d", len(final.Agents))
		}

		if final.Agents[0] != "github-copilot" {
			t.Errorf("Agents = %v, want [github-copilot]", final.Agents)
		}
	})
}

func TestRoundTrip(t *testing.T) {
	t.Run("load, modify, save, load config", func(t *testing.T) {
		tmpDir := t.TempDir()
		configPath := tmpDir + "/dnaspec.yaml"

		// Create initial config
		cfg := &ProjectConfig{
			Version: 1,
			Agents:  []string{"claude-code"},
			Sources: []ProjectSource{
				{
					Name: "source-1",
					Type: "git-repo",
					URL:  "https://github.com/test/repo",
				},
			},
		}

		// Save it
		if err := SaveProjectConfig(configPath, cfg); err != nil {
			t.Fatalf("SaveProjectConfig() error = %v", err)
		}

		// Load it back
		loaded, err := LoadProjectConfig(configPath)
		if err != nil {
			t.Fatalf("LoadProjectConfig() error = %v", err)
		}

		// Modify it
		newSource := ProjectSource{
			Name: "source-2",
			Type: "local-path",
			Path: "/path/to/local",
		}
		if err := AddSource(loaded, newSource); err != nil {
			t.Fatalf("AddSource() error = %v", err)
		}

		// Save again
		if err := AtomicWriteProjectConfig(configPath, loaded); err != nil {
			t.Fatalf("AtomicWriteProjectConfig() error = %v", err)
		}

		// Load final version
		final, err := LoadProjectConfig(configPath)
		if err != nil {
			t.Fatalf("LoadProjectConfig() error = %v", err)
		}

		// Verify
		if len(final.Sources) != 2 {
			t.Errorf("Expected 2 sources, got %d", len(final.Sources))
		}

		if final.Sources[1].Name != "source-2" {
			t.Errorf("Second source name = %s, want source-2", final.Sources[1].Name)
		}
	})
}
