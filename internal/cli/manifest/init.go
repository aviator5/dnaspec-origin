package manifest

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/aviator5/dnaspec/internal/core/config"
	"github.com/aviator5/dnaspec/internal/ui"
)

const manifestFileName = "dnaspec-manifest.yaml"

// NewInitCmd creates the manifest init subcommand
func NewInitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a new manifest file",
		Long: `Create a new dnaspec-manifest.yaml file with example structure.

This command creates a manifest file in the current directory with example
guidelines and prompts sections, including helpful comments to guide you.`,
		Example: `  # Create a new manifest in the current directory
  dnaspec manifest init`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runInit()
		},
	}

	return cmd
}

func runInit() error {
	// Check if manifest already exists
	if _, err := os.Stat(manifestFileName); err == nil {
		fmt.Println(ui.ErrorStyle.Render("✗ Error:"), "Manifest file already exists:", ui.CodeStyle.Render(manifestFileName))
		fmt.Println(ui.SubtleStyle.Render("  To create a new manifest, first remove or rename the existing file."))
		return fmt.Errorf("manifest file already exists")
	}

	// Create the manifest file
	if err := config.CreateExampleManifest(manifestFileName); err != nil {
		fmt.Println(ui.ErrorStyle.Render("✗ Error:"), "Failed to create manifest:", err)
		return err
	}

	// Success message
	fmt.Println(ui.SuccessStyle.Render("✓ Success:"), "Created", ui.CodeStyle.Render(manifestFileName))
	fmt.Println()
	fmt.Println(ui.InfoStyle.Render("Next steps:"))
	fmt.Println("  1. Edit the manifest file to add your guidelines and prompts")
	fmt.Println(
		"  2. Create the referenced files in", ui.CodeStyle.Render("guidelines/"),
		"and", ui.CodeStyle.Render("prompts/"), "directories",
	)
	fmt.Println("  3. Run", ui.CodeStyle.Render("dnaspec manifest validate"), "to check your manifest")

	return nil
}
