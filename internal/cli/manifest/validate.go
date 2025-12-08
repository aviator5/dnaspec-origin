package manifest

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/aviator5/dnaspec/internal/core/config"
	"github.com/aviator5/dnaspec/internal/core/validate"
	"github.com/aviator5/dnaspec/internal/ui"
)

// NewValidateCmd creates the manifest validate subcommand
func NewValidateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate the manifest file",
		Long: `Validate the dnaspec-manifest.yaml file in the current directory.

This command checks:
- Manifest structure and required fields
- Guideline and prompt definitions
- File references (files must exist)
- Cross-references (prompts referenced by guidelines must exist)
- Naming conventions (spinal-case)
- Path security (no absolute paths or path traversal)`,
		Example: `  # Validate the manifest in the current directory
  dnaspec manifest validate`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runValidate()
		},
	}

	return cmd
}

func runValidate() error {
	// Check if manifest exists
	if _, err := os.Stat(manifestFileName); os.IsNotExist(err) {
		fmt.Println(ui.ErrorStyle.Render("✗ Error:"), ui.CodeStyle.Render(manifestFileName), "not found")
		fmt.Println(
			ui.SubtleStyle.Render("  Run"), ui.CodeStyle.Render("dnaspec manifest init"),
			ui.SubtleStyle.Render("to create a new manifest."),
		)
		return fmt.Errorf("manifest file not found")
	}

	// Load the manifest
	manifest, err := config.LoadManifest(manifestFileName)
	if err != nil {
		fmt.Println(ui.ErrorStyle.Render("✗ Error:"), "Failed to load manifest:", err)
		return err
	}

	// Validate the manifest
	baseDir, _ := os.Getwd()
	errors := validate.ValidateManifest(manifest, baseDir)

	// Report results
	if errors.IsEmpty() {
		fmt.Println(ui.SuccessStyle.Render("✓ Manifest is valid"))
		return nil
	}

	// Display errors
	fmt.Println(ui.ErrorStyle.Render(fmt.Sprintf("✗ Found %d validation error(s):", len(errors))))
	fmt.Println()
	for _, err := range errors {
		fmt.Println("  •", ui.CodeStyle.Render(err.Field)+":", err.Message)
	}
	fmt.Println()
	fmt.Println(
		ui.SubtleStyle.Render("Fix these errors and run"),
		ui.CodeStyle.Render("dnaspec manifest validate"),
		ui.SubtleStyle.Render("again."),
	)

	return fmt.Errorf("validation failed")
}
