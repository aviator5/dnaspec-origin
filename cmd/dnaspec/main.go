package main

import (
	"os"

	"github.com/aviator5/dnaspec/internal/cli"
	"github.com/aviator5/dnaspec/internal/cli/manifest"
)

func main() {
	rootCmd := cli.NewRootCmd()
	rootCmd.AddCommand(manifest.NewManifestCmd())

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
