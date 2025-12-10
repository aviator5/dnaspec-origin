package project

import (
	"fmt"
	"os"

	"github.com/aviator5/dnaspec/internal/core/config"
	"github.com/aviator5/dnaspec/internal/ui"
	"github.com/spf13/cobra"
)

// NewSyncCmd creates the sync command for updating all sources and regenerating agent files
func NewSyncCmd() *cobra.Command {
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Update all sources and regenerate agent files",
		Long: `Update all DNA sources and regenerate agent files in a single operation.

This command is a convenience wrapper that:
1. Updates all sources from their origins (dnaspec update --all)
2. Regenerates all agent files (dnaspec update-agents --no-ask)

The sync command is non-interactive and safe for CI/CD pipelines. It uses saved
agent configurations and does not prompt for user input. New guidelines are NOT
added automatically (--add-new=none policy).`,
		Example: `  # Sync all sources and regenerate agent files
  dnaspec sync

  # Preview what would change without writing files
  dnaspec sync --dry-run`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSync(dryRun)
		},
	}

	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Preview changes without writing files")

	return cmd
}

func runSync(dryRun bool) error {
	// Check project config exists
	if _, err := os.Stat(projectConfigFileName); os.IsNotExist(err) {
		fmt.Println(ui.ErrorStyle.Render("✗ Error:"), "No project configuration found")
		fmt.Println(ui.SubtleStyle.Render("  Run"), ui.CodeStyle.Render("dnaspec init"), ui.SubtleStyle.Render("first to initialize a project"))
		return fmt.Errorf("project not initialized")
	}

	// Load project config
	cfg, err := config.LoadProjectConfig(projectConfigFileName)
	if err != nil {
		return fmt.Errorf("failed to load project config: %w", err)
	}

	// Check if there are any sources
	if len(cfg.Sources) == 0 {
		fmt.Println("No sources configured")
		return nil
	}

	fmt.Println(ui.InfoStyle.Render("Syncing all DNA sources..."))
	fmt.Printf("Updating %d sources...\n\n", len(cfg.Sources))

	// Update each source individually (non-interactive mode)
	flags := updateFlags{
		dryRun: dryRun,
	}

	var errors []error
	for i := range cfg.Sources {
		fmt.Printf("=== Updating %s ===\n", cfg.Sources[i].Name)

		if err := updateSingleSource(cfg, cfg.Sources[i].Name, flags); err != nil {
			errors = append(errors, fmt.Errorf("%s: %w", cfg.Sources[i].Name, err))
			fmt.Println(ui.ErrorStyle.Render("✗ Failed:"), err)
		}

		fmt.Println()
	}

	if len(errors) > 0 {
		return fmt.Errorf("failed to update %d sources", len(errors))
	}

	fmt.Println(ui.SuccessStyle.Render("✓ All sources updated"))

	// If dry-run, don't regenerate agents
	if dryRun {
		fmt.Println()
		fmt.Println(ui.InfoStyle.Render("=== Dry Run - Preview ==="))
		fmt.Println("No changes made (dry run)")
		return nil
	}

	fmt.Println()
	fmt.Println(ui.InfoStyle.Render("Regenerating agent files..."))

	// Regenerate agent files using non-interactive mode
	// Save old noAskFlag value and restore after
	oldNoAskFlag := noAskFlag
	noAskFlag = true
	defer func() { noAskFlag = oldNoAskFlag }()

	if err := runUpdateAgents(nil, nil); err != nil {
		return fmt.Errorf("failed to regenerate agent files: %w", err)
	}

	fmt.Println()
	fmt.Println(ui.SuccessStyle.Render("✓ Sync complete"))

	return nil
}
