// Package commands provides the CLI commands for Nexo.
package commands

import (
	"fmt"
	"os"

	"github.com/abdul-hamid-achik/nexo/internal/version"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "nexo",
	Short: "Nexo - A file-system based Go framework",
	Long: `Nexo is a file-system based Go framework for building APIs and websites.
Inspired by Next.js App Router, it brings convention over configuration to Go.

Quick Start:
  nexo new myapp      Create a new Nexo project
  nexo dev            Start development server with hot reload
  nexo build          Build for production
  nexo routes         List all registered routes
  nexo openapi        Generate OpenAPI specifications
  nexo upgrade        Upgrade to the latest version

Documentation: https://github.com/abdul-hamid-achik/nexo`,
	Version: version.GetVersion(),
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "Output in JSON format (for automation and LLM agents)")

	// Commands
	rootCmd.AddCommand(newCmd)
	rootCmd.AddCommand(devCmd)
	rootCmd.AddCommand(buildCmd)
	rootCmd.AddCommand(routesCmd)
	rootCmd.AddCommand(tailwindCmd)
}
