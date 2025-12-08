package cli

import (
	"github.com/spf13/cobra"
)

// NewRootCmd creates the root command
func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dnaspec",
		Short: "DNASpec - DNA repository management tool",
		Long: `DNASpec helps DNA repository maintainers create and validate manifest files,
and project developers integrate DNA guidelines into their projects.`,
	}

	return cmd
}
