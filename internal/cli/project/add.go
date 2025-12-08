package project

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/aviator5/dnaspec/internal/core/config"
	"github.com/aviator5/dnaspec/internal/core/files"
	"github.com/aviator5/dnaspec/internal/core/source"
	"github.com/aviator5/dnaspec/internal/ui"
)

type addFlags struct {
	gitRepo    string
	gitRef     string
	name       string
	all        bool
	guidelines []string
	dryRun     bool
}

// NewAddCmd creates the add command for adding DNA sources
func NewAddCmd() *cobra.Command {
	var flags addFlags

	cmd := &cobra.Command{
		Use:   "add [path]",
		Short: "Add a DNA source to the project",
		Long: `Add a DNA source (git repository or local directory) to your project.

This command fetches DNA guidelines from a source, lets you select which
guidelines to include, and copies them to your project's dnaspec/ directory.`,
		Example: `  # Add from git repository
  dnaspec add --git-repo https://github.com/company/dna

  # Add from git repository with specific tag
  dnaspec add --git-repo https://github.com/company/dna --git-ref v1.2.0

  # Add from local directory
  dnaspec add /path/to/local/dna

  # Add all guidelines without prompting
  dnaspec add --git-repo https://github.com/company/dna --all

  # Add specific guidelines
  dnaspec add --git-repo https://github.com/company/dna --guideline go-style --guideline rest-api

  # Specify custom source name
  dnaspec add --git-repo https://github.com/company/dna --name my-custom-name

  # Preview changes without writing
  dnaspec add --git-repo https://github.com/company/dna --dry-run`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAdd(flags, args)
		},
	}

	cmd.Flags().StringVar(&flags.gitRepo, "git-repo", "", "Git repository URL")
	cmd.Flags().StringVar(&flags.gitRef, "git-ref", "", "Git reference (branch, tag, or commit)")
	cmd.Flags().StringVar(&flags.name, "name", "", "Custom source name (auto-derived if not specified)")
	cmd.Flags().BoolVar(&flags.all, "all", false, "Add all guidelines without prompting")
	cmd.Flags().StringSliceVar(&flags.guidelines, "guideline", []string{}, "Add specific guideline by name (repeatable)")
	cmd.Flags().BoolVar(&flags.dryRun, "dry-run", false, "Preview changes without writing files")

	return cmd
}

func runAdd(flags addFlags, args []string) error {
	if err := validateAddFlags(flags, args); err != nil {
		return err
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

	// Fetch source
	sourceInfo, cleanup, err := fetchSource(flags, args)
	if err != nil {
		return err
	}
	defer cleanup()

	fmt.Println(ui.SuccessStyle.Render("✓"), "Source loaded successfully")

	// Select guidelines
	selectedGuidelines, err := selectGuidelines(flags, sourceInfo)
	if err != nil {
		return err
	}

	if len(selectedGuidelines) == 0 {
		fmt.Println(ui.SubtleStyle.Render("No guidelines selected. Exiting."))
		return nil
	}

	// Derive source name
	sourceName := flags.name
	if sourceName == "" {
		sourceName = config.DeriveSourceName(sourceInfo.URL, sourceInfo.Path)
	}

	// Check for duplicate source name
	for i := range cfg.Sources {
		if cfg.Sources[i].Name == sourceName {
			return fmt.Errorf("source with name '%s' already exists, use --name to specify a different name", sourceName)
		}
	}

	// Extract referenced prompts
	selectedPrompts := config.ExtractReferencedPrompts(selectedGuidelines, sourceInfo.Manifest.Prompts)

	// Build source entry
	newSource := config.ProjectSource{
		Name:       sourceName,
		Type:       sourceInfo.SourceType,
		URL:        sourceInfo.URL,
		Path:       sourceInfo.Path,
		Ref:        sourceInfo.Ref,
		Commit:     sourceInfo.Commit,
		Guidelines: config.ManifestGuidelinesToProject(selectedGuidelines),
		Prompts:    selectedPrompts,
	}

	// Preview mode
	if flags.dryRun {
		printDryRun(newSource)
		return nil
	}

	// Copy files
	destDir := filepath.Join("dnaspec", sourceName)
	fmt.Println(ui.InfoStyle.Render("⏳ Copying files to"), ui.CodeStyle.Render(destDir))

	if err := files.CopyGuidelineFiles(sourceInfo.SourceDir, destDir, selectedGuidelines, sourceInfo.Manifest.Prompts); err != nil {
		return fmt.Errorf("failed to copy files: %w", err)
	}

	// Update config
	if err := config.AddSource(cfg, newSource); err != nil {
		return fmt.Errorf("failed to add source: %w", err)
	}

	// Save config
	if err := config.AtomicWriteProjectConfig(projectConfigFileName, cfg); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	// Success message
	printSuccess(newSource, destDir)

	return nil
}

func validateAddFlags(flags addFlags, args []string) error {
	if flags.gitRepo == "" && len(args) == 0 {
		return fmt.Errorf("must specify either --git-repo or a local path")
	}

	if flags.gitRepo != "" && len(args) > 0 {
		return fmt.Errorf("cannot specify both --git-repo and a local path")
	}

	if flags.all && len(flags.guidelines) > 0 {
		return fmt.Errorf("cannot use both --all and --guideline flags")
	}
	return nil
}

func fetchSource(flags addFlags, args []string) (*source.SourceInfo, func(), error) {
	if flags.gitRepo != "" {
		fmt.Println(ui.InfoStyle.Render("⏳ Cloning repository..."))
		sourceInfo, cleanup, err := source.FetchGitSource(flags.gitRepo, flags.gitRef)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to fetch git source: %w", err)
		}
		return sourceInfo, cleanup, nil
	}

	localPath := args[0]
	fmt.Println(ui.InfoStyle.Render("⏳ Loading local source..."))
	sourceInfo, err := source.FetchLocalSource(localPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch local source: %w", err)
	}
	return sourceInfo, func() {}, nil
}

func selectGuidelines(flags addFlags, sourceInfo *source.SourceInfo) ([]config.ManifestGuideline, error) {
	if flags.all {
		selected := sourceInfo.Manifest.Guidelines
		fmt.Println(ui.InfoStyle.Render("ℹ"), "Selected all", len(selected), "guidelines")
		return selected, nil
	}

	if len(flags.guidelines) > 0 {
		selected, err := ui.SelectGuidelinesByName(sourceInfo.Manifest.Guidelines, flags.guidelines)
		if err != nil {
			return nil, fmt.Errorf("failed to select guidelines: %w", err)
		}
		fmt.Println(ui.InfoStyle.Render("ℹ"), "Selected", len(selected), "guidelines")
		return selected, nil
	}

	// Interactive selection
	selected, err := ui.SelectGuidelines(sourceInfo.Manifest.Guidelines)
	if err != nil {
		return nil, fmt.Errorf("failed to select guidelines: %w", err)
	}
	return selected, nil
}

func printDryRun(newSource config.ProjectSource) {
	fmt.Println()
	fmt.Println(ui.InfoStyle.Render("Dry run - would add source:"))
	fmt.Println("  Name:", ui.CodeStyle.Render(newSource.Name))
	fmt.Println("  Type:", newSource.Type)
	if newSource.URL != "" {
		fmt.Println("  URL:", newSource.URL)
	}
	if newSource.Path != "" {
		fmt.Println("  Path:", newSource.Path)
	}
	fmt.Println("  Guidelines:", len(newSource.Guidelines))
	fmt.Println("  Prompts:", len(newSource.Prompts))
}

func printSuccess(newSource config.ProjectSource, destDir string) {
	fmt.Println()
	fmt.Println(ui.SuccessStyle.Render("✓ Success:"), "Added source", ui.CodeStyle.Render(newSource.Name))
	fmt.Println("  Guidelines:", len(newSource.Guidelines))
	fmt.Println("  Prompts:", len(newSource.Prompts))
	fmt.Println("  Files copied to:", ui.CodeStyle.Render(destDir))
	fmt.Println()
	fmt.Println(ui.SubtleStyle.Render("Next steps:"))
	fmt.Println("  Run", ui.CodeStyle.Render("dnaspec update-agents"), "to configure AI agents")
}
