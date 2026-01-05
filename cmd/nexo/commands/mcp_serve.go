package commands

import (
	"fmt"
	"os"

	"github.com/abdul-hamid-achik/nexo/pkg/mcp"
	"github.com/spf13/cobra"
)

var mcpServeCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start MCP server over stdio",
	Long: `Start an MCP server over stdio for LLM agent integration.

The server exposes tools for:
  - Creating new Nexo projects
  - Generating routes, middleware, proxy, and pages
  - Listing routes and project info
  - Validating project structure

Usage with Claude Desktop (add to claude_desktop_config.json):

  {
    "mcpServers": {
      "nexo": {
        "command": "nexo",
        "args": ["mcp", "serve", "--workdir", "/path/to/project"]
      }
    }
  }

Usage with Claude Code or OpenCode:

  Configure your MCP settings to run:
    fuego mcp serve --workdir /path/to/project

Available tools:
  - nexo_new: Create a new Nexo project
  - nexo_generate_route: Generate a route file
  - nexo_generate_middleware: Generate middleware
  - nexo_generate_proxy: Generate proxy file
  - nexo_generate_page: Generate page template
  - nexo_list_routes: List all routes
  - nexo_info: Get project information
  - nexo_validate: Validate project structure`,
	Run: runMCPServe,
}

var mcpWorkdir string

func init() {
	mcpServeCmd.Flags().StringVarP(&mcpWorkdir, "workdir", "w", "", "Working directory for operations (default: current directory)")
	mcpCmd.AddCommand(mcpServeCmd)
}

func runMCPServe(cmd *cobra.Command, args []string) {
	workdir := mcpWorkdir
	if workdir == "" {
		var err error
		workdir, err = os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting working directory: %v\n", err)
			os.Exit(1)
		}
	}

	server := mcp.NewServer(workdir)
	if err := server.ServeStdio(); err != nil {
		fmt.Fprintf(os.Stderr, "MCP server error: %v\n", err)
		os.Exit(1)
	}
}
