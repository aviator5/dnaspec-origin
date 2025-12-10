package project

import (
	"fmt"
	"os"

	"github.com/aviator5/dnaspec/internal/core/agents"
	"github.com/aviator5/dnaspec/internal/core/config"
	"github.com/aviator5/dnaspec/internal/ui"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var (
	noAskFlag bool
)

// NewUpdateAgentsCmd creates the update-agents command
func NewUpdateAgentsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-agents",
		Short: "Configure AI agents and generate agent integration files",
		Long: `Configure AI agents and generate agent integration files.

This command allows you to select which AI agents to integrate with (Antigravity, Claude Code,
Cursor, GitHub Copilot, Windsurf) and generates the necessary files for each agent:

- AGENTS.md: Context-aware guideline references for all AI agents
- CLAUDE.md: Same as AGENTS.md, for Claude Code discovery
- Claude commands: Slash commands in .claude/commands/dnaspec/
- Copilot prompts: Prompt files in .github/prompts/
- Antigravity prompts: Workflow files in .agent/workflows/
- Windsurf workflows: Workflow files in .windsurf/workflows/
- Cursor commands: Command files in .cursor/commands/

Use --no-ask to skip agent selection and use saved configuration.`,
		RunE: runUpdateAgents,
	}

	cmd.Flags().BoolVar(&noAskFlag, "no-ask", false, "Skip agent selection, use saved configuration")

	return cmd
}

func runUpdateAgents(cmd *cobra.Command, args []string) error {
	// Load project configuration
	cfg, err := config.LoadProjectConfig("dnaspec.yaml")
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("dnaspec.yaml not found, run 'dnaspec init' first")
		}
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Check if sources are configured
	if len(cfg.Sources) == 0 {
		// No sources configured, cleanup any existing DNASPEC blocks
		fmt.Println(ui.InfoStyle.Render("No DNA sources configured."))
		fmt.Println(ui.InfoStyle.Render("Checking for DNASPEC blocks to remove..."))

		summary, err := agents.CleanupAgentFiles()
		if err != nil {
			return fmt.Errorf("failed to cleanup agent files: %w", err)
		}

		// Check if any files were cleaned
		if !summary.AgentsMDCleaned && !summary.ClaudeMDCleaned {
			fmt.Println(ui.InfoStyle.Render("No DNASPEC blocks found to remove."))
			fmt.Println(ui.InfoStyle.Render("Run 'dnaspec add' to add guidelines first."))
			return nil
		}

		// Display what was cleaned
		fmt.Println(ui.SuccessStyle.Render("\nRemoved DNASPEC blocks from:"))
		if summary.AgentsMDCleaned {
			fmt.Println(ui.SuccessStyle.Render("  ✓ AGENTS.md"))
		}
		if summary.ClaudeMDCleaned {
			fmt.Println(ui.SuccessStyle.Render("  ✓ CLAUDE.md"))
		}

		return nil
	}

	var selectedAgents []string

	if noAskFlag {
		// Non-interactive mode: use saved configuration
		if len(cfg.Agents) == 0 {
			return fmt.Errorf("no agents configured, run without --no-ask to select agents")
		}
		selectedAgents = cfg.Agents
		fmt.Println(ui.InfoStyle.Render(fmt.Sprintf("Using saved agents: %v", selectedAgents)))
	} else {
		// Interactive mode: prompt for agent selection
		selected, err := ui.SelectAgents(cfg.Agents)
		if err != nil {
			return fmt.Errorf("agent selection canceled: %w", err)
		}
		selectedAgents = selected

		// Save agents to config
		config.UpdateAgents(cfg, selectedAgents)
		if err := config.SaveProjectConfig("dnaspec.yaml", cfg); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		fmt.Println(ui.SuccessStyle.Render("✓ Updated dnaspec.yaml"))
	}

	// Generate agent files
	fmt.Println(ui.InfoStyle.Render("\nGenerating agent files..."))

	summary, err := agents.GenerateAgentFiles(cfg, selectedAgents)

	// Display summary
	displaySummary(summary)

	if err != nil {
		return err
	}

	fmt.Println(ui.SuccessStyle.Render("\n✓ Agent files generated successfully"))

	return nil
}

func displaySummary(summary *agents.GenerationSummary) {
	successStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("42"))

	if summary.AgentsMD {
		fmt.Println(successStyle.Render("  ✓ AGENTS.md"))
	}

	if summary.ClaudeMD {
		fmt.Println(successStyle.Render("  ✓ CLAUDE.md"))
	}

	if summary.ClaudeCommands > 0 {
		fmt.Println(successStyle.Render(fmt.Sprintf("  ✓ Generated %d Claude command(s)", summary.ClaudeCommands)))
	}

	if summary.CopilotPrompts > 0 {
		fmt.Println(successStyle.Render(fmt.Sprintf("  ✓ Generated %d Copilot prompt(s)", summary.CopilotPrompts)))
	}

	if summary.AntigravityPrompts > 0 {
		fmt.Println(successStyle.Render(fmt.Sprintf("  ✓ Generated %d Antigravity prompt(s)", summary.AntigravityPrompts)))
	}

	if summary.WindsurfWorkflows > 0 {
		fmt.Println(successStyle.Render(fmt.Sprintf("  ✓ Generated %d Windsurf workflow(s)", summary.WindsurfWorkflows)))
	}

	if summary.CursorCommands > 0 {
		fmt.Println(successStyle.Render(fmt.Sprintf("  ✓ Generated %d Cursor command(s)", summary.CursorCommands)))
	}

	if len(summary.Errors) > 0 {
		errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
		fmt.Println(errorStyle.Render(fmt.Sprintf("\n  %d error(s) occurred:", len(summary.Errors))))
		for _, err := range summary.Errors {
			fmt.Println(errorStyle.Render(fmt.Sprintf("    • %s", err.Error())))
		}
	}
}
