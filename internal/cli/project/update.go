package project

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/aviator5/dnaspec/internal/core/config"
	"github.com/aviator5/dnaspec/internal/core/files"
	"github.com/aviator5/dnaspec/internal/core/paths"
	"github.com/aviator5/dnaspec/internal/core/source"
	"github.com/aviator5/dnaspec/internal/ui"
	"github.com/spf13/cobra"
)

type updateFlags struct {
	dryRun bool
}

// NewUpdateCmd creates the update command for updating DNA sources
func NewUpdateCmd() *cobra.Command {
	var flags updateFlags

	cmd := &cobra.Command{
		Use:   "update [source-name]",
		Short: "Update source from its origin with interactive guideline selection",
		Long: `Update a DNA source from its origin (git repository or local directory).

This command fetches the latest manifest from the source and presents an interactive
multi-select UI to choose which guidelines to keep, add, or remove.`,
		Example: `  # Update a source with interactive selection
  dnaspec update my-company-dna

  # Preview changes without writing
  dnaspec update my-company-dna --dry-run`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runUpdate(flags, args)
		},
	}

	cmd.Flags().BoolVar(&flags.dryRun, "dry-run", false, "Preview changes without writing files")

	return cmd
}

func runUpdate(flags updateFlags, args []string) error {
	// Validate source name is provided
	if len(args) == 0 {
		return fmt.Errorf("must specify a source name")
	}

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

	return updateSingleSource(cfg, args[0], flags)
}

func updateSingleSource(cfg *config.ProjectConfig, sourceName string, flags updateFlags) error {
	// Find source by name
	src := config.FindSourceByName(cfg, sourceName)
	if src == nil {
		fmt.Println(ui.ErrorStyle.Render("✗ Error:"), "Source not found:", sourceName)
		fmt.Println("\nAvailable sources:")
		for i := range cfg.Sources {
			fmt.Printf("  - %s\n", cfg.Sources[i].Name)
		}
		return fmt.Errorf("source not found")
	}

	// Fetch latest from origin
	sourceInfo, cleanup, upToDate, err := fetchAndCheckSource(src)
	if err != nil {
		return err
	}
	if cleanup != nil {
		defer cleanup()
	}
	if upToDate {
		fmt.Println(ui.SuccessStyle.Render("✓ Current commit:"), src.Commit[:8])
		fmt.Println(ui.SuccessStyle.Render("✓ Already at latest commit"))
		fmt.Println("\nAll guidelines up to date.")
		return nil
	}

	// Compare current vs latest
	comparison := config.CompareGuidelines(src.Guidelines, sourceInfo.Manifest.Guidelines)

	// Dry run check - show preview without interactive selection
	if flags.dryRun {
		fmt.Println(ui.InfoStyle.Render("\n=== Dry Run - Preview ==="))
		fmt.Println("\nAvailable guidelines in source:")
		for _, g := range sourceInfo.Manifest.Guidelines {
			fmt.Printf("  - %s: %s\n", g.Name, g.Description)
		}

		if len(comparison.Unchanged) > 0 || len(comparison.Updated) > 0 {
			fmt.Println("\nAlready in config:")
			for _, name := range comparison.Unchanged {
				fmt.Println(ui.SuccessStyle.Render("  ✓"), name)
			}
			for _, name := range comparison.Updated {
				fmt.Println(ui.SuccessStyle.Render("  ✓"), name, ui.SubtleStyle.Render("(updated)"))
			}
		}

		if len(comparison.Removed) > 0 {
			fmt.Println("\nOrphaned (in config but not in source):")
			for _, name := range comparison.Removed {
				fmt.Printf(ui.SubtleStyle.Render("  ⚠️ %s (no longer in manifest)\n"), name)
			}
		}

		fmt.Println("\nNo changes made (dry run)")
		return nil
	}

	// Interactive guideline selection
	// Build lists for selection
	var existingNames []string
	existingNames = append(existingNames, comparison.Unchanged...)
	existingNames = append(existingNames, comparison.Updated...)

	var orphanedGuidelines []config.ProjectGuideline
	for _, g := range src.Guidelines {
		for _, removedName := range comparison.Removed {
			if g.Name == removedName {
				orphanedGuidelines = append(orphanedGuidelines, g)
				break
			}
		}
	}

	// Call interactive selection
	selectedNames, err := ui.SelectGuidelinesWithStatus(
		sourceInfo.Manifest.Guidelines,
		existingNames,
		orphanedGuidelines,
	)
	if err != nil {
		return fmt.Errorf("guideline selection canceled or failed: %w", err)
	}

	// Apply selection
	return applyUpdate(cfg, src, sourceInfo, selectedNames)
}

func fetchAndCheckSource(src *config.ProjectSource) (info *source.SourceInfo, cleanup func(), upToDate bool, err error) {
	if src.Type == config.SourceTypeGitRepo {
		fmt.Println(ui.InfoStyle.Render("⏳ Fetching latest from"), src.URL+"...")
		info, cleanup, err = source.FetchGitSource(src.URL, src.Ref)
		if err != nil {
			return nil, nil, false, fmt.Errorf("failed to fetch git source: %w", err)
		}

		// Check if commit changed
		if info.Commit == src.Commit {
			return nil, cleanup, true, nil
		}

		fmt.Println(ui.SuccessStyle.Render("✓ Current commit:"), src.Commit[:8])
		fmt.Println(ui.SuccessStyle.Render("✓ Latest commit:"), info.Commit[:8], ui.SubtleStyle.Render("(changed)"))
		return info, cleanup, false, nil
	}

	// Local path source
	fmt.Println(ui.InfoStyle.Render("⏳ Refreshing from local directory..."))

	// Resolve relative path if necessary
	sourcePath := src.Path
	if !filepath.IsAbs(src.Path) {
		projectRoot, err := filepath.Abs(filepath.Dir(projectConfigFileName))
		if err != nil {
			return nil, nil, false, fmt.Errorf("failed to resolve project root: %w", err)
		}
		absPath, resolveErr := paths.ResolveRelative(projectRoot, src.Path)
		if resolveErr != nil {
			return nil, nil, false, fmt.Errorf("failed to resolve relative path %s: %w", src.Path, resolveErr)
		}
		sourcePath = absPath
	}

	info, err = source.FetchLocalSource(sourcePath)
	if err != nil {
		return nil, nil, false, fmt.Errorf("failed to fetch local source: %w", err)
	}
	return info, nil, false, nil
}

func applyUpdate(cfg *config.ProjectConfig, src *config.ProjectSource, sourceInfo *source.SourceInfo, selectedNames []string) error {
	// Update guidelines in config based on selection
	updatedSource := *src

	// Build updated guidelines list from selected names
	var updatedGuidelines []config.ProjectGuideline
	for _, name := range selectedNames {
		manifestGuideline := findManifestGuideline(sourceInfo.Manifest, name)
		if manifestGuideline != nil {
			updatedGuidelines = append(updatedGuidelines, config.ProjectGuideline{
				Name:                manifestGuideline.Name,
				File:                manifestGuideline.File,
				Description:         manifestGuideline.Description,
				ApplicableScenarios: manifestGuideline.ApplicableScenarios,
				Prompts:             manifestGuideline.Prompts,
			})
		}
	}

	updatedSource.Guidelines = updatedGuidelines

	// Extract and update prompts
	manifestGuidelines := make([]config.ManifestGuideline, 0, len(updatedGuidelines))
	for _, g := range updatedGuidelines {
		manifestGuidelines = append(manifestGuidelines, config.ManifestGuideline(g))
	}
	updatedSource.Prompts = config.ExtractReferencedPrompts(manifestGuidelines, sourceInfo.Manifest.Prompts)

	// Update commit hash for git sources
	if src.Type == "git-repo" {
		updatedSource.Commit = sourceInfo.Commit
	}

	// Copy files
	destDir := filepath.Join("dnaspec", src.Name)
	if err := files.CopyGuidelineFiles(sourceInfo.SourceDir, destDir, manifestGuidelines, sourceInfo.Manifest.Prompts); err != nil {
		return fmt.Errorf("failed to copy files: %w", err)
	}

	// Update config
	if err := config.UpdateSourceInConfig(cfg, src.Name, updatedSource); err != nil {
		return fmt.Errorf("failed to update source in config: %w", err)
	}

	// Save config
	if err := config.AtomicWriteProjectConfig(projectConfigFileName, cfg); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Println(ui.SuccessStyle.Render("\n✓ Updated"), ui.CodeStyle.Render(projectConfigFileName))
	fmt.Println(
		ui.SubtleStyle.Render("\nRun"), ui.CodeStyle.Render("dnaspec update-agents"),
		ui.SubtleStyle.Render("to regenerate agent files"),
	)

	return nil
}

// Helper functions

func findManifestGuideline(manifest *config.Manifest, name string) *config.ManifestGuideline {
	for i := range manifest.Guidelines {
		if manifest.Guidelines[i].Name == name {
			return &manifest.Guidelines[i]
		}
	}
	return nil
}

func promptYesNo(question string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(question + " [y/N]: ")
	response, err := reader.ReadString('\n')
	if err != nil {
		return false
	}
	response = strings.TrimSpace(strings.ToLower(response))
	return response == "y" || response == "yes"
}
