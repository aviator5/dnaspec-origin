package project

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/aviator5/dnaspec/internal/core/config"
	"github.com/aviator5/dnaspec/internal/ui"
)

// NewListCmd creates the list command for displaying project configuration
func NewListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Display configured DNA sources, guidelines, prompts, and agents",
		Long: `Display all configured DNA sources with their metadata, including:
- Configured AI agents (Claude Code, GitHub Copilot, etc.)
- DNA sources with type-specific metadata (URL/path, ref, commit)
- Guidelines and prompts for each source

This command provides a quick overview of the current DNA configuration.`,
		Example: `  # Display current configuration
  dnaspec list`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runList()
		},
	}

	return cmd
}

func runList() error {
	// Load project configuration
	cfg, err := config.LoadProjectConfig(projectConfigFileName)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println(ui.ErrorStyle.Render("✗ Error:"), "Project configuration not found:", ui.CodeStyle.Render(projectConfigFileName))
			fmt.Println(
				ui.SubtleStyle.Render("  Run"), ui.CodeStyle.Render("dnaspec init"),
				ui.SubtleStyle.Render("to create a new project configuration."),
			)
			return fmt.Errorf("project configuration not found")
		}
		fmt.Println(ui.ErrorStyle.Render("✗ Error:"), "Failed to load project configuration:", err)
		return err
	}

	// Display agents
	displayAgents(cfg)

	// Display sources
	fmt.Println()
	displaySources(cfg)

	return nil
}

func displayAgents(cfg *config.ProjectConfig) {
	fmt.Println("Configured Agents (Phase 1):")
	if len(cfg.Agents) == 0 {
		fmt.Println("  None configured")
	} else {
		for _, agent := range cfg.Agents {
			// Map agent IDs to display names
			displayName := agent
			switch agent {
			case "claude-code":
				displayName = "Claude Code"
			case "github-copilot":
				displayName = "GitHub Copilot"
			}
			fmt.Printf("  - %s\n", displayName)
		}
	}
}

func displaySources(cfg *config.ProjectConfig) {
	fmt.Println("Sources:")
	if len(cfg.Sources) == 0 {
		fmt.Println("  No sources configured")
		return
	}

	fmt.Println()
	for i := range cfg.Sources {
		source := &cfg.Sources[i]
		// Display source name with type
		fmt.Printf("%s (%s)\n", source.Name, source.Type)

		// Display type-specific metadata
		switch source.Type {
		case config.SourceTypeGitRepo:
			fmt.Printf("  URL: %s\n", source.URL)
			if source.Ref != "" {
				fmt.Printf("  Ref: %s\n", source.Ref)
			}
			if source.Commit != "" {
				// Display short commit hash (first 8 characters)
				commitDisplay := source.Commit
				if len(commitDisplay) > 8 {
					commitDisplay = commitDisplay[:8]
				}
				fmt.Printf("  Commit: %s\n", commitDisplay)
			}
		case config.SourceTypeLocalPath:
			fmt.Printf("  Path: %s\n", source.Path)
		}

		// Display guidelines
		displayGuidelines(*source)

		// Display prompts
		displayPrompts(*source)

		// Add blank line between sources (except after the last one)
		if i < len(cfg.Sources)-1 {
			fmt.Println()
		}
	}
}

func displayGuidelines(source config.ProjectSource) {
	fmt.Println()
	fmt.Println("  Guidelines:")
	if len(source.Guidelines) == 0 {
		fmt.Println("    None")
	} else {
		for _, guideline := range source.Guidelines {
			fmt.Printf("    - %s: %s\n", guideline.Name, guideline.Description)
		}
	}
}

func displayPrompts(source config.ProjectSource) {
	fmt.Println()
	fmt.Println("  Prompts:")
	if len(source.Prompts) == 0 {
		fmt.Println("    None")
	} else {
		for _, prompt := range source.Prompts {
			fmt.Printf("    - %s: %s\n", prompt.Name, prompt.Description)
		}
	}
}
