package project

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/aviator5/dnaspec/internal/core/agents"
	"github.com/aviator5/dnaspec/internal/core/config"
	"github.com/aviator5/dnaspec/internal/ui"
	"github.com/spf13/cobra"
)

// NewRemoveCmd creates the remove command for removing DNA sources
func NewRemoveCmd() *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "remove <source-name>",
		Short: "Remove a DNA source from the project",
		Long: `Remove a DNA source from your project configuration.

This command removes the source from dnaspec.yaml, deletes the source
directory and all guideline files, and cleans up generated agent files
for all supported agents (Antigravity, Claude Code, Cursor, GitHub Copilot,
and Windsurf).

By default, this command will show what will be deleted and ask for
confirmation before proceeding. Use --force to skip the confirmation.`,
		Example: `  # Remove a source with confirmation
  dnaspec remove company-dna

  # Remove a source without confirmation
  dnaspec remove company-dna --force`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRemove(args[0], force)
		},
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false, "Skip confirmation prompt")

	return cmd
}

const responseYes = "yes"

func runRemove(sourceName string, force bool) error {
	cfg, sourceIndex, err := loadConfigAndFindSource(sourceName)
	if err != nil {
		return err
	}

	// Display impact
	displayImpact(sourceName)

	// Confirmation prompt (unless --force is set)
	if !force {
		confirmed, err := confirmRemoval()
		if err != nil {
			return err
		}
		if !confirmed {
			fmt.Println(ui.SubtleStyle.Render("\nCanceled. No changes made."))
			return nil
		}
	}

	fmt.Println()

	return performRemoval(cfg, sourceName, sourceIndex)
}

func loadConfigAndFindSource(sourceName string) (*config.ProjectConfig, int, error) {
	// Load project configuration
	cfg, err := config.LoadProjectConfig(projectConfigFileName)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println(ui.ErrorStyle.Render("✗ Error:"), "Project configuration not found:", ui.CodeStyle.Render(projectConfigFileName))
			fmt.Println(
				ui.SubtleStyle.Render("  Run"), ui.CodeStyle.Render("dnaspec init"),
				ui.SubtleStyle.Render("to create a new project configuration."),
			)
			return nil, -1, fmt.Errorf("project configuration not found")
		}
		fmt.Println(ui.ErrorStyle.Render("✗ Error:"), "Failed to load project configuration:", err)
		return nil, -1, err
	}

	// Find source by name
	sourceIndex := -1
	for i := range cfg.Sources {
		if cfg.Sources[i].Name == sourceName {
			sourceIndex = i
			break
		}
	}

	if sourceIndex == -1 {
		fmt.Println(ui.ErrorStyle.Render("✗ Error:"), "Source not found:", ui.CodeStyle.Render(sourceName))
		if len(cfg.Sources) > 0 {
			fmt.Println(ui.SubtleStyle.Render("\nAvailable sources:"))
			for i := range cfg.Sources {
				fmt.Println("  -", cfg.Sources[i].Name)
			}
		} else {
			fmt.Println(ui.SubtleStyle.Render("\nNo sources configured."))
		}
		return nil, -1, fmt.Errorf("source not found: %s", sourceName)
	}

	return cfg, sourceIndex, nil
}

func confirmRemoval() (bool, error) {
	fmt.Println()
	fmt.Print(ui.SubtleStyle.Render("This cannot be undone. Continue? [y/N]: "))

	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return false, fmt.Errorf("failed to read input: %w", err)
	}

	response = strings.TrimSpace(strings.ToLower(response))
	return response == "y" || response == responseYes, nil
}

func performRemoval(cfg *config.ProjectConfig, sourceName string, sourceIndex int) error {
	// Delete generated agent files
	deletedCount, err := deleteGeneratedFiles(sourceName)
	if err != nil {
		return fmt.Errorf("failed to delete generated files: %w", err)
	}

	// Delete source directory
	sourceDir := filepath.Join("dnaspec", sourceName)
	if err := os.RemoveAll(sourceDir); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete source directory %s: %w", sourceDir, err)
	}

	// Update configuration - remove source entry
	cfg.Sources = append(cfg.Sources[:sourceIndex], cfg.Sources[sourceIndex+1:]...)

	if err := config.AtomicWriteProjectConfig(projectConfigFileName, cfg); err != nil {
		fmt.Println(ui.ErrorStyle.Render("✗ Critical:"), "Failed to update configuration:", err)
		fmt.Println(ui.SubtleStyle.Render("  Files have been deleted but configuration update failed."))
		fmt.Println(ui.SubtleStyle.Render("  You may need to manually remove the source entry from"), ui.CodeStyle.Render(projectConfigFileName))
		return fmt.Errorf("failed to update config: %w", err)
	}

	// Success message
	fmt.Println(ui.SuccessStyle.Render("✓ Success:"), "Removed source", ui.CodeStyle.Render(sourceName))
	fmt.Println("  Cleaned up", deletedCount, "file(s)")
	fmt.Println()
	fmt.Println(ui.SubtleStyle.Render("Next steps:"))
	fmt.Println("  Run", ui.CodeStyle.Render("dnaspec update-agents"), "to regenerate AGENTS.md")

	return nil
}

func displayImpact(sourceName string) {
	fmt.Println(ui.SubtleStyle.Render("\nThe following will be deleted:"))

	// Config entry
	fmt.Println("  - dnaspec.yaml entry for", ui.CodeStyle.Render(sourceName))

	// Source directory
	sourceDir := filepath.Join("dnaspec", sourceName)
	guidelineCount := 0
	promptCount := 0

	if info, err := os.Stat(sourceDir); err == nil && info.IsDir() {
		// Count guidelines (files in guidelines/ subdirectory)
		guidelineFiles, err := filepath.Glob(filepath.Join(sourceDir, "guidelines", "*"))
		if err == nil {
			guidelineCount = len(guidelineFiles)
		}

		// Count prompts (files in prompts/ subdirectory)
		promptFiles, err := filepath.Glob(filepath.Join(sourceDir, "prompts", "*"))
		if err == nil {
			promptCount = len(promptFiles)
		}

		fmt.Printf("  - %s directory (%d guidelines, %d prompts)\n",
			ui.CodeStyle.Render(sourceDir), guidelineCount, promptCount)
	} else {
		fmt.Printf("  - %s directory (not found, will skip)\n", ui.CodeStyle.Render(sourceDir))
	}

	// Agent-generated files
	for _, pattern := range agents.AgentFilePatterns {
		globPattern := pattern.GetFilePatternForSource(sourceName)
		files, err := filepath.Glob(globPattern)
		if err == nil && len(files) > 0 {
			displayPattern := pattern.GetDisplayPatternForSource(sourceName)
			fmt.Printf("  - %s (%d files)\n", displayPattern, len(files))
		}
	}
}

func deleteGeneratedFiles(sourceName string) (int, error) {
	deletedCount := 0

	// Delete all agent-generated files
	for _, pattern := range agents.AgentFilePatterns {
		globPattern := pattern.GetFilePatternForSource(sourceName)
		files, err := filepath.Glob(globPattern)
		if err == nil {
			for _, file := range files {
				if err := os.Remove(file); err != nil && !os.IsNotExist(err) {
					return deletedCount, fmt.Errorf("failed to delete %s: %w", file, err)
				}
				deletedCount++
			}
		}
	}

	return deletedCount, nil
}
