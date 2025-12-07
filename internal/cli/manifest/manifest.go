package manifest

import (
	"github.com/spf13/cobra"
)

// NewManifestCmd creates the manifest command group
func NewManifestCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "manifest",
		Short: "Manage DNA repository manifest files",
		Long: `Commands for creating and validating dnaspec-manifest.yaml files.

The manifest file defines the guidelines and prompts available in a DNA repository.`,
	}

	// Add subcommands
	cmd.AddCommand(NewInitCmd())
	cmd.AddCommand(NewValidateCmd())

	return cmd
}
