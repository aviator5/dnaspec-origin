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
	rootCmd.AddCommand(project.NewUpdateAgentsCmd())

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
