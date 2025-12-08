package project

import (
	"fmt"
	"os"

	"github.com/aviator5/dnaspec/internal/core/config"
	"github.com/aviator5/dnaspec/internal/ui"
	"github.com/spf13/cobra"
)

const projectConfigFileName = "dnaspec.yaml"

// NewInitCmd creates the init command for initializing a project
func NewInitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a new project configuration",
		Long: `Create a new dnaspec.yaml file in the current directory.

This command creates an empty project configuration file with example structure
and helpful comments to guide you in adding DNA sources to your project.`,
		Example: `  # Initialize a new project configuration
  dnaspec init

  # Then add DNA sources
  dnaspec add --git-repo https://github.com/company/dna`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runInit()
		},
	}

	return cmd
}

func runInit() error {
	// Check if config already exists
	if _, err := os.Stat(projectConfigFileName); err == nil {
		fmt.Println(ui.ErrorStyle.Render("✗ Error:"), "Project configuration already exists:", ui.CodeStyle.Render(projectConfigFileName))
		fmt.Println(ui.SubtleStyle.Render("  To create a new configuration, first remove or rename the existing file."))
		return fmt.Errorf("project configuration already exists")
	}

	// Create the project config file
	if err := config.CreateExampleProjectConfig(projectConfigFileName); err != nil {
		fmt.Println(ui.ErrorStyle.Render("✗ Error:"), "Failed to create project configuration:", err)
		return err
	}

	// Success message
	fmt.Println(ui.SuccessStyle.Render("✓ Success:"), "Created", ui.CodeStyle.Render(projectConfigFileName))
	fmt.Println()
	fmt.Println(ui.InfoStyle.Render("Next steps:"))
	fmt.Println("  1. Run", ui.CodeStyle.Render("dnaspec add"), "to add DNA sources (git repositories or local directories)")
	fmt.Println("  2. Select which guidelines to include from each source")
	fmt.Println()
	fmt.Println(ui.SubtleStyle.Render("Examples:"))
	fmt.Println("  ", ui.CodeStyle.Render("dnaspec add --git-repo https://github.com/company/dna"))
	fmt.Println("  ", ui.CodeStyle.Render("dnaspec add /path/to/local/dna"))

	return nil
}
