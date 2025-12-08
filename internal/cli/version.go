package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// Version is the current version of the application
	// This will be set via ldflags during build
	Version = "dev"
	// Commit is the git commit hash
	Commit = "none"
	// Date is the build date
	Date = "unknown"
)

// NewVersionCmd creates the version command
func NewVersionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version information",
		Long:  `Print the version information for dnaspec`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("dnaspec version %s\n", Version)
			fmt.Printf("  commit: %s\n", Commit)
			fmt.Printf("  built at: %s\n", Date)
		},
	}

	return cmd
}
