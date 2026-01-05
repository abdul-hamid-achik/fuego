package commands

import (
	"fmt"
	"strings"

	"github.com/abdul-hamid-achik/nexo/pkg/generator"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var generateRouteCmd = &cobra.Command{
	Use:   "route <path>",
	Short: "Generate a new route",
	Long: `Generate a new route file with handler functions.

The path supports dynamic segments:
  [param]      - Dynamic parameter (e.g., users/[id])
  [...param]   - Catch-all parameter (e.g., docs/[...slug])
  [[...param]] - Optional catch-all (e.g., shop/[[...categories]])
  (group)      - Route group (doesn't affect URL)

Examples:
  fuego generate route users              # GET /api/users
  fuego generate route users/[id]         # Dynamic route /api/users/:id
  fuego generate route posts/[...slug]    # Catch-all /api/posts/*
  fuego generate route users/[id] --methods GET,PUT,DELETE`,
	Args: cobra.ExactArgs(1),
	Run:  runGenerateRoute,
}

var (
	routeMethods string
	routeAppDir  string
)

func init() {
	generateRouteCmd.Flags().StringVarP(&routeMethods, "methods", "m", "GET", "HTTP methods (comma-separated: GET,POST,PUT,DELETE)")
	generateRouteCmd.Flags().StringVarP(&routeAppDir, "app-dir", "d", "app", "App directory")
	generateCmd.AddCommand(generateRouteCmd)
}

func runGenerateRoute(cmd *cobra.Command, args []string) {
	path := args[0]
	methods := strings.Split(strings.ToUpper(routeMethods), ",")

	// Trim whitespace from methods
	for i, m := range methods {
		methods[i] = strings.TrimSpace(m)
	}

	result, err := generator.GenerateRoute(generator.RouteConfig{
		Path:    path,
		Methods: methods,
		AppDir:  routeAppDir,
	})

	if err != nil {
		if jsonOutput {
			printJSONError(err)
		} else {
			red := color.New(color.FgRed).SprintFunc()
			fmt.Printf("  %s %v\n", red("Error:"), err)
		}
		return
	}

	if jsonOutput {
		printSuccess(GenerateOutput{
			Command: "generate route",
			Path:    path,
			Files:   result.Files,
			Pattern: result.Pattern,
			Methods: methods,
		})
		return
	}

	green := color.New(color.FgGreen).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()

	fmt.Printf("\n  %s Generated route\n\n", green("✓"))
	for _, f := range result.Files {
		fmt.Printf("    Created: %s\n", cyan(f))
	}
	fmt.Printf("    Pattern: %s\n", result.Pattern)
	fmt.Printf("    Methods: %s\n\n", strings.Join(methods, ", "))
}
