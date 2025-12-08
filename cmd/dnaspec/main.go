package main

import (
	"os"

	"github.com/aviator5/dnaspec/internal/cli"
	"github.com/aviator5/dnaspec/internal/cli/manifest"
	"github.com/aviator5/dnaspec/internal/cli/project"
)

func main() {
	rootCmd := cli.NewRootCmd()
	rootCmd.AddCommand(manifest.NewManifestCmd())
	rootCmd.AddCommand(project.NewInitCmd())
	rootCmd.AddCommand(project.NewAddCmd())
	rootCmd.AddCommand(project.NewUpdateCmd())
	rootCmd.AddCommand(project.NewUpdateAgentsCmd())
	rootCmd.AddCommand(project.NewListCmd())
	rootCmd.AddCommand(project.NewRemoveCmd())
	rootCmd.AddCommand(project.NewValidateCmd())
	rootCmd.AddCommand(project.NewSyncCmd())
	rootCmd.AddCommand(cli.NewVersionCmd())

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
