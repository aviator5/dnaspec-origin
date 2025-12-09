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
	all    bool
	dryRun bool
	addNew string
}

// NewUpdateCmd creates the update command for updating DNA sources
func NewUpdateCmd() *cobra.Command {
	var flags updateFlags

	cmd := &cobra.Command{
		Use:   "update [source-name]",
		Short: "Update source(s) from their origin",
		Long: `Update DNA sources from their origins (git repositories or local directories).

This command fetches the latest manifest from the source, updates existing guidelines,
and optionally adds new guidelines or reports removed ones.`,
		Example: `  # Update a specific source
  dnaspec update my-company-dna

  # Update all sources
  dnaspec update --all

  # Preview changes without writing
  dnaspec update my-company-dna --dry-run

  # Add all new guidelines automatically
  dnaspec update my-company-dna --add-new=all

  # Skip new guidelines
  dnaspec update my-company-dna --add-new=none`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runUpdate(flags, args)
		},
	}

	cmd.Flags().BoolVar(&flags.all, "all", false, "Update all sources")
	cmd.Flags().BoolVar(&flags.dryRun, "dry-run", false, "Preview changes without writing files")
	cmd.Flags().StringVar(&flags.addNew, "add-new", "", "Policy for new guidelines (all|none)")

	return cmd
}

const (
	addNewAll  = "all"
	addNewNone = "none"
)

func runUpdate(flags updateFlags, args []string) error {
	// Validate flags
	if len(args) == 0 && !flags.all {
		return fmt.Errorf("must specify either a source name or --all flag")
	}

	if len(args) > 0 && flags.all {
		return fmt.Errorf("cannot specify both a source name and --all flag")
	}

	if flags.addNew != "" && flags.addNew != addNewAll && flags.addNew != addNewNone {
		return fmt.Errorf("--add-new must be either '%s' or '%s'", addNewAll, addNewNone)
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

	// Update all sources or single source
	if flags.all {
		return updateAllSources(cfg, flags)
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

	// Display updates
	if len(comparison.Updated) > 0 {
		fmt.Println("\nUpdated guidelines:")
		for _, name := range comparison.Updated {
			fmt.Println(ui.SuccessStyle.Render("  ✓"), name)
		}
	}

	// Handle new guidelines
	addedNew := handleNewGuidelines(sourceInfo.Manifest, comparison.New, flags.addNew)

	// Display removed guidelines
	if len(comparison.Removed) > 0 {
		fmt.Println("\nRemoved from source:")
		for _, name := range comparison.Removed {
			fmt.Printf(ui.SubtleStyle.Render("  - %s (no longer in manifest)\n"), name)
		}
	}

	// Dry run check
	if flags.dryRun {
		fmt.Println(ui.InfoStyle.Render("\n=== Dry Run - Preview ==="))
		fmt.Println("Would update:", len(comparison.Updated), "guidelines")
		fmt.Println("Would add:", len(addedNew), "guidelines")
		fmt.Println("Removed from source:", len(comparison.Removed), "guidelines")
		fmt.Println("\nNo changes made (dry run)")
		return nil
	}

	return applyUpdate(cfg, src, sourceInfo, addedNew)
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

func handleNewGuidelines(manifest *config.Manifest, newGuidelines []string, policy string) []string {
	if len(newGuidelines) == 0 {
		return nil
	}

	fmt.Println("\nNew guidelines available:")
	for _, name := range newGuidelines {
		guideline := findManifestGuideline(manifest, name)
		if guideline != nil {
			fmt.Printf("  - %s: %s\n", name, guideline.Description)
		}
	}

	if policy == "" {
		// Interactive mode
		if promptYesNo("\nAdd new guidelines?") {
			policy = addNewAll
		} else {
			policy = addNewNone
		}
	}

	if policy == addNewAll {
		return newGuidelines
	}
	return nil
}

func applyUpdate(cfg *config.ProjectConfig, src *config.ProjectSource, sourceInfo *source.SourceInfo, addedNew []string) error {
	// Update guidelines in config
	updatedSource := *src

	// Update existing and unchanged guidelines with latest metadata
	var updatedGuidelines []config.ProjectGuideline
	for _, g := range src.Guidelines {
		// Find in manifest
		manifestGuideline := findManifestGuideline(sourceInfo.Manifest, g.Name)
		if manifestGuideline != nil {
			// Update metadata from manifest
			updatedGuidelines = append(updatedGuidelines, config.ProjectGuideline{
				Name:                g.Name,
				File:                manifestGuideline.File,
				Description:         manifestGuideline.Description,
				ApplicableScenarios: manifestGuideline.ApplicableScenarios,
				Prompts:             manifestGuideline.Prompts,
			})
		}
	}

	// Add new guidelines
	for _, name := range addedNew {
		manifestGuideline := findManifestGuideline(sourceInfo.Manifest, name)
		if manifestGuideline != nil {
			updatedGuidelines = append(updatedGuidelines, config.ProjectGuideline{
				Name:                manifestGuideline.Name,
				File:                manifestGuideline.File,
				Description:         manifestGuideline.Description,
				ApplicableScenarios: manifestGuideline.ApplicableScenarios,
				Prompts:             manifestGuideline.Prompts,
			})
			fmt.Println(ui.SuccessStyle.Render("\n✓ Added"), name)
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

func updateAllSources(cfg *config.ProjectConfig, flags updateFlags) error {
	if len(cfg.Sources) == 0 {
		fmt.Println("No sources configured")
		return nil
	}

	fmt.Printf("Updating %d sources...\n\n", len(cfg.Sources))

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
