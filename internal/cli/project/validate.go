package project

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/aviator5/dnaspec/internal/core/config"
	"github.com/aviator5/dnaspec/internal/core/paths"
	"github.com/aviator5/dnaspec/internal/ui"
	"github.com/spf13/cobra"
)

// NewValidateCmd creates the validate command for validating project configuration
func NewValidateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate the project configuration",
		Long: `Validate the dnaspec.yaml file in the current directory.

This command checks:
- YAML syntax and schema structure
- Config version is supported (currently version 1)
- All sources have required fields
- File references exist in dnaspec/ directory (guidelines and prompts)
- Agent IDs are recognized (claude-code, github-copilot)
- No duplicate source names
- Symlinked sources with missing paths (warning only)`,
		Example: `  # Validate the project configuration
  dnaspec validate`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runValidate()
		},
	}

	return cmd
}

func runValidate() error {
	// Check if config exists
	if _, err := os.Stat(projectConfigFileName); os.IsNotExist(err) {
		fmt.Println(ui.ErrorStyle.Render("✗ Error:"), ui.CodeStyle.Render(projectConfigFileName), "not found")
		fmt.Println(ui.SubtleStyle.Render("  Run"), ui.CodeStyle.Render("dnaspec init"), ui.SubtleStyle.Render("first to initialize a project"))
		return fmt.Errorf("project configuration not found")
	}

	// Load the configuration
	cfg, err := config.LoadProjectConfig(projectConfigFileName)
	if err != nil {
		fmt.Println(ui.ErrorStyle.Render("✗ Error:"), "Failed to load configuration:", err)
		return err
	}

	fmt.Println(ui.InfoStyle.Render("Validating"), ui.CodeStyle.Render(projectConfigFileName)+"...")

	// Collect all validation errors and warnings
	var errors []string
	var warnings []string
	var validatedFiles []string

	// Validate config version
	if cfg.Version != 1 {
		errors = append(errors, fmt.Sprintf("Unsupported config version: %d (only version 1 is supported)", cfg.Version))
	} else {
		fmt.Println(ui.SuccessStyle.Render("✓"), "YAML syntax valid")
		fmt.Println(ui.SuccessStyle.Render("✓"), "Version 1 schema valid")
	}

	// Track source names for duplicate detection
	sourceNames := make(map[string]bool)

	// Validate sources
	fmt.Printf(ui.SuccessStyle.Render("✓")+" %d sources configured\n", len(cfg.Sources))

	for i := range cfg.Sources {
		src := &cfg.Sources[i]
		// Check for duplicate source names
		if sourceNames[src.Name] {
			errors = append(errors, fmt.Sprintf("Duplicate source name: '%s'", src.Name))
		}
		sourceNames[src.Name] = true

		sourceErrors, sourceWarnings, sourceFiles := validateSource(src)
		errors = append(errors, sourceErrors...)
		warnings = append(warnings, sourceWarnings...)
		validatedFiles = append(validatedFiles, sourceFiles...)
	}

	// Validate agent IDs
	recognizedAgents := map[string]bool{
		"claude-code":    true,
		"github-copilot": true,
	}

	for _, agentID := range cfg.Agents {
		if !recognizedAgents[agentID] {
			errors = append(errors, fmt.Sprintf("Unknown agent ID: '%s' (recognized: claude-code, github-copilot)", agentID))
		}
	}

	if len(cfg.Agents) > 0 {
		fmt.Println(ui.SuccessStyle.Render("✓"), "All agent IDs recognized:", formatList(cfg.Agents))
	}

	// Report results
	if len(errors) == 0 {
		// Success case
		fmt.Println(ui.SuccessStyle.Render("✓"), "All referenced files exist:")
		for _, file := range validatedFiles {
			fmt.Println("  -", ui.CodeStyle.Render(file))
		}

		// Display warnings if any
		if len(warnings) > 0 {
			fmt.Println()
			fmt.Println(ui.WarningStyle.Render("⚠"), "Found", len(warnings), "warning(s):")
			for _, warning := range warnings {
				fmt.Println("  -", warning)
			}
		}

		fmt.Println()
		if len(warnings) > 0 {
			fmt.Println(ui.SuccessStyle.Render("✓ Configuration is valid (with warnings)"))
		} else {
			fmt.Println(ui.SuccessStyle.Render("✓ Configuration is valid"))
		}
		return nil
	}

	// Failure case
	fmt.Println()
	fmt.Println(ui.ErrorStyle.Render("✗"), "Validation found", len(errors), "errors:")
	for _, err := range errors {
		fmt.Println("  -", err)
	}

	return fmt.Errorf("validation failed")
}

func validateSource(src *config.ProjectSource) (errors []string, warnings []string, validatedFiles []string) {
	// Check required fields based on source type
	if src.Name == "" {
		errors = append(errors, "Source missing required field: name")
	}

	if src.Type == "" {
		errors = append(errors, fmt.Sprintf("Source '%s' missing required field: type", src.Name))
	}

	switch src.Type {
	case config.SourceTypeGitRepo:
		if src.URL == "" {
			errors = append(errors, fmt.Sprintf("Source '%s' (%s) missing required field: url", src.Name, config.SourceTypeGitRepo))
		}
		if src.Commit == "" {
			errors = append(errors, fmt.Sprintf("Source '%s' (%s) missing required field: commit", src.Name, config.SourceTypeGitRepo))
		}
	case config.SourceTypeLocalPath:
		if src.Path == "" {
			errors = append(errors, fmt.Sprintf("Source '%s' (%s) missing required field: path", src.Name, config.SourceTypeLocalPath))
		} else {
			// Warn on absolute paths (not error - maintains backward compatibility)
			if filepath.IsAbs(src.Path) {
				warnings = append(warnings, fmt.Sprintf(
					"Source '%s' uses absolute path: %s\n"+
						"    Run 'dnaspec update %s' to auto-convert, or manually edit dnaspec.yaml",
					src.Name, src.Path, src.Name,
				))
			} else {
				// Validate relative path resolves within project
				projectRoot, err := filepath.Abs(filepath.Dir(projectConfigFileName))
				if err != nil {
					errors = append(errors, fmt.Sprintf("Failed to resolve project root: %v", err))
				} else if err := paths.ValidateLocalPath(projectRoot, src.Path); err != nil {
					errors = append(errors, fmt.Sprintf(
						"Source '%s' path validation failed: %v",
						src.Name, err,
					))
				}
			}
		}
	}

	// Validate guideline file references
	for _, guideline := range src.Guidelines {
		filePath := filepath.Join("dnaspec", src.Name, guideline.File)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			errors = append(errors, fmt.Sprintf("File not found: %s", filePath))
		} else {
			validatedFiles = append(validatedFiles, filePath)
		}
	}

	// Validate prompt file references
	for _, prompt := range src.Prompts {
		filePath := filepath.Join("dnaspec", src.Name, prompt.File)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			errors = append(errors, fmt.Sprintf("File not found: %s", filePath))
		} else {
			validatedFiles = append(validatedFiles, filePath)
		}
	}

	return errors, warnings, validatedFiles
}

// Helper function to format list
func formatList(items []string) string {
	if len(items) == 0 {
		return ""
	}
	if len(items) == 1 {
		return items[0]
	}
	result := ""
	for i, item := range items {
		if i > 0 {
			result += ", "
		}
		result += item
	}
	return result
}
