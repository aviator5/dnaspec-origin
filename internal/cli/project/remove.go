package project

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
(Claude commands and Copilot prompts).

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

func runRemove(sourceName string, force bool) error {
	// Load project configuration
	cfg, err := config.LoadProjectConfig(projectConfigFileName)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println(ui.ErrorStyle.Render("✗ Error:"), "Project configuration not found:", ui.CodeStyle.Render(projectConfigFileName))
			fmt.Println(ui.SubtleStyle.Render("  Run"), ui.CodeStyle.Render("dnaspec init"), ui.SubtleStyle.Render("to create a new project configuration."))
			return fmt.Errorf("project configuration not found")
		}
		fmt.Println(ui.ErrorStyle.Render("✗ Error:"), "Failed to load project configuration:", err)
		return err
	}

	// Find source by name
	sourceIndex := -1
	for i, source := range cfg.Sources {
		if source.Name == sourceName {
			sourceIndex = i
			break
		}
	}

	if sourceIndex == -1 {
		fmt.Println(ui.ErrorStyle.Render("✗ Error:"), "Source not found:", ui.CodeStyle.Render(sourceName))
		if len(cfg.Sources) > 0 {
			fmt.Println(ui.SubtleStyle.Render("\nAvailable sources:"))
			for _, source := range cfg.Sources {
				fmt.Println("  -", source.Name)
			}
		} else {
			fmt.Println(ui.SubtleStyle.Render("\nNo sources configured."))
		}
		return fmt.Errorf("source not found: %s", sourceName)
	}

	// Display impact
	if err := displayImpact(sourceName); err != nil {
		return fmt.Errorf("failed to display impact: %w", err)
	}

	// Confirmation prompt (unless --force is set)
	if !force {
		fmt.Println()
		fmt.Print(ui.SubtleStyle.Render("This cannot be undone. Continue? [y/N]: "))

		reader := bufio.NewReader(os.Stdin)
		response, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}

		response = strings.TrimSpace(strings.ToLower(response))
		if response != "y" && response != "yes" {
			fmt.Println(ui.SubtleStyle.Render("\nCancelled. No changes made."))
			return nil
		}
	}

	fmt.Println()

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

func displayImpact(sourceName string) error {
	fmt.Println(ui.SubtleStyle.Render("\nThe following will be deleted:"))

	// Config entry
	fmt.Println("  - dnaspec.yaml entry for", ui.CodeStyle.Render(sourceName))

	// Source directory
	sourceDir := filepath.Join("dnaspec", sourceName)
	guidelineCount := 0
	promptCount := 0

	if info, err := os.Stat(sourceDir); err == nil && info.IsDir() {
		// Count guidelines (*.md files in root)
		guidelineFiles, err := filepath.Glob(filepath.Join(sourceDir, "*.md"))
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

	// Claude command files
	claudePattern := filepath.Join(".claude", "commands", "dnaspec", sourceName+"-*.md")
	claudeFiles, err := filepath.Glob(claudePattern)
	if err == nil && len(claudeFiles) > 0 {
		fmt.Printf("  - .claude/commands/dnaspec/%s-*.md (%d files)\n", sourceName, len(claudeFiles))
	}

	// Copilot prompt files
	copilotPattern := filepath.Join(".github", "prompts", "dnaspec-"+sourceName+"-*.prompt.md")
	copilotFiles, err := filepath.Glob(copilotPattern)
	if err == nil && len(copilotFiles) > 0 {
		fmt.Printf("  - .github/prompts/dnaspec-%s-*.prompt.md (%d files)\n", sourceName, len(copilotFiles))
	}

	return nil
}

func deleteGeneratedFiles(sourceName string) (int, error) {
	deletedCount := 0

	// Delete Claude command files
	claudePattern := filepath.Join(".claude", "commands", "dnaspec", sourceName+"-*.md")
	claudeFiles, err := filepath.Glob(claudePattern)
	if err == nil {
		for _, file := range claudeFiles {
			if err := os.Remove(file); err != nil && !os.IsNotExist(err) {
				return deletedCount, fmt.Errorf("failed to delete %s: %w", file, err)
			}
			deletedCount++
		}
	}

	// Delete Copilot prompt files
	copilotPattern := filepath.Join(".github", "prompts", "dnaspec-"+sourceName+"-*.prompt.md")
	copilotFiles, err := filepath.Glob(copilotPattern)
	if err == nil {
		for _, file := range copilotFiles {
			if err := os.Remove(file); err != nil && !os.IsNotExist(err) {
				return deletedCount, fmt.Errorf("failed to delete %s: %w", file, err)
			}
			deletedCount++
		}
	}

	return deletedCount, nil
}
