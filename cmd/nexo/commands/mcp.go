package commands

import "github.com/spf13/cobra"

var mcpCmd = &cobra.Command{
	Use:   "mcp",
	Short: "MCP (Model Context Protocol) server",
	Long: `MCP server for LLM agent integration.

The MCP server exposes Nexo operations as tools that LLM agents can use
to create projects, generate routes, list routes, and more.

Example:
  nexo mcp serve --workdir /path/to/project`,
}

func init() {
	rootCmd.AddCommand(mcpCmd)
}
