package project

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/aviator5/dnaspec/internal/core/agents"
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
- Agent IDs are recognized
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
	cfg, err := loadAndCheckConfig()
	if err != nil {
		return err
	}

	fmt.Println(ui.InfoStyle.Render("Validating"), ui.CodeStyle.Render(projectConfigFileName)+"...")

	// Collect all validation errors and warnings
	var errors []string
	var warnings []string
	var validatedFiles []string

	// Validate config version
	errors = validateConfigVersion(cfg, errors)

	// Validate sources
	fmt.Printf(ui.SuccessStyle.Render("✓")+" %d sources configured\n", len(cfg.Sources))
	errors, warnings, validatedFiles = validateAllSources(cfg.Sources, errors, warnings, validatedFiles)

	// Validate agent IDs
	errors = validateAgentIDs(cfg.Agents, errors)

	// Report results
	return reportValidationResults(errors, warnings, validatedFiles)
}

func loadAndCheckConfig() (*config.ProjectConfig, error) {
	if _, err := os.Stat(projectConfigFileName); os.IsNotExist(err) {
		fmt.Println(ui.ErrorStyle.Render("✗ Error:"), ui.CodeStyle.Render(projectConfigFileName), "not found")
		fmt.Println(ui.SubtleStyle.Render("  Run"), ui.CodeStyle.Render("dnaspec init"), ui.SubtleStyle.Render("first to initialize a project"))
		return nil, fmt.Errorf("project configuration not found")
	}

	cfg, err := config.LoadProjectConfig(projectConfigFileName)
	if err != nil {
		fmt.Println(ui.ErrorStyle.Render("✗ Error:"), "Failed to load configuration:", err)
		return nil, err
	}
	return cfg, nil
}

func validateConfigVersion(cfg *config.ProjectConfig, errors []string) []string {
	if cfg.Version != 1 {
		return append(errors, fmt.Sprintf("Unsupported config version: %d (only version 1 is supported)", cfg.Version))
	}
	fmt.Println(ui.SuccessStyle.Render("✓"), "YAML syntax valid")
	fmt.Println(ui.SuccessStyle.Render("✓"), "Version 1 schema valid")
	return errors
}

func validateAllSources(
	sources []config.ProjectSource,
	errors, warnings, validatedFiles []string,
) (outErrors, outWarnings, outValidatedFiles []string) {
	sourceNames := make(map[string]bool)
	for i := range sources {
		src := &sources[i]
		if sourceNames[src.Name] {
			errors = append(errors, fmt.Sprintf("Duplicate source name: '%s'", src.Name))
		}
		sourceNames[src.Name] = true

		sourceErrors, sourceWarnings, sourceFiles := validateSource(src)
		errors = append(errors, sourceErrors...)
		warnings = append(warnings, sourceWarnings...)
		validatedFiles = append(validatedFiles, sourceFiles...)
	}
	return errors, warnings, validatedFiles
}

func validateAgentIDs(agentIDs []string, errors []string) []string {
	availableAgents := agents.GetAvailableAgents()
	recognizedAgents := make(map[string]bool, len(availableAgents))
	agentNames := make([]string, 0, len(availableAgents))
	for _, agent := range availableAgents {
		recognizedAgents[agent.ID] = true
		agentNames = append(agentNames, agent.ID)
	}

	for _, agentID := range agentIDs {
		if !recognizedAgents[agentID] {
			errors = append(errors, fmt.Sprintf("Unknown agent ID: '%s' (recognized: %s)", agentID, formatList(agentNames)))
		}
	}

	if len(agentIDs) > 0 {
		fmt.Println(ui.SuccessStyle.Render("✓"), "All agent IDs recognized:", formatList(agentIDs))
	}
	return errors
}

func reportValidationResults(errors, warnings, validatedFiles []string) error {
	if len(errors) == 0 {
		printSuccessResults(validatedFiles, warnings)
		return nil
	}

	fmt.Println()
	fmt.Println(ui.ErrorStyle.Render("✗"), "Validation found", len(errors), "errors:")
	for _, err := range errors {
		fmt.Println("  -", err)
	}
	return fmt.Errorf("validation failed")
}

func printSuccessResults(validatedFiles, warnings []string) {
	fmt.Println(ui.SuccessStyle.Render("✓"), "All referenced files exist:")
	for _, file := range validatedFiles {
		fmt.Println("  -", ui.CodeStyle.Render(file))
	}

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
						"    Consider manually editing dnaspec.yaml to use a relative path",
					src.Name, src.Path,
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
